package static_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	flow := static.NewFlow("f5263dca-469b-47c2-be4f-845d3a14eedf", "Catch All", []byte(`{}`))
	assert.Equal(t, assets.FlowUUID("f5263dca-469b-47c2-be4f-845d3a14eedf"), flow.UUID())
	assert.Equal(t, "Catch All", flow.Name())
	assert.Equal(t, []byte(`{}`), flow.Definition())

	definition := []byte(`{"uuid": "f5263dca-469b-47c2-be4f-845d3a14eedf", "name": "Registration", "nodes": []}`)
	f := &static.Flow{}
	err := jsonx.Unmarshal(definition, f)

	assert.NoError(t, err)
	assert.Equal(t, assets.FlowUUID("f5263dca-469b-47c2-be4f-845d3a14eedf"), f.UUID())
	assert.Equal(t, "Registration", f.Name())
	assert.Equal(t, definition, f.Definition())

	// can also read legacy definition with metadata section
	definition = []byte(`{"metadata": {"uuid": "834ab66a-cc95-4a4f-8a45-2ff9cd2ec4ab", "name": "Legacy"}}`)
	f = &static.Flow{}
	err = jsonx.Unmarshal(definition, f)

	assert.NoError(t, err)
	assert.Equal(t, assets.FlowUUID("834ab66a-cc95-4a4f-8a45-2ff9cd2ec4ab"), f.UUID())
	assert.Equal(t, "Legacy", f.Name())
	assert.Equal(t, definition, f.Definition())

	// sometimes new flows also have a metadata section
	definition = []byte(`{"uuid": "f5263dca-469b-47c2-be4f-845d3a14eedf", "name": "Registration", "nodes": [], "metadata": {"revision": 1}}`)
	f = &static.Flow{}
	err = jsonx.Unmarshal(definition, f)

	assert.NoError(t, err)
	assert.Equal(t, assets.FlowUUID("f5263dca-469b-47c2-be4f-845d3a14eedf"), f.UUID())
	assert.Equal(t, "Registration", f.Name())
	assert.Equal(t, definition, f.Definition())
}
