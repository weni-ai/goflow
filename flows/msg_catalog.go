package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
)

type MsgCatalog struct {
	BaseMsg

	Header_      string        `json:"header,omitempty"`
	Body_        string        `json:"body,omitempty"`
	Footer_      string        `json:"footer,omitempty"`
	Products_    []string      `json:"products,omitempty"`
	Action_      string        `json:"action,omitempty"`
	Topic_       MsgTopic      `json:"topic,omitempty"`
	Smart_       bool          `json:"smart"`
	TextLanguage envs.Language `json:"text_language,omitempty"`
}

func NewMsgCatalog(urn urns.URN, channel *assets.ChannelReference, header string, body string, footer string, products []string, action string, smart bool, topic MsgTopic) *MsgCatalog {
	return &MsgCatalog{
		BaseMsg: BaseMsg{
			UUID_:    MsgUUID(uuids.New()),
			URN_:     urn,
			Channel_: channel,
		},
		Header_:   header,
		Body_:     body,
		Footer_:   footer,
		Products_: products,
		Action_:   action,
		Smart_:    smart,
		Topic_:    topic,
	}
}

func (m *MsgCatalog) Header() string { return m.Header_ }

func (m *MsgCatalog) Body() string { return m.Body_ }

func (m *MsgCatalog) Footer() string { return m.Footer_ }

func (m *MsgCatalog) Products() []string { return m.Products_ }

func (m *MsgCatalog) Topic() MsgTopic { return m.Topic_ }

func (m *MsgCatalog) Action() string { return m.Action_ }

func (m *MsgCatalog) Smart() bool { return m.Smart_ }
