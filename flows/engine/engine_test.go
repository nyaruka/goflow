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
	webhookSvc := webhooks.NewService(&http.Client{}, "goflow", 1000)

	eng := engine.NewBuilder().
		WithWebhookServiceFactory(func(flows.Session) (flows.WebhookService, error) { return webhookSvc, nil }).
		WithMaxStepsPerSprint(123).
		Build()

	assert.Equal(t, 123, eng.MaxStepsPerSprint())

	svc, err := eng.Services().Webhook(nil)
	assert.NoError(t, err)
	assert.Equal(t, webhookSvc, svc)
}
