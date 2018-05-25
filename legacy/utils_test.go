package legacy_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslations(t *testing.T) {
	// can unmarshall from a single string
	translations := make(legacy.Translations)
	json.Unmarshal([]byte(`"hello"`), &translations)
	assert.Equal(t, legacy.Translations{"base": "hello"}, translations)

	// or a map
	translations = make(legacy.Translations)
	json.Unmarshal([]byte(`{"eng": "hello", "fra": "bonjour"}`), &translations)
	assert.Equal(t, legacy.Translations{"eng": "hello", "fra": "bonjour"}, translations)

	// and back to JSON
	data, err := json.Marshal(translations)
	require.NoError(t, err)
	assert.Equal(t, []byte(`{"eng":"hello","fra":"bonjour"}`), data)

	translationSet := []legacy.Translations{
		{"eng": "Yes", "fra": "Oui"},
		{"eng": "No", "fra": "Non"},
		{"eng": "Maybe"},
		{"eng": "Never", "fra": "Jamas"},
	}
	assert.Equal(t, map[utils.Language][]string{
		"eng": {"Yes", "No", "Maybe", "Never"},
		"fra": {"Oui", "Non", "", "Jamas"},
	}, legacy.TransformTranslations(translationSet))
}

func TestDecimalString(t *testing.T) {
	// can unmarshall from a string
	var decimal legacy.DecimalString
	json.Unmarshal([]byte(`"123.45"`), &decimal)
	assert.Equal(t, legacy.DecimalString("123.45"), decimal)

	// or a floating point (JSON number type)
	json.Unmarshal([]byte(`567.89`), &decimal)
	assert.Equal(t, legacy.DecimalString("567.89"), decimal)
}
