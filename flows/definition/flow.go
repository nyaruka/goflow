package definition

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type flow struct {
	uuid               assets.FlowUUID
	name               string
	specVersion        string
	language           utils.Language
	flowType           flows.FlowType
	revision           int
	expireAfterMinutes int
	localization       flows.Localization
	nodes              []flows.Node

	// only read for legacy flows which are being migrated
	ui flows.UI

	// internal state
	nodeMap map[flows.NodeUUID]flows.Node
	valid   bool
}

// NewFlow creates a new flow
func NewFlow(uuid assets.FlowUUID, name string, specVersion string, language utils.Language, flowType flows.FlowType, revision int, expireAfterMinutes int, localization flows.Localization, nodes []flows.Node, ui flows.UI) flows.Flow {
	f := &flow{
		uuid:               uuid,
		name:               name,
		specVersion:        specVersion,
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
func (f *flow) SpecVersion() string                    { return f.specVersion }
func (f *flow) Revision() int                          { return f.revision }
func (f *flow) Language() utils.Language               { return f.language }
func (f *flow) Type() flows.FlowType                   { return f.flowType }
func (f *flow) ExpireAfterMinutes() int                { return f.expireAfterMinutes }
func (f *flow) Nodes() []flows.Node                    { return f.nodes }
func (f *flow) Localization() flows.Localization       { return f.localization }
func (f *flow) UI() flows.UI                           { return f.ui }
func (f *flow) GetNode(uuid flows.NodeUUID) flows.Node { return f.nodeMap[uuid] }

// Validates that structurally we are sane. IE, all required fields are present and
// all exits with destinations point to valid endpoints.
func (f *flow) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	// check we haven't already started validating this flow to avoid an infinite loop
	if context.IsStarted(f) {
		return nil
	}
	context.Start(f)

	// if this flow has already been validated, don't need to do it again
	if f.valid {
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

		if err := node.Validate(assets, context, f, seenUUIDs); err != nil {
			return errors.Wrapf(err, "validation failed for node[uuid=%s]", node.UUID())
		}
	}

	f.valid = true
	return nil
}

// Resolve resolves the given key when this flow is referenced in an expression
func (f *flow) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	case "uuid":
		return types.NewXText(string(f.UUID()))
	case "name":
		return types.NewXText(f.name)
	case "revision":
		return types.NewXNumberFromInt(f.revision)
	}

	return types.NewXResolveError(f, key)
}

// Describe returns a representation of this type for error messages
func (f *flow) Describe() string { return "flow" }

// Reduce is called when this object needs to be reduced to a primitive
func (f *flow) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(f.name)
}

// ToXJSON is called when this type is passed to @(json(...))
func (f *flow) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, f, "uuid", "name", "revision").ToXJSON(env)
}

var _ flows.Flow = (*flow)(nil)

// Reference returns a reference to this flow asset
func (f *flow) Reference() *assets.FlowReference {
	return assets.NewFlowReference(f.uuid, f.name)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func init() {
	utils.Validator.RegisterAlias("flow_type", "eq=messaging|eq=messaging_offline|eq=voice")
}

type flowEnvelope struct {
	UUID               assets.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name               string          `json:"name" validate:"required"`
	SpecVersion        string          `json:"spec_version" validate:"required"`
	Language           utils.Language  `json:"language" validate:"required"`
	Type               flows.FlowType  `json:"type" validate:"required,flow_type"`
	Revision           int             `json:"revision"`
	ExpireAfterMinutes int             `json:"expire_after_minutes"`
	Localization       localization    `json:"localization"`
	Nodes              []*node         `json:"nodes"`
}

type flowEnvelopeWithUI struct {
	flowEnvelope
	UI flows.UI `json:"_ui,omitempty"`
}

// ReadFlow reads a single flow definition from the passed in byte array
func ReadFlow(data json.RawMessage) (flows.Flow, error) {
	e := &flowEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, errors.Wrap(err, "unable to read flow")
	}
	nodes := make([]flows.Node, len(e.Nodes))
	for n := range e.Nodes {
		nodes[n] = e.Nodes[n]
	}

	return NewFlow(e.UUID, e.Name, e.SpecVersion, e.Language, e.Type, e.Revision, e.ExpireAfterMinutes, e.Localization, nodes, nil), nil
}

// MarshalJSON marshals this flow into JSON
func (f *flow) MarshalJSON() ([]byte, error) {
	var fe = &flowEnvelopeWithUI{
		flowEnvelope: flowEnvelope{
			UUID:               f.uuid,
			Name:               f.name,
			SpecVersion:        f.specVersion,
			Language:           f.language,
			Type:               f.flowType,
			Revision:           f.revision,
			ExpireAfterMinutes: f.expireAfterMinutes,
		},
		UI: f.ui,
	}

	if f.localization != nil {
		fe.Localization = f.localization.(localization)
	}

	fe.Nodes = make([]*node, len(f.nodes))
	for i := range f.nodes {
		fe.Nodes[i] = f.nodes[i].(*node)
	}

	return json.Marshal(fe)
}

// implemention of FlowAssets which provides lazy loading and validation of flows
type flowAssets struct {
	byUUID map[assets.FlowUUID]flows.Flow

	mutex  sync.Mutex
	source assets.AssetSource
}

// NewFlowAssets creates a new flow assets
func NewFlowAssets(source assets.AssetSource) flows.FlowAssets {
	return &flowAssets{
		byUUID: make(map[assets.FlowUUID]flows.Flow),
		source: source,
	}
}

// Get returns the flow with the given UUID
func (a *flowAssets) Get(uuid assets.FlowUUID) (flows.Flow, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	flow := a.byUUID[uuid]
	if flow != nil {
		return flow, nil
	}

	asset, err := a.source.Flow(uuid)
	if err != nil {
		return nil, err
	}

	flow, err = ReadFlow(asset.Definition())
	if err != nil {
		return nil, err
	}

	a.byUUID[flow.UUID()] = flow
	return flow, nil
}
