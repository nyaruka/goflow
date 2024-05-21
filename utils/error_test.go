package utils_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
)

func TestRichError(t *testing.T) {
	e := utils.NewRichError("goats", "bad_attitude", "Bad Attitude").WithExtra("name", "billy")

	assert.Equal(t, "Bad Attitude", e.Error())
	assert.Equal(t, "goats", e.Domain)
	assert.Equal(t, "bad_attitude", e.Code)
	assert.Equal(t, map[string]string{"name": "billy"}, e.Extra)

	e1 := fmt.Errorf("wrapped twice: %w", fmt.Errorf("wrapped once: %w", e))
	isRichError, cause := utils.IsRichError(e1)
	assert.True(t, isRichError)
	assert.Equal(t, e, cause)

	e2 := fmt.Errorf("wrapped twice: %w", fmt.Errorf("wrapped once: %w", fmt.Errorf("not a query error")))
	isRichError, cause = utils.IsRichError(e2)
	assert.False(t, isRichError)
	assert.Nil(t, cause)
}
