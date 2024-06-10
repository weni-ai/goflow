package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

var (
	codeActionsURL   = "http://localhost:8050"
	codeActionsToken = ""
)

func init() {
	registerType(TypeRunCodeAction, func() flows.Action { return &RunCodeAction{} })
	codeActionsURL = os.Getenv("MAILROOM_FLOWS_CODE_ACTIONS_URL")
	codeActionsToken = os.Getenv("MAILROOM_FLOWS_CODE_ACTIONS_TOKEN")
}

const TypeRunCodeAction = "run_code_action"

type CodeAction struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RunCodeAction struct {
	baseAction
	onlineAction

	CodeAction CodeAction `json:"codeaction"`
	ResultName string     `json:"result_name,omitempty"`
}

func NewRunCodeAction(uuid flows.ActionUUID, codeActionName string, codeActionID string, resultName string) *RunCodeAction {
	return &RunCodeAction{
		baseAction: newBaseAction(TypeRunCodeAction, uuid),
		CodeAction: CodeAction{Name: codeActionName, ID: codeActionID},
		ResultName: resultName,
	}
}

func (a *RunCodeAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {

	runActionURL, err := url.JoinPath(codeActionsURL, "/run/"+a.CodeAction.ID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, runActionURL, nil)
	if err != nil {
		logEvent(events.NewError(err))
	}
	req.Header.Add("Authorization", codeActionsToken)

	call, err := a.Call(run.Session(), req)
	if err != nil {
		logEvent(events.NewError(err))
	}

	if call != nil {
		status := callStatus(call, err, false)
		logEvent(events.NewWebhookCalled(call, status, ""))

		if a.ResultName != "" {
			result, err := parseResult(call.ResponseJSON)
			if err != nil {
				return err
			}
			a.saveResult(run, step, a.ResultName, result.Result.Result, CategorySuccess, "", "", call.ResponseJSON, logEvent)
		}
	} else {
		a.saveResult(run, step, a.ResultName, fmt.Sprintf("%s", err), CategoryFailure, "", "", nil, logEvent)
	}

	return nil
}

func (a *RunCodeAction) Call(session flows.Session, req *http.Request) (*flows.WebhookCall, error) {
	httpRetries := httpx.NewExponentialRetries(
		time.Duration(5000)*time.Millisecond,
		2,
		0.5,
	)
	trace, err := httpx.DoTrace(http.DefaultClient, req, httpRetries, nil, (1024 * 1024))
	if trace != nil {
		call := &flows.WebhookCall{Trace: trace}
		if trace.Response == nil {
			return call, nil
		}

		if len(call.ResponseBody) > 0 {
			call.ResponseJSON, call.ResponseCleaned = ExtractJSON(call.ResponseBody)
		}
		return call, err
	}
	return nil, err
}

func (a *RunCodeAction) Results(include func(*flows.ResultInfo)) {
	include(flows.NewResultInfo(a.ResultName, []string{CategorySuccess, CategoryFailure}))
}

func ExtractJSON(body []byte) ([]byte, bool) {
	// we make a best effort to turn the body into JSON, so we strip out:
	//  1. any invalid UTF-8 sequences
	//  2. null chars
	//  3. escaped null chars (\u0000)
	cleaned := bytes.ToValidUTF8(body, nil)
	cleaned = bytes.ReplaceAll(cleaned, []byte{0}, nil)
	cleaned = utils.ReplaceEscapedNulls(cleaned, nil)

	if json.Valid(cleaned) {
		changed := !bytes.Equal(body, cleaned)
		return cleaned, changed
	}
	return nil, false
}

func parseResult(responseJSON []byte) (*RunResult, error) {
	result := &RunResult{}
	err := json.Unmarshal(responseJSON, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type RunResult struct {
	CodeID string `json:"code_id"`
	Result struct {
		ID        string `json:"_id"`
		CodeID    string `json:"code_id"`
		Status    string `json:"status"`
		Result    string `json:"result"`
		CretedAt  string `json:"creted_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"result"`
}
