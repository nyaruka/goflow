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
var CurrentSpecVersion = semver.MustParse("12.0")

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

	// properties set after validation
	dependencies *dependencies
	results      []*flows.ResultSpec
	waitingExits []flows.ExitUUID

	// internal state
	nodeMap   map[flows.NodeUUID]flows.Node
	validated bool
}

// NewFlow creates a new flow
func NewFlow(uuid assets.FlowUUID, name string, language utils.Language, flowType flows.FlowType, revision int, expireAfterMinutes int, localization flows.Localization, nodes []flows.Node, ui json.RawMessage) flows.Flow {
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

	return f
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

// Validates that we are structurally currect. The SessionAssets `sa` is optional but if provided,
// we will also check that all dependencies actually exist, and refresh their names.
func (f *flow) Validate(sa flows.SessionAssets) error {
	return f.validate(sa, false, nil)
}

// Validates that we are structurally currect, have all the dependencies we need, and all our flow dependencies are also valid
func (f *flow) ValidateRecursively(sa flows.SessionAssets, missing func(assets.Reference)) error {
	return f.validate(sa, true, missing)
}

func (f *flow) validate(sa flows.SessionAssets, recursive bool, missing func(assets.Reference)) error {
	// if this flow has already been validated, don't need to do it again - avoid unnecessary work
	// but also prevents looping if recursively validating flows
	if f.validated {
		return nil
	}

	// track UUIDs used by nodes and actions to ensure that they are unique
	seenUUIDs := make(map[utils.UUID]bool)

	for _, node := range f.nodes {
		uuidAlreadySeen := seenUUIDs[utils.UUID(node.UUID())]
		if uuidAlreadySeen {
			return errors.Errorf("node UUID %s isn't unique", node.UUID())
		}
		seenUUIDs[utils.UUID(node.UUID())] = true

		if err := node.Validate(f, seenUUIDs); err != nil {
			return errors.Wrapf(err, "validation failed for node[uuid=%s]", node.UUID())
		}
	}

	// extract all dependencies (assets, contacts)
	deps := newDependencies(f.ExtractDependencies())

	// and validate that all assets are available in the session assets
	if sa != nil {
		missingAssets := make([]assets.Reference, 0)
		deps.refresh(sa, func(r assets.Reference) { missingAssets = append(missingAssets, r) })

		if len(missingAssets) > 0 {
			// if we have callback for missing dependencies, call that
			if missing != nil {
				for _, dep := range missingAssets {
					missing(dep)
				}
			} else {
				// otherwise error
				depStrings := make([]string, len(missingAssets))
				for d := range missingAssets {
					depStrings[d] = missingAssets[d].String()
				}
				return errors.Errorf("missing dependencies: %s", strings.Join(depStrings, ","))
			}
		}
	}

	f.validated = true
	f.dependencies = deps
	f.results = f.ExtractResults()
	f.waitingExits = f.ExtractExitsFromWaits()

	if recursive {
		if sa == nil {
			panic("can't do recursive flow validation without session assets")
		}

		// go validate any non-missing flow dependencies
		for _, flowRef := range deps.Flows {
			flowDep, err := sa.Flows().Get(flowRef.UUID)
			if err == nil {
				flowDep.(*flow).validate(sa, true, missing)
			}
		}
	}

	return nil
}

// ToXValue returns a representation of this object for use in expressions
func (f *flow) ToXValue(env utils.Environment) types.XValue {
	return types.NewXDict(map[string]types.XValue{
		"uuid":     types.NewXText(string(f.UUID())),
		"name":     types.NewXText(f.name),
		"revision": types.NewXNumberFromInt(f.revision),
	})
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
	f.inspect(func(item flows.Inspectable) {
		item.EnumerateTemplates(f.Localization(), func(template string) {
			if template != "" {
				templates = append(templates, template)
			}
		})
	})
	return templates
}

// RewriteTemplates rewrites all templates
func (f *flow) RewriteTemplates(rewrite func(string) string) {
	f.inspect(func(item flows.Inspectable) {
		item.RewriteTemplates(f.Localization(), rewrite)
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

	f.inspect(func(item flows.Inspectable) {
		item.EnumerateTemplates(f.Localization(), func(template string) {
			fieldRefs := flows.ExtractFieldReferences(template)
			for _, f := range fieldRefs {
				addDependency(f)
			}
		})

		item.EnumerateDependencies(f.Localization(), func(r assets.Reference) {
			addDependency(r)
		})
	})

	return dependencies
}

// ExtractResultNames extracts all result names
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

// additional properties that a validated flow can have
type validatedFlowEnvelope struct {
	*flowEnvelope

	Dependencies *dependencies       `json:"_dependencies"`
	Results      []*flows.ResultSpec `json:"_results"`
	WaitingExits []flows.ExitUUID    `json:"_waiting_exits"`
}

// IsSpecVersionSupported determines if we can read the given flow version
func IsSpecVersionSupported(ver *semver.Version) bool {
	// major versions change flow schema
	return ver.Major() <= CurrentSpecVersion.Major()
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
	for n := range e.Nodes {
		nodes[n] = e.Nodes[n]
	}

	if e.Localization == nil {
		e.Localization = make(localization)
	}

	return NewFlow(e.UUID, e.Name, e.Language, e.Type, e.Revision, e.ExpireAfterMinutes, e.Localization, nodes, e.UI), nil
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

	if f.validated {
		return json.Marshal(&validatedFlowEnvelope{
			flowEnvelope: e,
			Dependencies: f.dependencies,
			Results:      f.results,
			WaitingExits: f.waitingExits,
		})
	}

	return json.Marshal(e)
}
