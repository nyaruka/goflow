package wit_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://api.wit.ai/message?v=20170307&q=book+flight+to+Quito": []httpx.MockResponse{
			httpx.NewMockResponse(200, `{"_text":"book flight to Quito","entities":{"intent":[{"confidence":0.84709152161066,"value":"book_flight"}]},"msg_id":"1M7fAcDWag76OmgDI"}`),
		},
	}))

	svc := wit.NewService(test.NewClassifier("Booking", "wit", []string{"book_flight", "book_hotel"}), "23532624376")

	classification, traces, err := svc.Classify(session, "book flight to Quito")
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		flows.ExtractedIntent{Name: "book_flight", Confidence: decimal.RequireFromString(`0.84709152161066`)},
	}, classification.Intents)
	assert.Equal(t, map[string][]flows.ExtractedEntity{}, classification.Entities)

	assert.Equal(t, 1, len(traces))
	assert.Equal(t, "https://api.wit.ai/message?v=20170307&q=book+flight+to+Quito", traces[0].Request.URL.String())
}
