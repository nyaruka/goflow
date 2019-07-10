package completion

import (
	"github.com/pkg/errors"
)

// Completion generates auto-complete paths
type Completion struct {
	Types []Type      `json:"types"`
	Root  []*Property `json:"root"`
}

// NewCompletion creates a new empty completion
func NewCompletion(types []Type, root []*Property) *Completion {
	return &Completion{Types: types, Root: root}
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

// EnumerateNodes walks the context to enumerate all posible nodes
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
	callback(path, p.Help)

	if p.Array {
		path += "[0]"
		callback(path, "first of "+p.Help)
	}

	for _, pp := range t.EnumerateProperties(context) {
		enumeratePaths(path, pp, types, context, callback)
	}
}
