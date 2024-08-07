package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/require"
)

func TestMsgWppOut(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	msg := flows.NewMsgWppOut(
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "Whatsapp Cloud Dummy Channel"),
		"",
		"",
		"",
		"Hi there.",
		"Footer.",
		flows.CTAMessage{},
		flows.ListMessage{},
		flows.FlowMessage{},
		[]utils.Attachment{
			utils.Attachment("image/jpeg:https://example.com/test.jpg"),
			utils.Attachment("audio/mp3:https://example.com/test.mp3"),
		},
		nil,
		flows.MsgTopicAgent,
	)

	// test marshaling our msg
	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"uuid": "1ae96956-4b34-433e-8d1a-f05fe6923d6d",
		"urn": "tel:+1234567890",
		"channel": {"uuid":"61f38f46-a856-4f90-899e-905691784159", "name":"Whatsapp Cloud Dummy Channel"},
		"text": "Hi there.",
		"footer": "Footer.",
		"cta_message": {},
		"list_message": {},
		"flow_message": {},
		"attachments": ["image/jpeg:https://example.com/test.jpg", "audio/mp3:https://example.com/test.mp3"],
		"topic": "agent"
	}`), marshaled, "JSON mismatch")
}
