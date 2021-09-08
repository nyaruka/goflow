package static

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	channel := assets.ChannelReference{
		Name: "Test Channel",
		UUID: assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"),
	}

	translation := NewTemplateTranslation(channel, envs.Language("eng"), envs.Country("US"), "Hello {{1}}", 1, "0162a7f4_dfe4_4c96_be07_854d5dba3b2b")
	assert.Equal(t, channel, translation.Channel())
	assert.Equal(t, envs.Language("eng"), translation.Language())
	assert.Equal(t, envs.Country("US"), translation.Country())
	assert.Equal(t, "Hello {{1}}", translation.Content())
	assert.Equal(t, 1, translation.VariableCount())
	assert.Equal(t, "0162a7f4_dfe4_4c96_be07_854d5dba3b2b", translation.Namespace())

	template := NewTemplate(assets.TemplateUUID("8a9c1f73-5059-46a0-ba4a-6390979c01d3"), "hello", []*TemplateTranslation{translation})
	assert.Equal(t, assets.TemplateUUID("8a9c1f73-5059-46a0-ba4a-6390979c01d3"), template.UUID())
	assert.Equal(t, "hello", template.Name())
	assert.Equal(t, 1, len(template.Translations()))

	// test json and back
	asJSON, err := jsonx.Marshal(template)
	assert.NoError(t, err)

	copy := Template{}
	err = jsonx.Unmarshal(asJSON, &copy)
	assert.NoError(t, err)

	assert.Equal(t, copy.Name(), template.Name())
	assert.Equal(t, copy.UUID(), template.UUID())
	assert.Equal(t, copy.Translations()[0].Content(), template.Translations()[0].Content())
	assert.Equal(t, copy.Translations()[0].Namespace(), template.Translations()[0].Namespace())
}
