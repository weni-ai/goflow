package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
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

// NewServiceFactory creates a new brain service factory
func NewServiceFactory(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, defaultHeaders map[string]string, maxBodyBytes int, token string, url string) engine.BrainServiceFactory {
	return func(flows.Session) (flows.BrainService, error) {
		return NewService(httpClient, httpRetries, httpAccess, defaultHeaders, maxBodyBytes, token, url), nil
	}
}

// NewService creates a new default webhook service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, httpAccess *httpx.AccessConfig, defaultHeaders map[string]string, maxBodyBytes int, token string, url string) flows.BrainService {
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

func (s *service) Call(session flows.Session, projectUUID uuids.UUID, text string, contactURN urns.URN, attachments []utils.Attachment) (*flows.BrainCall, error) {

	body := struct {
		ProjectUUID uuids.UUID         `json:"project_uuid"`
		Text        string             `json:"text"`
		ContactURN  urns.URN           `json:"contact_urn"`
		Attachments []utils.Attachment `json:"attachments"`
	}{
		ProjectUUID: projectUUID,
		Text:        text,
		ContactURN:  contactURN.Identity(),
		Attachments: attachments,
	}

	var b io.Reader
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	b = bytes.NewReader(bodyJSON)

	params := url.Values{}
	params.Add("token", s.token)
	url_ := fmt.Sprintf("%s/messages?%s", s.url, params.Encode())
	req, err := httpx.NewRequest("POST", url_, b, nil)
	if err != nil {
		return nil, err
	}

	// set any headers with defaults
	for k, v := range s.defaultHeaders {
		if req.Header.Get(k) == "" {
			req.Header.Set(k, v)
		}
	}

	trace, err := httpx.DoTrace(s.httpClient, req, s.httpRetries, s.httpAccess, s.maxBodyBytes)
	if trace != nil {
		call := &flows.BrainCall{Trace: trace}
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
