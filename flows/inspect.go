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

// NewDependencies creates a new dependency listing from the slice of references
func NewDependencies(refs []assets.Reference) *Dependencies {
	d := &Dependencies{}
	for _, r := range refs {
		switch typed := r.(type) {
		case *assets.ChannelReference:
			d.Channels = append(d.Channels, InspectedChannelReference{ChannelReference: typed})
		case *assets.ClassifierReference:
			d.Classifiers = append(d.Classifiers, InspectedClassifierReference{ClassifierReference: typed})
		case *ContactReference:
			d.Contacts = append(d.Contacts, InspectedContactReference{ContactReference: typed})
		case *assets.FieldReference:
			d.Fields = append(d.Fields, InspectedFieldReference{FieldReference: typed})
		case *assets.FlowReference:
			d.Flows = append(d.Flows, InspectedFlowReference{FlowReference: typed})
		case *assets.GlobalReference:
			d.Globals = append(d.Globals, InspectedGlobalReference{GlobalReference: typed})
		case *assets.GroupReference:
			d.Groups = append(d.Groups, InspectedGroupReference{GroupReference: typed})
		case *assets.LabelReference:
			d.Labels = append(d.Labels, InspectedLabelReference{LabelReference: typed})
		case *assets.TemplateReference:
			d.Templates = append(d.Templates, InspectedTemplateReference{TemplateReference: typed})
		default:
			panic(fmt.Sprintf("unknown dependency type reference: %v", r))
		}
	}
	return d
}

// Check checks the asset dependencies and notifies the caller of missing assets via the callback
func (d *Dependencies) Check(sa SessionAssets, missing assets.MissingCallback) error {
	callback := func(iref InspectedReference, ref assets.Reference, err error) {
		iref.Missing = true
		missing(ref, err)
	}

	for _, ref := range d.Channels {
		if sa.Channels().Get(ref.UUID) == nil {
			callback(ref.InspectedReference, ref.ChannelReference, nil)
		}
	}
	for _, ref := range d.Classifiers {
		if sa.Classifiers().Get(ref.UUID) == nil {
			callback(ref.InspectedReference, ref.ClassifierReference, nil)
		}
	}
	for _, ref := range d.Fields {
		if sa.Fields().Get(ref.Key) == nil {
			callback(ref.InspectedReference, ref.FieldReference, nil)
		}
	}
	for _, ref := range d.Flows {
		_, err := sa.Flows().Get(ref.UUID)
		if err != nil {
			callback(ref.InspectedReference, ref.FlowReference, err)
		}
	}
	for _, ref := range d.Globals {
		if sa.Globals().Get(ref.Key) == nil {
			callback(ref.InspectedReference, ref.GlobalReference, nil)
		}
	}
	for _, ref := range d.Groups {
		if sa.Groups().Get(ref.UUID) == nil {
			callback(ref.InspectedReference, ref.GroupReference, nil)
		}
	}
	for _, ref := range d.Labels {
		if sa.Labels().Get(ref.UUID) == nil {
			callback(ref.InspectedReference, ref.LabelReference, nil)
		}
	}
	for _, ref := range d.Templates {
		if sa.Templates().Get(ref.UUID) == nil {
			callback(ref.InspectedReference, ref.TemplateReference, nil)
		}
	}

	return nil
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
