package static_test

import (
	"testing"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/stretchr/testify/assert"
)

func TestOrgContext(t *testing.T) {
	oc := static.NewOrgContext("context 1", assets.ChannelUUID("e3c07fe2-7542-42c6-a394-18d968999f51"), uuids.UUID("80cb8631-b84a-4949-a483-a3103f00edc2"), false, false)

	assert.Equal(t, "context 1", oc.Context())
	assert.Equal(t, assets.ChannelUUID("e3c07fe2-7542-42c6-a394-18d968999f51"), oc.ChannelUUID())
	assert.Equal(t, uuids.UUID("80cb8631-b84a-4949-a483-a3103f00edc2"), oc.ProjectUUID())

}
