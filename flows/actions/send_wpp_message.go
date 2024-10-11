package actions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/httpx"
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
	HeaderType                string              `json:"header_type,omitempty"`
	HeaderText                string              `json:"header_text,omitempty"`
	Attachment                string              `json:"attachment,omitempty"`
	Text                      string              `json:"text,omitempty"`
	Footer                    string              `json:"footer,omitempty"`
	ListItems                 []flows.ListItems   `json:"list_items,omitempty"`
	ButtonText                string              `json:"button_text,omitempty"`
	QuickReplies              []string            `json:"quick_replies,omitempty"`
	InteractionType           string              `json:"interaction_type,omitempty"`
	ActionURL                 string              `json:"action_url,omitempty"`
	FlowID                    string              `json:"flow_id,omitempty"`
	FlowData                  flows.FlowData      `json:"flow_data,omitempty"`
	FlowScreen                string              `json:"flow_screen,omitempty"`
	FlowMode                  string              `json:"flow_mode,omitempty"`
	FlowDataAttachmentNameMap map[string]string   `json:"flow_data_attachment_name_map,omitempty"`
	OrderDetails              *flows.OrderDetails `json:"order_details,omitempty"`
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
	orderDetails *flows.OrderDetails,
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
				// check if the evaluated value is an attachment
				if _, ok := a.FlowDataAttachmentNameMap[k]; ok {
					// if the attachment is found, fetch it's content and save the base64 encoded content
					client := &http.Client{}
					req, err := http.NewRequest("GET", evaluatedValue, nil)
					if err != nil {
						run.LogError(step, err)
						continue
					}
					trace, err := httpx.DoTrace(client, req, nil, nil, -1)
					if err != nil {
						run.LogError(step, err)
						continue
					}
					base64Data := base64.StdEncoding.EncodeToString(trace.ResponseBody)
					evaluatedFlowData[k] = base64Data
				} else {
					evaluatedFlowData[k] = evaluatedValue
				}
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

		if evaluatedOrderItems == "" {
			logEvent(events.NewErrorf("order items evaluated to empty string"))
			return nil
		}

		tempOrderItems := []map[string]interface{}{}
		err := json.Unmarshal([]byte(evaluatedOrderItems), &tempOrderItems)
		if err != nil {
			logEvent(events.NewErrorf("error unmarshalling order items: %v", err))
			return nil
		}

		if len(tempOrderItems) == 0 {
			logEvent(events.NewErrorf("order items evaluated to empty array"))
			return nil
		}

		orderItems := []flows.MessageOrderItem{}
		for _, item := range tempOrderItems {
			if item["quantity"] == nil {
				logEvent(events.NewErrorf("order item quantity is nil"))
				return nil
			}
			convertedQuantity, isFloat := item["quantity"].(float64)
			if !isFloat {
				logEvent(events.NewErrorf("error reading order item quantity: %v", item["quantity"]))
				return nil
			}

			if item["amount"] == nil {
				logEvent(events.NewErrorf("order item amount is required"))
				return nil
			}
			var convertedAmount float64
			var convertedOffset float64
			itemAmount, ok := item["amount"].(map[string]interface{})
			if !ok {
				logEvent(events.NewErrorf("error reading order item amount: %v", item["amount"]))
				return nil
			}
			convertedAmount, isFloat = itemAmount["value"].(float64)
			if !isFloat {
				logEvent(events.NewErrorf("error reading order item amount: %v", itemAmount["value"]))
				return nil
			}

			convertedOffset, isFloat = itemAmount["offset"].(float64)
			if !isFloat {
				logEvent(events.NewErrorf("error reading order item amount offset: %v", itemAmount["offset"]))
				return nil
			}

			orderItem := flows.MessageOrderItem{
				RetailerID: item["retailer_id"].(string),
				Name:       item["name"].(string),
				Quantity:   int(convertedQuantity),
				Amount: flows.MessageOrderAmountWithOffset{
					Value:  int(convertedAmount),
					Offset: int(convertedOffset),
				},
			}

			if item["sale_amount"] != nil {
				itemSaleAmount, ok := item["sale_amount"].(map[string]interface{})
				if !ok {
					logEvent(events.NewErrorf("error reading order item sale amount: %v", item["sale_amount"]))
					return nil
				}
				convertedSaleAmount, isFloat := itemSaleAmount["value"].(float64)
				if !isFloat {
					logEvent(events.NewErrorf("error converting order item sale amount %s: %v", itemSaleAmount["value"], err))
					return nil
				}

				convertedSaleAmountOffset, isFloat := itemSaleAmount["offset"].(float64)
				if !isFloat {
					logEvent(events.NewErrorf("error converting order item sale amount offset %s: %v", itemSaleAmount["offset"], err))
					return nil
				}

				if convertedSaleAmount > 0 {
					orderItem.SaleAmount = &flows.MessageOrderAmountWithOffset{
						Value:  int(convertedSaleAmount),
						Offset: int(convertedSaleAmountOffset),
					}
				}

			}

			orderItems = append(orderItems, orderItem)
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

		var convertedOrderTax float64
		var convertedOrderShipping float64
		var convertedOrderDiscount float64
		err = nil
		if evaluatedOrderTax != "" {
			taxInRealFloatRepresentation := strings.Replace(evaluatedOrderTax, ",", ".", -1)
			convertedOrderTax, err = strconv.ParseFloat(taxInRealFloatRepresentation, 64)
			if err != nil {
				logEvent(events.NewErrorf("error converting order tax %s to int: %v", evaluatedOrderTax, err))
				return nil
			}
		}

		if evaluatedOrderShipping != "" {
			shippingInRealFloatRepresentation := strings.Replace(evaluatedOrderShipping, ",", ".", -1)
			convertedOrderShipping, err = strconv.ParseFloat(shippingInRealFloatRepresentation, 64)
			if err != nil {
				logEvent(events.NewErrorf("error converting order shipping %s to int: %v", evaluatedOrderShipping, err))
				return nil
			}
		}

		if evaluatedOrderDiscount != "" {
			discountInRealFloatRepresentation := strings.Replace(evaluatedOrderDiscount, ",", ".", -1)
			convertedOrderDiscount, err = strconv.ParseFloat(discountInRealFloatRepresentation, 64)
			if err != nil {
				logEvent(events.NewErrorf("error converting order discount %s to int: %v", evaluatedOrderDiscount, err))
				return nil
			}
		}

		subTotalValue := 0
		for _, item := range orderItems {
			if item.SaleAmount != nil && item.SaleAmount.Value > 0 {
				subTotalValue += item.SaleAmount.Value * item.Quantity
			} else {
				subTotalValue += item.Amount.Value * item.Quantity
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
