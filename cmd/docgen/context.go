package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/nyaruka/goflow/cmd/docgen/context"
	"github.com/nyaruka/goflow/utils"
	"github.com/pkg/errors"
)

func init() {
	registerGenerator("context map", generateContextMap)
}

func generateContextMap(baseDir string, outputDir string, items map[string][]*TaggedItem) error {
	ctx := context.NewContext()

	// the dynamic types in the context aren't described in the code so we add them manually here
	ctx.AddType(context.NewDynamicType("fields", "field-keys", context.NewProperty("{key}", "{key} for the contact", "any")))
	ctx.AddType(context.NewDynamicType("results", "result-keys", context.NewProperty("{key}", "{key} value for the run", "result")))
	ctx.AddType(context.NewDynamicType("urns", "urn-schemes", context.NewProperty("{key}", "the {key} URN for the contact", "text")))

	// now add the types from tagged docstrings
	for _, item := range items["context"] {
		// examples are actually property descriptors for context items
		properties := make([]*context.Property, len(item.examples))
		for i, propDesc := range item.examples {
			prop := context.ParseProperty(propDesc)
			if prop == nil {
				return errors.Errorf("invalid format for property description \"%s\"", propDesc)
			}
			properties[i] = prop
		}

		if item.tagValue == "root" {
			ctx.SetRoot(properties)
		} else {
			ctx.AddType(context.NewStaticType(item.tagValue, properties))
		}
	}

	if err := ctx.Validate(); err != nil {
		return err
	}

	path := path.Join(outputDir, "context.json")
	marshaled, _ := utils.JSONMarshalPretty(ctx)
	ioutil.WriteFile(path, marshaled, 0755)

	fmt.Printf(" > %d context types written to %s\n", len(items["context"]), path)

	return nil
}
