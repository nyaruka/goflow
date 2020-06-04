package main

// generate full docs with:
//
// go install github.com/nyaruka/goflow/cmd/docgen; docgen

import (
	"fmt"
	"os"

	"github.com/nyaruka/goflow/cmd/docgen/docs"
)

const (
	outputDir  = "docs"
	localesDir = "locales"
)

func main() {
	if err := docs.Generate(".", outputDir, localesDir); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
