package webhooks_test

import (
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type call struct {
	method string
	url    string
	body   string
}

func (c *call) String() string { return c.method + " " + c.url }

type webhook struct {
	request  string
	response string
	body     string
	bodyJSON string
}

func TestWebhookTimeoutHeader(t *testing.T) {
	tcs := []struct {
		name           string
		header         string
		delay          time.Duration
		clientTimeout  time.Duration
		expectResponse bool
		minElapsed     time.Duration
		maxElapsed     time.Duration
	}{
		{
			name:           "header 60 succeeds with 35s delay",
			header:         "60",
			delay:          35 * time.Second,
			clientTimeout:  30 * time.Second,
			expectResponse: true,
			minElapsed:     34 * time.Second,
			maxElapsed:     45 * time.Second,
		},
		{
			name:           "default client timeout fails with 35s delay",
			header:         "",
			delay:          35 * time.Second,
			clientTimeout:  30 * time.Second,
			expectResponse: false,
			minElapsed:     28 * time.Second,
			maxElapsed:     38 * time.Second,
		},
		{
			name:           "header 60 fails with 65s delay at ~60s not ~30s",
			header:         "60",
			delay:          65 * time.Second,
			clientTimeout:  30 * time.Second,
			expectResponse: false,
			minElapsed:     58 * time.Second,
			maxElapsed:     68 * time.Second,
		},
		{
			name:           "header 30 fails with 35s delay at ~30s",
			header:         "30",
			delay:          35 * time.Second,
			clientTimeout:  60 * time.Second,
			expectResponse: false,
			minElapsed:     28 * time.Second,
			maxElapsed:     38 * time.Second,
		},
		{
			name:           "header 45 succeeds with 40s delay",
			header:         "45",
			delay:          40 * time.Second,
			clientTimeout:  30 * time.Second,
			expectResponse: true,
			minElapsed:     39 * time.Second,
			maxElapsed:     48 * time.Second,
		},
		{
			name:           "invalid header uses client default",
			header:         "not-a-number",
			delay:          35 * time.Second,
			clientTimeout:  30 * time.Second,
			expectResponse: false,
			minElapsed:     28 * time.Second,
			maxElapsed:     38 * time.Second,
		},
		{
			name:           "header above 60 is clamped to 60",
			header:         "120",
			delay:          55 * time.Second,
			clientTimeout:  30 * time.Second,
			expectResponse: true,
			minElapsed:     54 * time.Second,
			maxElapsed:     62 * time.Second,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tc.delay)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"ok":true}`))
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: tc.clientTimeout}
			svc := webhooks.NewService(httpClient, nil, nil, nil, 10000)

			request, err := http.NewRequest("GET", server.URL, nil)
			require.NoError(t, err)
			if tc.header != "" {
				request.Header.Set("X-Weni-Webhook-Timeout", tc.header)
			}

			start := time.Now()
			call, err := svc.Call(nil, request)
			elapsed := time.Since(start)

			assert.NoError(t, err)
			require.NotNil(t, call)

			if tc.expectResponse {
				require.NotNil(t, call.Response)
				assert.Equal(t, http.StatusOK, call.Response.StatusCode)
			} else {
				assert.Nil(t, call.Response)
			}

			assert.GreaterOrEqual(t, elapsed, tc.minElapsed, "elapsed too short: %s", elapsed)
			assert.LessOrEqual(t, elapsed, tc.maxElapsed, "elapsed too long: %s", elapsed)

			assert.Empty(t, call.Request.Header.Get("X-Weni-Webhook-Timeout"))
		})
	}
}

func TestWebhookTimeoutHeaderConcurrent(t *testing.T) {
	const (
		shortDelay = 100 * time.Millisecond
		longDelay  = 2 * time.Second
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		delay, _ := time.ParseDuration(r.URL.Query().Get("delay"))
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 500 * time.Millisecond}
	svc := webhooks.NewService(httpClient, nil, nil, nil, 10000)

	const workers = 20
	errs := make(chan error, workers*2)
	done := make(chan struct{}, workers)

	for i := 0; i < workers; i++ {
		go func(useLongTimeout bool) {
			defer func() { done <- struct{}{} }()

			delay := shortDelay
			header := "30"
			if useLongTimeout {
				delay = longDelay
				header = "60"
			}

			request, err := http.NewRequest("GET", server.URL+"?delay="+delay.String(), nil)
			if err != nil {
				errs <- err
				return
			}
			request.Header.Set("X-Weni-Webhook-Timeout", header)

			call, err := svc.Call(nil, request)
			if err != nil {
				errs <- err
				return
			}
			if call.Response == nil {
				errs <- assert.AnError
			}
		}(i%2 == 0)
	}

	for i := 0; i < workers; i++ {
		<-done
	}
	close(errs)

	for err := range errs {
		assert.NoError(t, err)
	}
}

func TestClampWebhookTimeout(t *testing.T) {
	tcs := []struct {
		seconds  int
		expected time.Duration
	}{
		{10, 30 * time.Second},
		{30, 30 * time.Second},
		{45, 45 * time.Second},
		{60, 60 * time.Second},
		{120, 60 * time.Second},
	}

	for _, tc := range tcs {
		assert.Equal(t, tc.expected, webhooks.ClampWebhookTimeout(tc.seconds), "seconds=%d", tc.seconds)
	}
}

func TestWebhookParsing(t *testing.T) {
	server := test.NewTestHTTPServer(49994)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	testCases := []struct {
		call    call
		webhook webhook
		isError bool
	}{
		{
			// successful GET
			call: call{"GET", "http://127.0.0.1:49994/?cmd=success", ""},
			webhook: webhook{
				request:  "GET /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
				bodyJSON: `{ "ok": "true" }`,
			},
		}, {
			// successful GET with valid JSON response body
			call: call{"GET", "http://127.0.0.1:49994/?cmd=textjs", ""},
			webhook: webhook{
				request:  "GET /?cmd=textjs HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/javascript; charset=iso-8859-1\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
				bodyJSON: `{ "ok": "true" }`,
			},
		}, {
			// successful POST without request body and valid JSON response body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=success", ""},
			webhook: webhook{
				request:  "POST /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
				bodyJSON: `{ "ok": "true" }`,
			},
		}, {
			// successful POST with request body and valid JSON response body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=success", `{"contact": "Bob"}`},
			webhook: webhook{
				request:  "POST /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 18\r\nAccept-Encoding: gzip\r\n\r\n{\"contact\": \"Bob\"}",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
				bodyJSON: `{ "ok": "true" }`,
			},
		}, {
			// successful GET with JSON response body containing escaped null chars (actual escaped nulls should be replaced with \ufffd)
			call: call{"GET", "http://127.0.0.1:49994/?cmd=badjson", ""},
			webhook: webhook{
				request:  "GET /?cmd=badjson HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 67\r\nContent-Type: application/json\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     "{ \"bad\": \"null=\x00 escaped=\\u0000 double-escaped=\\\\u0000 badseq=\x80\x81\" }",
				bodyJSON: "{ \"bad\": \"null= escaped= double-escaped=\\\\u0000 badseq=\" }",
			},
		}, {
			// successful POST receiving gzipped non-JSON body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=gzipped&content=Hello", ``},
			webhook: webhook{
				request:  "POST /?cmd=gzipped&content=Hello HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Type: application/x-gzip\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `Hello`,
				bodyJSON: ``,
			},
		}, {
			// successful POST receiving gzipped JSON body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=gzipped&content=%7B%22contact%22%3A%20%22Bob%22%7D", ``},
			webhook: webhook{
				request:  "POST /?cmd=gzipped&content=%7B%22contact%22%3A%20%22Bob%22%7D HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Type: application/x-gzip\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{"contact": "Bob"}`,
				bodyJSON: `{"contact": "Bob"}`,
			},
		}, {
			// POST returning 503
			call: call{"POST", "http://127.0.0.1:49994/?cmd=unavailable", ""},
			webhook: webhook{
				request:  "POST /?cmd=unavailable HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 503 Service Unavailable\r\nContent-Length: 37\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "errors": ["service unavailable"] }`,
				bodyJSON: `{ "errors": ["service unavailable"] }`,
			},
		}, {
			// GET returning text body larger than allowed
			call:    call{"GET", "http://127.0.0.1:49994/?cmd=binary&size=11000", ""},
			isError: true,
		}, {
			// GET returning non-JSON body
			call: call{"GET", "http://127.0.0.1:49994/?cmd=typeless&content=kthxbai", ""},
			webhook: webhook{
				request:  "GET /?cmd=typeless&content=kthxbai HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 7\r\nContent-Type: \r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     "kthxbai",
				bodyJSON: ``,
			},
		}, {
			// GET returning JSON body but an empty content-type header
			call: call{"GET", "http://127.0.0.1:49994/?cmd=typeless&content=%7B%22msg%22%3A%20%22I%27m%20JSON%22%7D", ""},
			webhook: webhook{
				request:  "GET /?cmd=typeless&content=%7B%22msg%22%3A%20%22I%27m%20JSON%22%7D HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 19\r\nContent-Type: \r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{"msg": "I'm JSON"}`,
				bodyJSON: `{"msg": "I'm JSON"}`,
			},
		}, {
			// connection error
			call: call{"POST", "http://127.0.0.1:55555/", ""},
			webhook: webhook{
				request:  "POST / HTTP/1.1\r\nHost: 127.0.0.1:55555\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "",
				body:     "",
				bodyJSON: ``,
			},
		},
	}

	for _, tc := range testCases {
		request, err := http.NewRequest(tc.call.method, tc.call.url, strings.NewReader(tc.call.body))
		require.NoError(t, err)

		svc, _ := session.Engine().Services().Webhook(session)
		c, err := svc.Call(session, request)

		if tc.isError {
			assert.Error(t, err, "expected error for call %s", tc.call)
		} else {
			assert.NoError(t, err, "unexpected error fetching %s", tc.call)

			assert.Equal(t, tc.call.url, c.Request.URL.String(), "URL mismatch for call %s", tc.call)
			assert.Equal(t, tc.call.method, c.Request.Method, "method mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.request, string(c.RequestTrace), "request trace mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.response, string(c.ResponseTrace), "response mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.body, string(c.ResponseBody), "body mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.bodyJSON, string(c.ResponseJSON), "body JSON mismatch for call %s", tc.call)
		}
	}
}

