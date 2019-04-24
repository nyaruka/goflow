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

// TemplateIncluder is interface passed to EnumerateTemplates to include templates on flow entities
type TemplateIncluder interface {
	String(*string)
	Slice([]string)
	Map(map[string]string)
}

type templateEnumerator struct {
	include func(string)
}

// NewTemplateEnumerator creates a template includer for enumerating templates
func NewTemplateEnumerator(include func(string)) TemplateIncluder {
	return &templateEnumerator{include: include}
}

func (t *templateEnumerator) String(s *string) {
	t.include(*s)
}

func (t *templateEnumerator) Slice(a []string) {
	for s := range a {
		t.include(a[s])
	}
}

func (t *templateEnumerator) Map(m map[string]string) {
	for k := range m {
		t.include(m[k])
	}
}

type templateRewriter struct {
	rewrite func(string) string
}

// NewTemplateRewriter creates a template includer for rewriting templates
func NewTemplateRewriter(rewrite func(string) string) TemplateIncluder {
	return &templateRewriter{rewrite: rewrite}
}

func (t *templateRewriter) String(s *string) {
	*s = t.rewrite(*s)
}

func (t *templateRewriter) Slice(a []string) {
	for s := range a {
		a[s] = t.rewrite(a[s])
	}
}

func (t *templateRewriter) Map(m map[string]string) {
	for k := range m {
		m[k] = t.rewrite(m[k])
	}
}

// Inspectable is implemented by various flow components to allow walking the definition and extracting things like dependencies
type Inspectable interface {
	Inspect(func(Inspectable))
	EnumerateTemplates(Localization, TemplateIncluder)
	EnumerateDependencies(Localization, func(assets.Reference))
	EnumerateResults(func(*ResultSpec))
}
