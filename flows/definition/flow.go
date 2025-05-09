package definition

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/flows/inspect/issues"
	"github.com/nyaruka/goflow/utils"
)

// CurrentSpecVersion is the flow spec version supported by this library
var CurrentSpecVersion = semver.MustParse("14.3.0")
var CurrentSupportedSpecs, _ = semver.NewConstraint(">= 11.0.0, < 14.4.0")

// IsVersionSupported checks the given version is supported
func IsVersionSupported(v *semver.Version) bool {
	// can't do anything with a pre-11 flow or a newer minor version
	return CurrentSupportedSpecs.Check(v)
}

type flow struct {
	// spec properties
	uuid         assets.FlowUUID
	name         string
	specVersion  *semver.Version
	language     i18n.Language
	flowType     flows.FlowType
	revision     int
	expireAfter  time.Duration
	localization flows.Localization
	nodes        []flows.Node

	// optional properties not used by engine itself
	ui json.RawMessage

	// flows can be loaded from a assets or just a definition
	asset assets.Flow

	// internal state
	nodeMap map[flows.NodeUUID]flows.Node
}

// NewFlow creates a new flow
func NewFlow(uuid assets.FlowUUID, name string, language i18n.Language, flowType flows.FlowType, revision int, expireAfter time.Duration, localization flows.Localization, nodes []flows.Node, ui json.RawMessage, a assets.Flow) (flows.Flow, error) {
	f := &flow{
		uuid:         uuid,
		name:         name,
		specVersion:  CurrentSpecVersion,
		language:     language,
		flowType:     flowType,
		revision:     revision,
		expireAfter:  expireAfter,
		localization: localization,
		nodes:        nodes,
		nodeMap:      make(map[flows.NodeUUID]flows.Node, len(nodes)),
		ui:           ui,
		asset:        a,
	}

	for _, node := range f.nodes {
		f.nodeMap[node.UUID()] = node
	}

	if err := f.validate(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *flow) UUID() assets.FlowUUID                  { return f.uuid }
func (f *flow) Name() string                           { return f.name }
func (f *flow) SpecVersion() *semver.Version           { return f.specVersion }
func (f *flow) Revision() int                          { return f.revision }
func (f *flow) Language() i18n.Language                { return f.language }
func (f *flow) Type() flows.FlowType                   { return f.flowType }
func (f *flow) Nodes() []flows.Node                    { return f.nodes }
func (f *flow) Localization() flows.Localization       { return f.localization }
func (f *flow) UI() json.RawMessage                    { return f.ui }
func (f *flow) GetNode(uuid flows.NodeUUID) flows.Node { return f.nodeMap[uuid] }

func (f *flow) ExpireAfter() time.Duration {
	if f.expireAfter == 0 {
		if f.flowType == flows.FlowTypeMessaging {
			return 3 * 24 * time.Hour // 3 days
		} else if f.flowType == flows.FlowTypeVoice {
			return 5 * time.Minute
		}
	}
	return f.expireAfter
}

func (f *flow) validate() error {
	if len(f.nodes) > flows.MaxNodesPerFlow {
		return fmt.Errorf("flow can't have more than %d nodes (has %d)", flows.MaxNodesPerFlow, len(f.nodes))
	}

	// track UUIDs used by nodes and actions to ensure that they are unique
	seenUUIDs := make(map[uuids.UUID]bool)

	for _, node := range f.nodes {
		uuidAlreadySeen := seenUUIDs[uuids.UUID(node.UUID())]
		if uuidAlreadySeen {
			return fmt.Errorf("node UUID %s isn't unique", node.UUID())
		}
		seenUUIDs[uuids.UUID(node.UUID())] = true

		if err := node.Validate(f, seenUUIDs); err != nil {
			return fmt.Errorf("invalid node[uuid=%s]: %w", node.UUID(), err)
		}
	}

	if err := f.localization.Validate(); err != nil {
		return fmt.Errorf("invalid localization: %w", err)
	}

	return nil
}

// Context returns the properties available in expressions
//
//	__default__:text -> the name
//	uuid:text -> the UUID of the flow
//	name:text -> the name of the flow
//	revision:text -> the revision number of the flow
//
// @context flow
func (f *flow) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"__default__": types.NewXText(f.name),
		"uuid":        types.NewXText(string(f.UUID())),
		"name":        types.NewXText(f.name),
		"revision":    types.NewXNumberFromInt(f.revision),
	}
}

