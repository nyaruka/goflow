package flows

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/tools"
	"github.com/nyaruka/goflow/utils"
)

type FlowInfo struct {
	Dependencies *Dependencies `json:"_dependencies"`
	Results      []*ResultInfo `json:"_results"`
	WaitingExits []ExitUUID    `json:"_waiting_exits"`
}

type Dependencies struct {
	Channels  []*assets.ChannelReference  `json:"channels,omitempty"`
	Contacts  []*ContactReference         `json:"contacts,omitempty"`
	Fields    []*assets.FieldReference    `json:"fields,omitempty"`
	Flows     []*assets.FlowReference     `json:"flows,omitempty"`
	Groups    []*assets.GroupReference    `json:"groups,omitempty"`
	Labels    []*assets.LabelReference    `json:"labels,omitempty"`
	Templates []*assets.TemplateReference `json:"templates,omitempty"`
}

func NewDependencies(refs []assets.Reference) *Dependencies {
	d := &Dependencies{}
	for _, r := range refs {
		switch typed := r.(type) {
		case *assets.ChannelReference:
			d.Channels = append(d.Channels, typed)
		case *ContactReference:
			d.Contacts = append(d.Contacts, typed)
		case *assets.FieldReference:
			d.Fields = append(d.Fields, typed)
		case *assets.FlowReference:
			d.Flows = append(d.Flows, typed)
		case *assets.GroupReference:
			d.Groups = append(d.Groups, typed)
		case *assets.LabelReference:
			d.Labels = append(d.Labels, typed)
		case *assets.TemplateReference:
			d.Templates = append(d.Templates, typed)
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
	EnumerateResults(func(*ResultInfo))
}

// ResultInfo is possible result that a flow might generate
type ResultInfo struct {
	Key        string   `json:"key"`
	Name       string   `json:"name"`
	Categories []string `json:"categories,omitempty"`
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

// MergeResultInfos merges result specs based on key
func MergeResultInfos(specs []*ResultInfo) []*ResultInfo {
	merged := make([]*ResultInfo, 0, len(specs))
	byKey := make(map[string]*ResultInfo)

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
	for i := range a {
		t.include(a[i])
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
	for i := range a {
		a[i] = t.rewrite(a[i])
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
func (r inspectableReference) EnumerateResults(include func(*ResultInfo)) {}

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
