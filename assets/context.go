package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

type Context interface {
	Context() string
	ChannelUUID() ChannelUUID
}

type ContextReference struct {
	Context string `json:"context"`
	UUID    string `json:"uuid"`
}

func NewContextReference(context string) *ContextReference {
	return &ContextReference{Context: context}
}

func (r *ContextReference) Type() string {
	return "context"
}

func (r *ContextReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

func (r *ContextReference) Identity() string {
	return string(r.UUID)
}

func (r *ContextReference) Variable() bool {
	return false
}

func (r *ContextReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,context=%s]", r.Type(), r.Identity(), r.Context)
}

var _ UUIDReference = (*ContextReference)(nil)
