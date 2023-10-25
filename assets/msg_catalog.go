package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

//type MsgCatalogUUID uuids.UUID

type MsgCatalog interface {
	ChannelUUID() uuids.UUID
	Name() string
	Type() string
}

type MsgCatalogReference struct {
	UUID uuids.UUID `json:"uuid" validate:"required,uuid"`
	Name string     `json:"name"`
}

type MsgCatalogCallAction struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewMsgCatalogReference(uuid uuids.UUID, name string) *MsgCatalogReference {
	return &MsgCatalogReference{UUID: uuid, Name: name}
}

func (r *MsgCatalogReference) Type() string {
	return "msg_catalog"
}

func (r *MsgCatalogReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

func (r *MsgCatalogReference) Identity() string {
	return string(r.UUID)
}

func (r *MsgCatalogReference) Variable() bool {
	return false
}

func (r *MsgCatalogReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*MsgCatalogReference)(nil)
