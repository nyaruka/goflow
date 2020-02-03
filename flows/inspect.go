package flows

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// FlowInfo contains the results of flow inspection
type FlowInfo struct {
	Dependencies []InspectedReference `json:"dependencies"`
	Results      []*ResultInfo        `json:"results"`
	WaitingExits []ExitUUID           `json:"waiting_exits"`
	ParentRefs   []string             `json:"parent_refs"`
}

type ReferenceInfo struct {
	Type    string `json:"type"`
	Missing bool   `json:"missing,omitempty"`
}

type InspectedReference struct {
	assets.Reference
	ReferenceInfo
}

func (r InspectedReference) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(r.Reference)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(r.ReferenceInfo)
	if err != nil {
		return nil, err
	}
	b := append(b1[0:len(b1)-1], byte(','))
	b = append(b, b2[1:]...)
	return b, nil
}

// InspectReferences inspects a list of references. If a session assets is provided,
// each dependency is checked to see if it is available or missing.
func InspectReferences(refs []assets.Reference, sa SessionAssets) []InspectedReference {
	inspected := make([]InspectedReference, len(refs))

	for i, ref := range refs {
		var type_ string

		// TODO derive from type name
		switch ref.(type) {
		case *assets.ChannelReference:
			type_ = "channel"
		case *assets.ClassifierReference:
			type_ = "classifier"
		case *ContactReference:
			type_ = "contact"
		case *assets.FieldReference:
			type_ = "field"
		case *assets.FlowReference:
			type_ = "flow"
		case *assets.GlobalReference:
			type_ = "global"
		case *assets.GroupReference:
			type_ = "group"
		case *assets.LabelReference:
			type_ = "label"
		case *assets.TemplateReference:
			type_ = "template"
		default:
			panic(fmt.Sprintf("unknown dependency type reference: %T", ref))
		}

		missing := false
		if sa != nil {
			missing = !checkDependency(sa, ref)
		}

		inspected[i] = InspectedReference{Reference: ref, ReferenceInfo: ReferenceInfo{Type: type_, Missing: missing}}
	}
	return inspected
}

// determines whether the given dependency exists
func checkDependency(sa SessionAssets, ref assets.Reference) bool {
	switch typed := ref.(type) {
	case *assets.ChannelReference:
		return sa.Channels().Get(typed.UUID) != nil
	case *assets.ClassifierReference:
		return sa.Classifiers().Get(typed.UUID) != nil
	case *ContactReference:
		return true // have to assume contacts exist
	case *assets.FieldReference:
		return sa.Fields().Get(typed.Key) != nil
	case *assets.FlowReference:
		_, err := sa.Flows().Get(typed.UUID)
		return err == nil
	case *assets.GlobalReference:
		return sa.Globals().Get(typed.Key) != nil
	case *assets.GroupReference:
		return sa.Groups().Get(typed.UUID) != nil
	case *assets.LabelReference:
		return sa.Labels().Get(typed.UUID) != nil
	case *assets.TemplateReference:
		return sa.Templates().Get(typed.UUID) != nil
	default:
		panic(fmt.Sprintf("unknown dependency type reference: %T", ref))
	}
}

// ResultInfo is possible result that a flow might generate
type ResultInfo struct {
	Key        string     `json:"key"`
	Name       string     `json:"name"`
	Categories []string   `json:"categories"`
	NodeUUIDs  []NodeUUID `json:"node_uuids"`
}

// NewResultInfo creates a new result spec
func NewResultInfo(name string, categories []string, node Node) *ResultInfo {
	return &ResultInfo{
		Key:        utils.Snakify(name),
		Name:       name,
		Categories: categories,
		NodeUUIDs:  []NodeUUID{node.UUID()},
	}
}

func (r *ResultInfo) String() string {
	return fmt.Sprintf("key=%s|name=%s|categories=%s", r.Key, r.Name, strings.Join(r.Categories, ","))
}

// MergeResultInfos merges result specs based on key
func MergeResultInfos(specs []*ResultInfo) []*ResultInfo {
	merged := make([]*ResultInfo, 0, len(specs))
	byKey := make(map[string]*ResultInfo)

	for _, spec := range specs {
		existing := byKey[spec.Key]

		// merge if we already have a result info with this key
		if existing != nil {
			// merge categories
			for _, category := range spec.Categories {
				if !utils.StringSliceContains(existing.Categories, category, false) {
					existing.Categories = append(existing.Categories, category)
				}
			}

			// merge node UUIDs
			for _, nodeUUID := range spec.NodeUUIDs {
				uuidSeen := false
				for _, u := range existing.NodeUUIDs {
					if u == nodeUUID {
						uuidSeen = true
						break
					}
				}
				if !uuidSeen {
					existing.NodeUUIDs = append(existing.NodeUUIDs, nodeUUID)
				}
			}
		} else {
			// if not, add as new unique result spec
			merged = append(merged, spec)
			byKey[spec.Key] = spec
		}
	}
	return merged
}
