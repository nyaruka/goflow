package engine_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
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
		WithWebhookLimits(666, 555).
		Build()

	assert.Equal(t, 123, eng.Options().MaxStepsPerSprint)
	assert.Equal(t, 567, eng.Options().MaxSprintsPerSession)
	assert.Equal(t, 999, eng.Options().MaxTemplateChars)
	assert.Equal(t, 888, eng.Options().MaxFieldChars)
	assert.Equal(t, 777, eng.Options().MaxResultChars)
	assert.Equal(t, 666, eng.Options().MaxRequestBytes)
	assert.Equal(t, 555, eng.Options().MaxResponseBytes)

	_, err := eng.Services().Email(nil)
	assert.EqualError(t, err, "no email service factory configured")
	_, err = eng.Services().Airtime(nil)
	assert.EqualError(t, err, "no airtime service factory configured")
	_, err = eng.Services().Webhook(nil)
	assert.EqualError(t, err, "no webhook service factory configured")

	// include a webhook service
	webhookSvc := webhooks.NewService(&http.Client{}, map[string]string{"User-Agent": "goflow"}, 1000)

	eng = engine.NewBuilder().
		WithWebhookServiceFactory(func(flows.Engine, flows.SessionAssets) (flows.WebhookService, error) { return webhookSvc, nil }).
		Build()

	svc, err := eng.Services().Webhook(nil)
	assert.NoError(t, err)
	assert.Equal(t, webhookSvc, svc)

	// create engine with a tiny evaluation budget
	eng = engine.NewBuilder().WithEvaluationBudget(100).Build()

	env := envs.NewBuilder().Build()
	root := types.NewXObject(map[string]types.XValue{})

	_, _, err = eng.Evaluator().Template(t.Context(), env, root, `@(repeat("x", 50))`, nil)
	assert.NoError(t, err)
	_, _, err = eng.Evaluator().Template(t.Context(), env, root, `@(repeat("x", 200))`, nil)
	assert.ErrorContains(t, err, "expression is too complex to evaluate")
}
