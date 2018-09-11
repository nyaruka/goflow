package flows

import (
	"github.com/nyaruka/goflow/assets"
)

// Label represents a message label
type Label struct {
	assets.Label
}

// NewLabel creates a new label from the given asset
func NewLabel(asset assets.Label) *Label {
	return &Label{Label: asset}
}

// Asset returns the underlying asset
func (l *Label) Asset() assets.Label { return l.Label }

// Reference returns a reference to this label
func (l *Label) Reference() *LabelReference { return NewLabelReference(l.UUID(), l.Name()) }

var _ assets.Label = (*Label)(nil)
