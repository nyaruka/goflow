package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"log"
	"strings"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/utils"
)

// the set of assets loaded into the session that function examples are evaluated against
var functionsTestAssets = `
[
	{
		"type": "flow",
		"url": "http://testserver/assets/flow/50c3706e-fedb-42c0-8eab-dda3335714b7",
		"content": {
			"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
			"name": "EmptyFlow",
			"nodes": []
		}
	},
	{
		"type": "field",
		"url": "http://testserver/assets/field",
		"content": [
			{"key": "gender", "label": "Gender", "value_type": "text"}
		],
		"is_set": true
	},
	{
		"type": "group",
		"url": "http://testserver/assets/group",
		"content": [],
		"is_set": true
	},
	{
		"type": "group",
		"url": "http://testserver/assets/group",
		"content": [],
		"is_set": true
	},
	{
		"type": "location_hierarchy",
		"url": "http://testserver/assets/location_hierarchy",
		"content": {
			"id": "2342",
			"name": "Rwanda",
			"aliases": ["Ruanda"],		
			"children": [
				{
					"id": "234521",
					"name": "Kigali City",
					"aliases": ["Kigali", "Kigari"],
					"children": [
						{
							"id": "57735322",
							"name": "Gasabo",
							"children": [
								{
									"id": "575743222",
									"name": "Gisozi"
								},
								{
									"id": "457378732",
									"name": "Ndera"
								}
							]
						},
						{
							"id": "46547322",
							"name": "Nyarugenge",
							"children": []
						}
					]
				}
			]
		}
	}
]
`

func newFuncVisitor(funcType string, output *bytes.Buffer) ast.Visitor {
	session, err := createExampleSession(functionsTestAssets)
	if err != nil {
		log.Fatalf("Error creating example session: %s", err)
	}

	run := session.Runs()[0]

	return &funcVisitor{
		prefix:   "@" + funcType,
		env:      run.Environment(),
		resolver: run.Context(),
		output:   output,
	}
}

type funcVisitor struct {
	prefix   string
	env      utils.Environment
	resolver utils.VariableResolver
	output   *bytes.Buffer
}

func (v *funcVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		switch exp := node.(type) {
		case *ast.FuncDecl:
			if exp.Doc != nil && strings.Contains(exp.Doc.Text(), v.prefix) {
				lines := strings.Split(exp.Doc.Text(), "\n")
				signature := ""

				docs := make([]string, 0, len(lines))
				examples := make([]string, 0, len(lines))
				literalExamples := make([]string, 0, len(lines))
				for _, l := range lines {
					if strings.HasPrefix(l, v.prefix) {
						signature = l[len(v.prefix)+1:]
					} else if strings.HasPrefix(l, "  ") {
						examples = append(examples, l[2:])
					} else if strings.HasPrefix(l, " ") {
						literalExamples = append(literalExamples, l[1:])
					} else {
						docs = append(docs, l)
					}
				}

				if signature != "" {
					name := signature[0:strings.Index(signature, "(")]
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
						val, err := excellent.EvaluateTemplateAsString(v.env, v.resolver, test, false)
						if err != nil && expected != "ERROR" {
							log.Fatalf("Invalid example: %s  Error: %s", l, err)
						}
						if val != expected && expected != "ERROR" {
							log.Fatalf("Invalid example: %s  Got: '%s' Expected: '%s'", l, val, expected)
						}
					}

					v.output.WriteString(fmt.Sprintf("## %s\n\n", signature))
					v.output.WriteString(fmt.Sprintf("%s", strings.Join(docs, "\n")))
					v.output.WriteString(fmt.Sprintf("```objectivec\n"))
					if len(examples) > 0 {
						v.output.WriteString(fmt.Sprintf("%s\n", strings.Join(examples, "\n")))
					}
					if len(literalExamples) > 0 {
						v.output.WriteString(fmt.Sprintf("%s\n", strings.Join(literalExamples, "\n")))
					}
					v.output.WriteString(fmt.Sprintf("```\n"))
					v.output.WriteString(fmt.Sprintf("\n"))
				}
			}
		}
	}
	return v
}
