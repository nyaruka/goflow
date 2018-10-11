package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
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

// ReadResthooks reads a resthook set from the given JSON
func ReadResthooks(data json.RawMessage) ([]assets.Resthook, error) {
	var items []*Resthook
	if err := utils.UnmarshalAndValidate(data, &items); err != nil {
		return nil, err
	}

	asAssets := make([]assets.Resthook, len(items))
	for i := range items {
		asAssets[i] = items[i]
	}

	return asAssets, nil
}
