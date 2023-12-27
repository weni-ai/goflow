package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows"
	"github.com/stretchr/testify/assert"
)

func TestOrgContext(t *testing.T) {
	oc1 := static.NewOrgContext("context 1", assets.ChannelUUID("e3c07fe2-7542-42c6-a394-18d968999f51"))

	oc := flows.NewOrgContextAssets([]assets.OrgContext{oc1})

	o1 := oc.GetByChannelUUID()

	assert.Equal(t, "context 1", o1.Context())
	assert.Equal(t, assets.ChannelUUID("e3c07fe2-7542-42c6-a394-18d968999f51"), o1.ChannelUUID())
	assert.Equal(t, oc1, o1.Asset())
	assert.Equal(t, assets.NewOrgContextReference("context 1"), o1.Reference())

	// nil object returns nil reference
	assert.Nil(t, (*flows.OrgContext)(nil).Reference())
}
