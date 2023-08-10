package events

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeMsgCatalogCreated, func() flows.Event { return &MsgCatalogCreatedEvent{} })
}

// TypeMsgCreated is a constant for incoming messages
const TypeMsgCatalogCreated string = "msg_catalog_created"

// MsgCatalogCreatedEvent events are created when an action wants to send a catalog msg to the current contact.
//
//	{
//	  "type": "msg_catalog_created",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "msg": {
//	    "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
//	    "channel": {"uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf", "name": "Twilio"},
//	    "urn": "tel:+12065551212",
//	    "text": "hi there",
//	    "attachments": ["image/jpeg:https://s3.amazon.com/mybucket/attachment.jpg"]
//	  }
//	}
//
// @event msg_catalog_created
type MsgCatalogCreatedEvent struct {
	baseEvent

	Msg *flows.MsgCatalog `json:"msg" validate:"required,dive"`
}

// NewMsgCreated creates a new outgoing msg event to a single contact
func NewMsgCatalogCreated(msg *flows.MsgCatalog) *MsgCatalogCreatedEvent {
	return &MsgCatalogCreatedEvent{
		baseEvent: newBaseEvent(TypeMsgCatalogCreated),
		Msg:       msg,
	}
}
