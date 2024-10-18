package engine_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/services/webhooks"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	// create engine with no services
	eng := engine.NewBuilder().WithMaxStepsPerSprint(123).WithMaxResumesPerSession(567).Build()

	assert.Equal(t, 123, eng.MaxStepsPerSprint())
	assert.Equal(t, 567, eng.MaxResumesPerSession())

	_, err := eng.Services().Email(nil)
	assert.EqualError(t, err, "no email service factory configured")
	_, err = eng.Services().Airtime(nil)
	assert.EqualError(t, err, "no airtime service factory configured")
	_, err = eng.Services().Classification(nil, nil)
	assert.EqualError(t, err, "no classification service factory configured")
	_, err = eng.Services().Ticket(nil, nil)
	assert.EqualError(t, err, "no ticket service factory configured")
	_, err = eng.Services().Webhook(nil)
	assert.EqualError(t, err, "no webhook service factory configured")
	_, err = eng.Services().MsgCatalog(nil, nil)
	assert.EqualError(t, err, "no msg catalog service factory configured")
	_, err = eng.Services().OrgContext(nil, nil)
	assert.EqualError(t, err, "no org context service factory configured")
	_, err = eng.Services().Brain(nil)
	assert.EqualError(t, err, "no brain service factory configured")
	_, err = eng.Services().Meta(nil)
	assert.EqualError(t, err, "no meta service factory configured")

	// include a webhook service
	webhookSvc := webhooks.NewService(&http.Client{}, nil, nil, map[string]string{"User-Agent": "goflow"}, 1000)

	eng = engine.NewBuilder().
		WithWebhookServiceFactory(func(flows.Session) (flows.WebhookService, error) { return webhookSvc, nil }).
		Build()

	svc, err := eng.Services().Webhook(nil)
	assert.NoError(t, err)
	assert.Equal(t, webhookSvc, svc)
}
