package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocGeneration(t *testing.T) {
	// tests run from the same working directory as the test file, so two directories up is our goflow root
	output, err := buildDocs("../../")
	require.NoError(t, err)

	existingDocs, err := ioutil.ReadFile("../../docs/docs.md")
	require.NoError(t, err)

	//ioutil.WriteFile("../../docs/tests.md", []byte(output), 0666)

	assert.Equal(t, string(existingDocs), output, "changes have been made that require re-running docgen")
}
