package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMsgWppCreated, func() flows.Event { return &MsgWppCreatedEvent{} })
}

// TypeMsgWppCreated is a constant for incoming messages
const TypeMsgWppCreated string = "msg_wpp_created"

// MsgWppCreatedEvent events are created when an action wants to send a wpp msg to the current contact.
//
//	{
//	  "type": "msg_wpp_created",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "msg": {
//	    "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//	    "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//	    "urn": "tel:+12065551212",
//			"body": "body text",
//			"footer": "footer text",
//	  }
//	}
//
// @event msg_wpp_created
type MsgWppCreatedEvent struct {
	baseEvent

	Msg *flows.MsgWppOut `json:"msg" validate:"required,dive"`
}

// NewMsgWppCreated creates a new outgoing msg event to a single contact
func NewMsgWppCreated(msg *flows.MsgWppOut) *MsgWppCreatedEvent {
	return &MsgWppCreatedEvent{
		baseEvent: newBaseEvent(TypeMsgCatalogCreated),
		Msg:       msg,
	}
}
