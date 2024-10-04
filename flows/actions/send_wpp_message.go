package actions

import (
	"encoding/json"
	"strconv"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendWppMsg, func() flows.Action { return &SendWppMsgAction{} })
}

// TypeSendWppMsg is the type for the send message whatsapp action
const TypeSendWppMsg string = "send_whatsapp_msg"

type SendWppMsgAction struct {
	baseAction
	universalAction
	createWppMsgAction

	AllURNs bool           `json:"all_urns,omitempty"`
	Topic   flows.MsgTopic `json:"topic,omitempty" validate:"omitempty,msg_topic"`
}

type createWppMsgAction struct {
	HeaderType      string             `json:"header_type,omitempty"`
	HeaderText      string             `json:"header_text,omitempty"`
	Attachment      string             `json:"attachment,omitempty"`
	Text            string             `json:"text,omitempty"`
	Footer          string             `json:"footer,omitempty"`
	ListItems       []flows.ListItems  `json:"list_items,omitempty"`
	ButtonText      string             `json:"button_text,omitempty"`
	QuickReplies    []string           `json:"quick_replies,omitempty"`
	InteractionType string             `json:"interaction_type,omitempty"`
	ActionURL       string             `json:"action_url,omitempty"`
	FlowID          string             `json:"flow_id,omitempty"`
	FlowData        flows.FlowData     `json:"flow_data,omitempty"`
	FlowScreen      string             `json:"flow_screen,omitempty"`
	FlowMode        string             `json:"flow_mode,omitempty"`
	OrderDetails    flows.OrderDetails `json:"order_details,omitempty"`
}

