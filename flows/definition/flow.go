package definition

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	flows.SetFlowReader(ReadFlow)
}

type flow struct {
	uuid     assets.FlowUUID
	name     string
	language utils.Language
	flowType flows.FlowType

	revision           int
	expireAfterMinutes int
	localization       flows.Localization

	nodes   []flows.Node
	nodeMap map[flows.NodeUUID]flows.Node

	// only read for legacy flows which are being migrated
	ui map[string]interface{}
}

func NewFlow(uuid assets.FlowUUID, name string, language utils.Language, flowType flows.FlowType, revision int, expireAfterMinutes int, localization flows.Localization, nodes []flows.Node, ui map[string]interface{}) (flows.Flow, error) {
	f := &flow{
		uuid:               uuid,
		name:               name,
		language:           language,
		flowType:           flowType,
		revision:           revision,
		expireAfterMinutes: expireAfterMinutes,
		localization:       localization,
		nodes:              nodes,
		ui:                 ui,
	}
	if err := f.buildNodeMap(); err != nil {
		return nil, err
	}

	// go back through nodes and perform basic structural validation
	for _, node := range f.nodes {

		// check every exit has a valid destination
		for _, exit := range node.Exits() {
			if exit.DestinationNodeUUID() != "" && f.nodeMap[exit.DestinationNodeUUID()] == nil {
				return nil, fmt.Errorf("destination %s of exit[uuid=%s] isn't a known node", exit.DestinationNodeUUID(), exit.UUID())
			}
		}

		// and the router if there is one
		if node.Router() != nil {
			if err := node.Router().Validate(node.Exits()); err != nil {
				return nil, fmt.Errorf("router is invalid on node[uuid=%s]: %v", node.UUID(), err)
			}
		}
	}

	return f, nil
}

func (f *flow) UUID() assets.FlowUUID                  { return f.uuid }
func (f *flow) Name() string                           { return f.name }
func (f *flow) Revision() int                          { return f.revision }
func (f *flow) Language() utils.Language               { return f.language }
func (f *flow) ExpireAfterMinutes() int                { return f.expireAfterMinutes }
func (f *flow) Nodes() []flows.Node                    { return f.nodes }
func (f *flow) Localization() flows.Localization       { return f.localization }
func (f *flow) GetNode(uuid flows.NodeUUID) flows.Node { return f.nodeMap[uuid] }

// Validates that structurally we are sane. IE, all required fields are present and
// all exits with destinations point to valid endpoints.
func (f *flow) Validate(assets flows.SessionAssets) error {
	var err error

	// track UUIDs used by nodes and actions to ensure that they are unique
	seenUUIDs := make(map[utils.UUID]bool)

	for _, node := range f.nodes {
		uuidAlreadySeen := seenUUIDs[utils.UUID(node.UUID())]
		if uuidAlreadySeen {
			return fmt.Errorf("node UUID %s isn't unique", node.UUID())
		}
		seenUUIDs[utils.UUID(node.UUID())] = true

		// validate all the node's actions
		for _, action := range node.Actions() {

			// check that this action is valid for this flow type
			isValidInType := false
			for _, allowedType := range action.AllowedFlowTypes() {
				if f.flowType == allowedType {
					isValidInType = true
					break
				}
			}
			if !isValidInType {
				return fmt.Errorf("action type '%s' is not allowed in a flow of type '%s'", action.Type(), f.flowType)
			}

			uuidAlreadySeen := seenUUIDs[utils.UUID(action.UUID())]
			if uuidAlreadySeen {
				return fmt.Errorf("action UUID %s isn't unique", action.UUID())
			}
			seenUUIDs[utils.UUID(action.UUID())] = true

			if err := action.Validate(assets); err != nil {
				return fmt.Errorf("validation failed for action[uuid=%s, type=%s]: %v", action.UUID(), action.Type(), err)
			}
		}
	}
	return err
}

// Resolve resolves the given key when this flow is referenced in an expression
func (f *flow) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
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

func (f *flow) Reference() *flows.FlowReference {
	return flows.NewFlowReference(f.uuid, f.name)
}

func (f *flow) buildNodeMap() error {
	f.nodeMap = make(map[flows.NodeUUID]flows.Node)

	for _, node := range f.nodes {
		// make sure we haven't seen this node before
		if f.nodeMap[node.UUID()] != nil {
			return fmt.Errorf("duplicate node uuid: %s", node.UUID())
		}
		f.nodeMap[node.UUID()] = node
	}
	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type flowEnvelope struct {
	UUID               assets.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name               string          `json:"name" validate:"required"`
	Language           utils.Language  `json:"language" validate:"required"`
	Type               flows.FlowType  `json:"type" validate:"required"`
	Revision           int             `json:"revision"`
	ExpireAfterMinutes int             `json:"expire_after_minutes"`
	Localization       localization    `json:"localization"`
	Nodes              []*node         `json:"nodes"`
}

type flowEnvelopeWithUI struct {
	flowEnvelope
	UI map[string]interface{} `json:"_ui,omitempty"`
}

// ReadFlow reads a single flow definition from the passed in byte array
func ReadFlow(data json.RawMessage) (flows.Flow, error) {
	e := &flowEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, fmt.Errorf("unable to read flow: %s", err)
	}
	nodes := make([]flows.Node, len(e.Nodes))
	for n := range e.Nodes {
		nodes[n] = e.Nodes[n]
	}

	return NewFlow(e.UUID, e.Name, e.Language, e.Type, e.Revision, e.ExpireAfterMinutes, e.Localization, nodes, nil)
}

// MarshalJSON marshals this flow into JSON
func (f *flow) MarshalJSON() ([]byte, error) {
	var fe = &flowEnvelopeWithUI{
		flowEnvelope: flowEnvelope{
			UUID:               f.uuid,
			Name:               f.name,
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
