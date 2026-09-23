package core_test

import (
	"testing"

	"github.com/nyaruka/goflow/core"
	"github.com/stretchr/testify/assert"
)

func TestLLMClassificationCategoryProbability(t *testing.T) {
	cls := &core.LLMClassification{Category: "Hotels", Probabilities: map[string]float64{"Flights": 0.08, "Hotels": 0.92}}
	p, ok := cls.CategoryProbability()
	assert.True(t, ok)
	assert.Equal(t, 0.92, p)

	// generative models don't provide probabilities
	_, ok = (&core.LLMClassification{Category: "Hotels"}).CategoryProbability()
	assert.False(t, ok)

	_, ok = (*core.LLMClassification)(nil).CategoryProbability()
	assert.False(t, ok)
}
