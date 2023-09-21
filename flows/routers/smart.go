package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
)

func init() {
	registerType(TypeSmart, readSmartRouter)
}

// TypeSmart is the constant for our smart router
const TypeSmart string = "smart"

// Regex pattern for category and argument
const CategoryRegex string = `^[A-Za-zÀ-ÖØ-öø-ÿ]+$`
const ArgumentsRegex string = `^[A-Za-zÀ-ÖØ-öø-ÿ@.,]+$`

var apiUrl = "https://api.bothub.it"

// SmartCase represents a single case and test in our smart
type SmartCase struct {
	UUID         uuids.UUID         `json:"uuid"                   validate:"required"`
	Type         string             `json:"type"                   validate:"required"`
	Arguments    []string           `json:"arguments,omitempty"    engine:"localized,evaluated"`
	CategoryUUID flows.CategoryUUID `json:"category_uuid"          validate:"required"`
}

// NewSmartCase creates a new smart case
func NewSmartCase(uuid uuids.UUID, type_ string, arguments []string, categoryUUID flows.CategoryUUID) *SmartCase {
	return &SmartCase{
		UUID:         uuid,
		Type:         type_,
		Arguments:    arguments,
		CategoryUUID: categoryUUID,
	}
}

// LocalizationUUID gets the UUID which identifies this object for localization
func (c *SmartCase) LocalizationUUID() uuids.UUID { return uuids.UUID(c.UUID) }

// SmartRouter is a router that allows you to specify 0-n cases that will be sent to Zeroshot's API,
// following whichever case the API returns as a response, or if none do, then taking the default category
type SmartRouter struct {
	baseRouter

	operand             string
	cases               []*SmartCase
	defaultCategoryUUID flows.CategoryUUID
}

// NewSmart creates a new smart router
func NewSmart(wait flows.Wait, resultName string, categories []flows.Category, operand string, cases []*SmartCase, defaultCategoryUUID flows.CategoryUUID) *SmartRouter {
	return &SmartRouter{
		baseRouter:          newBaseRouter(TypeSmart, wait, resultName, categories),
		defaultCategoryUUID: defaultCategoryUUID,
		operand:             operand,
		cases:               cases,
	}
}

// SmartCases returns the cases for this smart router
func (r *SmartRouter) SmartCases() []*SmartCase { return r.cases }

// Validate validates the arguments for this router
func (r *SmartRouter) Validate(flow flows.Flow, exits []flows.Exit) error {
	// check the default category is valid
	if r.defaultCategoryUUID != "" && !r.isValidCategory(r.defaultCategoryUUID) {
		return errors.Errorf("default category %s is not a valid category", r.defaultCategoryUUID)
	}

	for _, c := range r.cases {
		// check each case points to a valid category
		if !r.isValidCategory(c.CategoryUUID) {
			return errors.Errorf("case category %s is not a valid category", c.CategoryUUID)
		}

		// and each case test is valid
		if c.Type != "has_any_word" && c.Type != "has_category" {
			return errors.Errorf("case must be of type 'has_any_words', not %s", c.Type)
		}
	}

	return r.validate(flow, exits)
}

// Route determines which exit to take from a node
func (r *SmartRouter) Route(run flows.FlowRun, step flows.Step, logEvent flows.EventCallback) (flows.ExitUUID, string, error) {
	env := run.Environment()

	// first evaluate our operand
	operand, err := run.EvaluateTemplateValue(r.operand)
	if err != nil {
		run.LogError(step, err)
	}

	var operandAsStr string

	if operand != nil {
		asText, _ := types.ToXText(env, operand)
		operandAsStr = asText.Native()
	}

	// classify text between categories
	categoryName, categoryUUID, err := r.classifyText(run, step, operandAsStr, logEvent)
	if err != nil && r.defaultCategoryUUID != "" {
		categoryName = "Failure"
		categoryUUID = r.defaultCategoryUUID
	} else if categoryUUID == "" && r.defaultCategoryUUID != "" {
		categoryName = "All Responses"
		categoryUUID = r.defaultCategoryUUID
	}

	exit, err := r.routeToCategory(run, step, categoryUUID, categoryName, operandAsStr, nil, logEvent)
	return exit, operandAsStr, err
}

var token string

func SetToken(t string) {
	token = t
}

func SetAPIURL(url string) {
	apiUrl = url
}

