package core

import (
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.RegisterValidatorAlias("classifier_confidence", "eq=none|eq=low|eq=medium|eq=high", func(validator.FieldError) string {
		return "is not a valid classifier confidence"
	})
	utils.RegisterValidatorAlias("required_confidence", "eq=any|eq=low|eq=medium|eq=high", func(validator.FieldError) string {
		return "is not a valid required confidence"
	})
}

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

// ModelTokens is the number of tokens a model call consumed
type ModelTokens struct {
	Input  int64 `json:"input"`
	Output int64 `json:"output"`
}

// ModelResponse is the response from a model service call
type ModelResponse struct {
	Output string
	Tokens ModelTokens
}

// ClassifierOption is an option that input can be classified as. The description helps the model decide whether the
// option fits.
type ClassifierOption struct {
	Name        string `json:"name"                  validate:"required,result_category"`
	Description string `json:"description,omitempty" validate:"max=1000"`
}

// ClassifierConfidence is a level of confidence in a classification. Model services map whatever measure of confidence
// their model provides onto these levels, so that a level means the same thing whichever model is used.
type ClassifierConfidence string

// possible classifier confidence levels, where any is only used as a requirement that any level meets
const (
	ClassifierConfidenceAny    ClassifierConfidence = "any"
	ClassifierConfidenceNone   ClassifierConfidence = "none"
	ClassifierConfidenceLow    ClassifierConfidence = "low"
	ClassifierConfidenceMedium ClassifierConfidence = "medium"
	ClassifierConfidenceHigh   ClassifierConfidence = "high"
)

var classifierConfidenceRanks = map[ClassifierConfidence]int{
	ClassifierConfidenceNone:   0,
	ClassifierConfidenceLow:    1,
	ClassifierConfidenceMedium: 2,
	ClassifierConfidenceHigh:   3,
}

// Meets returns whether this confidence level meets the given required level
func (c ClassifierConfidence) Meets(required ClassifierConfidence) bool {
	return required == ClassifierConfidenceAny || classifierConfidenceRanks[c] >= classifierConfidenceRanks[required]
}

// Classification is the result of a model service classification call
type Classification struct {
	Option        string               // the name of the chosen option, which must be one of the given options
	Confidence    ClassifierConfidence // confidence in the chosen option
	Probabilities map[string]float64   // per-option probabilities, if the model provides them, which may be partial
	Tokens        ModelTokens
}

// Translation is the result of translating items of text
type Translation struct {
	Items  map[string][]string // translated strings by item, with any untranslatable items omitted
	Tokens ModelTokens
}
