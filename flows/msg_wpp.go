package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

type MsgWppOut struct {
	BaseMsg
	InteractionType_ string             `json:"interaction_type,omitempty"`
	HeaderType_      string             `json:"header_type,omitempty"`
	HeaderText_      string             `json:"header_text,omitempty"`
	Text_            string             `json:"text,omitempty"`
	Footer_          string             `json:"footer,omitempty"`
	Topic_           MsgTopic           `json:"topic,omitempty"`
	ListMessage_     ListMessage        `json:"list_message,omitempty"`
	Attachments_     []utils.Attachment `json:"attachments,omitempty"`
	QuickReplies_    []string           `json:"quick_replies,omitempty"`
	TextLanguage     envs.Language      `json:"text_language,omitempty"`
	CTAMessage_      CTAMessage         `json:"cta_message,omitempty"`
	FlowMessage_     FlowMessage        `json:"flow_message,omitempty"`
}

type ListMessage struct {
	ButtonText string      `json:"button_text,omitempty"`
	ListItems  []ListItems `json:"list_items,omitempty"`
}

type CTAMessage struct {
	DisplayText_ string `json:"display_text,omitempty"`
	URL_         string `json:"url,omitempty"`
}

type FlowData map[string]string

type FlowMessage struct {
	FlowID     string   `json:"flow_id,omitempty"`
	FlowData   FlowData `json:"flow_data,omitempty"`
	FlowScreen string   `json:"flow_screen,omitempty"`
}

type ListItems struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

func NewMsgWppOut(urn urns.URN, channel *assets.ChannelReference, interactionType, headerType, headerText, text, footer string, ctaMessage CTAMessage, listMessage ListMessage, flowMessage FlowMessage, attachments []utils.Attachment, replyButtons []string, topic MsgTopic) *MsgWppOut {
	return &MsgWppOut{
		BaseMsg: BaseMsg{
			UUID_:    MsgUUID(uuids.New()),
			URN_:     urn,
			Channel_: channel,
		},
		HeaderType_:      headerType,
		InteractionType_: interactionType,
		HeaderText_:      headerText,
		Text_:            text,
		Footer_:          footer,
		ListMessage_:     listMessage,
		Attachments_:     attachments,
		QuickReplies_:    replyButtons,
		Topic_:           topic,
		CTAMessage_:      ctaMessage,
		FlowMessage_:     flowMessage,
	}
}

func (m *MsgWppOut) InteractionType() string { return m.InteractionType_ }

func (m *MsgWppOut) HeaderType() string { return m.HeaderType_ }

func (m *MsgWppOut) HeaderText() string { return m.HeaderText_ }

func (m *MsgWppOut) Text() string { return m.Text_ }

func (m *MsgWppOut) Footer() string { return m.Footer_ }

func (m *MsgWppOut) ListMessage() ListMessage { return m.ListMessage_ }

func (m *MsgWppOut) Attachments() []utils.Attachment { return m.Attachments_ }

func (m *MsgWppOut) Topic() MsgTopic { return m.Topic_ }

func (m *MsgWppOut) QuickReplies() []string { return m.QuickReplies_ }

func (m *MsgWppOut) CTAMessage() CTAMessage { return m.CTAMessage_ }

func (m *MsgWppOut) FlowMessage() FlowMessage { return m.FlowMessage_ }
