package static

import (
	"github.com/nyaruka/goflow/assets"
)

type MsgCatalog struct {
	UUID_ assets.MsgCatalogUUID `json:"uuid" validate:"required,uuid`
	Name_ string                `json:"name"`
	Type_ string                `json:"type"`
}

func NewMsgCatalog(uuid assets.MsgCatalogUUID, name string, type_ string) assets.MsgCatalog {
	return &MsgCatalog{
		UUID_: uuid,
		Name_: name,
		Type_: type_,
	}
}

func (mc *MsgCatalog) UUID() assets.MsgCatalogUUID { return mc.UUID_ }

func (mc *MsgCatalog) Name() string { return mc.Name_ }

func (mc *MsgCatalog) Type() string { return mc.Type_ }
