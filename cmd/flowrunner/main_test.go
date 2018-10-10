package main_test

import (
	"github.com/nyaruka/goflow/assets"
	"strings"
	"testing"

	main "github.com/nyaruka/goflow/cmd/flowrunner"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunFlow(t *testing.T) {
	// create an input than can be scanned for two answers
	in := strings.NewReader("I like red\npepsi\n")
	out := &strings.Builder{}

	err := main.RunFlow("testdata/two_questions.json", assets.FlowUUID("615b8a0f-588c-4d20-a05f-363b0b4ce6f4"), in, out)
	require.NoError(t, err)

	assert.Equal(t, "Starting flow 'Two Questions'....\n💬 Hi Ben Haggerty! What is your favorite color? (red/blue)\n> 💬 Red it is! What is your favorite soda? (pepsi/coke)\n> 💬 Great, you are done!\n", out.String())
}
