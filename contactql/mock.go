package contactql

import (
	"strings"

	"github.com/nyaruka/goflow/assets"
)

type mockResolver struct {
	fields map[string]assets.Field
	groups map[string]assets.Group
}

// NewMockResolver creates a new mock resolver for fields and groups
func NewMockResolver(fields map[string]assets.Field, groups map[string]assets.Group) Resolver {
	return &mockResolver{
		fields: fields,
		groups: groups,
	}
}

func (r *mockResolver) ResolveField(key string) assets.Field {
	field, found := r.fields[key]
	if !found {
		return nil
	}
	return field
}

func (r *mockResolver) ResolveGroup(name string) assets.Group {
	group, found := r.groups[strings.ToLower(name)]
	if !found {
		return nil
	}
	return group
}
