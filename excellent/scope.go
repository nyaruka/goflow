package excellent

import (
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
)

// Scope provides access to context objects
type Scope struct {
	get    func(string) (types.XValue, bool)
	parent *Scope
}

// NewScope creates a new evaluation scope with an optional parent
func NewScope(ctx *types.XObject, parent *Scope) *Scope {
	if parent == nil {
		parent = rootScope
	}
	return &Scope{get: ctx.Get, parent: parent}
}

// Get looks up a named value in the context
func (s *Scope) Get(name string) (types.XValue, bool) {
	v, exists := s.get(name)
	if exists {
		return v, true
	}

	if s.parent != nil {
		return s.parent.Get(name)
	}

	return nil, false
}

var rootScope = &Scope{
	get: func(name string) (types.XValue, bool) {
		// at the root, we only have functions
		function := functions.Lookup(name)
		if function != nil {
			return function, true
		}
		return nil, false
	},
	parent: nil,
}
