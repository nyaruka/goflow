package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/flows/inspect/issues"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

// CurrentSpecVersion is the flow spec version supported by this library
var CurrentSpecVersion = semver.MustParse("13.1.0")

// IsVersionSupported checks the given version is supported
func IsVersionSupported(v *semver.Version) bool {
	// can't do anything with a pre-11 flow or a newer major version
	return v.Major() >= 11 && v.Major() <= CurrentSpecVersion.Major()
}

type flow struct {
	// spec properties
	uuid               assets.FlowUUID
	name               string
	specVersion        *semver.Version
	language           envs.Language
	flowType           flows.FlowType
	revision           int
	expireAfterMinutes int
	localization       flows.Localization
	nodes              []flows.Node

	// optional properties not used by engine itself
	ui json.RawMessage

	// internal state
	nodeMap map[flows.NodeUUID]flows.Node
}

// NewFlow creates a new flow
func NewFlow(uuid assets.FlowUUID, name string, language envs.Language, flowType flows.FlowType, revision int, expireAfterMinutes int, localization flows.Localization, nodes []flows.Node, ui json.RawMessage) (flows.Flow, error) {
	f := &flow{
		uuid:               uuid,
		name:               name,
		specVersion:        CurrentSpecVersion,
		language:           language,
		flowType:           flowType,
		revision:           revision,
		expireAfterMinutes: expireAfterMinutes,
		localization:       localization,
		nodes:              nodes,
		nodeMap:            make(map[flows.NodeUUID]flows.Node, len(nodes)),
		ui:                 ui,
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
func (f *flow) Language() envs.Language                { return f.language }
func (f *flow) Type() flows.FlowType                   { return f.flowType }
func (f *flow) ExpireAfterMinutes() int                { return f.expireAfterMinutes }
func (f *flow) Nodes() []flows.Node                    { return f.nodes }
func (f *flow) Localization() flows.Localization       { return f.localization }
func (f *flow) UI() json.RawMessage                    { return f.ui }
func (f *flow) GetNode(uuid flows.NodeUUID) flows.Node { return f.nodeMap[uuid] }

func (f *flow) validate() error {
	// track UUIDs used by nodes and actions to ensure that they are unique
	seenUUIDs := make(map[uuids.UUID]bool)

	for _, node := range f.nodes {
		uuidAlreadySeen := seenUUIDs[uuids.UUID(node.UUID())]
		if uuidAlreadySeen {
			return errors.Errorf("node UUID %s isn't unique", node.UUID())
		}
		seenUUIDs[uuids.UUID(node.UUID())] = true

		if err := node.Validate(f, seenUUIDs); err != nil {
			return errors.Wrapf(err, "invalid node[uuid=%s]", node.UUID())
		}
	}

	return nil
}

// Inspect enumerates dependencies, results etc
func (f *flow) Inspect(sa flows.SessionAssets) *flows.FlowInfo {
	assetRefs, parentRefs := f.extractAssetAndParentRefs()

	return &flows.FlowInfo{
		Dependencies: flows.NewDependencies(assetRefs, sa),
		Results:      flows.NewResultSpecs(f.extractResults()),
		WaitingExits: f.extractExitsFromWaits(),
		ParentRefs:   parentRefs,
		Issues:       issues.Check(sa, f, assetRefs),
	}
}

// Context returns the properties available in expressions
//
//   __default__:text -> the name
//   uuid:text -> the UUID of the flow
//   name:text -> the name of the flow
//   revision:text -> the revision number of the flow
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
func (f *flow) Reference() *assets.FlowReference {
	if f == nil {
		return nil
	}
	return assets.NewFlowReference(f.uuid, f.name)
}

// ExtractTemplates extracts all non-empty templates
func (f *flow) ExtractTemplates() []string {
	templates := make([]string, 0)
	include := func(a flows.Action, r flows.Router, l envs.Language, t string) {
		if t != "" {
			templates = append(templates, t)
		}
	}

	for _, n := range f.nodes {
		n.EnumerateTemplates(f.Localization(), include)
	}

	return templates
}

// ExtractDependencies extracts all asset dependencies
func (f *flow) extractAssetAndParentRefs() ([]flows.ExtractedReference, []string) {
	assetRefs := make([]flows.ExtractedReference, 0)
	parentRefs := make(map[string]bool)

	recordDependency := func(n flows.Node, a flows.Action, r flows.Router, l envs.Language, ref assets.Reference) {
		if ref != nil && !ref.Variable() {
			er := flows.ExtractedReference{Node: n, Action: a, Router: r, Language: l, Reference: ref}
			assetRefs = append(assetRefs, er)
		}
	}

	for _, n := range f.nodes {
		n.EnumerateTemplates(f.Localization(), func(a flows.Action, r flows.Router, l envs.Language, t string) {
			ars, prs := inspect.ExtractFromTemplate(t)
			for _, ref := range ars {
				recordDependency(n, a, r, l, ref)
			}
			for _, r := range prs {
				parentRefs[r] = true
			}
		})
		n.EnumerateDependencies(f.Localization(), func(a flows.Action, r flows.Router, ref assets.Reference) {
			recordDependency(n, a, r, envs.NilLanguage, ref)
		})
	}

	return assetRefs, utils.StringSetKeys(parentRefs)
}

// extracts all result specs
func (f *flow) extractResults() []flows.ExtractedResult {
	results := make([]flows.ExtractedResult, 0)

	for _, n := range f.nodes {
		n.EnumerateResults(func(a flows.Action, r flows.Router, i *flows.ResultInfo) {
			results = append(results, flows.ExtractedResult{Node: n, Action: a, Router: r, Info: i})
		})
	}

	return results
}

// extracts all exits coming from nodes with waits
func (f *flow) extractExitsFromWaits() []flows.ExitUUID {
	exitUUIDs := make([]flows.ExitUUID, 0)
	include := func(e flows.ExitUUID) { exitUUIDs = append(exitUUIDs, e) }

	for _, n := range f.nodes {
		if n.Router() != nil && n.Router().Wait() != nil {
			for _, e := range n.Exits() {
				include(e.UUID())
			}
		}
	}
	return exitUUIDs
}

var _ flows.Flow = (*flow)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func init() {
	utils.Validator.RegisterAlias("flow_type", "eq=messaging|eq=messaging_offline|eq=voice")
}

type flowEnvelope struct {
	migrations.Header13

	Language           envs.Language   `json:"language" validate:"required"`
	Type               flows.FlowType  `json:"type" validate:"required,flow_type"`
	Revision           int             `json:"revision"`
	ExpireAfterMinutes int             `json:"expire_after_minutes"`
	Localization       localization    `json:"localization"`
	Nodes              []*node         `json:"nodes"`
	UI                 json.RawMessage `json:"_ui,omitempty"`
}

// ReadFlow a flow definition from the passed in byte array, migrating it to the spec version of the engine if necessary
func ReadFlow(data json.RawMessage, migrationConfig *migrations.Config) (flows.Flow, error) {
	var err error
	data, err = migrations.MigrateToLatest(data, migrationConfig)
	if err != nil {
		return nil, err
	}

	header := &migrations.Header13{}
	json.Unmarshal(data, header)

	if !IsVersionSupported(header.SpecVersion) {
		return nil, errors.Errorf("spec version %s is newer than this library (%s)", header.SpecVersion, CurrentSpecVersion)
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

	return NewFlow(e.UUID, e.Name, e.Language, e.Type, e.Revision, e.ExpireAfterMinutes, e.Localization, nodes, e.UI)
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
		ExpireAfterMinutes: f.expireAfterMinutes,
		Localization:       f.localization.(localization),
		Nodes:              make([]*node, len(f.nodes)),
		UI:                 f.ui,
	}

	for i := range f.nodes {
		e.Nodes[i] = f.nodes[i].(*node)
	}

	return json.Marshal(e)
}
