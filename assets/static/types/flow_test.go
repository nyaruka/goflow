package types_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	definition := json.RawMessage(`{"uuid": "f5263dca-469b-47c2-be4f-845d3a14eedf", "name": "Registration", "nodes": []}`)
	f := &types.Flow{}
	err := json.Unmarshal(definition, f)

	assert.NoError(t, err)
	assert.Equal(t, assets.FlowUUID("f5263dca-469b-47c2-be4f-845d3a14eedf"), f.UUID())
	assert.Equal(t, "Registration", f.Name())
	assert.Equal(t, definition, f.Definition())
}
