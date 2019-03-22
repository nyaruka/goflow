package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// ResultSpec is possible result that a flow might generate
type ResultSpec struct {
	Key        string   `json:"key"`
	Name       string   `json:"name"`
	Categories []string `json:"categories,omitempty"`
}

// NewResultSpec creates a new result spec
func NewResultSpec(name string, categories []string) *ResultSpec {
	return &ResultSpec{
		Key:        utils.Snakify(name),
		Name:       name,
		Categories: categories,
	}
}

func (r *ResultSpec) String() string {
	return fmt.Sprintf("key=%s|name=%s|categories=%s", r.Key, r.Name, strings.Join(r.Categories, ","))
}

// MergeResultSpecs merges result specs based on key
func MergeResultSpecs(specs []*ResultSpec) []*ResultSpec {
	merged := make([]*ResultSpec, 0, len(specs))
	byKey := make(map[string]*ResultSpec)

	for _, spec := range specs {
		existing := byKey[spec.Key]

		if existing != nil {
			// if we already have a result spec with this key, merge categories
			for _, category := range spec.Categories {
				if !utils.StringSliceContains(existing.Categories, category, false) {
					existing.Categories = append(existing.Categories, category)
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

// Inspectable is implemented by various flow components to allow walking the definition and extracting things like dependencies
type Inspectable interface {
	Inspect(func(Inspectable))
	EnumerateTemplates(Localization, func(string))
	RewriteTemplates(Localization, func(string) string)
	EnumerateDependencies(Localization, func(assets.Reference))
	EnumerateResults(func(*ResultSpec))
}
