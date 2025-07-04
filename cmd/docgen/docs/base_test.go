package docs_test

import (
	"os"
	"path"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/cmd/docgen/docs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDocs(t *testing.T) {
	// create a temporary directory to hold generated doc files
	outputDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	// create a temporary directory to hold generated locale files
	localeDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	defer os.RemoveAll(outputDir)
	defer os.RemoveAll(localeDir)

	// and setup the locale directory for en_US and es
	os.Mkdir(path.Join(localeDir, "en_US"), 0700)
	os.Mkdir(path.Join(localeDir, "es"), 0700)

	os.WriteFile(path.Join(localeDir, "en_US", "flows.po"), []byte(``), 0700)
	os.WriteFile(path.Join(localeDir, "es", "flows.po"), []byte(``), 0700)

	// tests run from the same working directory as the test file, so two directories up is our goflow root
	err = docs.Generate("../../../", outputDir, localeDir)
	require.NoError(t, err)

	// check other outputs
	completion := readJSONOutput(t, outputDir, "en-us", "editor.json").(map[string]any)
	assert.Contains(t, completion, "functions")
	assert.Contains(t, completion, "context")

	context := completion["context"].(map[string]any)
	functions := completion["functions"].([]any)

	assert.Equal(t, 90, len(functions))

	types := context["types"].([]any)
	assert.Equal(t, 22, len(types))

	root := context["root"].([]any)
	assert.Equal(t, 15, len(root))
}

func readJSONOutput(t *testing.T, file ...string) any {
	output, err := os.ReadFile(path.Join(file...))
	require.NoError(t, err)

	generic, err := jsonx.DecodeGeneric(output)
	require.NoError(t, err)

	return generic
}
