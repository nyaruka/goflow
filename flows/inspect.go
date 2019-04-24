package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/tools"
	"github.com/nyaruka/goflow/utils"
)

var fieldRefPaths = [][]string{
	{"fields"},
	{"contact", "fields"},
	{"parent", "fields"},
	{"parent", "contact", "fields"},
	{"child", "fields"},
	{"child", "contact", "fields"},
}

// Inspectable is implemented by various flow components to allow walking the definition and extracting things like dependencies
type Inspectable interface {
	Inspect(func(Inspectable))
	EnumerateTemplates(TemplateIncluder)
	EnumerateDependencies(Localization, func(assets.Reference))
	EnumerateResults(func(*ResultSpec))
}

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
	Translations(Localizable, string)
}

type templateEnumerator struct {
	localization Localization
	include      func(string)
}

// NewTemplateEnumerator creates a template includer for enumerating templates
func NewTemplateEnumerator(localization Localization, include func(string)) TemplateIncluder {
	return &templateEnumerator{localization: localization, include: include}
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

func (t *templateEnumerator) Translations(localizable Localizable, key string) {
	for _, lang := range t.localization.Languages() {
		translations := t.localization.GetTranslations(lang)
		t.Slice(translations.GetTextArray(localizable.LocalizationUUID(), key))
	}
}

type templateRewriter struct {
	localization Localization
	rewrite      func(string) string
}

// NewTemplateRewriter creates a template includer for rewriting templates
func NewTemplateRewriter(localization Localization, rewrite func(string) string) TemplateIncluder {
	return &templateRewriter{localization: localization, rewrite: rewrite}
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

func (t *templateRewriter) Translations(localizable Localizable, key string) {
	for _, lang := range t.localization.Languages() {
		translations := t.localization.GetTranslations(lang)
		t.Slice(translations.GetTextArray(localizable.LocalizationUUID(), key))
	}
}

// wrapper for an asset reference to make it inspectable
type inspectableReference struct {
	ref assets.Reference
}

// InspectReference inspects the given asset reference if it's non-nil
func InspectReference(ref assets.Reference, inspect func(Inspectable)) {
	if ref != nil {
		inspectableReference{ref: ref}.Inspect(inspect)
	}
}

// Inspect inspects this object and any children
func (r inspectableReference) Inspect(inspect func(Inspectable)) {
	inspect(r)
}

// EnumerateTemplates enumerates all expressions on this object and its children
func (r inspectableReference) EnumerateTemplates(include TemplateIncluder) {
	if r.ref != nil && r.ref.Variable() {
		switch typed := r.ref.(type) {
		case *assets.GroupReference:
			include.String(&typed.NameMatch)
		case *assets.LabelReference:
			include.String(&typed.NameMatch)
		}
	}
}

// EnumerateDependencies enumerates all dependencies on this object and its children
func (r inspectableReference) EnumerateDependencies(localization Localization, include func(assets.Reference)) {
	if r.ref != nil && !r.ref.Variable() {
		include(r.ref)
	}
}

// EnumerateResults enumerates all potential results on this object
// Asset references can't contain results.
func (r inspectableReference) EnumerateResults(include func(*ResultSpec)) {}

// ExtractFieldReferences extracts fields references from the given template
func ExtractFieldReferences(template string) []*assets.FieldReference {
	fieldRefs := make([]*assets.FieldReference, 0)
	tools.FindContextRefsInTemplate(template, RunContextTopLevels, func(path []string) {
		isField, fieldKey := isFieldRefPath(path)
		if isField {
			fieldRefs = append(fieldRefs, assets.NewFieldReference(fieldKey, ""))
		}
	})
	return fieldRefs
}

// checks whether the given context path is a reference to a contact field
func isFieldRefPath(path []string) (bool, string) {
	for _, possible := range fieldRefPaths {
		if len(path) == len(possible)+1 {
			matches := true
			for i := range possible {
				if strings.ToLower(path[i]) != possible[i] {
					matches = false
					break
				}
			}
			if matches {
				return true, strings.ToLower(path[len(possible)])
			}
		}
	}
	return false, ""
}
