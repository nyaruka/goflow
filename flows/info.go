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
	Results      []*ResultSpec `json:"results"`
	WaitingExits []ExitUUID    `json:"waiting_exits"`
	ParentRefs   []string      `json:"parent_refs"`
	Problems     []Problem     `json:"problems"`
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
	Key        string   `json:"key"`
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
}

// NewResultInfo creates a new result spec
func NewResultInfo(name string, categories []string) *ResultInfo {
	return &ResultInfo{
		Key:        utils.Snakify(name),
		Name:       name,
		Categories: categories,
	}
}

func (r *ResultInfo) String() string {
	return fmt.Sprintf("key=%s|name=%s|categories=%s", r.Key, r.Name, strings.Join(r.Categories, ","))
}

type ExtractedResult struct {
	Node   Node
	Action Action
	Router Router
	Info   *ResultInfo
}

type ResultSpec struct {
	ResultInfo
	NodeUUIDs []string `json:"node_uuids"`
}

// NewResultSpecs merges extracted results based on key
func NewResultSpecs(results []ExtractedResult) []*ResultSpec {
	specs := make([]*ResultSpec, 0)
	specsSeen := make(map[string]*ResultSpec)

	for _, result := range results {
		existing := specsSeen[result.Info.Key]
		nodeUUID := string(result.Node.UUID())

		// merge if we already have a result info with this key
		if existing != nil {
			// merge categories
			for _, category := range result.Info.Categories {
				if !utils.StringSliceContains(existing.Categories, category, false) {
					existing.Categories = append(existing.Categories, category)
				}
			}

			// merge this node UUID
			if !utils.StringSliceContains(existing.NodeUUIDs, nodeUUID, true) {
				existing.NodeUUIDs = append(existing.NodeUUIDs, nodeUUID)
			}
		} else {
			// if not, add as new unique result spec
			spec := &ResultSpec{
				ResultInfo: ResultInfo{
					Key:        result.Info.Key,
					Name:       result.Info.Name,
					Categories: result.Info.Categories,
				},
				NodeUUIDs: []string{nodeUUID},
			}
			specs = append(specs, spec)
			specsSeen[result.Info.Key] = spec
		}
	}
	return specs
}
