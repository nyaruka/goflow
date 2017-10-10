package flows

import (
	"encoding/json"
	"strings"

	"github.com/nyaruka/goflow/utils"
)

// Label represents a message label
type Label struct {
	uuid LabelUUID
	name string
}

// NewLabel creates a new label given the passed in uuid and name
func NewLabel(uuid LabelUUID, name string) *Label {
	return &Label{uuid, name}
}

// UUID returns the UUID of this label
func (l *Label) UUID() LabelUUID { return l.uuid }

// Name returns the name of this label
func (l *Label) Name() string { return l.name }

func (l *Label) Reference() *LabelReference { return NewLabelReference(l.uuid, l.name) }

// LabelSet defines the unordered set of all labels for a session
type LabelSet struct {
	labels       []*Label
	labelsByUUID map[LabelUUID]*Label
}

func NewLabelSet(labels []*Label) *LabelSet {
	s := &LabelSet{labels: labels, labelsByUUID: make(map[LabelUUID]*Label, len(labels))}
	for _, label := range s.labels {
		s.labelsByUUID[label.uuid] = label
	}
	return s
}

func (s *LabelSet) FindByUUID(uuid LabelUUID) *Label {
	return s.labelsByUUID[uuid]
}

// FindByName looks for a label with the given name (case-insensitive)
func (s *LabelSet) FindByName(name string) *Label {
	name = strings.ToLower(name)
	for _, label := range s.labels {
		if strings.ToLower(label.name) == name {
			return label
		}
	}
	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type labelEnvelope struct {
	UUID LabelUUID `json:"uuid" validate:"required,uuid4"`
	Name string    `json:"name"`
}

func ReadLabel(data json.RawMessage) (*Label, error) {
	var le labelEnvelope
	if err := utils.UnmarshalAndValidate(data, &le, "label"); err != nil {
		return nil, err
	}

	return NewLabel(le.UUID, le.Name), nil
}

func ReadLabelSet(data json.RawMessage) (*LabelSet, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	labels := make([]*Label, len(items))
	for d := range items {
		if labels[d], err = ReadLabel(items[d]); err != nil {
			return nil, err
		}
	}

	return NewLabelSet(labels), nil
}
