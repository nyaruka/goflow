package main

// generate full docs with:
//
// go install github.com/nyaruka/goflow/cmd/docgen; docgen

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	outputDir string = "docs"
)

type generatorFunc func(baseDir string, outputDir string, items map[string][]*TaggedItem) error
type generator struct {
	name     string
	function generatorFunc
}

var generators []generator

func registerGenerator(name string, fn generatorFunc) {
	generators = append(generators, generator{name, fn})
}

func main() {
	if err := GenerateDocs(".", outputDir); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// GenerateDocs generates all documentation outputs
func GenerateDocs(baseDir string, outputDir string) error {
	fmt.Println("Processing sources...")

	// extract all documented items from the source code
	taggedItems, err := FindAllTaggedItems(baseDir)
	if err != nil {
		return errors.Wrap(err, "error extracting tagged items")
	}

	for k, v := range taggedItems {
		fmt.Printf(" > Found %d tagged items with tag %s\n", len(v), k)
	}

	// invoke doc generators...

	for _, g := range generators {
		fmt.Printf("Invoking generator: %s...\n", g.name)

		if err := g.function(baseDir, outputDir, taggedItems); err != nil {
			return errors.Wrapf(err, "error invoking generator %s", g.name)
		}
	}

	return nil
}
