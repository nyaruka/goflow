package excellent_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	// context callback is optional
	exp, err := excellent.Parse(`foo + 1`, nil)
	assert.NoError(t, err)
	assert.IsType(t, &excellent.Addition{}, exp)

	var paths [][]string
	exp, err = excellent.Parse(`foo.bar + 1`, func(p []string) { paths = append(paths, p) })
	assert.NoError(t, err)
	assert.IsType(t, &excellent.Addition{}, exp)
	assert.Equal(t, [][]string{{"foo"}, {"foo", "bar"}}, paths)

	// if errors occur during parsing, first is returned
	_, err = excellent.Parse(`(foo +)`, nil)
	assert.EqualError(t, err, "syntax error at )")
}

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
