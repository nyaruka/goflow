package static

import (
	"github.com/nyaruka/goflow/assets"
)

// Global is a JSON serializable implementation of a global asset
type Global struct {
	Key_   string `json:"key" validate:"required"`
	Name_  string `json:"name"`
	Value_ string `json:"value"`
}

// NewGlobal creates a new global
func NewGlobal(key, name, value string) assets.Global {
	return &Global{
		Key_:   key,
		Name_:  name,
		Value_: value,
	}
}

// Key returns the key of this global
func (g *Global) Key() string { return g.Key_ }

// Name returns the name of this global
func (g *Global) Name() string { return g.Name_ }

// Value returns the type of this global
func (g *Global) Value() string { return g.Value_ }