func TestRetries(t *testing.T) {
	session, _ := test.NewSessionBuilder().MustBuild()

	defer httpx.SetRequestor(httpx.DefaultRequestor)

	mocks := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(502, nil, "a"),
			httpx.NewMockResponse(200, nil, "b"),
		},
	})
	httpx.SetRequestor(mocks)

	request, err := http.NewRequest("GET", "http://temba.io/", strings.NewReader("BODY"))
	require.NoError(t, err)

	svc, _ := session.Engine().Services().Webhook(session)
	c, err := svc.Call(session, request)
	require.NoError(t, err)

	assert.Equal(t, 200, c.Response.StatusCode)
	assert.Equal(t, "GET / HTTP/1.1\r\nHost: temba.io\r\nUser-Agent: goflow-testing\r\nContent-Length: 4\r\nAccept-Encoding: gzip\r\n\r\nBODY", string(c.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 1\r\n\r\n", string(c.ResponseTrace))
	assert.Equal(t, "b", string(c.ResponseBody))
}

func TestAccessRestrictions(t *testing.T) {
	retries := httpx.NewFixedRetries(5, 10)
	access := httpx.NewAccessConfig(10, []net.IP{net.IPv4(127, 0, 0, 1)}, nil)

	factory := webhooks.NewServiceFactory(http.DefaultClient, retries, access, map[string]string{"User-Agent": "Foo"}, 12345)
	svc, err := factory(nil)
	assert.NoError(t, err)

	request, _ := http.NewRequest("GET", "http://localhost/foo", nil)
	call, err := svc.Call(nil, request)

	// actual error becomes a call with a connection error
	assert.NoError(t, err)

	// should still have a trace.. just no response part
	assert.Equal(t, "GET /foo HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Foo\r\nAccept-Encoding: gzip\r\n\r\n", string(call.RequestTrace))
	assert.Equal(t, "", string(call.ResponseTrace))
}

func TestGzipEncoding(t *testing.T) {
	session, _ := test.NewSessionBuilder().MustBuild()

	defer dates.SetNowSource(dates.DefaultNowSource)

	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))

	server := test.NewTestHTTPServer(52025)

	request, err := http.NewRequest("GET", server.URL+"?cmd=gzipped&content=Hello", nil)
	require.NoError(t, err)

	request.Header.Set("Accept-Encoding", "gzip")

	svc, _ := session.Engine().Services().Webhook(session)
	c, err := svc.Call(session, request)
	require.NoError(t, err)

	// check that gzip decompression happens transparently
	assert.Equal(t, 200, c.Response.StatusCode)
	assert.Equal(t, "GET /?cmd=gzipped&content=Hello HTTP/1.1\r\nHost: 127.0.0.1:52025\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n", string(c.RequestTrace))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Type: application/x-gzip\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n", string(c.ResponseTrace))
	assert.Equal(t, "Hello", string(c.ResponseBody))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Type: application/x-gzip\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\nHello", string(c.SanitizedResponse("...")))
}

