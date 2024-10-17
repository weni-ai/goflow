package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMetaCalled, func() flows.Event { return &MetaCalledEvent{} })
}

// TypeMetaCalled is our type for calling meta services
const TypeMetaCalled string = "meta_called"

// MetaCalledEvent events are created when a meta service is called.
//
//   {
//     "type": "meta_called",
//     "created_on": "2006-01-02T15:04:05Z",
//     "http_logs": [
//       {
//         "url": "https://graph.facebook.com/v21.0/me",
//         "status": "success",
//         "request": "GET /me HTTP/1.1",
//         "response": "HTTP/1.1 200 OK\r\n\r\n{\"status\":\"ok\"}",
//         "created_on": "2006-01-02T15:04:05Z",
//         "elapsed_ms": 123
//       }
//     ]
//   }
//
// @event meta_called
type MetaCalledEvent struct {
	baseEvent
	HTTPLogs []*flows.HTTPLog `json:"http_logs"`
}

// NewClassifierCalled returns a service called event for a classifier
func NewMetaCalled(httpLogs []*flows.HTTPLog) *MetaCalledEvent {
	return &MetaCalledEvent{
		baseEvent: newBaseEvent(TypeMetaCalled),
		HTTPLogs:  httpLogs,
	}
}
