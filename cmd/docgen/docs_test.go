package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocGeneration(t *testing.T) {
	// tests run from the same working directory as the test file, so two directories up is our goflow root
	_, err := buildDocs("../../")
	require.NoError(t, err)
}
