package core

import (
	"slices"

	"github.com/nyaruka/goflow/assets"
)

// Model represents an AI model.
type Model struct {
	assets.Model
}

// NewModel returns a new model object from the given model asset
func NewModel(asset assets.Model) *Model {
	return &Model{Model: asset}
}

// Asset returns the underlying asset
func (m *Model) Asset() assets.Model { return m.Model }

// Reference returns a reference to this model
func (m *Model) Reference() *assets.ModelReference {
	return assets.NewModelReference(m.UUID(), m.Name())
}

// HasRole returns whether this model has the given role
func (m *Model) HasRole(role assets.ModelRole) bool {
	return slices.Contains(m.Roles(), role)
}

// ModelAssets provides access to all model assets
type ModelAssets struct {
	byUUID map[assets.ModelUUID]*Model
}

// NewModelAssets creates a new set of model assets
func NewModelAssets(models []assets.Model) *ModelAssets {
	s := &ModelAssets{
		byUUID: make(map[assets.ModelUUID]*Model, len(models)),
	}
	for _, asset := range models {
		s.byUUID[asset.UUID()] = NewModel(asset)
	}
	return s
}

// Get returns the model with the given UUID
func (s *ModelAssets) Get(uuid assets.ModelUUID) *Model {
	return s.byUUID[uuid]
}

// ModelResponse is the response from a model service call
type ModelResponse struct {
	Output       string
	TokensInput  int64
	TokensOutput int64
}

// ModelClassification is the result of a model service classification call
type ModelClassification struct {
	Category      string             // the chosen category, which must be one of the given categories
	Confidence    float64            // confidence in the chosen category, between 0 and 1
	Probabilities map[string]float64 // per-category probabilities, if the model provides them
	TokensInput   int64
	TokensOutput  int64
}
