package main

import "testing"

func TestDocGeneration(t *testing.T) {
	// tests run from the same working directory as the test file, so two directories up is our goflow root
	path := "../../"
	buildExcellentDocs(path)
	buildExampleDocs(path, "flows/actions", "@action", handleActionDoc)
	buildExampleDocs(path, "flows/events", "@event", handleEventDoc)
}
