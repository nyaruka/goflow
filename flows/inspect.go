package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// FlowInfo contains the results of flow inspection
type FlowInfo struct {
	Dependencies *Dependencies `json:"dependencies"`
	Results      []*ResultInfo `json:"results"`
	WaitingExits []ExitUUID    `json:"waiting_exits"`
	ParentRefs   []string      `json:"parent_refs"`
}

type InspectedReference struct {
	Missing bool `json:"missing,omitempty"`
}

type InspectedClassifierReference struct {
	*assets.ClassifierReference
	InspectedReference
}

type InspectedChannelReference struct {
	*assets.ChannelReference
	InspectedReference
}

type InspectedContactReference struct {
	*ContactReference
	InspectedReference
}

type InspectedFieldReference struct {
	*assets.FieldReference
	InspectedReference
}

type InspectedFlowReference struct {
	*assets.FlowReference
	InspectedReference
}

type InspectedGlobalReference struct {
	*assets.GlobalReference
	InspectedReference
}

type InspectedGroupReference struct {
	*assets.GroupReference
	InspectedReference
}

type InspectedLabelReference struct {
	*assets.LabelReference
	InspectedReference
}

type InspectedTemplateReference struct {
	*assets.TemplateReference
	InspectedReference
}

// Dependencies contains a flows dependencies grouped by type
type Dependencies struct {
	Classifiers []InspectedClassifierReference `json:"classifiers,omitempty"`
	Channels    []InspectedChannelReference    `json:"channels,omitempty"`
	Contacts    []InspectedContactReference    `json:"contacts,omitempty"`
	Fields      []InspectedFieldReference      `json:"fields,omitempty"`
	Flows       []InspectedFlowReference       `json:"flows,omitempty"`
	Globals     []InspectedGlobalReference     `json:"globals,omitempty"`
	Groups      []InspectedGroupReference      `json:"groups,omitempty"`
	Labels      []InspectedLabelReference      `json:"labels,omitempty"`
	Templates   []InspectedTemplateReference   `json:"templates,omitempty"`
}

// NewDependencies creates a new dependency listing from the slice of references. If a session assets is provided,
// each dependency is checked to see if it is available or missing.
func NewDependencies(refs []assets.Reference, sa SessionAssets) *Dependencies {
	d := &Dependencies{}
	for _, r := range refs {
		switch typed := r.(type) {
		case *assets.ChannelReference:
			missing := sa != nil && sa.Channels().Get(typed.UUID) == nil
			d.Channels = append(d.Channels, InspectedChannelReference{
				ChannelReference:   typed,
				InspectedReference: InspectedReference{Missing: missing},
			})
		case *assets.ClassifierReference:
			missing := sa != nil && sa.Classifiers().Get(typed.UUID) == nil
			d.Classifiers = append(d.Classifiers, InspectedClassifierReference{
				ClassifierReference: typed,
				InspectedReference:  InspectedReference{Missing: missing},
			})
		case *ContactReference:
			d.Contacts = append(d.Contacts, InspectedContactReference{
				ContactReference: typed,
			})
		case *assets.FieldReference:
			missing := sa != nil && sa.Fields().Get(typed.Key) == nil
			d.Fields = append(d.Fields, InspectedFieldReference{
				FieldReference:     typed,
				InspectedReference: InspectedReference{Missing: missing},
			})
		case *assets.FlowReference:
			missing := false
			if sa != nil {
				_, err := sa.Flows().Get(typed.UUID)
				missing = err != nil
			}
			d.Flows = append(d.Flows, InspectedFlowReference{
				FlowReference:      typed,
				InspectedReference: InspectedReference{Missing: missing},
			})
		case *assets.GlobalReference:
			missing := sa != nil && sa.Globals().Get(typed.Key) == nil
			d.Globals = append(d.Globals, InspectedGlobalReference{
				GlobalReference:    typed,
				InspectedReference: InspectedReference{Missing: missing},
			})
		case *assets.GroupReference:
			missing := sa != nil && sa.Groups().Get(typed.UUID) == nil
			d.Groups = append(d.Groups, InspectedGroupReference{
				GroupReference:     typed,
				InspectedReference: InspectedReference{Missing: missing},
			})
		case *assets.LabelReference:
			missing := sa != nil && sa.Labels().Get(typed.UUID) == nil
			d.Labels = append(d.Labels, InspectedLabelReference{
				LabelReference:     typed,
				InspectedReference: InspectedReference{Missing: missing},
			})
		case *assets.TemplateReference:
			missing := sa != nil && sa.Templates().Get(typed.UUID) == nil
			d.Templates = append(d.Templates, InspectedTemplateReference{
				TemplateReference:  typed,
				InspectedReference: InspectedReference{Missing: missing},
			})
		default:
			panic(fmt.Sprintf("unknown dependency type reference: %v", r))
		}
	}
	return d
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
