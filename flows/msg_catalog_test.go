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

func TestMsgCatalogOut(t *testing.T) {

	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	msg := flows.NewMsgCatalogOut(
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "My Android"),
		"header",
		"body",
		"footer",
		"action",
		"productSearch",
		[]string{"product_1", "product_2"},
		false,
		flows.MsgTopic("none"),
		false,
	)

	marshaled, err := jsonx.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"action": "action",
		"body": "body",
		"channel": {
			"name": "My Android",
			"uuid": "61f38f46-a856-4f90-899e-905691784159"
		},
		"footer": "footer",
		"header": "header",
		"product_search": "productSearch",
		"products": [
			"product_1",
			"product_2"
		],
		"smart": false,
		"text": "",
		"topic": "none",
		"urn": "tel:+1234567890",
		"uuid": "1ae96956-4b34-433e-8d1a-f05fe6923d6d"
	}`), marshaled, "JSON mismatch")
}
