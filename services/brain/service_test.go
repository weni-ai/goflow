package brain_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type call struct {
	projectUUID uuids.UUID
	text        string
	contactURN  urns.URN
	attachments []utils.Attachment
}

type brain struct {
	request  string
	response string
	body     string
	bodyJSON string
}

func TestBrain(t *testing.T) {
	server := test.NewTestHTTPServer(49994)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
	require.NoError(t, err)

	testCases := []struct {
		call  call
		brain brain
	}{
		{
			call: call{uuids.UUID("dac08478-75a5-465c-8ddd-c631197a8c94"), "What is a cookie?", urns.URN("whatsapp:123456789"), []utils.Attachment{"image:https://img.img"}},
			brain: brain{
				request:  "POST /messages?token=token HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 157\r\nAccept-Encoding: gzip\r\n\r\n{\"project_uuid\":\"dac08478-75a5-465c-8ddd-c631197a8c94\",\"text\":\"What is a cookie?\",\"contact_urn\":\"whatsapp:123456789\",\"attachments\":[\"image:https://img.img\"]}",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 0\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
				bodyJSON: `{ "ok": "true" }`,
			},
		},
	}

	for _, tc := range testCases {
		svc, _ := session.Engine().Services().Brain(session)
		c, err := svc.Call(session, tc.call.projectUUID, tc.call.text, tc.call.contactURN, tc.call.attachments)
		assert.NoError(t, err, nil)

		assert.Equal(t, tc.brain.request, string(c.RequestTrace), "request trace mismatch for call %s", tc.call)
		assert.Equal(t, tc.brain.response, string(c.ResponseTrace), "response mismatch for call %s", tc.call)

	}

}
