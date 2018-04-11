package main

// generate full docs with:
//
// go install github.com/nyaruka/goflow/cmd/docgen; docgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	markdownFile := "docs/docs.md"
	htmlFile := "docs/index.html"

	output, err := buildDocs(".")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// write the markdown file
	ioutil.WriteFile(markdownFile, []byte(output), 0666)

	fmt.Printf("Markdown written to '%s'\n", markdownFile)

	panDocArgs := []string{
		"--from=markdown",
		"--to=html",
		"-o", htmlFile,
		"--standalone",
		"--template=cmd/docgen/templates/template.html",
		"--toc",
		"--toc-depth=1",
		markdownFile,
	}

	cmd := exec.Command("pandoc", panDocArgs...)
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("HTML written to '%s'\n", htmlFile)
}
