package excellent_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent"

	"github.com/stretchr/testify/assert"
)

func TestHasExpressions(t *testing.T) {
	topLevels := []string{"foo"}

	assert.False(t, excellent.HasExpressions("", topLevels))
	assert.False(t, excellent.HasExpressions("hi there", topLevels))
	assert.False(t, excellent.HasExpressions("bob@gmail", topLevels))
	assert.False(t, excellent.HasExpressions("@(", topLevels))
	assert.True(t, excellent.HasExpressions("@foo", topLevels))
	assert.True(t, excellent.HasExpressions("hi @foo.x", topLevels))
	assert.True(t, excellent.HasExpressions("hi @(foo)", topLevels))
}
