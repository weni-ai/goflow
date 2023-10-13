package actions

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendMsgCatalog, func() flows.Action { return &SendMsgCatalogAction{} })
}

// TypeSendMsgCatalog is the type for the send message catalog action
const TypeSendMsgCatalog string = "send_msg_catalog"

// SendMsgAction can be used to send a catalog of products or services to the current contact in a flow. The header, body or footer fields may contain templates. The action
// will attempt to find pairs of URNs and channels which can be used for sending. If it can't find such a pair, it will
// create a message without a channel or URN.
//
// A [event:msg_catalog_created] event will be created with the evaluated fields.
//
//		{
//		  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//		  "type": "send_msg_catalog",
//		  "products": [
//	        {
//						  "product_retailer_id": "e3e84fc4-5320-4321-bd3c-0dd2bd068189"
//			    }
//			],
//			"productViewSettings": {
//				"header": "Header text",
//				"body": "Body text",
//				"footer": "Footer text",
//			  "action": "view",
//			},
//		  "topic": "event",
//			"automaticProductSearch": "false",
//			"productSearch": "some product",
//	    "result_name": "Result"
//		}
//
// @action send_msg_catalog
type SendMsgCatalogAction struct {
	baseAction
	universalAction
	createMsgCatalogAction

	AllURNs    bool           `json:"all_urns,omitempty"`
	Templating *Templating    `json:"templating,omitempty" validate:"omitempty,dive"`
	Topic      flows.MsgTopic `json:"topic,omitempty" validate:"omitempty,msg_topic"`
	ResultName string         `json:"result_name,omitempty"`
}

type createMsgCatalogAction struct {
	Products            []map[string]string `json:"products"`
	AutomaticSearch     bool                `json:"automaticProductSearch"`
	ProductSearch       string              `json:"productSearch"`
	ProductViewSettings ProductViewSettings `json:"productViewSettings"`
}

type ProductViewSettings struct {
	Header string `json:"header" engine:"evaluated"`
	Body   string `json:"body" engine:"evaluated"`
	Footer string `json:"footer" engine:"evaluated"`
	Action string `json:"action"`
}

// NewSendMsgCatalog creates a new send msg catalog action
func NewSendMsgCatalog(uuid flows.ActionUUID, header, body, footer, action, productSearch string, products []map[string]string, automaticSearch, allURNs bool) *SendMsgCatalogAction {
	return &SendMsgCatalogAction{
		baseAction: newBaseAction(TypeSendMsgCatalog, uuid),
		createMsgCatalogAction: createMsgCatalogAction{
			ProductViewSettings: ProductViewSettings{
				Header: header,
				Body:   body,
				Footer: footer,
				Action: action,
			},
			Products:        products,
			AutomaticSearch: automaticSearch,
			ProductSearch:   productSearch,
		},
		AllURNs: allURNs,
	}
}

// Execute runs this action
func (a *SendMsgCatalogAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	if run.Contact() == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	evaluatedSearch, err := run.EvaluateTemplate(a.ProductSearch)
	if err != nil {
		logEvent(events.NewError(err))
	}
	if evaluatedSearch == "" {
		logEvent(events.NewErrorf("search text evaluated to empty string"))
	}

	evaluatedHeader, evaluatedBody, evaluatedFooter := a.evaluateMessageCatalog(run, nil, a.ProductViewSettings.Header, a.ProductViewSettings.Body, a.ProductViewSettings.Footer, logEvent)

	destinations := run.Contact().ResolveDestinations(a.AllURNs)

	var products []string
	for _, p := range a.Products {
		v, found := p["product_retailer_id"]
		if found {
			products = append(products, v)
		}
	}

	// create a new message for each URN+channel destination
	for _, dest := range destinations {
		var channelRef *assets.ChannelReference
		if dest.Channel != nil {
			channelRef = assets.NewChannelReference(dest.Channel.UUID(), dest.Channel.Name())
		}

		msg := flows.NewMsgCatalog(dest.URN.URN(), channelRef, evaluatedHeader, evaluatedBody, evaluatedFooter, a.ProductViewSettings.Action, evaluatedSearch, products, a.AutomaticSearch, a.Topic)
		logEvent(events.NewMsgCatalogCreated(msg))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgCatalog(urns.NilURN, nil, evaluatedHeader, evaluatedBody, evaluatedFooter, a.ProductViewSettings.Action, evaluatedSearch, products, a.AutomaticSearch, a.Topic)
		logEvent(events.NewMsgCatalogCreated(msg))
	}

	// TODO: saveResult with CategoryFailure
	a.saveResult(run, step, a.ResultName, "SUCCESS RESULT", CategorySuccess, "", "SUCCESS RESULT", nil, logEvent)

	return nil
}

var msgCatalogCategories = []string{CategorySuccess, CategoryFailure}

func (a *SendMsgCatalogAction) Results(include func(*flows.ResultInfo)) {
	include(flows.NewResultInfo(a.ResultName, msgCatalogCategories))
}