func (r *SmartRouter) classifyText(run flows.FlowRun, step flows.Step, operand string, logEvent flows.EventCallback) (string, flows.CategoryUUID, error) {
	if len(r.categories) == 1 && len(r.cases) == 0 {
		return "", "", nil
	}

	url := apiUrl + "/v2/repository/nlp/zeroshot/zeroshot-fast-predict"
	status := flows.CallStatusSuccess
	body := struct {
		Text       string `json:"text"`
		Categories []struct {
			Option   string   `json:"option"`
			Synonyms []string `json:"synonyms"`
		} `json:"categories"`
		Token string `json:"token"`
	}{
		Text: operand,
	}

	args := make(map[string][]string)
	for _, c := range r.cases {
		if c.Type == "has_category" {
			continue
		}
		var evaluatedArgs []string
		localizedArgs, _ := run.GetTextArray(c.UUID, "arguments", c.Arguments)
		for i := range c.Arguments {
			test := localizedArgs[i]
			arg, err := run.EvaluateTemplateValue(test)
			if err != nil {
				run.LogError(step, err)
			}

			resultAsStr, xerr := types.ToXText(run.Environment(), arg)
			if xerr != nil {
				run.LogError(step, xerr)
			}

			results := customSplit(resultAsStr.Native())

			for _, result := range results {
				arg := strings.TrimSpace(result)
				if res, err := RegexMatch(arg, ArgumentsRegex, false); res && err == nil {
					evaluatedArgs = append(evaluatedArgs, arg)
				} else {
					run.LogError(step, err)
				}
			}
		}

		args[string(c.CategoryUUID)] = evaluatedArgs
	}

	for category, arg := range args {
		for _, c := range r.categories {
			if string(c.UUID()) == category {
				if res, err := RegexMatch(c.Name(), CategoryRegex, true); res && err == nil {
					body.Categories = append(body.Categories, struct {
						Option   string   "json:\"option\""
						Synonyms []string "json:\"synonyms\""
					}{Option: c.Name(), Synonyms: arg})
					break
				} else {
					run.LogError(step, err)
				}
			}
		}
	}

	if token != "" {
		body.Token = token
	} else {
		run.LogError(step, fmt.Errorf("validation token cannot be empty"))
		status = flows.CallStatusConnectionError
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		run.LogError(step, err)
	}
	// build our request
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJSON))
	if err != nil {
		run.LogError(step, err)
		status = flows.CallStatusConnectionError
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response := &struct {
		Output struct {
			Classification string `json:"classification"`
			Other          bool   `json:"other"`
		}
	}{Output: struct {
		Classification string "json:\"classification\""
		Other          bool   "json:\"other\""
	}{Classification: ""}}

	trace, err := httpx.DoTrace(client, req, nil, nil, -1)
	call := &flows.WebhookCall{
		Trace:           trace,
		ResponseJSON:    trace.ResponseBody,
		ResponseCleaned: false,
	}

	if trace.Response.StatusCode >= 400 {
		run.LogError(step, err)
		status = flows.CallStatusConnectionError
		logEvent(events.NewWebhookCalled(call, status, ""))
		return "", "", err
	} else {
		logEvent(events.NewWebhookCalled(call, status, ""))
	}

	err = jsonx.Unmarshal(trace.ResponseBody, response)
	if err != nil {
		run.LogError(step, err)
	}

	var categoryUUID flows.CategoryUUID
	categoryUUID = ""

	// case with 'other' option
	if response.Output.Other {
		for _, c := range r.cases {
			if c.Type == "has_category" {
				categoryUUID = c.CategoryUUID
			}
		}
		return "Other", categoryUUID, nil
	}

	for _, category := range r.categories {
		if category.Name() == response.Output.Classification {
			categoryUUID = category.UUID()
		}
	}

	return response.Output.Classification, categoryUUID, nil

}

func RegexMatch(input string, regexPattern string, isCategory bool) (bool, error) {
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return false, err
	}

	if regex.MatchString(input) {
		if len(input) <= 20 && isCategory || !isCategory {
			return true, nil
		} else {
			return false, fmt.Errorf("error when combining input \"%s\" with regex", input)
		}
	} else {
		return false, fmt.Errorf("error when combining input \"%s\" with regex", input)
	}
}

func customSplit(input string) []string {
	return strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == ' '
	})
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (r *SmartRouter) EnumerateTemplates(localization flows.Localization, include func(envs.Language, string)) {
	include(envs.NilLanguage, r.operand)

	inspect.Templates(r.cases, localization, include)
}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (r *SmartRouter) EnumerateDependencies(localization flows.Localization, include func(envs.Language, assets.Reference)) {
	inspect.Dependencies(r.cases, localization, include)
}

// EnumerateLocalizables enumerates all the localizable text on this object
func (r *SmartRouter) EnumerateLocalizables(include func(uuids.UUID, string, []string, func([]string))) {
	inspect.LocalizableText(r.cases, include)

	r.baseRouter.EnumerateLocalizables(include)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type smartRouterEnvelope struct {
	baseRouterEnvelope

	Operand             string             `json:"operand"               validate:"required"`
	Cases               []*SmartCase       `json:"cases"`
	DefaultCategoryUUID flows.CategoryUUID `json:"default_category_uuid" validate:"omitempty,uuid4"`
}

func readSmartRouter(data json.RawMessage) (flows.Router, error) {
	e := &smartRouterEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &SmartRouter{
		operand:             e.Operand,
		cases:               e.Cases,
		defaultCategoryUUID: e.DefaultCategoryUUID,
	}

	if err := r.unmarshal(&e.baseRouterEnvelope); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *SmartRouter) MarshalJSON() ([]byte, error) {
	e := &smartRouterEnvelope{
		Operand:             r.operand,
		Cases:               r.cases,
		DefaultCategoryUUID: r.defaultCategoryUUID,
	}

	if err := r.marshal(&e.baseRouterEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
