package actions

import (
	"fmt"
	"regexp"

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
	searchSettings

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
	SendCatalog         bool                `json:"sendCatalog"`
}

type ProductViewSettings struct {
	Header string `json:"header" engine:"evaluated"`
	Body   string `json:"body" engine:"evaluated"`
	Footer string `json:"footer" engine:"evaluated"`
	Action string `json:"action"`
}

type searchSettings struct {
	SearchUrl  string `json:"search_url,omitempty"`
	SearchType string `json:"search_type"`
	PostalCode string `json:"postal_code"`
	SellerId   string `json:"seller_id"`
}

// NewSendMsgCatalog creates a new send msg catalog action
func NewSendMsgCatalog(uuid flows.ActionUUID, header, body, footer, action, productSearch string, products []map[string]string, automaticSearch bool, searchUrl string, searchType string, postalCode string, sellerId string, allURNs bool) *SendMsgCatalogAction {
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
		searchSettings: searchSettings{
			SearchUrl:  searchUrl,
			SearchType: searchType,
			PostalCode: postalCode,
			SellerId:   sellerId,
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
	if evaluatedSearch == "" && a.AutomaticSearch {
		logEvent(events.NewErrorf("search text evaluated to empty string"))
	}

	ProductEntries := []flows.ProductEntry{}
	for _, p := range a.Products {
		v, found := p["product_retailer_id"]
		if found {
			productEntry := flows.ProductEntry{
				Product:            "product_retailer_id",
				ProductRetailerIDs: []string{v},
			}
			ProductEntries = append(ProductEntries, productEntry)
		}
	}

	evaluatedHeader, evaluatedBody, evaluatedFooter, evaluatedPostalCode, evaluatedURL, evaluatedSellerId := a.evaluateMessageCatalog(run, nil, a.ProductViewSettings.Header, a.ProductViewSettings.Body, a.ProductViewSettings.Footer, a.Products, a.SendCatalog, a.PostalCode, a.SearchUrl, a.SellerId, logEvent)

	destinations := run.Contact().ResolveDestinations(a.AllURNs)

	// create a new message for each URN+channel destination
	for _, dest := range destinations {
		var channelRef *assets.ChannelReference
		if dest.Channel != nil {
			channelRef = assets.NewChannelReference(dest.Channel.UUID(), dest.Channel.Name())
		}
		var apiType string
		if a.SearchType == "vtex" {
			regexLegacy := regexp.MustCompile(`^https:\/\/([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)\.com(\.br)?\/api\/catalog_system\/pub\/products\/search$`)
			regexLegacySeller := regexp.MustCompile(`^https:\/\/([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)\.com(\.br)?\/api\/checkout\/pub\/orderForms\/simulation$`)
			regexIntelligent := regexp.MustCompile(`^https:\/\/([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)\.com(\.br)?\/api\/io\/_v\/api\/intelligent-search\/product_search(\/)?([\\?&]([^&=]+)=([^&=]+))?$`)
			regexSponsored := regexp.MustCompile(`^https:\/\/([a-zA-Z0-9_-]+)\.([a-zA-Z0-9_-]+)\.com(\.br)?\/api\/io\/_v\/api\/intelligent-search\/sponsored_products(\/)?$`)
			if regexLegacy.MatchString(evaluatedURL) {
				apiType = "legacy"
			} else if regexIntelligent.MatchString(evaluatedURL) {
				apiType = "intelligent"
			} else if regexLegacySeller.MatchString(evaluatedURL) {
				apiType = "legacy"
			} else if regexSponsored.MatchString(evaluatedURL) {
				apiType = "sponsored"
			}
		}

		if a.createMsgCatalogAction.AutomaticSearch {
			msgCatalog := run.Session().Assets().MsgCatalogs()
			mc := msgCatalog.GetByChannelUUID(channelRef.UUID)
			if mc == nil {
				a.saveResult(run, step, a.ResultName, fmt.Sprintf("channel with uuid: %s, does not have an active catalog", channelRef.UUID), CategoryFailure, "", "", nil, logEvent)
				return nil
			}

			orgContext := run.Session().Assets().OrgContext()
			context := orgContext.GetHasVtexAdsByChannelUUID()
			var hasVtexAds bool
			if context != nil {
				hasVtexAds = context.OrgContext.HasVtexAds()
			}

			language := "eng"
			if run.Contact().Language() != envs.NilLanguage && run.Contact().Language() != "base" {
				language = string(run.Contact().Language())
			}

			params := assets.NewMsgCatalogParam(evaluatedSearch, uuids.UUID(dest.Channel.UUID()), a.SearchType, evaluatedURL, apiType, evaluatedPostalCode, evaluatedSellerId, hasVtexAds, language)
			c, err := a.call(run, step, params, mc, logEvent)
			if err != nil {
				if c != nil {
					for _, trace := range c.Traces {
						call := &flows.WebhookCall{Trace: trace}
						logEvent(events.NewWebhookCalled(call, callStatus(call, nil, true), ""))
					}
				}
				a.saveResult(run, step, a.ResultName, fmt.Sprintf("%s", err), CategoryFailure, "", "", nil, logEvent)
				return nil
			}
			for _, trace := range c.Traces {
				call := &flows.WebhookCall{Trace: trace}
				logEvent(events.NewWebhookCalled(call, callStatus(call, nil, true), ""))
			}
			a.saveResult(run, step, a.ResultName, string(c.ResponseJSON), CategorySuccess, "", "", c.ResponseJSON, logEvent)
			ProductEntries = c.ProductRetailerIDS
		} else {
			a.saveResult(run, step, a.ResultName, "", CategorySuccess, "", "", nil, logEvent)
		}

		msg := flows.NewMsgCatalogOut(dest.URN.URN(), channelRef, evaluatedHeader, evaluatedBody, evaluatedFooter, a.ProductViewSettings.Action, evaluatedSearch, ProductEntries, a.AutomaticSearch, a.Topic, a.SendCatalog)
		logEvent(events.NewMsgCatalogCreated(msg))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgCatalogOut(urns.NilURN, nil, evaluatedHeader, evaluatedBody, evaluatedFooter, a.ProductViewSettings.Action, evaluatedSearch, ProductEntries, a.AutomaticSearch, a.Topic, a.SendCatalog)
		logEvent(events.NewMsgCatalogCreated(msg))
	}
	return nil
}

var msgCatalogCategories = []string{CategorySuccess, CategoryFailure}

func (a *SendMsgCatalogAction) Results(include func(*flows.ResultInfo)) {
	include(flows.NewResultInfo(a.ResultName, msgCatalogCategories))
}

func (a *SendMsgCatalogAction) call(run flows.FlowRun, step flows.Step, params assets.MsgCatalogParam, msgCatalog *flows.MsgCatalog, logEvent flows.EventCallback) (*flows.MsgCatalogCall, error) {
	var call *flows.MsgCatalogCall

	if msgCatalog == nil {
		logEvent(events.NewDependencyError(a.MsgCatalog))
		return call, fmt.Errorf("msgCatalog cannot be nil")
	}

	svc, err := run.Session().Engine().Services().MsgCatalog(run.Session(), msgCatalog)
	if err != nil {
		logEvent(events.NewError(err))
		return call, err
	}

	httpLogger := &flows.HTTPLogger{}

	call, err = svc.Call(run.Session(), params, httpLogger.Log)
	if err != nil {
		logEvent(events.NewError(err))
		return call, err
	}

	if len(call.ProductRetailerIDS) == 0 {
		err = fmt.Errorf("product out of stock")
		logEvent(events.NewError(err))
		return call, err
	}

	return call, nil
}