// Reference returns a reference to this flow asset
func (f *flow) Reference(withRevision bool) *assets.FlowReference {
	if f == nil {
		return nil
	}
	if withRevision {
		return assets.NewFlowReferenceWithRevision(f.uuid, f.name, f.revision)
	}
	return assets.NewFlowReference(f.uuid, f.name)
}

// Inspect enumerates dependencies, results etc
func (f *flow) Inspect(sa flows.SessionAssets) *flows.Info {
	assetRefs := make([]flows.ExtractedReference, 0)
	locals := make(map[string]bool)
	results := make([]flows.ExtractedResult, 0)
	templates := make([]flows.ExtractedTemplate, 0)
	parentRefs := make(map[string]bool)

	recordAssetRef := func(n flows.Node, a flows.Action, r flows.Router, l i18n.Language, ref assets.Reference) {
		if ref != nil && !ref.Variable() {
			er := flows.NewExtractedReference(n, a, r, l, ref)
			assetRefs = append(assetRefs, er)
		}
	}

	for _, n := range f.nodes {
		n.Inspect(func(a flows.Action, r flows.Router, ref assets.Reference) {
			recordAssetRef(n, a, r, i18n.NilLanguage, ref)
		}, func(a flows.Action, r flows.Router, l string) {
			if l != "" {
				locals[l] = true
			}
		}, func(a flows.Action, r flows.Router, i *flows.ResultInfo) {
			results = append(results, flows.ExtractedResult{Node: n, Action: a, Router: r, Info: i})
		})

		n.EnumerateTemplates(f.Localization(), func(a flows.Action, r flows.Router, l i18n.Language, t string) {
			templates = append(templates, flows.NewExtractedTemplate(n, a, r, l, t))
			ars, prs := inspect.ExtractFromTemplate(t)
			for _, ref := range ars {
				recordAssetRef(n, a, r, l, ref)
			}
			for _, r := range prs {
				parentRefs[r] = true
			}
		})
	}

	return &flows.Info{
		Counts:       map[string]int{"nodes": len(f.nodes), "languages": len(f.localization.Languages())},
		Dependencies: inspect.NewDependencies(assetRefs, sa),
		Locals:       utils.EnsureNonNil(slices.Sorted(maps.Keys(locals))),
		Results:      flows.NewResultSpecs(results),
		ParentRefs:   utils.EnsureNonNil(slices.Sorted(maps.Keys(parentRefs))),
		Issues:       issues.Check(sa, f, templates, assetRefs),
	}
}

// ExtractTemplates extracts all non-empty templates
func (f *flow) ExtractTemplates() []string {
	templates := make([]string, 0)
	include := func(a flows.Action, r flows.Router, l i18n.Language, t string) {
		if t != "" {
			templates = append(templates, t)
		}
	}

	for _, n := range f.nodes {
		n.EnumerateTemplates(f.Localization(), include)
	}

	return templates
}

// ExtractLocalizables extracts all localizable text
func (f *flow) ExtractLocalizables() []string {
	texts := make([]string, 0)
	include := func(uuid uuids.UUID, property string, ts []string, w func([]string)) {
		for _, t := range ts {
			if t != "" {
				texts = append(texts, t)
			}
		}
	}

	for _, n := range f.nodes {
		n.EnumerateLocalizables(include)
	}

	return texts
}

