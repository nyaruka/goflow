package types_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestXError(t *testing.T) {
	err1 := types.NewXError(fmt.Errorf("I failed"))
	assert.Equal(t, types.NewXText("I failed"), err1.ToXText())
	assert.Equal(t, types.NewXText(`"I failed"`), err1.ToXJSON())
	assert.Equal(t, types.XBooleanFalse, err1.ToXBoolean())
	assert.Equal(t, "I failed", err1.String())
	assert.Equal(t, "I failed", err1.Error())

	err2 := types.NewXResolveError(nil, "foo")
	assert.Equal(t, "unable to resolve 'foo'", err2.Error())
}
