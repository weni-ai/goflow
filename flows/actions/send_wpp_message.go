package actions

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendWppMsg, func() flows.Action { return &SendWppMsgAction{} })
}

// TypeSendWppMsg is the type for the send message whatsapp action
const TypeSendWppMsg string = "send_wpp_msg"

type SendWppMsgAction struct {
	baseAction
	universalAction
	createWppMsgAction

	AllURNs bool           `json:"all_urns,omitempty"`
	Topic   flows.MsgTopic `json:"topic,omitempty" validate:"omitempty,msg_topic"`
}

type createWppMsgAction struct {
	Header       Header       `json:"header,omitempty"`
	Body         string       `json:"body,omitempty"`
	Footer       string       `json:"footer,omitempty"`
	ListMessages ListMessages `json:"list_messages,omitempty"`
	ReplyButtons []string     `json:"reply_buttons,omitempty"`
}

type Header struct {
	Type        string   `json:"type,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
	Text        string   `json:"text,omitempty"`
}

type ListMessages struct {
	Title   string `json:"title,omitempty"`
	Footer  string `json:"footer,omitempty"`
	Options []struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"options,omitempty"`
}

// NewSendWppMsg creates a new send msg whatsapp action
func NewSendWppMsg(uuid flows.ActionUUID, header Header, body string, footer string, listMessage ListMessages, replyButtons []string, allURNs bool) *SendWppMsgAction {
	return &SendWppMsgAction{
		baseAction: newBaseAction(TypeSendMsgCatalog, uuid),
		createWppMsgAction: createWppMsgAction{
			Header:       header,
			Body:         body,
			Footer:       footer,
			ListMessages: listMessage,
			ReplyButtons: replyButtons,
		},
		AllURNs: allURNs,
	}
}

// Execute runs this action
func (a *SendWppMsgAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	if run.Contact() == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	evaluatedHeaderText, evaluatedAttachments, evaluatedBody, evaluatedFooter, evaluatedReplyMessage := a.evaluateMessageWpp(run, nil, a.Header, a.Body, a.Footer, a.ListMessages, a.ReplyButtons, logEvent)

	evaluatedHeader := flows.Header{
		Type:        a.Header.Type,
		Text:        evaluatedHeaderText,
		Attachments: evaluatedAttachments,
	}

	destinations := run.Contact().ResolveDestinations(a.AllURNs)

	for _, dest := range destinations {
		var channelRef *assets.ChannelReference
		if dest.Channel != nil {
			channelRef = assets.NewChannelReference(dest.Channel.UUID(), dest.Channel.Name())
		}

		msg := flows.NewMsgWppOut(dest.URN.URN(), channelRef, evaluatedHeader, evaluatedBody, evaluatedFooter, flows.ListMessages(a.ListMessages), evaluatedReplyMessage, a.Topic)
		logEvent(events.NewMsgWppCreated(msg))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgWppOut(urns.NilURN, nil, evaluatedHeader, evaluatedBody, evaluatedFooter, flows.ListMessages(a.ListMessages), evaluatedReplyMessage, flows.NilMsgTopic)
		logEvent(events.NewMsgWppCreated(msg))
	}

	return nil
}
