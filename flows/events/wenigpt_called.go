package events

import "github.com/nyaruka/goflow/flows"

func init() {
	registerType(TypeWeniGPTCalled, func() flows.Event { return &WebhookCalledEvent{} })
}

// TypeWeniGPTCalled is the type for our weniGPT events
const TypeWeniGPTCalled string = "wenigpt_called"

// NewMsgCreated creates a new outgoing msg event to a single contact
func NewWeniGPTCalled(call *flows.WeniGPTCall, status flows.CallStatus, resthook string) *WebhookCalledEvent {
	extraction := ExtractionNone
	if len(call.ResponseBody) > 0 {
		if len(call.ResponseJSON) > 0 {
			if call.ResponseCleaned {
				extraction = ExtractionCleaned
			} else {
				extraction = ExtractionValid
			}
		} else {
			extraction = ExtractionIgnored
		}
	}

	return &WebhookCalledEvent{
		baseEvent: newBaseEvent(TypeWeniGPTCalled),
		HTTPTrace: flows.NewHTTPTrace(call.Trace, status),
		// Resthook:   resthook,
		Extraction: extraction,
	}
}
