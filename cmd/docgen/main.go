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

	var contextDocs, functionDocs, testDocs, actionDocs, eventDocs string

	if contextDocs, err = buildDocSet(baseDir, "flows", "@context", handleContextDoc, session); err != nil {
		return "", err
	}
	if functionDocs, err = buildDocSet(baseDir, "excellent/functions", "@function", handleFunctionDoc, session); err != nil {
		return "", err
	}
	if testDocs, err = buildDocSet(baseDir, "flows/routers/tests", "@test", handleFunctionDoc, session); err != nil {
		return "", err
	}
	if actionDocs, err = buildDocSet(baseDir, "flows/actions", "@action", handleActionDoc, session); err != nil {
		return "", err
	}
	if eventDocs, err = buildDocSet(baseDir, "flows/events", "@event", handleEventDoc, session); err != nil {
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

func buildDocSet(goflowPath string, subdir string, tag string, handler handleFunc, session flows.Session) (string, error) {
	buffer := bytes.Buffer{}
	examplePath := path.Join(goflowPath, subdir)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, examplePath, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	for _, f := range pkgs {
		p := doc.New(f, "./", 0)
		for _, t := range p.Types {
			if strings.Contains(t.Doc, tag) {
				if err := handler(&buffer, tag, t.Name, t.Doc, session); err != nil {
					return "", fmt.Errorf("error parsing %s docstrings: %s", tag, err)
				}
			}
		}
		for _, t := range p.Funcs {
			if strings.Contains(t.Doc, tag) {
				if err := handler(&buffer, tag, t.Name, t.Doc, session); err != nil {
					return "", fmt.Errorf("error parsing %s docstrings: %s", tag, err)
				}
			}
		}
	}

	output := buffer.String()
	if output == "" {
		return "", fmt.Errorf("found 0 docstrings for tag %s", tag)
	}

	return output, nil
}
