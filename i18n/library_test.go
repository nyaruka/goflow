package i18n_test

import (
	"testing"

	"github.com/nyaruka/goflow/i18n"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLibrary(t *testing.T) {
	library := i18n.NewLibrary("testdata/locales", "en")

	assert.Equal(t, "testdata/locales", library.Path())
	assert.Equal(t, "en", library.SrcLanguage())
	assert.Equal(t, []string{"en", "es"}, library.Languages())

	es, err := library.Load("es", "simple")
	require.NoError(t, err)

	assert.Equal(t, "Azul", es.GetText("", "Blue"))
}
