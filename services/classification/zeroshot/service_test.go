package zeroshot_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/classification/zeroshot"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	session, _ := test.NewSessionBuilder().MustBuild()

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://zeroshot.it/v2/repository/nlp/zeroshot/zeroshot-predict": {
			httpx.NewMockResponse(200, nil, `{
				"text": "book_flight"
			  }`),
		},
	}))

	svc := zeroshot.NewService(
		http.DefaultClient,
		nil,
		test.NewClassifier("Booking", "zeroshot", []string{"book_flight", "book_hotel"}),
		"f96abf2f-3b53-4766-8ea6-09a655222a02",
		"bc9641ed-6467-4737-8d66-7887edaadb94",
	)

	httpLogger := &flows.HTTPLogger{}

	classification, err := svc.Classify(session, "book flight to Quito", httpLogger.Log)
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		{Name: "book_flight"},
	}, classification.Intents)

	assert.Equal(t, 1, len(httpLogger.Logs))
	assert.Equal(t, "https://zeroshot.it/v2/repository/nlp/zeroshot/zeroshot-predict", httpLogger.Logs[0].URL)
	assert.Equal(t, "POST /v2/repository/nlp/zeroshot/zeroshot-predict HTTP/1.1\r\nHost: zeroshot.it\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 135\r\nAccept-Encoding: gzip\r\n\r\n{\"repository_uuid\":\"bc9641ed-6467-4737-8d66-7887edaadb94\",\"text\":\"book flight to Quito\",\"token\":\"****************\"}", httpLogger.Logs[0].Request)
}
