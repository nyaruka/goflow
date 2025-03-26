package static

import (
	"github.com/nyaruka/goflow/assets"
)

// LLM is a JSON serializable implementation of an LLM asset
type LLM struct {
	UUID_ assets.LLMUUID `json:"uuid" validate:"required,uuid"`
	Name_ string         `json:"name"`
	Type_ string         `json:"type"`
}

// NewLLM creates a new LLM
func NewLLM(uuid assets.LLMUUID, name string, type_ string) assets.LLM {
	return &LLM{
		UUID_: uuid,
		Name_: name,
		Type_: type_,
	}
}

// UUID returns the UUID of this LLM
func (l *LLM) UUID() assets.LLMUUID { return l.UUID_ }

// Name returns the name of this LLM
func (l *LLM) Name() string { return l.Name_ }

// Type returns the type of this LLM
func (l *LLM) Type() string { return l.Type_ }
