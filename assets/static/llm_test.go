package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestLLM(t *testing.T) {
	llm := static.NewLLM(
		assets.LLMUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"),
		"GPT-4",
		"openai",
	)
	assert.Equal(t, assets.LLMUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"), llm.UUID())
	assert.Equal(t, "GPT-4", llm.Name())
	assert.Equal(t, "openai", llm.Type())
}
