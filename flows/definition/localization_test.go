package definition_test

import (
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/stretchr/testify/assert"
)

func TestLocalization(t *testing.T) {
	l8n, err := definition.ReadLocalization([]byte(`{
		"spa": {
			"ac110f56-a66c-4462-921c-b2c6d1c6dadb": {
				"text": [
					"Hola @contact.name"
				],
				"quick_replies": [
					"Yes", "No"
				],
				"empty": [],
				"bad1": [""],
				"bad2": [{}],
				"_ui": {
					"auto_translated": [
						"text"
					]
				}
			}
		},
		"fra": {
			"ac110f56-a66c-4462-921c-b2c6d1c6dadb": {
				"text": [
					"Bonjour @contact.name"
				]
			}
		}
	}`))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []i18n.Language{"fra", "spa"}, l8n.Languages())
	assert.Equal(t, []string{"Hola @contact.name"}, l8n.GetItemTranslation("spa", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "text"))
	assert.Equal(t, []string{"Bonjour @contact.name"}, l8n.GetItemTranslation("fra", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "text"))
	assert.Equal(t, []string{"Yes", "No"}, l8n.GetItemTranslation("spa", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "quick_replies"))
	assert.Nil(t, l8n.GetItemTranslation("spa", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "empty"))
	assert.Nil(t, l8n.GetItemTranslation("spa", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "bad1"))
	assert.Nil(t, l8n.GetItemTranslation("spa", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "bad2"))
	assert.Nil(t, l8n.GetItemTranslation("spa", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "xxx"))

	l8n.SetItemTranslation("spa", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "text", []string{"Hola @contact"})
	l8n.SetItemTranslation("kin", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "text", []string{"Bite @contact"})
	assert.Equal(t, []string{"Hola @contact"}, l8n.GetItemTranslation("spa", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "text"))
	assert.Equal(t, []string{"Bite @contact"}, l8n.GetItemTranslation("kin", "ac110f56-a66c-4462-921c-b2c6d1c6dadb", "text"))
}
