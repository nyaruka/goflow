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
func (f *flow) Validate(assets flows.Assets) error {
	var err error
	f.nodeMap = make(map[flows.NodeUUID]flows.Node)

	// build up a list of all our node ids
	for i, node := range f.nodes {
		// make sure we haven't seen this node before
		if f.nodeMap[node.UUID()] != nil {
			return fmt.Errorf("duplicate node uuid: '%s'", node.UUID())
		}
		f.nodeMap[node.UUID()] = f.nodes[i]
	}

	for _, node := range f.nodes {
		// validate all our actions
		for _, action := range node.Actions() {
			err = action.Validate(assets)
			if err != nil {
				asJSON, jerr := json.MarshalIndent(action, "", "  ")
				if jerr != nil {
					return fmt.Errorf("%+v: %v", action, err)
				}
				return fmt.Errorf("%s: %v", asJSON, err)
			}
		}

		// and our router if we have one
		router := node.Router()
		if router != nil {
			err = router.Validate(node.Exits())
			if err != nil {
				asJSON, jerr := json.MarshalIndent(node.Router(), "", "  ")
				if jerr != nil {
					return fmt.Errorf("%+v: %v", node.Router(), err)
				}
				return fmt.Errorf("%s: %v", asJSON, err)
			}
		}

		// make sure all our exits have valid destinations
		for _, exit := range node.Exits() {
			if exit.DestinationNodeUUID() != "" && f.nodeMap[exit.DestinationNodeUUID()] == nil {
				return fmt.Errorf("exit '%s' on node '%s' has invalid destination '%s'", exit.UUID(), node.UUID(), exit.DestinationNodeUUID())
			}
		}
	}
	return err
}

func (f *flow) Resolve(key string) interface{} {
	switch key {

	case "name":
		return f.Name()

	case "uuid":
		return f.UUID()

	}

	return fmt.Errorf("No field '%s' on flow", key)
}

func (f *flow) Default() interface{} {
	return f
}

// String returns the default string value for this flow, which is just our name
func (f *flow) String() string {
	return f.name
}

var _ utils.VariableResolver = (*flow)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadFlow reads a single flow definition from the passed in byte array
func ReadFlow(data json.RawMessage) (flows.Flow, error) {
	flow := &flow{}
	err := json.Unmarshal(data, flow)
	return flow, err
}

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

func (f *flow) UnmarshalJSON(data []byte) error {
	var envelope flowEnvelope
	err := utils.UnmarshalAndValidate(data, &envelope, "flow")
	if err != nil {
		return err
	}

	f.uuid = envelope.UUID
	f.name = envelope.Name
	f.language = envelope.Language
	f.expireAfterMinutes = envelope.ExpireAfterMinutes

	f.translations = &envelope.Localization

	// for each node
	f.nodes = make([]flows.Node, len(envelope.Nodes))
	for i := range envelope.Nodes {
		f.nodes[i] = envelope.Nodes[i]
	}

	return err
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
