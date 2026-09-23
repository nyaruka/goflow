package services_test

import (
	"testing"

	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/test/services"
	"github.com/stretchr/testify/assert"
)

func TestLLMService(t *testing.T) {
	svc := services.NewLLM()
	ctx := t.Context()

	// plain input is echoed with instructions
	resp, err := svc.Response(ctx, "Summarize", "Hello", 100)
	assert.NoError(t, err)
	assert.Equal(t, "You asked:\n\nSummarize\n\nHello", resp.Output)

	// instructions starting with "Translate" leetify the whole input
	resp, err = svc.Response(ctx, "Translate to Spanish", "Hello", 100)
	assert.NoError(t, err)
	assert.Equal(t, "H3110", resp.Output)

	// lower-case "translate" does NOT trigger — only the capitalized prefix
	resp, err = svc.Response(ctx, "please translate this", "Hello", 100)
	assert.NoError(t, err)
	assert.Equal(t, "You asked:\n\nplease translate this\n\nHello", resp.Output)

	// "Translate" instructions mentioning "JSON" parse the input as a string->[]string object and leetify values
	resp, err = svc.Response(ctx, "Translate as JSON", `{"greeting":["Hello","Hi"],"name":["World"]}`, 100)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"greeting":["H3110","H1"],"name":["W0r1d"]}`, resp.Output)

	// values exactly equal to "untranslatable" become "<CANT>"
	resp, err = svc.Response(ctx, "Translate to Spanish", "untranslatable", 100)
	assert.NoError(t, err)
	assert.Equal(t, "<CANT>", resp.Output)

	resp, err = svc.Response(ctx, "Translate as JSON", `{"a":["Hi","untranslatable"]}`, 100)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"a":["H1","<CANT>"]}`, resp.Output)

	// values exactly equal to "error" cause the service to error
	_, err = svc.Response(ctx, "Translate to Spanish", "error", 100)
	assert.EqualError(t, err, "simulated LLM error")

	_, err = svc.Response(ctx, "Translate as JSON", `{"a":["Hi","error"]}`, 100)
	assert.EqualError(t, err, "simulated LLM error")

	// invalid JSON input with a JSON translate instruction errors
	_, err = svc.Response(ctx, "Translate as JSON", "not json", 100)
	assert.Error(t, err)

	// Categorize instructions pick the last word
	resp, err = svc.Response(ctx, "Categorize into [A, B, C]", "input", 100)
	assert.NoError(t, err)
	assert.Equal(t, "C", resp.Output)

}

func TestLLMServiceClassify(t *testing.T) {
	svc := services.NewLLM()
	ctx := t.Context()
	categories := []string{"Flights", "Hotels"}

	// plain input chooses the last category
	cls, err := svc.Classify(ctx, "I want to book a room", categories)
	assert.NoError(t, err)
	assert.Equal(t, "Hotels", cls.Category)
	assert.Equal(t, 0.8, cls.Confidence)
	assert.Equal(t, map[string]float64{"Flights": 0.1, "Hotels": 0.9}, cls.Probabilities)
}

func TestMockLLM(t *testing.T) {
	ctx := t.Context()

	svc := services.NewMockLLM(
		&services.MockLLMResult{Output: "Bonjour", TokensInput: 12, TokensOutput: 3},
		&services.MockLLMResult{Category: "Flights", Confidence: 0.7},
		&services.MockLLMResult{Error: "boom"},
	)
	assert.True(t, svc.HasUnused())

	resp, err := svc.Response(ctx, "Translate to French", "Hello", 100)
	assert.NoError(t, err)
	assert.Equal(t, &core.LLMResponse{Output: "Bonjour", TokensInput: 12, TokensOutput: 3}, resp)

	cls, err := svc.Classify(ctx, "I want to fly to Paris", []string{"Flights", "Hotels"})
	assert.NoError(t, err)
	assert.Equal(t, &core.LLMClassification{Category: "Flights", Confidence: 0.7}, cls)

	_, err = svc.Classify(ctx, "Hi", []string{"Flights", "Hotels"})
	assert.EqualError(t, err, "boom")

	assert.False(t, svc.HasUnused())
	assert.Equal(t, []*services.LLMCall{
		{Instructions: "Translate to French", Input: "Hello", MaxTokens: 100},
		{Input: "I want to fly to Paris", Categories: []string{"Flights", "Hotels"}},
		{Input: "Hi", Categories: []string{"Flights", "Hotels"}},
	}, svc.Calls())

	// running out of results, or a result that doesn't fit the call, is a test setup mistake
	assert.PanicsWithValue(t, "missing mock LLM result for call with input 'Hi'", func() { svc.Response(ctx, "Summarize", "Hi", 100) })
	assert.Panics(t, func() {
		services.NewMockLLM(&services.MockLLMResult{Category: "Cars"}).Classify(ctx, "Hi", []string{"Flights"})
	})
	assert.Panics(t, func() {
		services.NewMockLLM(&services.MockLLMResult{Category: "Cars"}).Response(ctx, "Summarize", "Hi", 100)
	})
}
