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
	eng := engine.NewBuilder().
		WithMaxStepsPerSprint(123).
		WithMaxSprintsPerSession(567).
		WithMaxTemplateChars(999).
		WithMaxFieldChars(888).
		WithMaxResultChars(777).
		Build()

	assert.Equal(t, 123, eng.Options().MaxStepsPerSprint)
	assert.Equal(t, 567, eng.Options().MaxSprintsPerSession)
	assert.Equal(t, 999, eng.Options().MaxTemplateChars)
	assert.Equal(t, 888, eng.Options().MaxFieldChars)
	assert.Equal(t, 777, eng.Options().MaxResultChars)

	_, err := eng.Services().Email(nil)
	assert.EqualError(t, err, "no email service factory configured")
	_, err = eng.Services().Airtime(nil)
	assert.EqualError(t, err, "no airtime service factory configured")
	_, err = eng.Services().Classification(nil)
	assert.EqualError(t, err, "no classification service factory configured")
	_, err = eng.Services().Webhook(nil)
	assert.EqualError(t, err, "no webhook service factory configured")

	// include a webhook service
	webhookSvc := webhooks.NewService(&http.Client{}, nil, nil, map[string]string{"User-Agent": "goflow"}, 1000)

	eng = engine.NewBuilder().
		WithWebhookServiceFactory(func(flows.SessionAssets) (flows.WebhookService, error) { return webhookSvc, nil }).
		Build()

	svc, err := eng.Services().Webhook(nil)
	assert.NoError(t, err)
	assert.Equal(t, webhookSvc, svc)
}
