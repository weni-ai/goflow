package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
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
	Header        string        `json:"header,omitempty"`
	Body          string        `json:"body,omitempty"`
	Footer        string        `json:"footer,omitempty"`
	Products      []string      `json:"products,omitempty"`
	Action        string        `json:"action,omitempty"`
	Topic         string        `json:"topic,omitempty"`
	Smart         bool          `json:"smart"`
	ProductSearch string        `json:"product_search,omitempty"`
	TextLanguage  envs.Language `json:"text_language,omitempty"`
	ChannelUUID   uuids.UUID    `json:"channel_uuid,omitempty"`
}

func NewMsgCatalogParam(header string, body string, footer string, products []string, action string, topic string, smart bool, productSearch string, textLanguage envs.Language, channelUUID uuids.UUID) MsgCatalogParam {
	p := MsgCatalogParam{
		Header:        header,
		Body:          body,
		Footer:        footer,
		Products:      products,
		Action:        action,
		Topic:         topic,
		Smart:         smart,
		ProductSearch: productSearch,
		TextLanguage:  textLanguage,
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
