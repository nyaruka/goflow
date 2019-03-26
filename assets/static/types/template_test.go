package types

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	channel := assets.ChannelReference{
		Name: "Test Channel",
		UUID: assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"),
	}

	translation := NewTemplateTranslation(channel, utils.Language("eng"), "Hello {{1}}", 1)
	assert.Equal(t, channel, translation.Channel())
	assert.Equal(t, utils.Language("eng"), translation.Language())
	assert.Equal(t, "Hello {{1}}", translation.Content())
	assert.Equal(t, 1, translation.VariableCount())

	template := NewTemplate(assets.TemplateUUID("8a9c1f73-5059-46a0-ba4a-6390979c01d3"), "hello", []*TemplateTranslation{translation})
	assert.Equal(t, assets.TemplateUUID("8a9c1f73-5059-46a0-ba4a-6390979c01d3"), template.UUID())
	assert.Equal(t, "hello", template.Name())
	assert.Equal(t, 1, len(template.Translations()))

	// test json and back
	asJSON, err := json.Marshal(template)
	assert.NoError(t, err)

	copy := Template{}
	err = json.Unmarshal(asJSON, &copy)

	assert.Equal(t, copy.Name(), template.Name())
	assert.Equal(t, copy.UUID(), template.UUID())
	assert.Equal(t, *copy.Translations()[0], *template.Translations()[0])
}
