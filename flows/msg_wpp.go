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
	ListMessage_     ListMessage        `json:"list_messages,omitempty"`
	Attachments_     []utils.Attachment `json:"attachments,omitempty"`
	QuickReplies_    []string           `json:"quick_replies,omitempty"`
	TextLanguage     envs.Language      `json:"text_language,omitempty"`
}

type ListMessage struct {
	ListTitle  string      `json:"list_title,omitempty"`
	ListFooter string      `json:"list_footer,omitempty"`
	ListItems  []ListItems `json:"list_items,omitempty"`
}

type ListItems struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

func NewMsgWppOut(urn urns.URN, channel *assets.ChannelReference, interactionType string, headerType string, headerText string, text string, footer string, listMessage ListMessage, attachments []utils.Attachment, replyButtons []string, topic MsgTopic) *MsgWppOut {
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
	}
}

func (m *MsgWppOut) InteractionType() string { return m.InteractionType_ }

func (m *MsgWppOut) HeaderType() string { return m.HeaderType_ }

func (m *MsgWppOut) HeaderText() string { return m.HeaderText_ }

func (m *MsgWppOut) Text() string { return m.Text_ }

func (m *MsgWppOut) Footer() string { return m.Footer_ }

func (m *MsgWppOut) ListMessages() ListMessage { return m.ListMessage_ }

func (m *MsgWppOut) Attachments() []utils.Attachment { return m.Attachments_ }

func (m *MsgWppOut) Topic() MsgTopic { return m.Topic_ }

func (m *MsgWppOut) QuickReplies() []string { return m.QuickReplies_ }
