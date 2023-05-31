package zeroshot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/utils"
)

const (
	apiBaseURL = "https://zeroshot.it"
)

// Intent is possible intent match
type Intent struct {
	ID     int    `json:"id"`
	Text   string `json:"text" validate:"required"`
	Option string `json:"option"`
}

type PredictionResponse struct {
	Intent Intent
}

// Client is a basic zeroshot client
type Client struct {
	httpClient  *http.Client
	httpRetries *httpx.RetryConfig
	token       string
	repository  string
	//saber se é preciso add repository tbm
}

// NewClient creates a new client
func NewClient(httpClient *http.Client, httpRetries *httpx.RetryConfig, token string, repository string) *Client {
	return &Client{
		httpClient:  httpClient,
		httpRetries: httpRetries,
		token:       token,
		repository:  repository,
	}
}

// Predict
func (c *Client) Predict(q string) (*Intent, *httpx.Trace, error) {
	endpoint := fmt.Sprintf("%s/v2/repository/nlp/zeroshot/zeroshot-predict", apiBaseURL)

	body, _ := json.Marshal(map[string]string{
		"repository_uuid": c.repository, //repository_uuid
		"text":            q,
		"token":           c.token,
	})

	request, err := httpx.NewRequest("POST", endpoint, bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, nil, err
	}

	trace, err := httpx.DoTrace(c.httpClient, request, c.httpRetries, nil, -1)
	if err != nil {
		return nil, trace, err
	}

	if trace.Response != nil && trace.Response.StatusCode == 200 {
		response := &Intent{}
		if err := utils.UnmarshalAndValidate(trace.ResponseBody, response); err != nil {
			return nil, trace, err
		}

		return response, trace, nil
	}

	return nil, trace, errors.New("zeroshot API request failed")
}
