package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocGeneration(t *testing.T) {
	session, err := createExampleSession(nil)
	require.NoError(t, err)

	// tests run from the same working directory as the test file, so two directories up is our goflow root
	path := "../../"
	buildDocSet(path, "excellent", "@function", handleFunctionDoc, session)
	buildDocSet(path, "flows/routers/tests", "@test", handleFunctionDoc, session)
	buildDocSet(path, "flows/actions", "@action", handleActionDoc, session)
	buildDocSet(path, "flows/events", "@event", handleEventDoc, session)
}
