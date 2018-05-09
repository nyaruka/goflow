package definition

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type flow struct {
	uuid               flows.FlowUUID
	name               string
	language           utils.Language
	expireAfterMinutes int

	localization flows.Localization

	nodes   []flows.Node
	nodeMap map[flows.NodeUUID]flows.Node

	// only read for legacy flows which are being migrated
	ui map[string]interface{}
}

type FlowObj = flow

func NewFlow(uuid flows.FlowUUID, name string, language utils.Language, expireAfterMinutes int, localization flows.Localization, nodes []flows.Node, ui map[string]interface{}) (flows.Flow, error) {
	f := &flow{
		uuid:               uuid,
		name:               name,
		language:           language,
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

func (f *flow) UUID() flows.FlowUUID                   { return f.uuid }
func (f *flow) Name() string                           { return f.name }
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
func (f *flow) Resolve(key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(f.UUID()))
	case "name":
		return types.NewXText(f.name)
	}

	return types.NewXResolveError(f, key)
}

// Repr returns the representation of this type
func (f *flow) Repr() string { return "flow" }

// Reduce is called when this object needs to be reduced to a primitive
func (f *flow) Reduce() types.XPrimitive {
	return types.NewXText(f.name)
}

// ToXJSON is called when this type is passed to @(json(...))
func (f *flow) ToXJSON() types.XText {
	return types.ResolveKeys(f, "uuid", "name").ToXJSON()
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
	UUID               flows.FlowUUID `json:"uuid"               validate:"required,uuid4"`
	Name               string         `json:"name"               validate:"required"`
	Language           utils.Language `json:"language"`
	ExpireAfterMinutes int            `json:"expire_after_minutes"`
	Localization       localization   `json:"localization"`
	Nodes              []*node        `json:"nodes"`
}

type flowEnvelopeWithUI struct {
	flowEnvelope
	UI map[string]interface{} `json:"_ui,omitempty"`
}

// ReadFlow reads a single flow definition from the passed in byte array
func ReadFlow(data json.RawMessage) (flows.Flow, error) {
	envelope := flowEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &envelope, "flow"); err != nil {
		return nil, err
	}
	nodes := make([]flows.Node, len(envelope.Nodes))
	for n := range envelope.Nodes {
		nodes[n] = envelope.Nodes[n]
	}

	return NewFlow(envelope.UUID, envelope.Name, envelope.Language, envelope.ExpireAfterMinutes, envelope.Localization, nodes, nil)
}

// MarshalJSON marshals this flow into JSON
func (f *flow) MarshalJSON() ([]byte, error) {
	var fe = &flowEnvelopeWithUI{
		flowEnvelope: flowEnvelope{
			UUID:               f.uuid,
			Name:               f.name,
			Language:           f.language,
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
