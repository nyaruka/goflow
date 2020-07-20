package contactql_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/contactql"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestQueryError(t *testing.T) {
	e := contactql.NewQueryErrorf("bad query")

	assert.Equal(t, "bad query", e.Error())

	e1 := errors.Wrap(errors.Wrap(e, "wrapped once"), "wrapped twice")
	isQueryError, cause := contactql.IsQueryError(e1)
	assert.True(t, isQueryError)
	assert.Equal(t, e, cause)

	e2 := errors.Wrap(errors.Wrap(fmt.Errorf("not a query error"), "wrapped once"), "wrapped twice")
	isQueryError, cause = contactql.IsQueryError(e2)
	assert.False(t, isQueryError)
	assert.Nil(t, cause)
}
