package flows

import "github.com/nyaruka/goflow/assets"

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
