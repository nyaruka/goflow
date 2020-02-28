package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/cmd/docgen/completion"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/pkg/errors"
)

func init() {
	registerGenerator("completion map", generateCompletionMap)
}

func generateCompletionMap(baseDir string, outputDir string, items map[string][]*TaggedItem) error {
	types := []completion.Type{
		// the dynamic types in the context aren't described in the code so we add them manually here
		completion.NewDynamicType("fields", "fields", completion.NewProperty("{key}", "{key} for the contact", "any")),
		completion.NewDynamicType("results", "results", completion.NewProperty("{key}", "the result for {key}", "result")),
		completion.NewDynamicType("globals", "globals", completion.NewProperty("{key}", "the global value {key}", "text")),

		// the urns type also added here as it's "dynamic" in sense that keys are known at build time
		createURNsType(),
	}

	// now collect the types from tagged docstrings
	var root []*completion.Property

	for _, item := range items["context"] {
		// examples are actually property descriptors for context items
		properties := make([]*completion.Property, len(item.examples))
		for i, propDesc := range item.examples {
			prop := completion.ParseProperty(propDesc)
			if prop == nil {
				return errors.Errorf("invalid format for property description \"%s\"", propDesc)
			}
			properties[i] = prop
		}

		if item.tagValue == "root" {
			root = properties
		} else {
			types = append(types, completion.NewStaticType(item.tagValue, properties))
		}
	}

	c := completion.NewCompletion(types, root)

	if err := c.Validate(); err != nil {
		return err
	}

	mapPath := path.Join(outputDir, "completion.json")
	marshaled, _ := jsonx.MarshalPretty(c)
	ioutil.WriteFile(mapPath, marshaled, 0755)

	fmt.Printf(" > %d completion map written to %s\n", len(items["context"]), mapPath)

	context := completion.NewContext(map[string][]string{
		"fields":  []string{"age", "gender"},
		"globals": []string{"org_name"},
		"results": []string{"response_1"},
	})
	nodes := c.EnumerateNodes(context)

	nodeOutput := &strings.Builder{}
	for _, n := range nodes {
		nodeOutput.WriteString(fmt.Sprintf("%s -> %s\n", n.Path, n.Help))
	}

	listPath := path.Join(outputDir, "completion.txt")
	ioutil.WriteFile(listPath, []byte(nodeOutput.String()), 0755)

	return nil
}

func createURNsType() completion.Type {
	properties := make([]*completion.Property, 0, len(urns.ValidSchemes))
	for k := range urns.ValidSchemes {
		name := strings.Title(k)
		properties = append(properties, completion.NewProperty(k, name+" URN for the contact", "text"))
	}
	sort.SliceStable(properties, func(i, j int) bool { return properties[i].Key < properties[j].Key })

	return completion.NewStaticType("urns", properties)
}