// ChangeLanguage changes the language of the flow saving the current flow text as a translation and replacing it with
// the specified translation. It returns an error if there are missing translations.
func (f *flow) ChangeLanguage(lang i18n.Language) (flows.Flow, error) {
	// make a copy of the flow
	copy, err := f.copy()
	if err != nil {
		return nil, err
	}

	ll := copy.localization.(localization)

	outTranslation := make(languageTranslation) // current flow text extracted out as a translation
	inTranslation := ll[lang]                   // translation being imported in as new flow text
	if inTranslation == nil {
		inTranslation = languageTranslation{}
	}

	include := func(uuid uuids.UUID, property string, oldValues []string, w func([]string)) {
		// save current flow text into a translation
		if len(oldValues) > 0 {
			outTranslation.setTextArray(uuid, property, oldValues)
		}

		newValues := inTranslation.getTextArray(uuid, property)

		// if we have a translation, update flow text in the definition, if not then leave it as is
		if len(newValues) > 0 {
			w(newValues)
		}
	}

	for _, n := range copy.nodes {
		n.EnumerateLocalizables(include)
	}

	ll[copy.language] = outTranslation
	delete(ll, lang)
	copy.language = lang

	return copy, nil
}

// makes a copy of this flow which this differs from cloning as UUIDs are preserved
func (f *flow) copy() (*flow, error) {
	// by marshaling and unmarshaling...
	marshaled, err := jsonx.Marshal(f)
	if err != nil {
		return nil, err
	}
	cp, err := ReadFlow(marshaled, nil)
	if err != nil {
		return nil, err
	}
	return cp.(*flow), nil
}

// Asset returns the asset from which this flow was created if there was one
func (f *flow) Asset() assets.Flow { return f.asset }

var _ flows.Flow = (*flow)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type flowEnvelope struct {
	migrations.Header13

	Language           i18n.Language   `json:"language"             validate:"required,language"`
	Type               flows.FlowType  `json:"type"                 validate:"required,flow_type"`
	Revision           int             `json:"revision"`
	ExpireAfterMinutes int             `json:"expire_after_minutes"`
	Localization       localization    `json:"localization"`
	Nodes              []*node         `json:"nodes"                validate:"dive,required"`
	UI                 json.RawMessage `json:"_ui,omitempty"`
}

// ReadFlow reads a flow definition from the passed in byte array, migrating it to the spec version of the engine if necessary
func ReadFlow(data []byte, mc *migrations.Config) (flows.Flow, error) {
	return readFlow(data, mc, nil)
}

// ReadAsset reads a flow definition from the passed in flow asset, migrating it to the spec version of the engine if necessary
func ReadAsset(a assets.Flow, mc *migrations.Config) (flows.Flow, error) {
	return readFlow(a.Definition(), mc, a)
}

func readFlow(data []byte, mc *migrations.Config, a assets.Flow) (flows.Flow, error) {
	if mc == nil {
		mc = migrations.DefaultConfig
	}

	var err error
	data, err = migrations.MigrateToLatest(data, mc)
	if err != nil {
		return nil, err
	}

	header := &migrations.Header13{}
	jsonx.Unmarshal(data, header)

	if !IsVersionSupported(header.SpecVersion) {
		return nil, fmt.Errorf("spec version %s is newer than this library (%s)", header.SpecVersion, CurrentSpecVersion)
	}

	e := &flowEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	nodes := make([]flows.Node, len(e.Nodes))
	for i := range e.Nodes {
		nodes[i] = e.Nodes[i]
	}

	if e.Localization == nil {
		e.Localization = make(localization)
	}

	return NewFlow(e.UUID, e.Name, e.Language, e.Type, e.Revision, time.Duration(e.ExpireAfterMinutes)*time.Minute, e.Localization, nodes, e.UI, a)
}

// MarshalJSON marshals this flow into JSON
func (f *flow) MarshalJSON() ([]byte, error) {
	e := &flowEnvelope{
		Header13: migrations.Header13{
			UUID:        f.uuid,
			Name:        f.name,
			SpecVersion: f.specVersion,
		},
		Language:           f.language,
		Type:               f.flowType,
		Revision:           f.revision,
		ExpireAfterMinutes: int(f.expireAfter / time.Minute),
		Localization:       f.localization.(localization),
		Nodes:              make([]*node, len(f.nodes)),
		UI:                 f.ui,
	}

	for i := range f.nodes {
		e.Nodes[i] = f.nodes[i].(*node)
	}

	return jsonx.Marshal(e)
}
