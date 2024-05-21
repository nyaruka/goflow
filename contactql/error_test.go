package contactql_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/contactql"
	"github.com/stretchr/testify/assert"
)

func TestQueryError(t *testing.T) {
	e := contactql.NewQueryError("bad_query", "Bad Query")

	assert.Equal(t, "Bad Query", e.Error())

	e1 := fmt.Errorf("wrapped twice: %w", fmt.Errorf("wrapped once: %w", e))
	isQueryError, cause := contactql.IsQueryError(e1)
	assert.True(t, isQueryError)
	assert.Equal(t, e, cause)

	e2 := fmt.Errorf("wrapped twice: %w", fmt.Errorf("wrapped once: %w", fmt.Errorf("not a query error")))
	isQueryError, cause = contactql.IsQueryError(e2)
	assert.False(t, isQueryError)
	assert.Nil(t, cause)
}
