package contactql

import (
	"strings"

	"github.com/nyaruka/goflow/assets"
)

type mockResolver struct {
	fields []assets.Field
	flows  []assets.Flow
	groups []assets.Group
}

// NewMockResolver creates a new mock resolver for fields and groups
func NewMockResolver(fields []assets.Field, flows []assets.Flow, groups []assets.Group) Resolver {
	return &mockResolver{
		fields: fields,
		flows:  flows,
		groups: groups,
	}
}

func (r *mockResolver) ResolveField(key string) assets.Field {
	for _, f := range r.fields {
		if f.Key() == key {
			return f
		}
	}
	return nil
}

func (r *mockResolver) ResolveGroup(name string) assets.Group {
	for _, g := range r.groups {
		if strings.EqualFold(g.Name(), name) {
			return g
		}
	}
	return nil
}

func (r *mockResolver) ResolveFlow(name string) assets.Flow {
	for _, f := range r.flows {
		if strings.EqualFold(f.Name(), name) {
			return f
		}
	}
	return nil
}
