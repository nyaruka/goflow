package completion

import (
	"fmt"

	"github.com/pkg/errors"
)

// types available in root of context even without a session
var rootNoSessionTypes = map[string]bool{
	"contact": true,
	"fields":  true,
	"globals": true,
	"urns":    true,
}

// Completion generates auto-complete paths
type Completion struct {
	Types         []Type      `json:"types"`
	Root          []*Property `json:"root"`
	RootNoSession []*Property `json:"root_no_session"`
}

// NewCompletion creates a new empty completion
func NewCompletion(types []Type, root []*Property) *Completion {
	// extract types which are available in a context without a session
	rootNoSession := make([]*Property, 0, len(rootNoSessionTypes))
	for _, p := range root {
		if rootNoSessionTypes[p.Type] {
			rootNoSession = append(rootNoSession, p)
		}
	}

	return &Completion{Types: types, Root: root, RootNoSession: rootNoSession}
}

// Validate checks that all type references are valid
func (c *Completion) Validate() error {
	knownTypes := make(map[string]bool, len(c.Types))
	for _, t := range primitiveTypes {
		knownTypes[t.Name()] = true
	}
	for _, t := range c.Types {
		knownTypes[t.Name()] = true
	}

	for _, t := range c.Types {
		for _, ref := range t.TypeRefs() {
			if !knownTypes[ref] {
				return errors.Errorf("context type %s references unknown type %s", t.Name(), ref)
			}
		}
	}
	for _, p := range c.Root {
		if !knownTypes[p.Type] {
			return errors.Errorf("context root references unknown type %s", p.Type)
		}
	}
	return nil
}

// Node represents a part of the context that can be referenced
type Node struct {
	Path string
	Help string
}

// EnumerateNodes walks the context to enumerate all possible nodes
func (c *Completion) EnumerateNodes(context *Context) []Node {
	// make a lookup of all types by their name
	types := make(map[string]Type, len(c.Types))
	for _, t := range primitiveTypes {
		types[t.Name()] = t
	}
	for _, t := range c.Types {
		types[t.Name()] = t
	}

	nodes := make([]Node, 0)

	callback := func(path, help string) {
		nodes = append(nodes, Node{path, help})
	}
	for _, p := range c.Root {
		enumeratePaths("", p, types, context, callback)
	}

	return nodes
}

func enumeratePaths(base string, p *Property, types map[string]Type, context *Context, callback func(path, help string)) {
	t := types[p.Type]

	path := p.Key
	if base != "" {
		path = base + "." + path
	}

	props := t.EnumerateProperties(context)
	help := p.Help

	// if this type has a default value, append that to the help
	for _, pp := range props {
		if pp.Key == "__default__" {
			help += fmt.Sprintf(" (defaults to %s)", pp.Help)
			break
		}
	}

	callback(path, help)

	if p.Array {
		path += "[0]"
		callback(path, "first of "+help)
	}

	for _, pp := range props {
		if pp.Key != "__default__" {
			enumeratePaths(path, pp, types, context, callback)
		}
	}
}
