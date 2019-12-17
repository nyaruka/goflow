package definition

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

// CurrentSpecVersion is the flow spec version supported by this library
var CurrentSpecVersion = semver.MustParse("13.1.0")

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
func (f *flow) Inspect() *flows.FlowInfo {
	return &flows.FlowInfo{
		Dependencies: flows.NewDependencies(f.ExtractDependencies()),
		Results:      f.ExtractResults(),
		WaitingExits: f.ExtractExitsFromWaits(),
	}
}

// Validate checks that all of this flow's dependencies exist
func (f *flow) Validate(sa flows.SessionAssets, missing func(assets.Reference)) error {
	return f.validateAssets(sa, false, nil, missing)
}

// Validate checks that all of this flow's dependencies exist, and all our flow dependencies are also valid
func (f *flow) ValidateRecursive(sa flows.SessionAssets, missing func(assets.Reference)) error {
	seen := make(map[assets.FlowUUID]bool)

	return f.validateAssets(sa, true, seen, missing)
}

type brokenDependency struct {
	dependency assets.Reference
	reason     error
}

func (f *flow) validateAssets(sa flows.SessionAssets, recursive bool, seen map[assets.FlowUUID]bool, missing func(assets.Reference)) error {
	// prevent looping if recursive
	if recursive && seen[f.UUID()] {
		return nil
	}

	// extract all dependencies (assets, contacts)
	deps := flows.NewDependencies(f.ExtractDependencies())

	// check dependencies actually exist
	missingAssets := make([]brokenDependency, 0)
	err := deps.Check(sa, func(r assets.Reference, err error) {
		missingAssets = append(missingAssets, brokenDependency{r, err})
	})
	if err != nil {
		return err
	}

	if len(missingAssets) > 0 {
		// if we have callback for missing dependencies, call that
		if missing != nil {
			for _, ma := range missingAssets {
				missing(ma.dependency)
			}
		} else {
			// otherwise error
			depStrings := make([]string, len(missingAssets))
			for i, ma := range missingAssets {
				depStrings[i] = ma.dependency.String()
				if ma.reason != nil {
					depStrings[i] += fmt.Sprintf(" (%s)", ma.reason)
				}
			}
			return errors.Errorf("missing dependencies: %s", strings.Join(depStrings, ","))
		}
	}

	if recursive {
		seen[f.UUID()] = true

		// go check any non-missing flow dependencies
		for _, flowRef := range deps.Flows {
			flowDep, err := sa.Flows().Get(flowRef.UUID)
			if err == nil {
				if err := flowDep.(*flow).validateAssets(sa, true, seen, missing); err != nil {
					return errors.Wrapf(err, "invalid child %s", flowRef)
				}
			}
		}
	}

	return nil
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
	include := func(template string) {
		if template != "" {
			templates = append(templates, template)
		}
	}

	for _, n := range f.nodes {
		n.EnumerateTemplates(f.Localization(), include)
	}

	return templates
}

// ExtractDependencies extracts all asset dependencies
func (f *flow) ExtractDependencies() []assets.Reference {
	dependencies := make([]assets.Reference, 0)
	dependenciesSeen := make(map[string]bool)
	addDependency := func(r assets.Reference) {
		if !utils.IsNil(r) && !r.Variable() {
			key := fmt.Sprintf("%s:%s", r.Type(), r.Identity())
			if !dependenciesSeen[key] {
				dependencies = append(dependencies, r)
				dependenciesSeen[key] = true
			}

			// TODO replace if we saw a field ref without a name but now have same field with a name
		}
	}

	include := func(template string) {
		refs := inspect.ExtractFromTemplate(template)
		for _, f := range refs {
			addDependency(f)
		}
	}

	for _, n := range f.nodes {
		n.EnumerateTemplates(f.Localization(), include)
		n.EnumerateDependencies(f.Localization(), func(r assets.Reference) {
			addDependency(r)
		})
	}

	return dependencies
}

// ExtractResults extracts all result specs
func (f *flow) ExtractResults() []*flows.ResultInfo {
	specs := make([]*flows.ResultInfo, 0)

	for _, n := range f.nodes {
		n.EnumerateResults(n, func(spec *flows.ResultInfo) {
			specs = append(specs, spec)
		})
	}

	return flows.MergeResultInfos(specs)
}

// ExtractExitsFromWaits extracts all exits coming from nodes with waits
func (f *flow) ExtractExitsFromWaits() []flows.ExitUUID {
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

	// can't do anything with a newer major version than this library supports
	if header.SpecVersion.Major() > CurrentSpecVersion.Major() {
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

type flowWithInfoEnvelope struct {
	*flowEnvelope

	Dependencies *flows.Dependencies `json:"_dependencies"`
	Results      []*flows.ResultInfo `json:"_results"`
	WaitingExits []flows.ExitUUID    `json:"_waiting_exits"`
}

// MarshalWithInfo is temporary workaround to help us with the response format of the current /flow/validate endpoint
func (f *flow) MarshalWithInfo() ([]byte, error) {
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

	info := f.Inspect()

	return json.Marshal(&flowWithInfoEnvelope{
		flowEnvelope: e,
		Dependencies: info.Dependencies,
		Results:      info.Results,
		WaitingExits: info.WaitingExits,
	})
}
