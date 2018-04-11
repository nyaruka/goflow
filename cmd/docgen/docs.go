package main

import (
	"bytes"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"path"
	"strings"
	"text/template"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type taggedType struct {
	docString string
	typeName  string
}

type handleFunc func(output *strings.Builder, prefix string, docString string, typeName string, session flows.Session) error

// builds all documentation from the given base directory
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

	if contextDocs, err = buildDocSet(baseDir, []string{"flows"}, "@context", handleContextDoc, session); err != nil {
		return "", err
	}
	if functionDocs, err = buildDocSet(baseDir, []string{"excellent/functions"}, "@function", handleFunctionDoc, session); err != nil {
		return "", err
	}
	if testDocs, err = buildDocSet(baseDir, []string{"flows/routers/tests"}, "@test", handleFunctionDoc, session); err != nil {
		return "", err
	}
	if actionDocs, err = buildDocSet(baseDir, []string{"flows/actions"}, "@action", handleActionDoc, session); err != nil {
		return "", err
	}
	if eventDocs, err = buildDocSet(baseDir, []string{"flows/events"}, "@event", handleEventDoc, session); err != nil {
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

func buildDocSet(baseDir string, searchDirs []string, tag string, handler handleFunc, session flows.Session) (string, error) {
	taggedTypes := make([]taggedType, 0)
	for _, searchDir := range searchDirs {
		fromDir, err := findTaggedTypes(baseDir, searchDir, tag)
		if err != nil {
			return "", err
		}
		taggedTypes = append(taggedTypes, fromDir...)
	}

	buffer := &strings.Builder{}

	for _, taggedType := range taggedTypes {
		if err := handler(buffer, tag, taggedType.typeName, taggedType.docString, session); err != nil {
			return "", fmt.Errorf("error parsing %s docstrings: %s", tag, err)
		}
	}

	output := buffer.String()
	if output == "" {
		return "", fmt.Errorf("found 0 docstrings for tag %s", tag)
	}

	return output, nil
}

// finds all tagged types in go files in the given directory
func findTaggedTypes(baseDir string, searchDir string, tag string) ([]taggedType, error) {
	taggedTypes := make([]taggedType, 0)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path.Join(baseDir, searchDir), nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, f := range pkgs {
		p := doc.New(f, "./", 0)
		for _, t := range p.Types {
			if strings.Contains(t.Doc, tag) {
				taggedTypes = append(taggedTypes, taggedType{docString: t.Doc, typeName: t.Name})
			}
		}
		for _, t := range p.Funcs {
			if strings.Contains(t.Doc, tag) {
				taggedTypes = append(taggedTypes, taggedType{docString: t.Doc, typeName: t.Name})
			}
		}
	}

	return taggedTypes, nil
}
