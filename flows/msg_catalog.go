package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
)

type MsgCatalogOut struct {
	BaseMsg

	Header_        string        `json:"header,omitempty"`
	Body_          string        `json:"body,omitempty"`
	Footer_        string        `json:"footer,omitempty"`
	Products_      []string      `json:"products,omitempty"`
	Action_        string        `json:"action,omitempty"`
	Topic_         MsgTopic      `json:"topic,omitempty"`
	Smart_         bool          `json:"smart"`
	ProductSearch_ string        `json:"product_search,omitempty"`
	TextLanguage   envs.Language `json:"text_language,omitempty"`
}

type MsgCatalog struct {
	assets.MsgCatalog
}

func NewMsgCatalogOut(urn urns.URN, channel *assets.ChannelReference, header, body, footer, action, productSearch string, products []string, smart bool, topic MsgTopic) *MsgCatalogOut {
	return &MsgCatalogOut{
		BaseMsg: BaseMsg{
			UUID_:    MsgUUID(uuids.New()),
			URN_:     urn,
			Channel_: channel,
		},
		Header_:        header,
		Body_:          body,
		Footer_:        footer,
		Products_:      products,
		Action_:        action,
		Smart_:         smart,
		ProductSearch_: productSearch,
		Topic_:         topic,
	}
}

func NewMsgCatalog(asset assets.MsgCatalog) *MsgCatalog {
	return &MsgCatalog{
		MsgCatalog: asset,
	}
}

func (s *MsgCatalogAssets) Get(uuid assets.MsgCatalogUUID) *MsgCatalog {
	return s.byUUID[uuid]
}

func (e *MsgCatalog) Asset() assets.MsgCatalog { return e.MsgCatalog }

// Reference returns a reference to this external service
func (e *MsgCatalog) Reference() *assets.MsgCatalogReference {
	return assets.NewMsgCatalogReference(e.UUID(), e.Name())
}

type MsgCatalogAssets struct {
	byUUID map[assets.MsgCatalogUUID]*MsgCatalog
}

func NewMsgCatalogAssets(msgCatalogs []assets.MsgCatalog) *MsgCatalogAssets {
	s := &MsgCatalogAssets{
		byUUID: make(map[assets.MsgCatalogUUID]*MsgCatalog, len(msgCatalogs)),
	}
	for _, asset := range msgCatalogs {
		s.byUUID[asset.UUID()] = NewMsgCatalog(asset)
	}
	return s
}

func (m *MsgCatalogOut) Header() string { return m.Header_ }

func (m *MsgCatalogOut) Body() string { return m.Body_ }

func (m *MsgCatalogOut) Footer() string { return m.Footer_ }

func (m *MsgCatalogOut) Products() []string { return m.Products_ }

func (m *MsgCatalogOut) Topic() MsgTopic { return m.Topic_ }

func (m *MsgCatalogOut) Action() string { return m.Action_ }

func (m *MsgCatalogOut) Smart() bool { return m.Smart_ }

func (m *MsgCatalogOut) ProductSearch() string { return m.ProductSearch_ }
