package actions

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendMsg, func() flows.Action { return &SendMsgAction{} })
}

// TypeSendMsg is the type for the send message action
const TypeSendMsg string = "send_msg"

// SendMsgAction can be used to reply to the current contact in a flow. The text field may contain templates. The action
// will attempt to find pairs of URNs and channels which can be used for sending. If it can't find such a pair, it will
// create a message without a channel or URN.
//
// A [event:msg_created] event will be created with the evaluated text.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "send_msg",
//	  "text": "Hi @contact.name, are you ready to complete today's survey?",
//	  "attachments": [],
//	  "all_urns": false,
//	  "templating": {
//	    "uuid": "32c2ead6-3fa3-4402-8e27-9cc718175c5a",
//	    "template": {
//	      "uuid": "3ce100b7-a734-4b4e-891b-350b1279ade2",
//	      "name": "revive_issue"
//	    },
//	    "variables": ["@contact.name"]
//	  },
//	  "topic": "event"
//	  "ig_comment": "0123456789",
//	  "ig_response_type": "comment"
//	}
//
// @action send_msg
type SendMsgAction struct {
	baseAction
	universalAction
	createMsgAction

	AllURNs           bool               `json:"all_urns,omitempty"`
	Templating        *Templating        `json:"templating,omitempty" validate:"omitempty,dive"`
	Topic             flows.MsgTopic     `json:"topic,omitempty" validate:"omitempty,msg_topic"`
	InstagramSettings *InstagramSettings `json:"instagram_settings,omitempty"`
}

type InstagramSettings struct {
	ResponseType string `json:"response_type,omitempty"`
	CommentID    string `json:"comment_id,omitempty"`
	Tag          string `json:"tag,omitempty"`
}

// Templating represents the templating that should be used if possible
type Templating struct {
	UUID          uuids.UUID                `json:"uuid" validate:"required,uuid4"`
	Template      *assets.TemplateReference `json:"template" validate:"required"`
	Variables     []string                  `json:"variables" engine:"localized,evaluated"`
	CarouselCards []CarouselCard            `json:"carousel_cards,omitempty"`
}

type CarouselCard struct {
	Body    string               `json:"body,omitempty"`
	Index   int                  `json:"index,omitempty"`
	Buttons []CarouselCardButton `json:"buttons,omitempty"`
}

type CarouselCardButton struct {
	SubType   string `json:"sub_type"`  // quick_reply, url, phone_number
	Parameter string `json:"parameter"` // payload for quick_reply, url variable for url, phone number for phone_number
}

// LocalizationUUID gets the UUID which identifies this object for localization
func (t *Templating) LocalizationUUID() uuids.UUID { return t.UUID }

// NewSendMsg creates a new send msg action
func NewSendMsg(uuid flows.ActionUUID, text string, attachments []string, quickReplies []string, commentID string, responseType string, tag string, allURNs bool) *SendMsgAction {
	return &SendMsgAction{
		baseAction: newBaseAction(TypeSendMsg, uuid),
		createMsgAction: createMsgAction{
			Text:         text,
			Attachments:  attachments,
			QuickReplies: quickReplies,
		},
		InstagramSettings: &InstagramSettings{
			CommentID:    commentID,
			ResponseType: responseType,
			Tag:          tag,
		},
		AllURNs: allURNs,
	}
}

// Execute runs this action
func (a *SendMsgAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	if run.Contact() == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	evaluatedText, evaluatedAttachments, evaluatedQuickReplies := a.evaluateMessage(run, nil, a.Text, a.Attachments, a.QuickReplies, logEvent)

	var evaluatedIGComment string
	var IGresponseType string
	var IGTag string
	if a.InstagramSettings != nil {
		evaluatedIGComment = a.evaluateMessageIG(run, nil, a.InstagramSettings.CommentID, logEvent)
		IGresponseType = a.InstagramSettings.ResponseType
		IGTag = a.InstagramSettings.Tag
	}

	destinations := run.Contact().ResolveDestinations(a.AllURNs)

	sa := run.Session().Assets()

	// create a new message for each URN+channel destination
	for _, dest := range destinations {
		var channelRef *assets.ChannelReference
		if dest.Channel != nil {
			channelRef = assets.NewChannelReference(dest.Channel.UUID(), dest.Channel.Name())
		}

		var templating *flows.MsgTemplating

		// do we have a template defined?
		if a.Templating != nil {
			// looks for a translation in these locales
			locales := []envs.Locale{
				run.Contact().Locale(run.Environment()),
				run.Environment().DefaultLocale(),
			}

			translation := sa.Templates().FindTranslation(a.Templating.Template.UUID, channelRef, locales)
			if translation != nil {
				localizedVariables, _ := run.GetTextArray(uuids.UUID(a.Templating.UUID), "variables", a.Templating.Variables)

				// evaluate our variables
				evaluatedVariables := make([]string, len(localizedVariables))
				for i, variable := range localizedVariables {
					sub, err := run.EvaluateTemplate(variable)
					if err != nil {
						logEvent(events.NewError(err))
					}
					evaluatedVariables[i] = sub
				}
				// Build evaluated carousel cards
				var evaluatedCarouselCards []flows.CarouselCard
				if len(a.Templating.CarouselCards) > 0 {
					evaluatedCarouselCards = make([]flows.CarouselCard, len(a.Templating.CarouselCards))
					for idx, carouselCard := range a.Templating.CarouselCards {
						// Evaluate body variables
						localizedCarouselCardsBody := run.GetText(uuids.UUID(a.Templating.UUID), "carousel_cards.body", carouselCard.Body)
						evaluatedBody, err := run.EvaluateTemplate(localizedCarouselCardsBody)
						if err != nil {
							logEvent(events.NewError(err))
						}

						// Evaluate button text variables
						evaluatedButtons := make([]flows.CarouselCardButton, len(carouselCard.Buttons))
						for i, button := range carouselCard.Buttons {
							localizedCarouselCardsButtonsParameter := run.GetText(uuids.UUID(a.Templating.UUID), "carousel_cards.buttons.parameter", button.Parameter)
							evaluatedButtonParameter, err := run.EvaluateTemplate(localizedCarouselCardsButtonsParameter)
							if err != nil {
								logEvent(events.NewError(err))
							}
							evaluatedButtons[i] = flows.CarouselCardButton{
								SubType:   button.SubType,
								Parameter: evaluatedButtonParameter,
							}
						}

						evaluatedCarouselCards[idx] = flows.CarouselCard{
							Body:    evaluatedBody,
							Index:   carouselCard.Index,
							Buttons: evaluatedButtons,
						}
					}
				}

				evaluatedText = translation.Substitute(evaluatedVariables)
				template := sa.Templates().Get(a.Templating.Template.UUID)
				templating = flows.NewMsgTemplating(template.Reference(), translation.Language(), translation.Country(), evaluatedVariables, translation.Namespace(), evaluatedCarouselCards)
			}
		}

		msg := flows.NewMsgOut(dest.URN.URN(), channelRef, evaluatedText, evaluatedAttachments, evaluatedQuickReplies, templating, a.Topic, evaluatedIGComment, IGresponseType, IGTag)
		logEvent(events.NewMsgCreated(msg))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgOut(urns.NilURN, nil, evaluatedText, evaluatedAttachments, evaluatedQuickReplies, nil, flows.NilMsgTopic, evaluatedIGComment, IGresponseType, IGTag)
		logEvent(events.NewMsgCreated(msg))
	}

	return nil
}
