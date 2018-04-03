package main

// generate full docs with:
//
// go install github.com/nyaruka/goflow/cmd/docgen
// $GOPATH/bin/docgen . | pandoc --from=markdown --to=html -o docs/index.html --standalone --template=cmd/docgen/templates/template.html --toc --toc-depth=1

import (
	"bytes"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/nyaruka/goflow/flows"
)

func buildDocSet(goflowPath string, subdir string, tag string, handler handleFunc, session flows.Session) string {
	output := bytes.Buffer{}
	examplePath := path.Join(goflowPath, subdir)

	fset := token.NewFileSet()
	pkgs, e := parser.ParseDir(fset, examplePath, nil, parser.ParseComments)
	if e != nil {
		log.Fatal(e)
	}

	for _, f := range pkgs {
		p := doc.New(f, "./", 0)
		for _, t := range p.Types {
			if strings.Contains(t.Doc, tag) {
				handler(&output, tag, t.Name, t.Doc, session)
			}
		}
		for _, t := range p.Funcs {
			if strings.Contains(t.Doc, tag) {
				handler(&output, tag, t.Name, t.Doc, session)
			}
		}
	}
	return output.String()
}

type handleFunc func(output *bytes.Buffer, prefix string, typeName string, docString string, session flows.Session)

func main() {
	path := os.Args[1]

	session, err := createExampleSession(nil)
	if err != nil {
		log.Fatalf("Error creating example session: %s", err)
	}

	context := struct {
		FunctionDocs string
		TestDocs     string
		ActionDocs   string
		EventDocs    string
	}{
		FunctionDocs: buildDocSet(path, "excellent", "@function", handleFunctionDoc, session),
		TestDocs:     buildDocSet(path, "flows/tests", "@test", handleFunctionDoc, session),
		ActionDocs:   buildDocSet(path, "flows/actions", "@action", handleActionDoc, session),
		EventDocs:    buildDocSet(path, "flows/events", "@event", handleEventDoc, session),
	}

	// generate our complete docs
	docTpl, err := template.ParseFiles("cmd/docgen/templates/docs.md")
	if err != nil {
		log.Fatalf("Error reading template file: %s", err)
	}

	output := bytes.Buffer{}
	err = docTpl.Execute(&output, context)
	if err != nil {
		log.Fatalf("Error executing template: %s", err)
	}

	fmt.Println(output.String())
}
