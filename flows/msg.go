package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// BaseMsg represents a incoming or outgoing message with the session contact
type BaseMsg struct {
	UUID_        MsgUUID                  `json:"uuid"`
	ID_          MsgID                    `json:"id,omitempty"`
	URN_         urns.URN                 `json:"urn,omitempty" validate:"omitempty,urn"`
	Channel_     *assets.ChannelReference `json:"channel,omitempty"`
	Text_        string                   `json:"text"`
	Attachments_ []Attachment             `json:"attachments,omitempty"`
}

// MsgIn represents a incoming message from the session contact
type MsgIn struct {
	BaseMsg

	ExternalID_ string `json:"external_id,omitempty"`
}

// MsgOut represents a outgoing message to the session contact
type MsgOut struct {
	BaseMsg

	QuickReplies_ []string     `json:"quick_replies,omitempty"`
	Template_     *MsgTemplate `json:"template,omitempty"`
}

// NewMsgIn creates a new incoming message
func NewMsgIn(uuid MsgUUID, urn urns.URN, channel *assets.ChannelReference, text string, attachments []Attachment) *MsgIn {
	return &MsgIn{
		BaseMsg: BaseMsg{
			UUID_:        uuid,
			URN_:         urn,
			Channel_:     channel,
			Text_:        text,
			Attachments_: attachments,
		},
	}
}

// NewMsgOut creates a new outgoing message
func NewMsgOut(urn urns.URN, channel *assets.ChannelReference, text string, attachments []Attachment, quickReplies []string, template *MsgTemplate) *MsgOut {
	return &MsgOut{
		BaseMsg: BaseMsg{
			UUID_:        MsgUUID(utils.NewUUID()),
			URN_:         urn,
			Channel_:     channel,
			Text_:        text,
			Attachments_: attachments,
		},
		QuickReplies_: quickReplies,
		Template_:     template,
	}
}

// UUID returns the UUID of this message
func (m *BaseMsg) UUID() MsgUUID { return m.UUID_ }

// ID returns the internal ID of this message
func (m *BaseMsg) ID() MsgID { return m.ID_ }

// SetID sets the internal ID of this message
func (m *BaseMsg) SetID(id MsgID) { m.ID_ = id }

// URN returns the URN of this message
func (m *BaseMsg) URN() urns.URN { return m.URN_ }

// Channel returns the channel of this message
func (m *BaseMsg) Channel() *assets.ChannelReference { return m.Channel_ }

// Text returns the text of this message
func (m *BaseMsg) Text() string { return m.Text_ }

// Attachments returns the attachments of this message
func (m *BaseMsg) Attachments() []Attachment { return m.Attachments_ }

// ExternalID returns the optional external ID of this incoming message
func (m *MsgIn) ExternalID() string { return m.ExternalID_ }

// SetExternalID sets the external ID of this message
func (m *MsgIn) SetExternalID(id string) { m.ExternalID_ = id }

// QuickReplies returns the quick replies of this outgoing message
func (m *MsgOut) QuickReplies() []string { return m.QuickReplies_ }

// Template returns the template to use to send this message (if any)
func (m *MsgOut) Template() *MsgTemplate { return m.Template_ }

// MsgTemplate represents a substituted message template, containing the uuid and name of the template that should be used as well
// as the language and variables that should be substituted
type MsgTemplate struct {
	UUID_      assets.TemplateUUID `json:"uuid"`
	Name_      string              `json:"name"`
	Language_  utils.Language      `json:"language"`
	Variables_ []string            `json:"template_variables"`
}

// UUID returns the uuid of the template that should be used in the template
func (t MsgTemplate) UUID() assets.TemplateUUID { return t.UUID_ }

// Name returns the name of the template that should be used in the template
func (t MsgTemplate) Name() string { return t.Name_ }

// Language returns the language that should be used for the template
func (t MsgTemplate) Language() utils.Language { return t.Language_ }

// Variables returns the variables that should be substituted in the template
func (t MsgTemplate) Variables() []string { return t.Variables_ }

// NewMsgTemplate creates and returns a new msg template
func NewMsgTemplate(uuid assets.TemplateUUID, name string, language utils.Language, variables []string) *MsgTemplate {
	return &MsgTemplate{
		UUID_:      uuid,
		Name_:      name,
		Language_:  language,
		Variables_: variables,
	}
}
