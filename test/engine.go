package test

import (
	"net/http"
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/test/services"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	retries := httpx.NewFixedRetries(1*time.Millisecond, 2*time.Millisecond)
	llmResponses := map[string]string{} // TODO

	return engine.NewBuilder().
		WithMaxFieldChars(256).
		WithEmailServiceFactory(func(s flows.SessionAssets) (flows.EmailService, error) {
			return services.NewEmail(), nil
		}).
		WithWebhookServiceFactory(webhooks.NewServiceFactory(http.DefaultClient, retries, nil, map[string]string{"User-Agent": "goflow-testing"}, 10000)).
		WithClassificationServiceFactory(func(c *flows.Classifier) (flows.ClassificationService, error) {
			return services.NewClassification(c), nil
		}).
		WithLLMServiceFactory(func(l *flows.LLM) (flows.LLMService, error) {
			return services.NewLLM(llmResponses), nil
		}).
		WithAirtimeServiceFactory(func(flows.SessionAssets) (flows.AirtimeService, error) {
			return services.NewAirtime("RWF"), nil
		}).
		Build()
}
