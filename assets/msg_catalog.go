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
	ProductSearch   string     `json:"product_search,omitempty"`
	ChannelUUID     uuids.UUID `json:"channel_uuid,omitempty"`
	SearchType      string     `json:"search_type,omitempty"`
	SearchUrl       string     `json:"search_url,omitempty"`
	ApiType         string     `json:"api_type,omitempty"`
	PostalCode      string     `json:"postal_code,omitempty"`
	SellerId        string     `json:"seller_id,omitempty"`
	HasVtexAds      bool       `json:"vtex_ads,omitempty"`
	HideUnavailable bool       `json:"hide_unavailable"`
	Language        string     `json:"language"`
}

func NewMsgCatalogParam(productSearch string, channelUUID uuids.UUID, searchType string, searchUrl string, apiType string, postalCode string, sellerId string, hasVtexAds bool, hideUnavailable bool, language string) MsgCatalogParam {
	p := MsgCatalogParam{
		ProductSearch:   productSearch,
		ChannelUUID:     channelUUID,
		SearchType:      searchType,
		SearchUrl:       searchUrl,
		ApiType:         apiType,
		PostalCode:      postalCode,
		SellerId:        sellerId,
		HasVtexAds:      hasVtexAds,
		HideUnavailable: hideUnavailable,
		Language:        language,
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
