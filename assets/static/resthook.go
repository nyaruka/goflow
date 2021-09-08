package static

import (
	"github.com/nyaruka/goflow/assets"
)

// Resthook is a JSON serializable implementation of a resthook asset
type Resthook struct {
	Slug_        string   `json:"slug" validate:"required"`
	Subscribers_ []string `json:"subscribers" validate:"required,dive,url"`
}

// NewResthook creates a new resthook
func NewResthook(slug string, subscribers []string) assets.Resthook {
	return &Resthook{Slug_: slug, Subscribers_: subscribers}
}

// Slug returns the slug of the resthook
func (r *Resthook) Slug() string { return r.Slug_ }

// Subscribers returns the subscribers to the resthook
func (r *Resthook) Subscribers() []string { return r.Subscribers_ }
