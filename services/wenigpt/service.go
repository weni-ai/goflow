package wenigpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
)

type service struct {
	httpClient     *http.Client
	httpRetries    *httpx.RetryConfig
	httpAccess     *httpx.AccessConfig
	defaultHeaders map[string]string
	maxBodyBytes   int
	token          string
	url            string
}

// NewServiceFactory creates a new wenigpt service factory
func NewServiceFactory(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, defaultHeaders map[string]string, maxBodyBytes int, token string, url string) engine.WeniGPTServiceFactory {
	return func(flows.Session) (flows.WeniGPTService, error) {
		return NewService(httpClient, httpRetries, httpAccess, defaultHeaders, maxBodyBytes, token, url), nil
	}
}

// NewService creates a new default webhook service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, defaultHeaders map[string]string, maxBodyBytes int, token string, url string) flows.WeniGPTService {
	return &service{
		httpClient:     httpClient,
		httpRetries:    httpRetries,
		httpAccess:     httpAccess,
		defaultHeaders: defaultHeaders,
		maxBodyBytes:   maxBodyBytes,
		token:          token,
		url:            url,
	}
}

func (s *service) Call(session flows.Session, input string, contentBaseUUID string) (*flows.WeniGPTCall, error) {

	body := struct {
		Text            string `json:"text"`
		ContentBaseUUID string `json:"content_base_uuid"`
	}{
		Text:            input,
		ContentBaseUUID: contentBaseUUID,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// build our request
	req, err := http.NewRequest("POST", s.url+"/api/v1/wenigpt_question", bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, errors.Wrapf(err, "error building request")
	}

	// set any headers with defaults
	for k, v := range s.defaultHeaders {
		if req.Header.Get(k) == "" {
			req.Header.Set(k, v)
		}
	}

	req.Header.Set("Content-Type", "application/json")

	if s.token != "" {
		req.Header.Set("Authorization", "Bearer "+s.token)
	} else {
		return nil, fmt.Errorf("validation token cannot be empty")
	}

	trace, err := httpx.DoTrace(s.httpClient, req, s.httpRetries, s.httpAccess, s.maxBodyBytes)
	if trace != nil {
		call := &flows.WeniGPTCall{Trace: trace}
		// throw away any error that happened prior to getting a response.. these will be surfaced to the user
		// as connection_error status on the response
		if trace.Response == nil {
			return call, err
		}

		if len(call.ResponseBody) > 0 {
			call.ResponseJSON, call.ResponseCleaned = ExtractJSON(call.ResponseBody)
		}

		return call, err
	}

	return nil, errors.Wrapf(err, "")
}

func ExtractJSON(body []byte) ([]byte, bool) {
	// we make a best effort to turn the body into JSON, so we strip out:
	//  1. any invalid UTF-8 sequences
	//  2. null chars
	//  3. escaped null chars (\u0000)
	cleaned := bytes.ToValidUTF8(body, nil)
	cleaned = bytes.ReplaceAll(cleaned, []byte{0}, nil)
	cleaned = utils.ReplaceEscapedNulls(cleaned, nil)

	if json.Valid(cleaned) {
		changed := !bytes.Equal(body, cleaned)
		return cleaned, changed
	}
	return nil, false
}
