package types

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// json serializable implementation of a resthook asset
type resthook struct {
	Slug_        string   `json:"slug" validate:"required"`
	Subscribers_ []string `json:"subscribers" validate:"required,dive,url"`
}

func NewResthook(slug string, subscribers []string) assets.Resthook {
	return &resthook{Slug_: slug, Subscribers_: subscribers}
}

// Slug returns the slug of the resthook
func (r *resthook) Slug() string { return r.Slug_ }

// Subscribers returns the subscribers to the resthook
func (r *resthook) Subscribers() []string { return r.Subscribers_ }

// ReadResthook reads a resthook from the given JSON
func ReadResthook(data json.RawMessage) (assets.Resthook, error) {
	r := &resthook{}
	if err := utils.UnmarshalAndValidate(data, r); err != nil {
		return nil, fmt.Errorf("unable to read resthook: %s", err)
	}

	return r, nil
}

// ReadResthooks reads a resthook set from the given JSON
func ReadResthooks(data json.RawMessage) ([]assets.Resthook, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	resthooks := make([]assets.Resthook, len(items))
	for d := range items {
		if resthooks[d], err = ReadResthook(items[d]); err != nil {
			return nil, err
		}
	}

	return resthooks, nil
}
