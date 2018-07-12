package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/utils"
)

type exit struct {
	uuid        flows.ExitUUID
	destination flows.NodeUUID
	name        string
}

// NewExit creates a new exit
func NewExit(uuid flows.ExitUUID, destination flows.NodeUUID, name string) flows.Exit {
	return &exit{uuid: uuid, destination: destination, name: name}
}

func (e *exit) UUID() flows.ExitUUID                { return e.uuid }
func (e *exit) DestinationNodeUUID() flows.NodeUUID { return e.destination }
func (e *exit) Name() string                        { return e.name }

type node struct {
	uuid flows.NodeUUID

	actions []flows.Action
	router  flows.Router
	exits   []flows.Exit
	wait    flows.Wait
}

// NewNode creates a new flow node
func NewNode(uuid flows.NodeUUID, actions []flows.Action, router flows.Router, exits []flows.Exit, wait flows.Wait) flows.Node {
	return &node{
		uuid:    uuid,
		actions: actions,
		router:  router,
		exits:   exits,
		wait:    wait,
	}
}

func (n *node) UUID() flows.NodeUUID    { return n.uuid }
func (n *node) Router() flows.Router    { return n.router }
func (n *node) Actions() []flows.Action { return n.actions }
func (n *node) Exits() []flows.Exit     { return n.exits }
func (n *node) Wait() flows.Wait        { return n.wait }

func (n *node) AddAction(action flows.Action) {
	n.actions = append(n.actions, action)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type nodeEnvelope struct {
	UUID    flows.NodeUUID         `json:"uuid"                  validate:"required,uuid4"`
	Actions []*utils.TypedEnvelope `json:"actions,omitempty"`
	Router  *utils.TypedEnvelope   `json:"router,omitempty"`
	Exits   []*exit                `json:"exits"`
	Wait    *utils.TypedEnvelope   `json:"wait,omitempty"`
}

// UnmarshalJSON unmarshals a flow node from the given JSON
func (n *node) UnmarshalJSON(data []byte) error {
	var envelope nodeEnvelope
	err := utils.UnmarshalAndValidate(data, &envelope, "node")
	if err != nil {
		return err
	}

	n.uuid = envelope.UUID

	// instantiate the right kind of router
	if envelope.Router != nil {
		n.router, err = routers.ReadRouter(envelope.Router)
		if err != nil {
			return err
		}
	}

	// and the right kind of wait
	if envelope.Wait != nil {
		n.wait, err = waits.ReadWait(envelope.Wait)
		if err != nil {
			return err
		}
	}

	// and the right kind of actions
	n.actions = make([]flows.Action, len(envelope.Actions))
	for i := range envelope.Actions {
		n.actions[i], err = actions.ReadAction(envelope.Actions[i])
		if err != nil {
			return err
		}
	}

	// populate our exits
	n.exits = make([]flows.Exit, len(envelope.Exits))
	for i := range envelope.Exits {
		n.exits[i] = envelope.Exits[i]
	}

	return nil
}

// MarshalJSON marshals this flow node into JSON
func (n *node) MarshalJSON() ([]byte, error) {
	envelope := nodeEnvelope{}
	var err error

	envelope.UUID = n.uuid

	envelope.Router, err = utils.EnvelopeFromTyped(n.router)
	if err != nil {
		return nil, err
	}

	envelope.Wait, err = utils.EnvelopeFromTyped(n.wait)
	if err != nil {
		return nil, err
	}

	// and the right kind of actions
	envelope.Actions = make([]*utils.TypedEnvelope, len(n.actions))
	for i := range n.actions {
		envelope.Actions[i], err = utils.EnvelopeFromTyped(n.actions[i])
		if err != nil {
			return nil, err
		}
	}

	envelope.Exits = make([]*exit, len(n.exits))
	for i := range n.exits {
		envelope.Exits[i] = n.exits[i].(*exit)
	}

	return json.Marshal(envelope)
}

type exitEnvelope struct {
	UUID                flows.ExitUUID `json:"uuid"                               validate:"required,uuid4"`
	DestinationNodeUUID flows.NodeUUID `json:"destination_node_uuid,omitempty"    validate:"omitempty,uuid4"`
	Name                string         `json:"name,omitempty"`
}

// UnmarshalJSON unmarshals a node exit from the given JSON
func (e *exit) UnmarshalJSON(data []byte) error {
	var envelope exitEnvelope
	err := utils.UnmarshalAndValidate(data, &envelope, "exit")
	if err != nil {
		return err
	}

	e.uuid = envelope.UUID
	e.destination = envelope.DestinationNodeUUID
	e.name = envelope.Name

	return nil
}

// MarshalJSON marshals this node exit into JSON
func (e *exit) MarshalJSON() ([]byte, error) {
	envelope := exitEnvelope{e.uuid, e.destination, e.name}
	return json.Marshal(envelope)
}
