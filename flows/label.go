package flows

import (
	"strings"

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
func (l *Label) Reference() *assets.LabelReference {
	if l == nil {
		return nil
	}
	return assets.NewLabelReference(l.UUID(), l.Name())
}

var _ assets.Label = (*Label)(nil)

// LabelAssets provides access to all label assets
type LabelAssets struct {
	all    []*Label
	byUUID map[assets.LabelUUID]*Label
}

// NewLabelAssets creates a new set of label assets
func NewLabelAssets(labels []assets.Label) *LabelAssets {
	s := &LabelAssets{
		all:    make([]*Label, len(labels)),
		byUUID: make(map[assets.LabelUUID]*Label, len(labels)),
	}
	for i, asset := range labels {
		label := NewLabel(asset)
		s.all[i] = label
		s.byUUID[label.UUID()] = label
	}
	return s
}

// All returns all the labels
func (s *LabelAssets) All() []*Label {
	return s.all
}

// Get returns the label with the given UUID
func (s *LabelAssets) Get(uuid assets.LabelUUID) *Label {
	return s.byUUID[uuid]
}

// FindByName looks for a label with the given name (case-insensitive)
func (s *LabelAssets) FindByName(name string) *Label {
	name = strings.ToLower(name)
	for _, label := range s.all {
		if strings.ToLower(label.Name()) == name {
			return label
		}
	}
	return nil
}
