package static_test

import (
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	channel := assets.NewChannelReference("Test Channel", "ffffffff-9b24-92e1-ffff-ffffb207cdb4")

	v1 := static.NewTemplateVariable("text")
	v2 := static.NewTemplateVariable("text")
	assert.Equal(t, "text", v1.Type())
	assert.Equal(t, "text", v2.Type())

	c1 := static.NewTemplateComponent("body", "body", "Hello {{1}}", "", map[string]int{"1": 0})
	c2 := static.NewTemplateComponent("button/url", "button.0", "http://google.com?q={{1}}", "Go", map[string]int{"1": 1})

	assert.Equal(t, "body", c1.Type())
	assert.Equal(t, "body", c1.Name())
	assert.Equal(t, "Hello {{1}}", c1.Content())
	assert.Equal(t, "", c1.Display())
	assert.Equal(t, map[string]int{"1": 0}, c1.Variables())

	translation := static.NewTemplateTranslation(channel, i18n.Locale("eng-US"), "0162a7f4_dfe4_4c96_be07_854d5dba3b2b", []*static.TemplateComponent{c1, c2}, []*static.TemplateVariable{v1, v2})
	assert.Equal(t, channel, translation.Channel())
	assert.Equal(t, i18n.Locale("eng-US"), translation.Locale())
	assert.Equal(t, "0162a7f4_dfe4_4c96_be07_854d5dba3b2b", translation.Namespace())
	assert.Equal(t, []assets.TemplateComponent{
		(assets.TemplateComponent)(c1),
		(assets.TemplateComponent)(c2),
	}, translation.Components())

	template := static.NewTemplate(assets.TemplateUUID("8a9c1f73-5059-46a0-ba4a-6390979c01d3"), "hello", []*static.TemplateTranslation{translation})
	assert.Equal(t, assets.TemplateUUID("8a9c1f73-5059-46a0-ba4a-6390979c01d3"), template.UUID())
	assert.Equal(t, "hello", template.Name())
	assert.Equal(t, 1, len(template.Translations()))

	// test json and back
	asJSON, err := jsonx.Marshal(template)
	assert.NoError(t, err)

	copy := &static.Template{}
	err = jsonx.Unmarshal(asJSON, copy)
	assert.NoError(t, err)

	assert.Equal(t, copy.Name(), template.Name())
	assert.Equal(t, copy.UUID(), template.UUID())
	assert.Equal(t, copy.Translations()[0].Namespace(), template.Translations()[0].Namespace())
}
