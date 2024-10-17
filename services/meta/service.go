package meta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/utils"
)

type service struct {
	httpClient      *http.Client
	httpRetries     *httpx.RetryConfig
	httpAccess      *httpx.AccessConfig
	defaultHeaders  map[string]string
	maxBodyBytes    int
	redactor        utils.Redactor
	systemUserToken string
	url             string
}

// NewServiceFactory creates a new meta service factory
func NewServiceFactory(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, defaultHeaders map[string]string, maxBodyBytes int, systemUserToken string, url string) engine.MetaServiceFactory {
	return func(flows.Session) (flows.MetaService, error) {
		return NewService(httpClient, httpRetries, httpAccess, defaultHeaders, maxBodyBytes, systemUserToken, url), nil
	}
}

// NewService creates a new default meta service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, defaultHeaders map[string]string, maxBodyBytes int, systemUserToken string, url string) flows.MetaService {
	return &service{
		httpClient:      httpClient,
		httpRetries:     httpRetries,
		httpAccess:      httpAccess,
		defaultHeaders:  defaultHeaders,
		maxBodyBytes:    maxBodyBytes,
		redactor:        utils.NewRedactor(flows.RedactionMask, systemUserToken),
		url:             url,
		systemUserToken: systemUserToken,
	}
}

// Filter represents the structure of the filter for the Meta products API request
type Filter struct {
	Or []OrCondition `json:"or"`
}

// OrCondition represents an OR condition
type OrCondition struct {
	And []AndCondition `json:"and"`
}

// AndCondition represents an AND condition
type AndCondition struct {
	RetailerID   map[string]string `json:"retailer_id,omitempty"`
	Availability map[string]string `json:"availability,omitempty"`
	Visibility   map[string]string `json:"visibility,omitempty"`
}

// OrderProductsSearch searches for products in a catalog based on the order items
func (s *service) OrderProductsSearch(order flows.Order, logHTTP flows.HTTPLogCallback) ([]byte, error) {
	filter, err := createFilter(order.ProductItems)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("fields", "[\"category\",\"name\",\"retailer_id\",\"availability\"]")
	params.Add("summary", "true")
	params.Add("access_token", s.systemUserToken)
	params.Add("filter", filter)

	url := fmt.Sprintf("%s/%s/products?%s", s.url, order.CatalogID, params.Encode())

	request, err := httpx.NewRequest("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}

	// set any headers with defaults
	for k, v := range s.defaultHeaders {
		if request.Header.Get(k) == "" {
			request.Header.Set(k, v)
		}
	}

	trace, err := httpx.DoTrace(s.httpClient, request, s.httpRetries, s.httpAccess, s.maxBodyBytes)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode, s.redactor))

		// throw away any error that happened prior to getting a response.. these will be surfaced to the user
		// as connection_error status on the response
		if trace.Response == nil {
			return nil, err
		}

		var responseJSON []byte
		if len(trace.ResponseBody) > 0 {
			responseJSON, _ = webhooks.ExtractJSON(trace.ResponseBody)
		}

		return responseJSON, err
	}

	return nil, err
}

func createFilter(orderitems []flows.ProductItem) (string, error) {
	var filter Filter

	for _, item := range orderitems {
		andCondition := []AndCondition{
			{
				RetailerID: map[string]string{"i_contains": item.ProductRetailerID},
			},
		}
		filter.Or = append(filter.Or, OrCondition{And: andCondition})
	}

	filterJSON, err := json.Marshal(filter)
	if err != nil {
		return "", err
	}

	return string(filterJSON), nil
}

var _ flows.MetaService = (*service)(nil)
