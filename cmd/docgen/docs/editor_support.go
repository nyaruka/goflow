package docs

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/cmd/docgen/completion"

	"github.com/pkg/errors"
)

func init() {
	RegisterGenerator(&editorSupportGenerator{})
}

type functionExample struct {
	Template string `json:"template"`
	Output   string `json:"output"`
}

type functionListing struct {
	Signature string             `json:"signature"`
	Summary   string             `json:"summary"`
	Detail    string             `json:"detail"`
	Examples  []*functionExample `json:"examples"`
}

type editorSupport struct {
	Context   *completion.Completion `json:"context"`
	Functions []*functionListing     `json:"functions"`
}

type editorSupportGenerator struct{}

func (g *editorSupportGenerator) Name() string {
	return "editor support files"
}

func (g *editorSupportGenerator) Generate(baseDir, outputDir string, items map[string][]*TaggedItem, gettext func(string) string) error {
	es := &editorSupport{}
	var err error

	es.Context, err = g.buildContextCompletion(items, gettext)
	if err != nil {
		return err
	}

	es.Functions = g.buildFunctionListing(items, gettext)

	outputPath := path.Join(outputDir, "editor.json")
	marshaled, err := jsonx.MarshalPretty(es)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(outputPath, marshaled, 0755); err != nil {
		return err
	}
	fmt.Printf(" > editor support file written to %s\n", outputPath)

	// also output list of context paths.. not used by the editor but useful for checking
	if err := createContextPathListFile(outputDir, es.Context); err != nil {
		return err
	}

	return nil
}

func (g *editorSupportGenerator) buildContextCompletion(items map[string][]*TaggedItem, gettext func(string) string) (*completion.Completion, error) {
	types := []completion.Type{
		// the dynamic types in the context aren't described in the code so we add them manually here
		completion.NewDynamicType("fields", "fields", completion.NewProperty("{key}", gettext("{key} for the contact"), "any")),
		completion.NewDynamicType("results", "results", completion.NewProperty("{key}", gettext("the result for {key}"), "result")),
		completion.NewDynamicType("globals", "globals", completion.NewProperty("{key}", gettext("the global value {key}"), "text")),

		// the urns type also added here as it's "dynamic" in sense that keys are known at build time
		createURNsType(gettext),
	}

	// now collect the types from tagged docstrings
	var root []*completion.Property

	for _, item := range items["context"] {
		// examples are actually property descriptors for context items
		properties := make([]*completion.Property, len(item.examples))
		for i, propDesc := range item.examples {
			prop := completion.ParseProperty(propDesc)
			if prop == nil {
				return nil, errors.Errorf("invalid format for property description \"%s\"", propDesc)
			}
			prop.Help = gettext(prop.Help)
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
		return nil, err
	}

	return c, nil
}

func (g *editorSupportGenerator) buildFunctionListing(items map[string][]*TaggedItem, gettext func(string) string) []*functionListing {
	funcItems := items["function"]
	listings := make([]*functionListing, len(funcItems))

	for i, funcItem := range funcItems {
		summary := funcItem.description[0]
		detail := strings.TrimSpace(strings.Join(funcItem.description[1:len(funcItem.description)-1], "\n"))

		examples := make([]*functionExample, len(funcItem.examples))
		for j := range funcItem.examples {
			parts := strings.Split(funcItem.examples[j], "→")
			examples[j] = &functionExample{Template: strings.TrimSpace(parts[0]), Output: strings.TrimSpace(parts[1])}
		}

		listings[i] = &functionListing{
			Signature: funcItem.tagValue + funcItem.tagExtra,
			Summary:   gettext(summary),
			Detail:    gettext(detail),
			Examples:  examples,
		}
	}

	return listings
}

// creates a text file which lists all the context paths using example fields
func createContextPathListFile(outputDir string, c *completion.Completion) error {
	context := completion.NewContext(map[string][]string{
		"fields":  {"age", "gender"},
		"globals": {"org_name"},
		"results": {"response_1"},
	})
	nodes := c.EnumerateNodes(context)

	nodeOutput := &strings.Builder{}
	for _, n := range nodes {
		nodeOutput.WriteString(fmt.Sprintf("%s -> %s\n", n.Path, n.Help))
	}

	listPath := path.Join(outputDir, "completion.txt")
	return ioutil.WriteFile(listPath, []byte(nodeOutput.String()), 0755)
}

func createURNsType(gettext func(string) string) completion.Type {
	properties := make([]*completion.Property, 0, len(urns.ValidSchemes))
	for k := range urns.ValidSchemes {
		name := strings.Title(k)
		help := strings.ReplaceAll(gettext("{type} URN for the contact"), "{type}", name)
		properties = append(properties, completion.NewProperty(k, help, "text"))
	}
	sort.SliceStable(properties, func(i, j int) bool { return properties[i].Key < properties[j].Key })

	return completion.NewStaticType("urns", properties)
}
