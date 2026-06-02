package test

import (
	"context"
	"net/http"
	"strings"
	"text/template"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/test/services"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	return newEngine(http.DefaultClient)
}

// NewMockedEngine creates an engine instance for testing whose webhook calls are answered from the given mocks
func NewMockedEngine(mocks map[string][]*httpx.MockResponse) flows.Engine {
	client, _ := MockedHTTP(mocks)
	return newEngine(client)
}

func newEngine(httpClient *http.Client) flows.Engine {
	return engine.NewBuilder().
		WithHTTPClient(httpClient).
		WithMaxFieldChars(256).
		WithLLMPrompts(map[string]*template.Template{
			"categorize": template.Must(template.New("").Parse("Categorize the following text into one of the following: {{ .arg1 }}")),
		}).
		WithWebhookServiceFactory(webhooks.NewServiceFactory(map[string]string{"User-Agent": "goflow-testing"}, 10000)).
		WithEmailServiceFactory(func(s flows.SessionAssets) (flows.EmailService, error) {
			return services.NewEmail(), nil
		}).
		WithLLMServiceFactory(func(l *flows.LLM) (flows.LLMService, error) {
			return services.NewLLM(), nil
		}).
		WithAirtimeServiceFactory(func(flows.SessionAssets) (flows.AirtimeService, error) {
			return services.NewAirtime("RWF"), nil
		}).
		WithCheckSendable(func(ctx context.Context, sa flows.SessionAssets, c *flows.Contact, mc *flows.MsgContent) (flows.UnsendableReason, error) {
			return "", nil
		}).
		WithClaimURN(func(ctx context.Context, sa flows.SessionAssets, c *flows.Contact, u urns.URN) (bool, error) {
			return !strings.Contains(u.Path(), "taken"), nil
		}).
		Build()
}