type Header struct {
	Type        string   `json:"type,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
	Text        string   `json:"text,omitempty"`
}

// NewSendWppMsg creates a new send msg whatsapp action
func NewSendWppMsg(
	uuid flows.ActionUUID,
	headerType string,
	headerText string,
	attachment string,
	text string,
	footer string,
	listItems []flows.ListItems,
	buttonText string,
	quickReplies []string,
	interactionType string,
	actionURL string,
	flowID string,
	flowData flows.FlowData,
	flowScreen string,
	flowMode string,
	orderDetails flows.OrderDetails,
	allURNs bool) *SendWppMsgAction {
	return &SendWppMsgAction{
		baseAction: newBaseAction(TypeSendWppMsg, uuid),
		createWppMsgAction: createWppMsgAction{
			HeaderType:      headerType,
			HeaderText:      headerText,
			Attachment:      attachment,
			Text:            text,
			Footer:          footer,
			ListItems:       listItems,
			ButtonText:      buttonText,
			QuickReplies:    quickReplies,
			InteractionType: interactionType,
			ActionURL:       actionURL,
			FlowID:          flowID,
			FlowData:        flowData,
			FlowScreen:      flowScreen,
			FlowMode:        flowMode,
			OrderDetails:    orderDetails,
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

	evaluatedHeaderText, evaluatedFooter, evaluatedText, evaluatedListItems, evaluatedButtonText, evaluatedAttachments, evaluatedReplyMessage := a.evaluateMessageWpp(run, nil, a.HeaderType, a.InteractionType, a.HeaderText, a.Footer, a.Text, a.ListItems, a.ButtonText, a.Attachment, a.QuickReplies, logEvent)

	listMessage := flows.ListMessage{}
	if len(evaluatedListItems) > 0 {
		listMessage = flows.ListMessage{
			ButtonText: evaluatedButtonText,
			ListItems:  evaluatedListItems,
		}
	}

	ctaMessage := flows.CTAMessage{}
	if a.ActionURL != "" {
		evaluatedActionURL, _ := run.EvaluateTemplate(a.ActionURL)
		ctaMessage = flows.CTAMessage{
			DisplayText_: evaluatedButtonText,
			URL_:         evaluatedActionURL,
		}
	}

	flowMessage := flows.FlowMessage{}
	if a.FlowID != "" {
		evaluatedFlowID, _ := run.EvaluateTemplate(a.FlowID)
		evaluatedFlowScreen, _ := run.EvaluateTemplate(a.FlowScreen)

		evaluatedFlowData := make(flows.FlowData)
		for k, v := range a.FlowData {
			evaluatedValue, _ := run.EvaluateTemplate(v.(string))

			// check if the evalutaed value is a valid JSON, and if so do not convert to string
			var jsonValue json.RawMessage
			err := json.Unmarshal([]byte(evaluatedValue), &jsonValue)
			if err == nil {
				evaluatedFlowData[k] = jsonValue
			} else {
				evaluatedFlowData[k] = evaluatedValue
			}
		}

		flowMessage = flows.FlowMessage{
			FlowID:     evaluatedFlowID,
			FlowData:   evaluatedFlowData,
			FlowScreen: evaluatedFlowScreen,
			FlowCTA:    evaluatedButtonText,
			FlowMode:   a.FlowMode,
		}
	}

	orderDetailsMessage := flows.OrderDetailsMessage{}
	if a.InteractionType == "order_details" {
		evaluatedReferenceID, _ := run.EvaluateTemplate(a.OrderDetails.ReferenceID)
		evaluatedOrderItems, _ := run.EvaluateTemplate(a.OrderDetails.Items)

		orderItems := []flows.MessageOrderItem{}
		tempOrderItems := []map[string]interface{}{}
		err := json.Unmarshal([]byte(evaluatedOrderItems), &tempOrderItems)
		if err != nil {
			logEvent(events.NewErrorf("error unmarshalling order items: %v", err))
			return nil
		}

		evaluatedOrderTax, _ := run.EvaluateTemplate(a.OrderDetails.Tax.Value)
		evaluatedOrderTaxDescription, _ := run.EvaluateTemplate(a.OrderDetails.Tax.Description)

		evaluatedOrderShipping, _ := run.EvaluateTemplate(a.OrderDetails.Shipping.Value)
		evaluatedOrderShippingDescription, _ := run.EvaluateTemplate(a.OrderDetails.Shipping.Description)

		evaluatedOrderDiscount, _ := run.EvaluateTemplate(a.OrderDetails.Discount.Value)
		evaluatedOrderDiscountDescription, _ := run.EvaluateTemplate(a.OrderDetails.Discount.Description)
		evaluatedOrderDiscountProgramName, _ := run.EvaluateTemplate(a.OrderDetails.Discount.ProgramName)

		evaluatedOrderPaymentType, _ := run.EvaluateTemplate(a.OrderDetails.PaymentSettings.Type)
		evaluatedOrderPaymentLink, _ := run.EvaluateTemplate(a.OrderDetails.PaymentSettings.PaymentLink)

		evaluatedOrderPixKey, _ := run.EvaluateTemplate(a.OrderDetails.PaymentSettings.PixConfig.Key)
		evaluatedOrderPixKeyType, _ := run.EvaluateTemplate(a.OrderDetails.PaymentSettings.PixConfig.KeyType)
		evaluatedOrderPixMerchantName, _ := run.EvaluateTemplate(a.OrderDetails.PaymentSettings.PixConfig.MerchantName)
		evaluatedOrderPixCode, _ := run.EvaluateTemplate(a.OrderDetails.PaymentSettings.PixConfig.Code)

		convertedOrderTax, err := strconv.ParseFloat(evaluatedOrderTax, 64)
		if err != nil {
			logEvent(events.NewErrorf("error converting order tax %s to int: %v", evaluatedOrderTax, err))
			return nil
		}

		convertedOrderShipping, err := strconv.ParseFloat(evaluatedOrderShipping, 64)
		if err != nil {
			logEvent(events.NewErrorf("error converting order shipping %s to int: %v", evaluatedOrderShipping, err))
			return nil
		}

		convertedOrderDiscount, err := strconv.ParseFloat(evaluatedOrderDiscount, 64)
		if err != nil {
			logEvent(events.NewErrorf("error converting order discount %s to int: %v", evaluatedOrderDiscount, err))
			return nil
		}

		subTotalValue := 0
		for _, item := range orderItems {
			if item.SaleAmount != 0 {
				subTotalValue += item.SaleAmount * item.Quantity
			} else {
				subTotalValue += item.Amount * item.Quantity
			}
		}

		taxValue := int(convertedOrderTax * 100)
		shippingValue := int(convertedOrderShipping * 100)
		discountValue := int(convertedOrderDiscount * 100)
		totalValue := subTotalValue + taxValue + shippingValue - discountValue

		orderDetailsMessage = flows.OrderDetailsMessage{
			ReferenceID: evaluatedReferenceID,
			PaymentSettings: &flows.OrderPaymentSettings{
				Type:        evaluatedOrderPaymentType,
				PaymentLink: evaluatedOrderPaymentLink,
				PixConfig: &flows.OrderPixConfig{
					Key:          evaluatedOrderPixKey,
					KeyType:      evaluatedOrderPixKeyType,
					MerchantName: evaluatedOrderPixMerchantName,
					Code:         evaluatedOrderPixCode,
				},
			},
			TotalAmount: totalValue,
			Order: &flows.MessageOrder{
				Items:    &orderItems,
				Subtotal: subTotalValue,
				Tax: &flows.MessageOrderAmountWithDescription{
					Value:       taxValue,
					Description: evaluatedOrderTaxDescription,
				},
				Shipping: &flows.MessageOrderAmountWithDescription{
					Value:       shippingValue,
					Description: evaluatedOrderShippingDescription,
				},
				Discount: &flows.MessageOrderDiscount{
					Value:       discountValue,
					Description: evaluatedOrderDiscountDescription,
					ProgramName: evaluatedOrderDiscountProgramName,
				},
			},
		}

	}

	destinations := run.Contact().ResolveDestinations(a.AllURNs)

	for _, dest := range destinations {
		var channelRef *assets.ChannelReference
		if dest.Channel != nil {
			channelRef = assets.NewChannelReference(dest.Channel.UUID(), dest.Channel.Name())
		}

		msg := flows.NewMsgWppOut(dest.URN.URN(), channelRef, a.InteractionType, a.HeaderType, evaluatedHeaderText, evaluatedText, evaluatedFooter, ctaMessage, listMessage, flowMessage, orderDetailsMessage, evaluatedAttachments, evaluatedReplyMessage, a.Topic)
		logEvent(events.NewMsgWppCreated(msg))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgWppOut(urns.NilURN, nil, a.InteractionType, a.HeaderType, evaluatedHeaderText, evaluatedText, evaluatedFooter, ctaMessage, listMessage, flowMessage, orderDetailsMessage, evaluatedAttachments, evaluatedReplyMessage, flows.NilMsgTopic)
		logEvent(events.NewMsgWppCreated(msg))
	}

	return nil
}
