package actions

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeCallExternalService, func() flows.Action { return &CallExternalServiceAction{} })
}

const TypeCallExternalService string = "call_external_service"

type CallExternalServiceAction struct {
	baseAction
	onlineAction
	ExternalService *assets.ExternalServiceReference
	Input           string
	ResultName      string
}

func (a *CallExternalServiceAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	externalServices := run.Session().Assets().ExternalServices()
	externalService := externalServices.Get(a.ExternalService.UUID)
	return a.call(run, step, "", externalService, logEvent)
}

func (a *CallExternalServiceAction) call(run flows.FlowRun, step flows.Step, input string, externalService *flows.ExternalService, logEvent flows.EventCallback) error {
	return nil
}
