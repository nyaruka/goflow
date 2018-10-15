package definition

import (
	"encoding/json"
	"fmt"

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
	uuid    flows.NodeUUID
	actions []flows.Action
	wait    flows.Wait
	router  flows.Router
	exits   []flows.Exit
}

// NewNode creates a new flow node
func NewNode(uuid flows.NodeUUID, actions []flows.Action, wait flows.Wait, router flows.Router, exits []flows.Exit) flows.Node {
	return &node{
		uuid:    uuid,
		actions: actions,
		wait:    wait,
		router:  router,
		exits:   exits,
	}
}

func (n *node) UUID() flows.NodeUUID    { return n.uuid }
func (n *node) Actions() []flows.Action { return n.actions }
func (n *node) Wait() flows.Wait        { return n.wait }
func (n *node) Router() flows.Router    { return n.router }
func (n *node) Exits() []flows.Exit     { return n.exits }

func (n *node) AddAction(action flows.Action) {
	n.actions = append(n.actions, action)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type nodeEnvelope struct {
	UUID    flows.NodeUUID    `json:"uuid" validate:"required,uuid4"`
	Actions []json.RawMessage `json:"actions,omitempty"`
	Wait    json.RawMessage   `json:"wait,omitempty"`
	Router  json.RawMessage   `json:"router,omitempty"`
	Exits   []*exit           `json:"exits"`
}

// UnmarshalJSON unmarshals a flow node from the given JSON
func (n *node) UnmarshalJSON(data []byte) error {
	e := &nodeEnvelope{}
	err := utils.UnmarshalAndValidate(data, e)
	if err != nil {
		return fmt.Errorf("unable to read node: %s", err)
	}

	n.uuid = e.UUID

	// instantiate the right kind of router
	if e.Router != nil {
		n.router, err = routers.ReadRouter(e.Router)
		if err != nil {
			return fmt.Errorf("unable to read router: %s", err)
		}
	}

	// and the right kind of wait
	if e.Wait != nil {
		n.wait, err = waits.ReadWait(e.Wait)
		if err != nil {
			return fmt.Errorf("unable to read wait: %s", err)
		}
	}

	// and the right kind of actions
	n.actions = make([]flows.Action, len(e.Actions))
	for i := range e.Actions {
		n.actions[i], err = actions.ReadAction(e.Actions[i])
		if err != nil {
			return fmt.Errorf("unable to read action: %s", err)
		}
	}

	// populate our exits
	n.exits = make([]flows.Exit, len(e.Exits))
	for i := range e.Exits {
		n.exits[i] = e.Exits[i]
	}

	return nil
}

// MarshalJSON marshals this flow node into JSON
func (n *node) MarshalJSON() ([]byte, error) {
	var err error

	e := &nodeEnvelope{
		UUID: n.uuid,
	}

	if n.router != nil {
		e.Router, err = json.Marshal(n.router)
		if err != nil {
			return nil, err
		}
	}

	if n.wait != nil {
		e.Wait, err = json.Marshal(n.wait)
		if err != nil {
			return nil, err
		}
	}

	// and the right kind of actions
	e.Actions = make([]json.RawMessage, len(n.actions))
	for i := range n.actions {
		e.Actions[i], err = json.Marshal(n.actions[i])
		if err != nil {
			return nil, err
		}
	}

	e.Exits = make([]*exit, len(n.exits))
	for i := range n.exits {
		e.Exits[i] = n.exits[i].(*exit)
	}

	return json.Marshal(e)
}

type exitEnvelope struct {
	UUID                flows.ExitUUID `json:"uuid"                               validate:"required,uuid4"`
	DestinationNodeUUID flows.NodeUUID `json:"destination_node_uuid,omitempty"    validate:"omitempty,uuid4"`
	Name                string         `json:"name,omitempty"`
}

// UnmarshalJSON unmarshals a node exit from the given JSON
func (e *exit) UnmarshalJSON(data []byte) error {
	var envelope exitEnvelope
	err := utils.UnmarshalAndValidate(data, &envelope)
	if err != nil {
		return fmt.Errorf("unable to read exit: %s", err)
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
