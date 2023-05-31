package zeroshot_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/services/classification/zeroshot"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

func TestPredict(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://zeroshot.it/v2/repository/nlp/zeroshot/zeroshot-predict": {
			httpx.NewMockResponse(200, nil, `xx`), // non-JSON response
			httpx.NewMockResponse(200, nil, `{}`), // invalid JSON response
			httpx.NewMockResponse(200, nil, `{
				"text": "book_flight"
			  }`),
		},
	}))

	client := zeroshot.NewClient(http.DefaultClient, nil, "123e4567-e89b-12d3-a456-426655440000", "656b0715-9071-4782-8cce-4f33d0bf7c38")

	response, trace, err := client.Predict("book flight to Quito")
	assert.EqualError(t, err, `invalid character 'x' looking for beginning of value`)
	test.AssertSnapshot(t, "predict_request", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, "xx", string(trace.ResponseBody))
	assert.Nil(t, response)

	response, trace, err = client.Predict("book flight to Quito")
	fmt.Println(err)
	assert.EqualError(t, err, `field 'text' is required`)
	assert.NotNil(t, trace)
	assert.Nil(t, response)

	response, trace, err = client.Predict("book flight to Quito")
	assert.NoError(t, err)
	assert.NotNil(t, trace)
	assert.Equal(t, "book_flight", response.Text)

}
