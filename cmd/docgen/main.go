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
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type handleFunc func(output *bytes.Buffer, prefix string, typeName string, docString string, session flows.Session) error

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

func buildDocs(baseDir string) (string, error) {
	server, err := utils.NewTestHTTPServer()
	if err != nil {
		return "", fmt.Errorf("error starting mock HTTP server: %s", err)
	}
	server.Start()
	defer server.Close()

	session, err := createExampleSession(nil)
	if err != nil {
		return "", fmt.Errorf("error creating example session: %s", err)
	}

	contextDocs, _, err := buildDocSet(baseDir, "flows", "@context", handleContextDoc, session)
	if err != nil {
		return "", err
	}
	functionDocs, _, err := buildDocSet(baseDir, "excellent/functions", "@function", handleFunctionDoc, session)
	if err != nil {
		return "", err
	}
	testDocs, _, err := buildDocSet(baseDir, "flows/routers/tests", "@test", handleFunctionDoc, session)
	if err != nil {
		return "", err
	}
	actionDocs, _, err := buildDocSet(baseDir, "flows/actions", "@action", handleActionDoc, session)
	if err != nil {
		return "", err
	}
	eventDocs, _, err := buildDocSet(baseDir, "flows/events", "@event", handleEventDoc, session)
	if err != nil {
		return "", err
	}

	context := struct {
		ContextDocs  string
		FunctionDocs string
		TestDocs     string
		ActionDocs   string
		EventDocs    string
	}{
		ContextDocs:  contextDocs,
		FunctionDocs: functionDocs,
		TestDocs:     testDocs,
		ActionDocs:   actionDocs,
		EventDocs:    eventDocs,
	}

	// generate our complete docs
	docTpl, err := template.ParseFiles(path.Join(baseDir, "cmd/docgen/templates/docs.md"))
	if err != nil {
		return "", fmt.Errorf("Error reading template file: %s", err)
	}

	output := bytes.Buffer{}
	err = docTpl.Execute(&output, context)
	if err != nil {
		return "", fmt.Errorf("Error executing template: %s", err)
	}

	return output.String(), nil
}

func buildDocSet(goflowPath string, subdir string, tag string, handler handleFunc, session flows.Session) (string, int, error) {
	output := bytes.Buffer{}
	examplePath := path.Join(goflowPath, subdir)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, examplePath, nil, parser.ParseComments)
	if err != nil {
		return "", 0, err
	}

	itemsFound := 0

	for _, f := range pkgs {
		p := doc.New(f, "./", 0)
		for _, t := range p.Types {
			if strings.Contains(t.Doc, tag) {
				itemsFound++
				if err := handler(&output, tag, t.Name, t.Doc, session); err != nil {
					return "", 0, fmt.Errorf("error parsing %s docstrings: %s", tag, err)
				}
			}
		}
		for _, t := range p.Funcs {
			if strings.Contains(t.Doc, tag) {
				itemsFound++
				if err := handler(&output, tag, t.Name, t.Doc, session); err != nil {
					return "", 0, fmt.Errorf("error parsing %s docstrings: %s", tag, err)
				}
			}
		}
	}
	return output.String(), itemsFound, nil
}
