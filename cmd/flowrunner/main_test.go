package main_test

import (
	"strings"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/test"

	main "github.com/nyaruka/goflow/cmd/flowrunner"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunFlow(t *testing.T) {
	// create an input than can be scanned for two answers
	in := strings.NewReader("I like red\npepsi\n")
	out := &strings.Builder{}

	_, err := main.RunFlow(test.NewEngine(), "testdata/two_questions.json", assets.FlowUUID("615b8a0f-588c-4d20-a05f-363b0b4ce6f4"), "", "eng", in, out)
	require.NoError(t, err)

	// remove input prompts and split output by line to get each event
	lines := strings.Split(strings.Replace(out.String(), "> ", "", -1), "\n")

	assert.Equal(t, []string{
		"Starting flow 'Two Questions'....",
		"---------------------------------------",
		"💬 message created \"Hi Ben Haggerty! What is your favorite color? (red/blue)\"",
		"⏳ waiting for message (600 sec timeout, type /timeout to simulate)....",
		"📥 message received \"I like red\"",
		"📈 run result 'Favorite Color' changed to 'red' with category 'Red'",
		"🌐 language changed to 'fra'",
		"💬 message created \"Red it is! What is your favorite soda? (pepsi/coke)\"",
		"⏳ waiting for message....",
		"📥 message received \"pepsi\"",
		"📈 run result 'Soda' changed to 'pepsi' with category 'Pepsi'",
		"💬 message created \"Great, you are done!\"",
		"",
	}, lines)

	// run again but don't specify the flow
	in = strings.NewReader("I like red\npepsi\n")
	out = &strings.Builder{}
	_, err = main.RunFlow(test.NewEngine(), "testdata/two_questions.json", "", "", "eng", in, out)
	require.NoError(t, err)

	assert.Contains(t, out.String(), "Starting flow 'Two Questions'")
}
