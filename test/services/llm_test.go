package services_test

import (
	"context"
	"testing"

	"github.com/nyaruka/goflow/test/services"
	"github.com/stretchr/testify/assert"
)

func TestLLMService(t *testing.T) {
	svc := services.NewLLM()
	ctx := context.Background()

	// plain input is echoed with instructions
	resp, err := svc.Response(ctx, "Translate", "Hello", 100)
	assert.NoError(t, err)
	assert.Equal(t, "You asked:\n\nTranslate\n\nHello", resp.Output)

	// \return directive returns what follows
	resp, err = svc.Response(ctx, "whatever", "\\return foo", 100)
	assert.NoError(t, err)
	assert.Equal(t, "foo", resp.Output)

	// \error directive returns an error
	_, err = svc.Response(ctx, "whatever", "\\error boom", 100)
	assert.EqualError(t, err, "boom")

	// Categorize instructions pick the last word
	resp, err = svc.Response(ctx, "Categorize into [A, B, C]", "input", 100)
	assert.NoError(t, err)
	assert.Equal(t, "C", resp.Output)

	// directive in the first element of a JSON array input fires
	resp, err = svc.Response(ctx, "whatever", `["\\return [\"T-Hi\"]"]`, 100)
	assert.NoError(t, err)
	assert.Equal(t, `["T-Hi"]`, resp.Output)

	_, err = svc.Response(ctx, "whatever", `["\\error boom","ignored"]`, 100)
	assert.EqualError(t, err, "boom")

	// JSON array without a directive falls through to echo
	resp, err = svc.Response(ctx, "whatever", `["Hi","Bye"]`, 100)
	assert.NoError(t, err)
	assert.Equal(t, "You asked:\n\nwhatever\n\n[\"Hi\",\"Bye\"]", resp.Output)

	// non-JSON input starting with [ falls through
	resp, err = svc.Response(ctx, "whatever", "[not json", 100)
	assert.NoError(t, err)
	assert.Equal(t, "You asked:\n\nwhatever\n\n[not json", resp.Output)
}
