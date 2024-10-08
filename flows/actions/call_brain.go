package actions

import (
	"strings"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeCallBrain, func() flows.Action { return &CallBrainAction{} })
}

// TypeCallBrain is the type for the call brain action
const TypeCallBrain string = "call_brain"

type CallBrainAction struct {
	baseAction
	onlineAction
	Entry      string `json:"entry"`
	ResultName string `json:"result_name,omitempty"`
}

// NewCallBrain creates a new call brain action
func NewCallBrain(uuid flows.ActionUUID, entry string, resultName string) *CallBrainAction {
	return &CallBrainAction{
		baseAction: newBaseAction(TypeCallBrain, uuid),
		Entry:      entry,
		ResultName: resultName,
	}
}

// Validate validates our action is valid
func (a *CallBrainAction) Validate() error {
	return nil
}

// Execute runs this action
func (a *CallBrainAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	return a.call(run, step, logEvent)
}

// Execute runs this action
func (a *CallBrainAction) call(run flows.FlowRun, step flows.Step, logEvent flows.EventCallback) error {
	var entry = a.Entry
	if len(entry) == 0 {
		entry = "@input.text"
	}

	var attachments []string = nil
	if entry == "@input.text" {
		attachmentsString, _ := run.EvaluateTemplate("@input.attachments")
		trimmedString := strings.Trim(attachmentsString, "[]")
		attachments = strings.Split(trimmedString, ", ")
		if len(attachments) == 1 && strings.Trim(attachments[0], " ") == "" {
			attachments = nil
		}
	}

	evaluatedText, evaluatedAttachment, _ := a.evaluateMessage(run, nil, entry, attachments, nil, logEvent)

	contactURN := run.Contact().PreferredURN()
	svc, err := run.Session().Engine().Services().Brain(run.Session())
	if err != nil {
		logEvent(events.NewError(err))
		return nil
	}

	orgContext := run.Session().Assets().OrgContext()
	c := orgContext.GetProjectUUIDByChannelUUID()
	var projectUUID uuids.UUID
	if c != nil {
		projectUUID = c.OrgContext.ProjectUUID()
	}

	call, err := svc.Call(run.Session(), projectUUID, evaluatedText, contactURN.URN(), evaluatedAttachment)

	if err != nil {
		logEvent(events.NewError(err))
	}

	if call != nil {
		a.updateBrain(run, call)

		status := callStatusBrain(call, err)

		c := &flows.WebhookCall{
			Trace:           call.Trace,
			ResponseJSON:    call.ResponseBody,
			ResponseCleaned: false,
		}

		logEvent(events.NewWebhookCalled(c, status, ""))
	}

	return nil
}

// determines the brain status from the HTTP status code
func callStatusBrain(call *flows.BrainCall, err error) flows.CallStatus {
	if call.Response == nil || err != nil {
		return flows.CallStatusConnectionError
	}
	if call.Response.StatusCode/100 == 2 {
		return flows.CallStatusSuccess
	}

	return flows.CallStatusResponseError
}
