package types

import (
	"github.com/nyaruka/goflow/assets"
)

// Classifier is a JSON serializable implementation of a classifier asset
type Classifier struct {
	UUID_     assets.ClassifierUUID `json:"uuid" validate:"required,uuid"`
	Name_     string                `json:"name"`
	Provider_ string                `json:"provider"`
}

// NewClassifier creates a new classifier
func NewClassifier(uuid assets.ClassifierUUID, name string, provider string) assets.Classifier {
	return &Classifier{
		UUID_:     uuid,
		Name_:     name,
		Provider_: provider,
	}
}

// UUID returns the UUID of this channel
func (c *Classifier) UUID() assets.ClassifierUUID { return c.UUID_ }

// Name returns the name of this channel
func (c *Classifier) Name() string { return c.Name_ }

// Provider returns the provider of this classifier
func (c *Classifier) Provider() string { return c.Provider_ }
