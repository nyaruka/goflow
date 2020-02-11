package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
)

// ExtractedReference is a reference and its location in a flow
type ExtractedReference struct {
	Node      Node
	Action    Action
	Router    Router
	Reference assets.Reference
}

// FlowInfo contains the results of flow inspection
type FlowInfo struct {
	Dependencies []*Dependency `json:"dependencies"`
	Results      []*ResultInfo `json:"results"`
	WaitingExits []ExitUUID    `json:"waiting_exits"`
	ParentRefs   []string      `json:"parent_refs"`
}

type nodeLocations map[NodeUUID][]string

func (l nodeLocations) add(n Node, a Action, r Router) {
	var loc string
	if a != nil {
		loc = string(a.UUID())
	} else {
		loc = "router"
	}

	locs := l[n.UUID()]
	if !utils.StringSliceContains(locs, loc, true) {
		locs = append(locs, loc)
	}
	l[n.UUID()] = locs
}

type Dependency struct {
	Reference assets.Reference `json:"-"`
	Type      string           `json:"type"`
	Missing   bool             `json:"missing,omitempty"`
	Nodes     nodeLocations    `json:"nodes"`
}

func (d Dependency) MarshalJSON() ([]byte, error) {
	type dependency Dependency // need to alias type to avoid circular calls to this method
	return jsonx.MarshalMerged(d.Reference, dependency(d))
}

// NewDependencies inspects a list of references. If a session assets is provided,
// each dependency is checked to see if it is available or missing.
func NewDependencies(refs []ExtractedReference, sa SessionAssets) []*Dependency {
	deps := make([]*Dependency, 0)
	depsSeen := make(map[string]*Dependency, 0)

	for _, er := range refs {
		key := fmt.Sprintf("%s:%s", er.Reference.Type(), er.Reference.Identity())

		// if we already created a dependency for this reference, add this location
		if dep, seen := depsSeen[key]; seen {
			dep.Nodes.add(er.Node, er.Action, er.Router)
		} else {
			// check if this dependency is accessible
			missing := false
			if sa != nil {
				missing = !checkDependency(sa, er.Reference)
			}

			nodes := nodeLocations{}
			nodes.add(er.Node, er.Action, er.Router)

			dep := &Dependency{
				Reference: er.Reference,
				Type:      er.Reference.Type(),
				Missing:   missing,
				Nodes:     nodes,
			}
			deps = append(deps, dep)
			depsSeen[key] = dep
		}
	}

	return deps
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
