package main

import (
	"fmt"
	"go/ast"
	"log"
	"strings"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/utils"
)

func newFuncVisitor(funcType string) ast.Visitor {
	return &funcVisitor{
		prefix:   "@" + funcType,
		env:      utils.NewDefaultEnvironment(),
		resolver: utils.NewMapResolver(make(map[string]interface{})),
	}
}

type funcVisitor struct {
	prefix   string
	env      utils.Environment
	resolver utils.VariableResolver
}

func (v *funcVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch exp := node.(type) {
		case *ast.FuncDecl:
			if exp.Doc != nil && strings.Contains(exp.Doc.Text(), v.prefix) {
				lines := strings.Split(exp.Doc.Text(), "\n")
				name := ""

				docs := make([]string, 0, len(lines))
				examples := make([]string, 0, len(lines))
				literalExamples := make([]string, 0, len(lines))
				for _, l := range lines {
					if strings.HasPrefix(l, v.prefix) {
						name = l[len(v.prefix)+1:]
					} else if strings.HasPrefix(l, "  ") {
						examples = append(examples, l[2:])
					} else if strings.HasPrefix(l, " ") {
						literalExamples = append(literalExamples, l[1:])
					} else {
						docs = append(docs, l)
					}
				}

				if name != "" {
					if len(docs) > 0 && strings.HasPrefix(docs[0], exp.Name.String()) {
						docs[0] = strings.Replace(docs[0], exp.Name.String(), name, 1)
					}

					// check our examples
					for _, l := range examples {
						pieces := strings.Split(l, "->")
						if len(pieces) != 2 {
							log.Fatalf("Invalid example: %s", l)
						}
						test, expected := strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1])

						if expected[0] == '"' && expected[len(expected)-1] == '"' {
							expected = expected[1 : len(expected)-1]
						}

						// evaluate our expression
						val, err := excellent.EvaluateTemplateAsString(v.env, v.resolver, test)
						if err != nil && expected != "ERROR" {
							log.Fatalf("Invalid example: %s  Error: %s", l, err)
						}
						if val != expected && expected != "ERROR" {
							log.Fatalf("Invalid example: %s  Got: '%s' Expected: '%s'", l, val, expected)
						}
					}

					fmt.Printf("# %s\n\n", name)
					fmt.Printf("%s", strings.Join(docs, "\n"))
					fmt.Printf("```objective-c\n")
					if len(examples) > 0 {
						fmt.Printf("%s\n", strings.Join(examples, "\n"))
					}
					if len(literalExamples) > 0 {
						fmt.Printf("%s\n", strings.Join(literalExamples, "\n"))
					}
					fmt.Printf("```\n")
					fmt.Printf("\n")
				}
			}
		}
	}
	return v
}
