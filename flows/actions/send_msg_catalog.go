package actions

import (
	"fmt"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
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

	MsgCatalog *assets.MsgCatalogReference `json:"msg_catalog,omitempty"`
	AllURNs    bool                        `json:"all_urns,omitempty"`
	Templating *Templating                 `json:"templating,omitempty" validate:"omitempty,dive"`
	Topic      flows.MsgTopic              `json:"topic,omitempty" validate:"omitempty,msg_topic"`
	ResultName string                      `json:"result_name,omitempty"`
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

		/////////////////////

		msgCatalog := run.Session().Assets().MsgCatalogs()
		mc := msgCatalog.Get(a.MsgCatalog.UUID)
		params := assets.NewMsgCatalogParam(evaluatedHeader, evaluatedBody, evaluatedFooter, products, a.ProductViewSettings.Action, string(a.Topic), a.AutomaticSearch, evaluatedSearch, envs.NilLanguage, uuids.UUID(dest.Channel.UUID()))
		c, err := a.call(run, step, params, mc, logEvent)
		if err != nil {
			a.saveResult(run, step, a.ResultName, fmt.Sprintf("%s", err), CategoryFailure, "", "", nil, logEvent)
		}
		a.saveResult(run, step, a.ResultName, string(c.ResponseJSON), CategorySuccess, "", "", c.ResponseJSON, logEvent)

		////////////////////

		products = c.ProductRetailerIDS

		msg := flows.NewMsgCatalogOut(dest.URN.URN(), channelRef, evaluatedHeader, evaluatedBody, evaluatedFooter, a.ProductViewSettings.Action, evaluatedSearch, products, a.AutomaticSearch, a.Topic)
		logEvent(events.NewMsgCatalogCreated(msg))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgCatalogOut(urns.NilURN, nil, evaluatedHeader, evaluatedBody, evaluatedFooter, a.ProductViewSettings.Action, evaluatedSearch, products, a.AutomaticSearch, a.Topic)
		logEvent(events.NewMsgCatalogCreated(msg))
	}
	return nil
}

var msgCatalogCategories = []string{CategorySuccess, CategoryFailure}

func (a *SendMsgCatalogAction) Results(include func(*flows.ResultInfo)) {
	include(flows.NewResultInfo(a.ResultName, msgCatalogCategories))
}

func (a *SendMsgCatalogAction) call(run flows.FlowRun, step flows.Step, params assets.MsgCatalogParam, msgCatalog *flows.MsgCatalog, logEvent flows.EventCallback) (*flows.MsgCatalogCall, error) {
	if msgCatalog == nil {
		logEvent(events.NewDependencyError(a.MsgCatalog))
		return nil, nil
	}

	svc, err := run.Session().Engine().Services().MsgCatalog(run.Session(), msgCatalog)
	if err != nil {
		logEvent(events.NewError(err))
		return nil, nil
	}

	httpLogger := &flows.HTTPLogger{}

	call, err := svc.Call(run.Session(), params, httpLogger.Log)
	if err != nil {
		logEvent(events.NewError(err))
		return nil, err
	}

	return call, nil
}
