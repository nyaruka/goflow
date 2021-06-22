package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Global represents a constant value available in expressions.
type Global struct {
	assets.Global
}

// NewGlobal returns a new global object from the given global asset
func NewGlobal(asset assets.Global) *Global {
	return &Global{Global: asset}
}

// Asset returns the underlying asset
func (g *Global) Asset() assets.Global { return g.Global }

// Reference returns a reference to this global
func (g *Global) Reference() *assets.GlobalReference {
	return assets.NewGlobalReference(g.Key(), g.Name())
}

// GlobalAssets provides access to all global assets
type GlobalAssets struct {
	all   []*Global
	byKey map[string]*Global
}

// NewGlobalAssets creates a new set of global assets
func NewGlobalAssets(globals []assets.Global) *GlobalAssets {
	s := &GlobalAssets{
		all:   make([]*Global, len(globals)),
		byKey: make(map[string]*Global, len(globals)),
	}
	for i, asset := range globals {
		global := NewGlobal(asset)
		s.all[i] = global
		s.byKey[global.Key()] = global
	}
	return s
}

// Get returns the global with the given key
func (s *GlobalAssets) Get(key string) *Global {
	return s.byKey[key]
}

// Context returns the properties available in expressions
func (s *GlobalAssets) Context(env envs.Environment) map[string]types.XValue {
	entries := make(map[string]types.XValue, len(s.all)+1)
	entries["__default__"] = types.NewXText(s.format())

	for _, g := range s.all {
		entries[g.Key()] = types.NewXText(g.Value())
	}
	return entries
}

func (s *GlobalAssets) format() string {
	lines := make([]string, 0, len(s.all))
	for _, g := range s.all {
		lines = append(lines, fmt.Sprintf("%s: %s", g.Name(), g.Value()))
	}
	return strings.Join(lines, "\n")
}
