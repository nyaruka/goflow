package wit_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/nlu/wit"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*http.Response{
		"https://api.wit.ai/message?v=20170307&q=book+flight+to+Quito": []*http.Response{
			httpx.NewMockResponse(200, `{"_text":"book flight to Quito","entities":{"intent":[{"confidence":0.84709152161066,"value":"book_flight"}]},"msg_id":"1M7fAcDWag76OmgDI"}`),
		},
	}))
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	session, _, err := test.CreateTestSession("", nil, envs.RedactionPolicyNone)
	require.NoError(t, err)

	svc := wit.NewService(test.NewClassifier("Booking", "wit", []string{"book_flight", "book_hotel"}), "23532624376")

	events := make([]flows.Event, 0)
	log := func(e flows.Event) { events = append(events, e) }

	classification, err := svc.Classify(session, "book flight to Quito", log)
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		flows.ExtractedIntent{Name: "book_flight", Confidence: decimal.RequireFromString(`0.84709152161066`)},
	}, classification.Intents)
	assert.Equal(t, map[string][]flows.ExtractedEntity{}, classification.Entities)
}
