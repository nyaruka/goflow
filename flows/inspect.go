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

// Dependencies contains a flows dependencies grouped by type
type Dependencies struct {
	Classifiers []*assets.ClassifierReference `json:"classifiers,omitempty"`
	Channels    []*assets.ChannelReference    `json:"channels,omitempty"`
	Contacts    []*ContactReference           `json:"contacts,omitempty"`
	Fields      []*assets.FieldReference      `json:"fields,omitempty"`
	Flows       []*assets.FlowReference       `json:"flows,omitempty"`
	Globals     []*assets.GlobalReference     `json:"globals,omitempty"`
	Groups      []*assets.GroupReference      `json:"groups,omitempty"`
	Labels      []*assets.LabelReference      `json:"labels,omitempty"`
	Templates   []*assets.TemplateReference   `json:"templates,omitempty"`
}

// NewDependencies creates a new dependency listing from the slice of references
func NewDependencies(refs []assets.Reference) *Dependencies {
	d := &Dependencies{}
	for _, r := range refs {
		switch typed := r.(type) {
		case *assets.ChannelReference:
			d.Channels = append(d.Channels, typed)
		case *assets.ClassifierReference:
			d.Classifiers = append(d.Classifiers, typed)
		case *ContactReference:
			d.Contacts = append(d.Contacts, typed)
		case *assets.FieldReference:
			d.Fields = append(d.Fields, typed)
		case *assets.FlowReference:
			d.Flows = append(d.Flows, typed)
		case *assets.GlobalReference:
			d.Globals = append(d.Globals, typed)
		case *assets.GroupReference:
			d.Groups = append(d.Groups, typed)
		case *assets.LabelReference:
			d.Labels = append(d.Labels, typed)
		case *assets.TemplateReference:
			d.Templates = append(d.Templates, typed)
		default:
			panic(fmt.Sprintf("unknown dependency type reference: %v", r))
		}
	}
	return d
}

// Check checks the asset dependencies and notifies the caller of missing assets via the callback
func (d *Dependencies) Check(sa SessionAssets, missing assets.MissingCallback) error {
	for _, ref := range d.Channels {
		if sa.Channels().Get(ref.UUID) == nil {
			missing(ref, nil)
		}
	}
	for _, ref := range d.Classifiers {
		if sa.Classifiers().Get(ref.UUID) == nil {
			missing(ref, nil)
		}
	}
	for _, ref := range d.Fields {
		if sa.Fields().Get(ref.Key) == nil {
			missing(ref, nil)
		}
	}
	for _, ref := range d.Flows {
		_, err := sa.Flows().Get(ref.UUID)
		if err != nil {
			missing(ref, err)
		}
	}
	for _, ref := range d.Globals {
		if sa.Globals().Get(ref.Key) == nil {
			missing(ref, nil)
		}
	}
	for _, ref := range d.Groups {
		if sa.Groups().Get(ref.UUID) == nil {
			missing(ref, nil)
		}
	}
	for _, ref := range d.Labels {
		if sa.Labels().Get(ref.UUID) == nil {
			missing(ref, nil)
		}
	}
	for _, ref := range d.Templates {
		if sa.Templates().Get(ref.UUID) == nil {
			missing(ref, nil)
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
