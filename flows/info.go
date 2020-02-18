package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

type baseExtractedItem struct {
	Node     Node
	Action   Action
	Router   Router
	Language envs.Language
}

// ExtractedTemplate is a template and its location in a flow
type ExtractedTemplate struct {
	baseExtractedItem

	Template string
}

// NewExtractedTemplate creates a new extracted template
func NewExtractedTemplate(n Node, a Action, r Router, l envs.Language, t string) ExtractedTemplate {
	return ExtractedTemplate{
		baseExtractedItem: baseExtractedItem{Node: n, Action: a, Router: r, Language: l},
		Template:          t,
	}
}

// ExtractedReference is a reference and its location in a flow
type ExtractedReference struct {
	baseExtractedItem

	Reference assets.Reference
}

// NewExtractedReference creates a new extracted reference
func NewExtractedReference(n Node, a Action, r Router, l envs.Language, ref assets.Reference) ExtractedReference {
	return ExtractedReference{
		baseExtractedItem: baseExtractedItem{Node: n, Action: a, Router: r, Language: l},
		Reference:         ref,
	}
}

// Inspection contains the results of flow inspection
type Inspection struct {
	Dependencies []Dependency  `json:"dependencies"`
	Issues       []Issue       `json:"issues"`
	Results      []*ResultSpec `json:"results"`
	WaitingExits []ExitUUID    `json:"waiting_exits"`
	ParentRefs   []string      `json:"parent_refs"`
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
