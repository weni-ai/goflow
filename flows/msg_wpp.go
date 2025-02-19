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
	InteractionType_     string              `json:"interaction_type,omitempty"`
	HeaderType_          string              `json:"header_type,omitempty"`
	HeaderText_          string              `json:"header_text,omitempty"`
	Text_                string              `json:"text,omitempty"`
	Footer_              string              `json:"footer,omitempty"`
	Topic_               MsgTopic            `json:"topic,omitempty"`
	ListMessage_         ListMessage         `json:"list_message,omitempty"`
	Attachments_         []utils.Attachment  `json:"attachments,omitempty"`
	QuickReplies_        []string            `json:"quick_replies,omitempty"`
	TextLanguage         envs.Language       `json:"text_language,omitempty"`
	CTAMessage_          CTAMessage          `json:"cta_message,omitempty"`
	FlowMessage_         FlowMessage         `json:"flow_message,omitempty"`
	OrderDetailsMessage_ OrderDetailsMessage `json:"order_details_message,omitempty"`
	Templating_          *MsgTemplating      `json:"templating,omitempty"`
	Buttons_             []ButtonComponent   `json:"buttons,omitempty"`

	// fields for msg_catalog
	Header_        string         `json:"header,omitempty"`
	Body_          string         `json:"body,omitempty"`
	Products_      []ProductEntry `json:"products,omitempty"`
	Action_        string         `json:"action,omitempty"`
	Smart_         bool           `json:"smart"`
	ProductSearch_ string         `json:"product_search,omitempty"`
	SendCatalog_   bool           `json:"send_catalog,omitempty"`
}

type ButtonComponent struct {
	SubType    string        `json:"sub_type"`
	Parameters []ButtonParam `json:"parameters"`
}

type ButtonParam struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type ListMessage struct {
	ButtonText string      `json:"button_text,omitempty"`
	ListItems  []ListItems `json:"list_items,omitempty"`
}

type CTAMessage struct {
	DisplayText_ string `json:"display_text,omitempty"`
	URL_         string `json:"url,omitempty"`
}

type FlowData map[string]interface{}

type FlowMessage struct {
	FlowID     string   `json:"flow_id,omitempty"`
	FlowData   FlowData `json:"flow_data,omitempty"`
	FlowScreen string   `json:"flow_screen,omitempty"`
	FlowCTA    string   `json:"flow_cta,omitempty"`
	FlowMode   string   `json:"flow_mode,omitempty"`
}

type ListItems struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

// Message order details structs, with string-like attributes to be evaluated and calculated
type OrderAmountWithDescription struct {
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type OrderDiscount struct {
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	ProgramName string `json:"program_name,omitempty"`
}

type OrderPixConfig struct {
	Key          string `json:"key,omitempty"`
	KeyType      string `json:"key_type,omitempty"`
	MerchantName string `json:"merchant_name,omitempty"`
	Code         string `json:"code,omitempty"`
}

type OrderPaymentSettings struct {
	Type        string          `json:"type,omitempty"`
	PaymentLink string          `json:"payment_link,omitempty"`
	PixConfig   *OrderPixConfig `json:"pix_config,omitempty"`
}

type OrderDetails struct {
	ReferenceID     string                      `json:"reference_id,omitempty"`
	Items           string                      `json:"item_list,omitempty"`
	Tax             *OrderAmountWithDescription `json:"tax,omitempty"`
	Shipping        *OrderAmountWithDescription `json:"shipping,omitempty"`
	Discount        *OrderDiscount              `json:"discount,omitempty"`
	PaymentSettings *OrderPaymentSettings       `json:"payment_settings,omitempty"`
}

// Message for order details, with attribute types defined such as int values
type OrderDetailsMessage struct {
	ReferenceID     string                `json:"reference_id,omitempty"`
	PaymentSettings *OrderPaymentSettings `json:"payment_settings,omitempty"`
	TotalAmount     int                   `json:"total_amount,omitempty"`
	Order           *MessageOrder         `json:"order,omitempty"`
}

type MessageOrder struct {
	Items    *[]MessageOrderItem                `json:"items,omitempty"`
	Subtotal int                                `json:"subtotal,omitempty"`
	Tax      *MessageOrderAmountWithDescription `json:"tax,omitempty"`
	Shipping *MessageOrderAmountWithDescription `json:"shipping,omitempty"`
	Discount *MessageOrderDiscount              `json:"discount,omitempty"`
}

type MessageOrderItem struct {
	RetailerID string                        `json:"retailer_id,omitempty"`
	Name       string                        `json:"name,omitempty"`
	Quantity   int                           `json:"quantity,omitempty"`
	Amount     MessageOrderAmountWithOffset  `json:"amount,omitempty"`
	SaleAmount *MessageOrderAmountWithOffset `json:"sale_amount,omitempty"`
}

type MessageOrderAmountWithDescription struct {
	Value       int    `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type MessageOrderDiscount struct {
	Value       int    `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	ProgramName string `json:"program_name,omitempty"`
}

type MessageOrderAmountWithOffset struct {
	Value  int `json:"value"`
	Offset int `json:"offset"`
}

func NewMsgWppOut(urn urns.URN, channel *assets.ChannelReference, interactionType, headerType, headerText, text, footer string, ctaMessage CTAMessage, listMessage ListMessage, flowMessage FlowMessage, orderDetailsMessage OrderDetailsMessage, attachments []utils.Attachment, replyButtons []string, buttons []ButtonComponent, templating *MsgTemplating, topic MsgTopic, header string, body string, products []ProductEntry, action string, smart bool, productSearch string, sendCatalog bool) *MsgWppOut {
	return &MsgWppOut{
		BaseMsg: BaseMsg{
			UUID_:    MsgUUID(uuids.New()),
			URN_:     urn,
			Channel_: channel,
		},
		HeaderType_:          headerType,
		InteractionType_:     interactionType,
		HeaderText_:          headerText,
		Text_:                text,
		Footer_:              footer,
		ListMessage_:         listMessage,
		Attachments_:         attachments,
		QuickReplies_:        replyButtons,
		Topic_:               topic,
		CTAMessage_:          ctaMessage,
		FlowMessage_:         flowMessage,
		OrderDetailsMessage_: orderDetailsMessage,
		Templating_:          templating,
		Buttons_:             buttons,
		Header_:              header,
		Body_:                body,
		Products_:            products,
		Action_:              action,
		Smart_:               smart,
		ProductSearch_:       productSearch,
		SendCatalog_:         sendCatalog,
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

func (m *MsgWppOut) OrderDetailsMessage() OrderDetailsMessage { return m.OrderDetailsMessage_ }

func (m *MsgWppOut) Templating() *MsgTemplating { return m.Templating_ }

func (m *MsgWppOut) Buttons() []ButtonComponent { return m.Buttons_ }

func (m *MsgWppOut) Header() string { return m.Header_ }

func (m *MsgWppOut) Body() string { return m.Body_ }

func (m *MsgWppOut) Products() []ProductEntry { return m.Products_ }

func (m *MsgWppOut) Action() string { return m.Action_ }

func (m *MsgWppOut) Smart() bool { return m.Smart_ }

func (m *MsgWppOut) ProductSearch() string { return m.ProductSearch_ }

func (m *MsgWppOut) SendCatalog() bool { return m.SendCatalog_ }
