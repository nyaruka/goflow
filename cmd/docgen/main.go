package main

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
)

func buildExcellentDocs(goflowPath string) {
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

	for _, f := range astf {
		ast.Walk(newFuncVisitor("function"), f)
	}
	for _, f := range astf {
		ast.Walk(newFuncVisitor("test"), f)
	}
}

func buildExampleDocs(goflowPath string, subdir string, tag string, handler handleExampleFunc) {
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
				handler(tag, t.Name, t.Doc)
			}
		}
	}
}

type handleExampleFunc func(prefix string, typeName string, docString string)

func main() {
	path := os.Args[1]
	buildExcellentDocs(path)
	buildExampleDocs(path, "flows/actions", "@action", handleActionDoc)
	buildExampleDocs(path, "flows/events", "@event", handleEventDoc)
}
