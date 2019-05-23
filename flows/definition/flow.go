package definition

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

// CurrentSpecVersion is the flow spec version supported by this library
var CurrentSpecVersion = semver.MustParse("13.0")

type flow struct {
	// spec properties
	uuid               assets.FlowUUID
	name               string
	specVersion        *semver.Version
	language           utils.Language
	flowType           flows.FlowType
	revision           int
	expireAfterMinutes int
	localization       flows.Localization
	nodes              []flows.Node

	// optional properties not used by engine itself
	ui json.RawMessage

	// properties set after inspection
	dependencies *dependencies
	results      []*flows.ResultSpec
	waitingExits []flows.ExitUUID

	// internal state
	nodeMap map[flows.NodeUUID]flows.Node
}

// NewFlow creates a new flow
func NewFlow(uuid assets.FlowUUID, name string, language utils.Language, flowType flows.FlowType, revision int, expireAfterMinutes int, localization flows.Localization, nodes []flows.Node, ui json.RawMessage) (flows.Flow, error) {
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
func (f *flow) Language() utils.Language               { return f.language }
func (f *flow) Type() flows.FlowType                   { return f.flowType }
func (f *flow) ExpireAfterMinutes() int                { return f.expireAfterMinutes }
func (f *flow) Nodes() []flows.Node                    { return f.nodes }
func (f *flow) Localization() flows.Localization       { return f.localization }
func (f *flow) UI() json.RawMessage                    { return f.ui }
func (f *flow) GetNode(uuid flows.NodeUUID) flows.Node { return f.nodeMap[uuid] }

func (f *flow) validate() error {
	// track UUIDs used by nodes and actions to ensure that they are unique
	seenUUIDs := make(map[utils.UUID]bool)

	for _, node := range f.nodes {
		uuidAlreadySeen := seenUUIDs[utils.UUID(node.UUID())]
		if uuidAlreadySeen {
			return errors.Errorf("node UUID %s isn't unique", node.UUID())
		}
		seenUUIDs[utils.UUID(node.UUID())] = true

		if err := node.Validate(f, seenUUIDs); err != nil {
			return errors.Wrapf(err, "invalid node[uuid=%s]", node.UUID())
		}
	}

	return nil
}

// Inspect enumerates dependencies, checks that they exist
func (f *flow) Inspect(sa flows.SessionAssets) error {
	return f.doInspect(sa, false, nil, nil)
}

// InspectRecursively checks that all of this flow's dependencies exist, and all our flow dependencies are also valid
func (f *flow) InspectRecursively(sa flows.SessionAssets, missing func(assets.Reference)) error {
	seen := make(map[assets.FlowUUID]bool)

	return f.doInspect(sa, true, seen, missing)
}

type brokenDependency struct {
	dependency assets.Reference
	reason     error
}

func (f *flow) doInspect(sa flows.SessionAssets, recursive bool, seen map[assets.FlowUUID]bool, missing func(assets.Reference)) error {
	// prevent looping if recursively inspecting flows
	if recursive && seen[f.UUID()] {
		return nil
	}

	// extract all dependencies (assets, contacts)
	deps := newDependencies(f.ExtractDependencies())

	// if we have session assets, check dependencies actually exist
	if sa != nil {
		missingAssets := make([]brokenDependency, 0)
		err := deps.refresh(sa, func(r assets.Reference, err error) {
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
	}

	f.dependencies = deps
	f.results = f.ExtractResults()
	f.waitingExits = f.ExtractExitsFromWaits()

	if recursive {
		seen[f.UUID()] = true

		// go check any non-missing flow dependencies
		for _, flowRef := range deps.Flows {
			flowDep, err := sa.Flows().Get(flowRef.UUID)
			if err == nil {
				if err := flowDep.(*flow).doInspect(sa, true, seen, missing); err != nil {
					return errors.Wrapf(err, "invalid child %s", flowRef)
				}
			}
		}
	}

	return nil
}

// Context returns the properties available in expressions
func (f *flow) Context(env utils.Environment) map[string]types.XValue {
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

func (f *flow) inspect(inspect func(flows.Inspectable)) {
	// inspect each node
	for _, n := range f.Nodes() {
		n.Inspect(inspect)
	}
}

// ExtractTemplates extracts all non-empty templates
func (f *flow) ExtractTemplates() []string {
	templates := make([]string, 0)
	include := flows.NewTemplateEnumerator(f.Localization(), func(template string) {
		if template != "" {
			templates = append(templates, template)
		}
	})

	f.inspect(func(item flows.Inspectable) {
		item.EnumerateTemplates(include)
	})
	return templates
}

// RewriteTemplates rewrites all templates
func (f *flow) RewriteTemplates(rewrite func(string) string) {
	f.inspect(func(item flows.Inspectable) {
		item.EnumerateTemplates(flows.NewTemplateRewriter(f.Localization(), rewrite))
	})
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
		}
	}

	include := flows.NewTemplateEnumerator(f.Localization(), func(template string) {
		fieldRefs := flows.ExtractFieldReferences(template)
		for _, f := range fieldRefs {
			addDependency(f)
		}
	})

	f.inspect(func(item flows.Inspectable) {
		item.EnumerateTemplates(include)
		item.EnumerateDependencies(f.Localization(), func(r assets.Reference) {
			addDependency(r)
		})
	})

	return dependencies
}

// ExtractResults extracts all result specs
func (f *flow) ExtractResults() []*flows.ResultSpec {
	specs := make([]*flows.ResultSpec, 0)
	f.inspect(func(item flows.Inspectable) {
		item.EnumerateResults(func(spec *flows.ResultSpec) {
			specs = append(specs, spec)
		})
	})
	return flows.MergeResultSpecs(specs)
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

// Generic returns this flow's data modelled as a hierarchy of generic maps and slices
func (f *flow) Generic() map[string]interface{} {
	marshaled, err := json.Marshal(f)
	if err != nil {
		panic(err.Error())
	}

	g, err := utils.JSONDecodeGeneric(marshaled)
	if err != nil {
		panic(err.Error())
	}

	return g.(map[string]interface{})
}

// Clone clones this flow replacing all UUIDs using the provided mapping and generating new
// random UUIDs if they aren't in the mapping
func (f *flow) Clone(depMapping map[utils.UUID]utils.UUID) flows.Flow {
	generic := f.Generic()
	remapUUIDs(generic, depMapping)

	// read back as a real flow in the current spec.. since we control this it can't error in theory
	return MustReadFlowFromGeneric(generic)
}

var _ flows.Flow = (*flow)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func init() {
	utils.Validator.RegisterAlias("flow_type", "eq=messaging|eq=messaging_offline|eq=voice")
}

// the set of fields common to all new flow spec versions
type flowHeader struct {
	UUID        assets.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name        string          `json:"name"`
	SpecVersion *semver.Version `json:"spec_version" validate:"required"`
}

type flowEnvelope struct {
	flowHeader

	Language           utils.Language  `json:"language" validate:"required"`
	Type               flows.FlowType  `json:"type" validate:"required,flow_type"`
	Revision           int             `json:"revision"`
	ExpireAfterMinutes int             `json:"expire_after_minutes"`
	Localization       localization    `json:"localization"`
	Nodes              []*node         `json:"nodes"`
	UI                 json.RawMessage `json:"_ui,omitempty"`
}

// additional properties that an inspected flow can have
type inspectedFlowEnvelope struct {
	*flowEnvelope

	Dependencies *dependencies       `json:"_dependencies"`
	Results      []*flows.ResultSpec `json:"_results"`
	WaitingExits []flows.ExitUUID    `json:"_waiting_exits"`
}

// IsSpecVersionSupported determines if we can read the given flow version
func IsSpecVersionSupported(ver *semver.Version) bool {
	// major versions change flow schema
	// we currently have no support for migrations but that will change in future
	return ver.Major() == CurrentSpecVersion.Major()
}

// ReadFlow reads a single flow definition from the passed in byte array
func ReadFlow(data json.RawMessage) (flows.Flow, error) {
	header := &flowHeader{}
	if err := utils.UnmarshalAndValidate(data, header); err != nil {
		return nil, errors.Wrap(err, "unable to read flow header")
	}

	// can't do anything with a newer major version than this library supports
	if !IsSpecVersionSupported(header.SpecVersion) {
		return nil, errors.Errorf("spec version %s is newer than this library (%s)", header.SpecVersion, CurrentSpecVersion)
	}

	// TODO flow spec migrations

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
		flowHeader: flowHeader{
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

	if f.dependencies != nil {
		return json.Marshal(&inspectedFlowEnvelope{
			flowEnvelope: e,
			Dependencies: f.dependencies,
			Results:      f.results,
			WaitingExits: f.waitingExits,
		})
	}

	return json.Marshal(e)
}
