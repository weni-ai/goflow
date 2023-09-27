package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/require"
)

func TestCatalogMsg(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	msg := flows.NewMsgCatalog(
		urns.URN("whatsapp:+558299988776655"),
		assets.NewChannelReference(assets.ChannelUUID("1fbd661a-d519-42c5-9bbe-d12a6da9621b"), "My WPP Cloud"),
		"Header text",
		"Body text",
		"Footer text",
		[]string{
			"524580ca-406d-491b-b97e-b07113e322db",
			"ee48c9ed-6e52-4a76-8e7b-70c5046c559a",
		},
		"View Products",
		false,
		flows.MsgTopic("none"),
	)

	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"action": "View Products",
		"body": "Body text",
		"channel": {
				"name": "My WPP Cloud",
				"uuid": "1fbd661a-d519-42c5-9bbe-d12a6da9621b"
		},
		"footer": "Footer text",
		"header": "Header text",
		"products": [
				"524580ca-406d-491b-b97e-b07113e322db",
				"ee48c9ed-6e52-4a76-8e7b-70c5046c559a"
		],
		"smart": false,
		"text": "",
		"topic": "none",
		"urn": "whatsapp:+558299988776655",
		"uuid": "1ae96956-4b34-433e-8d1a-f05fe6923d6d"
	}`), marshaled)
}

func TestCatalogMsgSmart(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	msg := flows.NewMsgCatalog(
		urns.URN("whatsapp:+558299988776655"),
		assets.NewChannelReference(assets.ChannelUUID("1fbd661a-d519-42c5-9bbe-d12a6da9621b"), "My WPP Cloud"),
		"Header text",
		"Body text",
		"Footer text",
		[]string{},
		"View Products",
		true,
		flows.MsgTopic("none"),
	)

	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"action": "View Products",
		"body": "Body text",
		"channel": {
				"name": "My WPP Cloud",
				"uuid": "1fbd661a-d519-42c5-9bbe-d12a6da9621b"
		},
		"footer": "Footer text",
		"header": "Header text",
		"smart": true,
		"text": "",
		"topic": "none",
		"urn": "whatsapp:+558299988776655",
		"uuid": "1ae96956-4b34-433e-8d1a-f05fe6923d6d"
	}`), marshaled)
}
