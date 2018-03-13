package definition

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type flow struct {
	uuid               flows.FlowUUID
	name               string
	language           utils.Language
	expireAfterMinutes int

	translations flows.FlowTranslations

	nodes   []flows.Node
	nodeMap map[flows.NodeUUID]flows.Node
}

func (f *flow) UUID() flows.FlowUUID                   { return f.uuid }
func (f *flow) Name() string                           { return f.name }
func (f *flow) Language() utils.Language               { return f.language }
func (f *flow) ExpireAfterMinutes() int                { return f.expireAfterMinutes }
func (f *flow) Nodes() []flows.Node                    { return f.nodes }
func (f *flow) Translations() flows.FlowTranslations   { return f.translations }
func (f *flow) GetNode(uuid flows.NodeUUID) flows.Node { return f.nodeMap[uuid] }

// Validates that structurally we are sane. IE, all required fields are present and
// all exits with destinations point to valid endpoints.
func (f *flow) Validate(assets flows.SessionAssets) error {
	var err error

	for _, node := range f.nodes {
		// validate all the node's actions
		for _, action := range node.Actions() {
			if err := action.Validate(assets); err != nil {
				return fmt.Errorf("validation failed for action[uuid=%s, type=%s]: %v", action.UUID(), action.Type(), err)
			}
		}
	}
	return err
}

// Resolve resolves the given key when this flow is referenced in an expression
func (f *flow) Resolve(key string) interface{} {
	switch key {
	case "uuid":
		return f.UUID()
	case "name":
		return f.Name()
	}

	return fmt.Errorf("no field '%s' on flow", key)
}

// Default returns the value of this flow when it is the result of an expression
func (f *flow) Default() interface{} {
	return f
}

// String returns the default string value for this flow, which is just our name
func (f *flow) String() string {
	return f.name
}

func (f *flow) Reference() *flows.FlowReference {
	return flows.NewFlowReference(f.uuid, f.name)
}

var _ utils.VariableResolver = (*flow)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type flowEnvelope struct {
	UUID               flows.FlowUUID   `json:"uuid"               validate:"required,uuid4"`
	Name               string           `json:"name"               validate:"required"`
	Language           utils.Language   `json:"language"`
	ExpireAfterMinutes int              `json:"expire_after_minutes"`
	Localization       flowTranslations `json:"localization"`
	Nodes              []*node          `json:"nodes"`

	// only for writing out, optional
	Metadata map[string]interface{} `json:"_ui,omitempty"`
}

// ReadFlow reads a single flow definition from the passed in byte array
func ReadFlow(data json.RawMessage) (flows.Flow, error) {
	envelope := flowEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &envelope, "flow"); err != nil {
		return nil, err
	}

	f := &flow{}
	f.uuid = envelope.UUID
	f.name = envelope.Name
	f.language = envelope.Language
	f.expireAfterMinutes = envelope.ExpireAfterMinutes
	f.translations = envelope.Localization

	f.nodes = make([]flows.Node, len(envelope.Nodes))
	f.nodeMap = make(map[flows.NodeUUID]flows.Node)

	// for each node...
	for n, node := range envelope.Nodes {
		f.nodes[n] = node

		// make sure we haven't seen this node before
		if f.nodeMap[node.UUID()] != nil {
			return nil, fmt.Errorf("duplicate node uuid: %s", node.UUID())
		}
		f.nodeMap[node.UUID()] = node
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

func (f *flow) MarshalJSON() ([]byte, error) {

	var fe = flowEnvelope{}
	fe.UUID = f.uuid
	fe.Name = f.name
	fe.Language = f.language
	fe.ExpireAfterMinutes = f.expireAfterMinutes

	if f.translations != nil {
		fe.Localization = *f.translations.(*flowTranslations)
	}

	fe.Nodes = make([]*node, len(f.nodes))
	for i := range f.nodes {
		fe.Nodes[i] = f.nodes[i].(*node)
	}

	return json.Marshal(&fe)
}
