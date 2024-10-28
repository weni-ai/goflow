package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// TopicUUID is the UUID of a topic
type TopicUUID uuids.UUID

// QueueUUID is the UUID of a queue related to wenichats
type QueueUUID string

// Topic categorizes tickets
//
//	{
//	  "uuid": "cd48bd11-08b9-44e3-9778-8e26adf08a7a",
//	  "name": "Weather",
//	  "queue_uuid": "3108d8c2-2a93-4db8-b7c1-d4b1b9c812b3"
//	}
//
// @asset topic
type Topic interface {
	UUID() TopicUUID
	Name() string
	QueueUUID() QueueUUID
}

// TopicReference is used to reference a topic
type TopicReference struct {
	UUID      TopicUUID `json:"uuid" validate:"required,uuid"`
	Name      string    `json:"name"`
	QueueUUID QueueUUID `json:"queue_uuid"`
}

// NewTopicReference creates a new topic reference with the given UUID and name
func NewTopicReference(uuid TopicUUID, name string, queueUUID QueueUUID) *TopicReference {
	return &TopicReference{UUID: uuid, Name: name, QueueUUID: queueUUID}
}

// Type returns the name of the asset type
func (r *TopicReference) Type() string {
	return "topic"
}

// GenericUUID returns the untyped UUID
func (r *TopicReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *TopicReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *TopicReference) Variable() bool {
	return false
}

func (r *TopicReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*TopicReference)(nil)
