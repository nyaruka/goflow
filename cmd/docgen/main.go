package main

// generate full docs with:
//
// go install github.com/nyaruka/goflow/cmd/docgen
// $GOPATH/bin/docgen . | pandoc --from=markdown --to=html -o docs/index.html --standalone --template=cmd/docgen/templates/template.html --toc --toc-depth=1

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: docgen <basedir>")
		os.Exit(1)
	}

	output, err := buildDocs(os.Args[1])

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// write output to stdout so it can be piped elsewhere
	fmt.Println(output)
}
