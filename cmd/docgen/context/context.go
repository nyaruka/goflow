package context

import (
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
