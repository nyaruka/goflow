package context

import (
	"strings"

	"github.com/pkg/errors"
)

type Context struct {
	Types []Type      `json:"types"`
	Root  []*Property `json:"root"`
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) SetRoot(root []*Property) {
	c.Root = root
}

func (c *Context) AddType(t Type) {
	c.Types = append(c.Types, t)
}

func (c *Context) Validate() error {
	knownTypes := make(map[string]bool, len(c.Types))
	for _, t := range primitiveTypes {
		knownTypes[t.TypeName()] = true
	}
	for _, t := range c.Types {
		knownTypes[t.TypeName()] = true
	}

	for _, t := range c.Types {
		for _, ref := range t.TypeRefs() {
			if !knownTypes[ref] {
				return errors.Errorf("context type %s references unknown type %s", t.TypeName(), ref)
			}
		}
	}
	for _, p := range c.Root {
		if !knownTypes[p.TypeRef] {
			return errors.Errorf("context root references unknown type %s", p.TypeRef)
		}
	}
	return nil
}

type Node struct {
	Path        string
	Description string
}

func (c *Context) EnumerateNodes(sources map[string][]string) []Node {
	// make a lookup of all types by their name
	types := make(map[string]Type, len(c.Types))
	for _, t := range primitiveTypes {
		types[t.TypeName()] = t
	}
	for _, t := range c.Types {
		types[t.TypeName()] = t
	}

	nodes := make([]Node, 0)

	callback := func(path, desc string) {
		nodes = append(nodes, Node{path, desc})
	}
	for _, p := range c.Root {
		enumeratePaths("", p, types, sources, callback)
	}

	return nodes
}

func enumeratePaths(base string, p *Property, types map[string]Type, sources map[string][]string, callback func(path, desc string)) {
	t := types[p.TypeRef]

	path := p.Name
	if base != "" {
		path = base + "." + path
	}
	callback(path, p.Description)

	if p.Array {
		path += "[0]"
		callback(path, "first of "+p.Description)
	}

	switch typed := t.(type) {
	case *primitiveType:
		// primitives don't have properties
	case *staticType:
		for _, pp := range typed.Properties {
			enumeratePaths(path, pp, types, sources, callback)
		}
	case *dynamicType:
		nameTemplate := typed.PropertyTemplate.Name
		descTemplate := typed.PropertyTemplate.Description

		for _, key := range sources[typed.Source] {
			name := strings.Replace(nameTemplate, "{key}", key, -1)
			description := strings.Replace(descTemplate, "{key}", key, -1)
			pp := NewProperty(name, description, typed.PropertyTemplate.TypeRef)
			enumeratePaths(path, pp, types, sources, callback)
		}
	}
}
