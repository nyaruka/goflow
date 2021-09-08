package static

import (
	"github.com/nyaruka/goflow/assets"
)

// Classifier is a JSON serializable implementation of a classifier asset
type Classifier struct {
	UUID_    assets.ClassifierUUID `json:"uuid" validate:"required,uuid"`
	Name_    string                `json:"name"`
	Type_    string                `json:"type"`
	Intents_ []string              `json:"intents"`
}

// NewClassifier creates a new classifier
func NewClassifier(uuid assets.ClassifierUUID, name string, type_ string, intents []string) assets.Classifier {
	return &Classifier{
		UUID_:    uuid,
		Name_:    name,
		Type_:    type_,
		Intents_: intents,
	}
}

// UUID returns the UUID of this channel
func (c *Classifier) UUID() assets.ClassifierUUID { return c.UUID_ }

// Name returns the name of this channel
func (c *Classifier) Name() string { return c.Name_ }

// Type returns the type of this classifier
func (c *Classifier) Type() string { return c.Type_ }

// Intents returns the intents of this classifier
func (c *Classifier) Intents() []string { return c.Intents_ }
