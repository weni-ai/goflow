package wenigpt_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type call struct {
	kb       string
	input    string
	language string
}

type wenigpt struct {
	request  string
	response string
	body     string
	bodyJSON string
}

func TestWeniGPT(t *testing.T) {
	server := test.NewTestHTTPServerWeniGPT(49994)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	testCases := []struct {
		call    call
		wenigpt wenigpt
	}{
		{
			call: call{"59a246a0-775f-49df-a734-ac95220d96ca", "What is a cookie?", "por"},
			wenigpt: wenigpt{
				request:  "POST /api/v1/wenigpt_question HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 104\r\nAuthorization: Bearer token\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"text\":\"What is a cookie?\",\"content_base_uuid\":\"59a246a0-775f-49df-a734-ac95220d96ca\",\"language\":\"por\"}",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 46\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{"answers":[{"text":"Is a cookie."}],"id":"0"}`,
				bodyJSON: `{"answers":[{"text":"Is a cookie."}],"id":"0"}`,
			},
		},
	}

	for _, tc := range testCases {
		svc, _ := session.Engine().Services().WeniGPT(session)
		c, err := svc.Call(session, tc.call.input, tc.call.kb, tc.call.language)

		assert.NoError(t, err, nil)

		assert.Equal(t, tc.wenigpt.request, string(c.RequestTrace), "request trace mismatch for call %s", tc.call)
		assert.Equal(t, tc.wenigpt.response, string(c.ResponseTrace), "response mismatch for call %s", tc.call)
		assert.Equal(t, tc.wenigpt.body, string(c.ResponseBody), "body mismatch for call %s", tc.call)
		assert.Equal(t, tc.wenigpt.bodyJSON, string(c.ResponseJSON), "body JSON mismatch for call %s", tc.call)

	}

}
