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

	// if the docs we just generated don't match the existing ones, someone needs to run docgen
	assert.Equal(t, string(existingDocs), output, "changes have been made that require re-running docgen (go install github.com/nyaruka/goflow/cmd/docgen; docgen)")
}
