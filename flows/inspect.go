package flows

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// FlowInfo contains the results of flow inspection
type FlowInfo struct {
	Dependencies []Dependency  `json:"dependencies"`
	Results      []*ResultInfo `json:"results"`
	WaitingExits []ExitUUID    `json:"waiting_exits"`
	ParentRefs   []string      `json:"parent_refs"`
}

type DependencyInfo struct {
	Type    string `json:"type"`
	Missing bool   `json:"missing,omitempty"`
}

type Dependency struct {
	Reference assets.Reference
	Info      DependencyInfo
}

func (d Dependency) MarshalJSON() ([]byte, error) {
	return utils.JSONMarshalMerged(d.Reference, d.Info)
}

// NewDependencies inspects a list of references. If a session assets is provided,
// each dependency is checked to see if it is available or missing.
func NewDependencies(refs []assets.Reference, sa SessionAssets) []Dependency {
	deps := make([]Dependency, len(refs))

	for i, ref := range refs {
		missing := false
		if sa != nil {
			missing = !checkDependency(sa, ref)
		}

		deps[i] = Dependency{
			Reference: ref,
			Info:      DependencyInfo{Type: referenceTypeName(ref), Missing: missing},
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

// derives a dependency type name (e.g. group) from a reference
func referenceTypeName(ref assets.Reference) string {
	t := reflect.TypeOf(ref).String()
	t = strings.Split(t, ".")[1]
	if strings.HasSuffix(t, "Reference") {
		t = t[0 : len(t)-9]
	}
	return strings.ToLower(t)
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
