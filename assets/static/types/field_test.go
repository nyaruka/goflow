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

	fields, err := types.ReadFields([]byte(`[{"key": "gender", "name": "Gender", "value_type": "text"}]`))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(fields))
	assert.Equal(t, "gender", fields[0].Key())
	assert.Equal(t, "Gender", fields[0].Name())
	assert.Equal(t, assets.FieldTypeText, fields[0].Type())
}
