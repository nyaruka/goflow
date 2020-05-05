package types

import (
	"github.com/nyaruka/goflow/assets"
)

// Ticketer is a JSON serializable implementation of a ticketer asset
type Ticketer struct {
	UUID_ assets.TicketerUUID `json:"uuid" validate:"required,uuid"`
	Name_ string              `json:"name"`
	Type_ string              `json:"type"`
}

// NewTicketer creates a new ticketer
func NewTicketer(uuid assets.TicketerUUID, name string, type_ string) assets.Ticketer {
	return &Ticketer{
		UUID_: uuid,
		Name_: name,
		Type_: type_,
	}
}

// UUID returns the UUID of this ticketer
func (t *Ticketer) UUID() assets.TicketerUUID { return t.UUID_ }

// Name returns the name of this ticketer
func (t *Ticketer) Name() string { return t.Name_ }

// Type returns the type of this ticketer
func (t *Ticketer) Type() string { return t.Type_ }
