package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"

	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	field := types.NewField(assets.FieldUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), "age", "Age", assets.FieldTypeNumber)
	assert.Equal(t, assets.FieldUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), field.UUID())
	assert.Equal(t, "age", field.Key())
	assert.Equal(t, "Age", field.Name())
	assert.Equal(t, assets.FieldTypeNumber, field.Type())
}
