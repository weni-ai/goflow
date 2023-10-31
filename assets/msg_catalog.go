package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

type MsgCatalogUUID uuids.UUID

type MsgCatalog interface {
	UUID() MsgCatalogUUID
	Name() string
	Type() string
	ChannelUUID() ChannelUUID
}

type MsgCatalogReference struct {
	UUID MsgCatalogUUID `json:"uuid" validate:"required,uuid"`
	Name string         `json:"name"`
}

type MsgCatalogParam struct {
	ProductSearch string     `json:"product_search,omitempty"`
	ChannelUUID   uuids.UUID `json:"channel_uuid,omitempty"`
}

func NewMsgCatalogParam(productSearch string, channelUUID uuids.UUID) MsgCatalogParam {
	p := MsgCatalogParam{
		ProductSearch: productSearch,
		ChannelUUID:   channelUUID,
	}
	return p
}

func NewMsgCatalogReference(uuid MsgCatalogUUID, name string) *MsgCatalogReference {
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
