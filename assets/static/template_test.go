package static

import (
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	channel := assets.NewChannelReference("Test Channel", "ffffffff-9b24-92e1-ffff-ffffb207cdb4")

	tp1 := NewTemplateParam("text", "1")
	assert.Equal(t, "text", tp1.Type())
	assert.Equal(t, "1", tp1.Name())

	tc1 := NewTemplateComponent("body", "body", "Hello {{1}}", "", []*TemplateParam{tp1})
	tc2 := NewTemplateComponent("button/url", "button.0", "http://google.com?q={{1}}", "Go {{1}", []*TemplateParam{NewTemplateParam("text", "1"), NewTemplateParam("text", "2")})

	assert.Equal(t, "body", tc1.Type())
	assert.Equal(t, "body", tc1.Name())
	assert.Equal(t, "Hello {{1}}", tc1.Content())
	assert.Equal(t, "", tc1.Display())
	assert.Equal(t, []assets.TemplateParam{tp1}, tc1.Params())

	translation := NewTemplateTranslation(channel, i18n.Locale("eng-US"), "0162a7f4_dfe4_4c96_be07_854d5dba3b2b", []*TemplateComponent{tc1, tc2})
	assert.Equal(t, channel, translation.Channel())
	assert.Equal(t, i18n.Locale("eng-US"), translation.Locale())
	assert.Equal(t, "0162a7f4_dfe4_4c96_be07_854d5dba3b2b", translation.Namespace())
	assert.Equal(t, []assets.TemplateComponent{
		(assets.TemplateComponent)(tc1),
		(assets.TemplateComponent)(tc2),
	}, translation.Components())

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
	assert.Equal(t, copy.Translations()[0].Namespace(), template.Translations()[0].Namespace())
}
