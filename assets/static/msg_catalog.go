package static

import (
	"github.com/nyaruka/goflow/assets"
)

type MsgCatalog struct {
	UUID_        assets.MsgCatalogUUID `json:"uuid" validate:"required,uuid`
	Name_        string                `json:"name"`
	Type_        string                `json:"type"`
	ChannelUUID_ assets.ChannelUUID    `json:"channel_uuid" validate:"required,channel_uuid`
}

func NewMsgCatalog(uuid assets.MsgCatalogUUID, name string, type_ string, channelUUID assets.ChannelUUID) assets.MsgCatalog {
	return &MsgCatalog{
		UUID_:        uuid,
		Name_:        name,
		Type_:        type_,
		ChannelUUID_: channelUUID,
	}
}

func (mc *MsgCatalog) UUID() assets.MsgCatalogUUID { return mc.UUID_ }

func (mc *MsgCatalog) ChannelUUID() assets.ChannelUUID { return mc.ChannelUUID_ }

func (mc *MsgCatalog) Name() string { return mc.Name_ }

func (mc *MsgCatalog) Type() string { return mc.Type_ }
