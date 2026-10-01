package brain_test

import (
	"testing"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type call struct {
	projectUUID uuids.UUID
	text        string
	contact     *flows.Contact
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
			call: call{
				uuids.UUID("dac08478-75a5-465c-8ddd-c631197a8c94"),
				"What is a cookie?",
				session.Contact(),
				[]utils.Attachment{"image:https://img.img"},
			},
			brain: brain{
				request:  "POST /messages?token=token HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 479\r\nAccept-Encoding: gzip\r\n\r\n{\"project_uuid\":\"dac08478-75a5-465c-8ddd-c631197a8c94\",\"contact_urn\":\"tel:+12024561111\",\"text\":\"What is a cookie?\",\"attachments\":[\"image:https://img.img\"],\"channel_uuid\":\"57f1078f-88aa-46f4-a59a-948a5739c03d\",\"contact_fields\":{\"activation_token\":{\"value\":\"AACC55\",\"type\":\"text\"},\"age\":{\"value\":23,\"type\":\"number\"},\"gender\":{\"value\":\"Male\",\"type\":\"text\"},\"join_date\":{\"value\":\"2017-12-02T00:00:00-02:00\",\"type\":\"datetime\"},\"not_set\":null,\"state\":null},\"contact_name\":\"Ryan Lewis\"}",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 0\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
				bodyJSON: `{ "ok": "true" }`,
			},
		},
	}

	for _, tc := range testCases {
		svc, _ := session.Engine().Services().Brain(session)
		c, err := svc.Call(session, tc.call.projectUUID, tc.call.text, tc.call.contact, tc.call.attachments)
		assert.NoError(t, err, nil)

		assert.Equal(t, tc.brain.request, string(c.RequestTrace), "request trace mismatch for call %s", tc.call)
		assert.Equal(t, tc.brain.response, string(c.ResponseTrace), "response mismatch for call %s", tc.call)

	}

}
