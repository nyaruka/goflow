package flows

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"
)

// Resthook represents a named event and a set of subscribers
type Resthook struct {
	slug        string
	subscribers []string
}

// NewResthook returns a new resthook object
func NewResthook(slug string, subscribers []string) *Resthook {
	return &Resthook{slug: slug, subscribers: subscribers}
}

// Slug returns the slug of the resthook
func (r *Resthook) Slug() string { return r.slug }

// Subscribers returns the subscribers to the resthook
func (r *Resthook) Subscribers() []string { return r.subscribers }

// ResthookSet defines the unordered set of all resthooks for a session
type ResthookSet struct {
	resthooksBySlug map[string]*Resthook
}

// NewResthookSet creates a new resthook set from the given list of resthooks
func NewResthookSet(resthooks []*Resthook) *ResthookSet {
	s := &ResthookSet{
		resthooksBySlug: make(map[string]*Resthook, len(resthooks)),
	}

	for _, resthook := range resthooks {
		s.resthooksBySlug[resthook.slug] = resthook
	}

	return s
}

// FindBySlug finds the group with the given UUID
func (s *ResthookSet) FindBySlug(slug string) *Resthook {
	return s.resthooksBySlug[slug]
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type resthookEnvelope struct {
	Slug        string   `json:"slug" validate:"required"`
	Subscribers []string `json:"subscribers" validate:"required,dive,url"`
}

// ReadResthook reads a resthook from the given JSON
func ReadResthook(data json.RawMessage) (*Resthook, error) {
	var e resthookEnvelope
	if err := utils.UnmarshalAndValidate(data, &e, "resthook"); err != nil {
		return nil, err
	}

	return NewResthook(e.Slug, e.Subscribers), nil
}

// ReadResthookSet reads a resthook set from the given JSON
func ReadResthookSet(data json.RawMessage) (*ResthookSet, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	resthooks := make([]*Resthook, len(items))
	for d := range items {
		if resthooks[d], err = ReadResthook(items[d]); err != nil {
			return nil, err
		}
	}

	return NewResthookSet(resthooks), nil
}
