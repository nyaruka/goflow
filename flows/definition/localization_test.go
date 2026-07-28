package definition_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
}

func TestLocalizationWithTooManyItems(t *testing.T) {
	// a language can't have translations for more items than a max sized flow could possibly contain
	b := &strings.Builder{}
	b.WriteString(`{"spa": {`)
	for i := 0; i <= flows.MaxItemsPerLanguage; i++ {
		if i > 0 {
			b.WriteString(",")
		}
		fmt.Fprintf(b, `"e0323e14-9dcf-4a8d-b721-%012d": {}`, i)
	}
	b.WriteString(`}}`)

	l8n, err := definition.ReadLocalization([]byte(b.String()))
	require.NoError(t, err)

	assert.EqualError(t, l8n.Validate(), fmt.Sprintf("invalid translation for 'spa': can't have more than %d item translations (has %d)", flows.MaxItemsPerLanguage, flows.MaxItemsPerLanguage+1))
}
