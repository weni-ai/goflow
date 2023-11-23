package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

type MsgWppOut struct {
	BaseMsg

	Header_       Header       `json:"header,omitempty"`
	Body_         string       `json:"body,omitempty"`
	Footer_       string       `json:"footer,omitempty"`
	Topic_        MsgTopic     `json:"topic,omitempty"`
	ListMessages_ ListMessages `json:"list_messages,omitempty"`
	ReplyButtons_ []string     `json:"reply_buttons,omitempty"`
}

type Header struct {
	Type        string             `json:"type,omitempty"`
	Attachments []utils.Attachment `json:"attachments,omitempty"`
	Text        string             `json:"text,omitempty"`
}

type ListMessages struct {
	Title   string `json:"title,omitempty"`
	Footer  string `json:"footer,omitempty"`
	Options []struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"options,omitempty"`
}

func NewMsgWppOut(urn urns.URN, channel *assets.ChannelReference, header Header, body string, footer string, listMessage ListMessages, replyButtons []string, topic MsgTopic) *MsgWppOut {
	return &MsgWppOut{
		BaseMsg: BaseMsg{
			UUID_:    MsgUUID(uuids.New()),
			URN_:     urn,
			Channel_: channel,
		},
		Header_:       header,
		Body_:         body,
		Footer_:       footer,
		ListMessages_: listMessage,
		ReplyButtons_: replyButtons,
		Topic_:        topic,
	}
}

func (m *MsgWppOut) Header() Header { return m.Header_ }

func (m *MsgWppOut) Body() string { return m.Body_ }

func (m *MsgWppOut) Footer() string { return m.Footer_ }

func (m *MsgWppOut) ListMessages() ListMessages { return m.ListMessages_ }

func (m *MsgWppOut) Topic() MsgTopic { return m.Topic_ }

func (m *MsgWppOut) ReplyButtons() []string { return m.ReplyButtons_ }
