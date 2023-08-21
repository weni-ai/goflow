package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeZeroshotCalled, func() flows.Event { return &ZeroshotCalledEvent{} })
}

const TypeZeroshotCalled string = "zeroshot_called"

type ZeroshotCalledEvent struct {
	baseEvent

	*flows.HTTPTrace

	Resthook   string     `json:"resthook,omitempty"`
	Extraction Extraction `json:"extraction"`
}

func NewZeroshotCalled(call *flows.ZeroshotCall, status flows.CallStatus, resthook string) *ZeroshotCalledEvent {
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

	return &ZeroshotCalledEvent{
		baseEvent:  newBaseEvent(TypeZeroshotCalled),
		HTTPTrace:  flows.NewHTTPTrace(call.Trace, status),
		Resthook:   resthook,
		Extraction: extraction,
	}
}
