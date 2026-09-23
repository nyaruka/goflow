package services_test

import (
	"testing"

	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/test/services"
	"github.com/stretchr/testify/assert"
)

func TestModelService(t *testing.T) {
	svc := services.NewModel()
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
	assert.EqualError(t, err, "simulated model error")

	_, err = svc.Response(ctx, "Translate as JSON", `{"a":["Hi","error"]}`, 100)
	assert.EqualError(t, err, "simulated model error")

	// invalid JSON input with a JSON translate instruction errors
	_, err = svc.Response(ctx, "Translate as JSON", "not json", 100)
	assert.Error(t, err)

	// Categorize instructions pick the last word
	resp, err = svc.Response(ctx, "Categorize into [A, B, C]", "input", 100)
	assert.NoError(t, err)
	assert.Equal(t, "C", resp.Output)

}

func TestModelServiceClassify(t *testing.T) {
	svc := services.NewModel()
	ctx := t.Context()
	options := []*core.ClassifierOption{{Name: "Flights"}, {Name: "Hotels"}}

	// plain input chooses the last option
	cls, err := svc.Classify(ctx, "I want to book a room", options)
	assert.NoError(t, err)
	assert.Equal(t, "Hotels", cls.Option)
	assert.Equal(t, 0.8, cls.Confidence)
	assert.Equal(t, map[string]float64{"Flights": 0.1, "Hotels": 0.9}, cls.Probabilities)
}

func TestMockModel(t *testing.T) {
	ctx := t.Context()

	svc := services.NewMockModel(
		&services.MockModelResult{Output: "Bonjour", TokensInput: 12, TokensOutput: 3},
		&services.MockModelResult{Option: "Flights", Confidence: 0.7},
		&services.MockModelResult{Error: "boom"},
	)
	assert.True(t, svc.HasUnused())

	resp, err := svc.Response(ctx, "Translate to French", "Hello", 100)
	assert.NoError(t, err)
	assert.Equal(t, &core.ModelResponse{Output: "Bonjour", TokensInput: 12, TokensOutput: 3}, resp)

	options := []*core.ClassifierOption{{Name: "Flights"}, {Name: "Hotels"}}

	cls, err := svc.Classify(ctx, "I want to fly to Paris", options)
	assert.NoError(t, err)
	assert.Equal(t, &core.Classification{Option: "Flights", Confidence: 0.7}, cls)

	_, err = svc.Classify(ctx, "Hi", options)
	assert.EqualError(t, err, "boom")

	assert.False(t, svc.HasUnused())
	assert.Equal(t, []*services.ModelCall{
		{Instructions: "Translate to French", Input: "Hello", MaxTokens: 100},
		{Input: "I want to fly to Paris", Options: options},
		{Input: "Hi", Options: options},
	}, svc.Calls())

	// running out of results, or a result that doesn't fit the call, is a test setup mistake
	assert.PanicsWithValue(t, "missing mock model result for call with input 'Hi'", func() { svc.Response(ctx, "Summarize", "Hi", 100) })
	assert.Panics(t, func() {
		services.NewMockModel(&services.MockModelResult{Option: "Cars"}).Classify(ctx, "Hi", options)
	})
	assert.Panics(t, func() {
		services.NewMockModel(&services.MockModelResult{Option: "Cars"}).Response(ctx, "Summarize", "Hi", 100)
	})
}
