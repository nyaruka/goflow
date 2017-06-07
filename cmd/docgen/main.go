package main

// generate full docs with:
// go install github.com/nyaruka/goflow/cmd/docgen && $GOPATH/bin/docgen . | pandoc --from=markdown --to=html -o docs.html --standalone --template=cmd/docgen/templates/template.html

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

func buildExcellentDocs(goflowPath string) (string, string) {
	excellentPath := path.Join(goflowPath, "excellent")

	fset := token.NewFileSet()
	pkgs, e := parser.ParseDir(fset, excellentPath, nil, parser.ParseComments)
	if e != nil {
		log.Fatal(e)
	}

	astf := make([]*ast.File, 0)
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			astf = append(astf, f)
		}
	}

	functionOutput := bytes.Buffer{}
	for _, f := range astf {
		ast.Walk(newFuncVisitor("function", &functionOutput), f)
	}

	testOutput := bytes.Buffer{}
	for _, f := range astf {
		ast.Walk(newFuncVisitor("test", &testOutput), f)
	}
	return functionOutput.String(), testOutput.String()
}

func buildExampleDocs(goflowPath string, subdir string, tag string, handler handleExampleFunc) string {
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
				handler(&output, tag, t.Name, t.Doc)
			}
		}
	}
	return output.String()
}

type handleExampleFunc func(output *bytes.Buffer, prefix string, typeName string, docString string)

type docContext struct {
	ExcellentFunctionDocs string
	ExcellentTestDocs     string
	ActionDocs            string
	EventDocs             string
}

func main() {
	path := os.Args[1]

	context := docContext{}

	context.ExcellentFunctionDocs, context.ExcellentTestDocs = buildExcellentDocs(path)
	context.ActionDocs = buildExampleDocs(path, "flows/actions", "@action", handleActionDoc)
	context.EventDocs = buildExampleDocs(path, "flows/events", "@event", handleEventDoc)

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
