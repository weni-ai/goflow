package engine_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/engine"
	"github.com/stretchr/testify/assert"
)

func TestEmptyServices(t *testing.T) {
	// default engine configation provides no services for anything
	eng := engine.NewBuilder().Build()

	webhookSvc, err := eng.Services().Webhook(nil)
	assert.EqualError(t, err, "no webhook service factory configured")
	assert.Nil(t, webhookSvc)

	brainSvc, err := eng.Services().Brain(nil)
	assert.EqualError(t, err, "no brain service factory configured")
	assert.Nil(t, brainSvc)

	classificationSvc, err := eng.Services().Classification(nil, nil)
	assert.EqualError(t, err, "no classification service factory configured")
	assert.Nil(t, classificationSvc)

	airtimeSvc, err := eng.Services().Airtime(nil)
	assert.EqualError(t, err, "no airtime service factory configured")
	assert.Nil(t, airtimeSvc)

	msgCatalogSvc, err := eng.Services().MsgCatalog(nil, nil)
	assert.EqualError(t, err, "no msg catalog service factory configured")
	assert.Nil(t, msgCatalogSvc)

	orgContextSvc, err := eng.Services().OrgContext(nil, nil)
	assert.EqualError(t, err, "no org context service factory configured")
	assert.Nil(t, orgContextSvc)

	metaSvc, err := eng.Services().Meta(nil)
	assert.EqualError(t, err, "no meta service factory configured")
	assert.Nil(t, metaSvc)
}
