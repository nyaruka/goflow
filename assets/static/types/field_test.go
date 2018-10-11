package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"

	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	field := types.NewField("age", "Age", assets.FieldTypeNumber)
	assert.Equal(t, "age", field.Key())
	assert.Equal(t, "Age", field.Name())
	assert.Equal(t, assets.FieldTypeNumber, field.Type())
}
