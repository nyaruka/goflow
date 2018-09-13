package flows

import (
	"github.com/nyaruka/goflow/assets"
)

// Resthook represents a named event and a set of subscribers
type Resthook struct {
	assets.Resthook
}

// NewResthook returns a new resthook object
func NewResthook(asset assets.Resthook) *Resthook {
	return &Resthook{Resthook: asset}
}

// Asset returns the underlying asset
func (r *Resthook) Asset() assets.Resthook { return r.Resthook }

// ResthookAssets provides access to all resthook assets
type ResthookAssets struct {
	bySlug map[string]*Resthook
}

// NewResthookAssets creates a new set of resthook assets
func NewResthookAssets(resthooks []assets.Resthook) *ResthookAssets {
	s := &ResthookAssets{
		bySlug: make(map[string]*Resthook, len(resthooks)),
	}
	for _, asset := range resthooks {
		s.bySlug[asset.Slug()] = NewResthook(asset)
	}
	return s
}

// FindBySlug finds the group with the given UUID
func (s *ResthookAssets) FindBySlug(slug string) *Resthook {
	return s.bySlug[slug]
}
