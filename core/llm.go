package core

import (
	"slices"

	"github.com/nyaruka/goflow/assets"
)

// LLM represents a large language model.
type LLM struct {
	assets.LLM
}

// NewLLM returns a new LLM object from the given LLM asset
func NewLLM(asset assets.LLM) *LLM {
	return &LLM{LLM: asset}
}

// Asset returns the underlying asset
func (l *LLM) Asset() assets.LLM { return l.LLM }

// Reference returns a reference to this LLM
func (l *LLM) Reference() *assets.LLMReference {
	return assets.NewLLMReference(l.UUID(), l.Name())
}

// HasRole returns whether this LLM has the given role
func (l *LLM) HasRole(role assets.LLMRole) bool {
	return slices.Contains(l.Roles(), role)
}

// LLMAssets provides access to all LLM assets
type LLMAssets struct {
	byUUID map[assets.LLMUUID]*LLM
}

// NewLLMAssets creates a new set of LLM assets
func NewLLMAssets(llms []assets.LLM) *LLMAssets {
	s := &LLMAssets{
		byUUID: make(map[assets.LLMUUID]*LLM, len(llms)),
	}
	for _, asset := range llms {
		s.byUUID[asset.UUID()] = NewLLM(asset)
	}
	return s
}

// Get returns the LLM with the given UUID
func (s *LLMAssets) Get(uuid assets.LLMUUID) *LLM {
	return s.byUUID[uuid]
}

// LLMResponse is the response from an LLM service call
type LLMResponse struct {
	Output       string
	TokensInput  int64
	TokensOutput int64
}

// LLMClassification is the result of an LLM service classification call
type LLMClassification struct {
	Category      string             // the chosen category, which must be one of the given categories
	Probabilities map[string]float64 // per-category probabilities, if the model provides them
	TokensInput   int64
	TokensOutput  int64
}

// CategoryProbability returns the probability of the chosen category, if the model provided probabilities
func (c *LLMClassification) CategoryProbability() (float64, bool) {
	if c == nil {
		return 0, false
	}
	p, ok := c.Probabilities[c.Category]
	return p, ok
}
