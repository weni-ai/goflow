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
	KnowledgeBase string `json:"knowledge_base"`
	Input         string `json:"input"`
}

// NewServiceFactory creates a new wenigpt service factory
func NewServiceFactory(kb string, input string) engine.WeniGPTServiceFactory {
	return func(flows.Session) (flows.WeniGPTService, error) {
		return NewService(kb, input), nil
	}
}

// NewService creates a new default webhook service
func NewService(kb string, input string) flows.WeniGPTService {
	return &service{
		KnowledgeBase: kb,
		Input:         input,
	}
}

func (s *service) Call(session flows.Session, input string, contentBaseUUID string, token string, url string) (*flows.WeniGPTCall, error) {

	body := struct {
		Message         string `json:"message"`
		ContentBaseUUID string `json:"content_base_uuid"`
	}{
		Message:         input,
		ContentBaseUUID: contentBaseUUID,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// build our request
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, errors.Wrapf(err, "error building request")
	}
	req.Header.Add("Content-Type", "application/json")

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	} else {
		return nil, fmt.Errorf("validation token cannot be empty")
	}

	client := &http.Client{}

	trace, err := httpx.DoTrace(client, req, nil, nil, -1)
	if trace != nil {
		call := &flows.WeniGPTCall{Trace: trace}
		// throw away any error that happened prior to getting a response.. these will be surfaced to the user
		// as connection_error status on the response
		if trace.Response == nil {
			return call, nil
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
