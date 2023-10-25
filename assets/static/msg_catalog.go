package static

import (
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
)

type MsgCatalog struct {
	ChannelUUID_ uuids.UUID `json:"uuid" validate:"required,uuid`
	Name_        string     `json:"name"`
	Type_        string     `json:"type"`
}

func NewMsgCatalog(uuid uuids.UUID, name string, type_ string) assets.MsgCatalog {
	return &MsgCatalog{
		ChannelUUID_: uuid,
		Name_:        name,
		Type_:        type_,
	}
}

func (mc *MsgCatalog) ChannelUUID() uuids.UUID { return mc.ChannelUUID_ }

func (mc *MsgCatalog) Name() string { return mc.Name_ }

func (mc *MsgCatalog) Type() string { return mc.Type_ }