func TestExtractJSON(t *testing.T) {
	tcs := []struct {
		body []byte
		json []byte
	}{
		{[]byte(`{`), nil}, // invalid JSON
		{[]byte(`"x"`), []byte(`"x"`)},
		{[]byte(`{"foo": ["x"]}`), []byte(`{"foo": ["x"]}`)},
		{[]byte("\"a\x80\x81b\""), []byte(`"ab"`)},                     // invalid UTF-8 sequences stripped
		{[]byte("\u0000{\"foo\": 123\u0000}"), []byte(`{"foo": 123}`)}, // null chars stripped
		{[]byte(`"a\u0000b"`), []byte(`"ab"`)},                         // escaped null chars stripped
		{[]byte(`"01\02\03"`), nil},                                    // \0 not valid JSON escape
		{[]byte(`"01\\02\\03"`), []byte(`"01\\02\\03"`)},
	}

	for _, tc := range tcs {
		actual, changed := webhooks.ExtractJSON(tc.body)
		assert.Equal(t, string(tc.json), string(actual), "extracted JSON mismatch for %s", string(tc.body))
		if len(actual) > 0 {
			assert.Equal(t, !bytes.Equal(tc.body, tc.json), changed)
		}
	}

	asXValue := types.JSONToXValue([]byte(`{"foo": "01\\02\\03"}`))
	asXObject := asXValue.(*types.XObject)
	foo, _ := asXObject.Get("foo")
	assert.Equal(t, types.NewXText(`01\02\03`), foo)
	assert.Equal(t, `"01\\02\\03"`, string(jsonx.MustMarshal(foo)))
}

func TestWebhookResponseWithEscapes(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	mocks := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://cheapcontactlookups.com": {
			httpx.NewMockResponse(200, nil, `{"name": "01\\02\\03", "joined": "04\\05\\06"}`),
		},
	})
	httpx.SetRequestor(mocks)

	session, _ := test.NewSessionBuilder().
		WithAssetsPath("testdata/webhook_flow.json").
		WithFlow("bb38eefb-3cd9-4f80-9867-9c84ae276f7a").MustBuild()

	joined := session.Assets().Fields().Get("joined")

	assert.Equal(t, flows.SessionStatusCompleted, session.Status())
	assert.Equal(t, `01\02\03`, session.Contact().Name())
	assert.Equal(t, types.NewXText(`04\05\06`), session.Contact().Fields().Get(joined).Text)

	// check nothing became an escaped NULL
	assert.NotContains(t, string(jsonx.MustMarshal(session)), `\u0000`)
}
