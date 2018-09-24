package events

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	RegisterType(TypeInputLabelsAdded, func() flows.Event { return &InputLabelsAddedEvent{} })
}

// TypeInputLabelsAdded is the type of our add label action
const TypeInputLabelsAdded string = "input_labels_added"

// InputLabelsAddedEvent events are created when an action wants to add labels to the current input.
//
//   {
//     "type": "input_labels_added",
//     "created_on": "2006-01-02T15:04:05Z",
//     "input_uuid": "4aef4050-1895-4c80-999a-70368317a4f5",
//     "labels": [{"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Spam"}]
//   }
//
// @event input_labels_added
type InputLabelsAddedEvent struct {
	BaseEvent

	InputUUID flows.InputUUID          `json:"input_uuid" validate:"required,uuid4"`
	Labels    []*assets.LabelReference `json:"labels" validate:"required,min=1,dive"`
}

// NewInputLabelsAddedEvent returns a new add to group event
func NewInputLabelsAddedEvent(inputUUID flows.InputUUID, labels []*assets.LabelReference) *InputLabelsAddedEvent {
	return &InputLabelsAddedEvent{
		BaseEvent: NewBaseEvent(),
		InputUUID: inputUUID,
		Labels:    labels,
	}
}

// Type returns the type of this event
func (e *InputLabelsAddedEvent) Type() string { return TypeInputLabelsAdded }
