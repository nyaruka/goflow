package i18n_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/nyaruka/goflow/utils/i18n"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLibrary(t *testing.T) {
	// create a temporary directory to hold a library
	libraryDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	defer os.RemoveAll(libraryDir)

	// and setup the locale directory for en_US and es
	os.Mkdir(path.Join(libraryDir, "en"), 0700)
	os.Mkdir(path.Join(libraryDir, "es"), 0700)

	// copy sample PO files from testdata
	poEN, err := ioutil.ReadFile(path.Join("testdata", "locale", "en", "simple.po"))
	require.NoError(t, err)
	poES, err := ioutil.ReadFile(path.Join("testdata", "locale", "es", "simple.po"))
	require.NoError(t, err)
	ioutil.WriteFile(path.Join(libraryDir, "en", "simple.po"), poEN, 0700)
	ioutil.WriteFile(path.Join(libraryDir, "es", "simple.po"), poES, 0700)

	library := i18n.NewLibrary(libraryDir, "en")

	assert.Equal(t, libraryDir, library.Path())
	assert.Equal(t, "en", library.SrcLanguage())
	assert.Equal(t, []string{"en", "es"}, library.Locales())

	es, err := library.Load("es", "simple")
	require.NoError(t, err)
	assert.Equal(t, 5, len(es.Entries))
	assert.Equal(t, "Azul", es.GetText("", "Blue"))

	en, err := library.Load("en", "simple")
	require.NoError(t, err)
	assert.Equal(t, 5, len(es.Entries))
	assert.Equal(t, "Blue", en.GetText("", "Blue"))

	// add new entry
	en.AddEntry(&i18n.POEntry{MsgID: "Green"})

	err = library.Update("simple", en)
	require.NoError(t, err)

	// the new entry will have been saved into the English PO
	en, err = library.Load("en", "simple")
	require.NoError(t, err)
	assert.Equal(t, 5, len(en.Entries))

	// the new entry will have been merged into the Spanish PO
	es, err = library.Load("es", "simple")
	require.NoError(t, err)
	assert.Equal(t, 5, len(es.Entries))
}
