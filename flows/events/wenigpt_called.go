package events

import "github.com/nyaruka/goflow/flows"

func init() {
	registerType(TypeWeniGPTCalled, func() flows.Event { return &WeniGPTCalledEvent{} })
}

// TypeWeniGPTCalled is the type for our weniGPT events
const TypeWeniGPTCalled string = "wenigpt_called"

// WeniGPTCalledEvent events are created when a weniGPT is called. The event contains
// the URL and the status of the response, as well as a full dump of the
// request and response.
//
//	{
//	  "type": "weniGPT_called",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "url": "http://localhost:49998/?cmd=success",
//	  "status": "success",
//	  "status_code": 200,
//	  "elapsed_ms": 123,
//	  "retries": 0,
//	  "request": "GET /?format=json HTTP/1.1",
//	  "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}",
//	  "extraction": "valid"
//	}
//
// @event wenigpt_called
type WeniGPTCalledEvent struct {
	baseEvent

	*flows.HTTPTrace

	// Resthook   string     `json:"resthook,omitempty"`
	Extraction Extraction `json:"extraction"`
}

// NewMsgCreated creates a new outgoing msg event to a single contact
func NewWeniGPTCalled(call *flows.WeniGPTCall, status flows.CallStatus, resthook string) *WeniGPTCalledEvent {
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

	return &WeniGPTCalledEvent{
		baseEvent: newBaseEvent(TypeWeniGPTCalled),
		HTTPTrace: flows.NewHTTPTrace(call.Trace, status),
		// Resthook:   resthook,
		Extraction: extraction,
	}
}
