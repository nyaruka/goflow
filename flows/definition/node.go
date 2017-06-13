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

func (n *node) UUID() flows.NodeUUID    { return n.uuid }
func (n *node) Router() flows.Router    { return n.router }
func (n *node) Actions() []flows.Action { return n.actions }
func (n *node) Exits() []flows.Exit     { return n.exits }
func (n *node) Wait() flows.Wait        { return n.wait }

func (n *node) Resolve(key string) interface{} {
	switch key {
	case "uuid":
		return n.uuid
	}
	return fmt.Errorf("No field '%s' on node", key)
}

func (n *node) Default() interface{} { return n.uuid }
func (n *node) String() string       { return n.uuid.String() }

var _ utils.VariableResolver = (*node)(nil)

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

func (n *node) UnmarshalJSON(data []byte) error {
	envelope := nodeEnvelope{}

	err := json.Unmarshal(data, &envelope)
	err = utils.ValidateAllUnlessErr(err, &envelope)
	if err != nil {
		return err
	}

	n.uuid = envelope.UUID

	// instantiate the right kind of router
	if envelope.Router != nil {
		n.router, err = routers.RouterFromEnvelope(envelope.Router)
		if err != nil {
			return err
		}
	}

	// and the right kind of wait
	if envelope.Wait != nil {
		n.wait, err = waits.WaitFromEnvelope(envelope.Wait)
		if err != nil {
			return err
		}
	}

	// and the right kind of actions
	n.actions = make([]flows.Action, len(envelope.Actions))
	for i := range envelope.Actions {
		n.actions[i], err = actions.ActionFromEnvelope(envelope.Actions[i])
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

func (e *exit) UnmarshalJSON(data []byte) error {
	envelope := exitEnvelope{}

	err := json.Unmarshal(data, &envelope)
	err = utils.ValidateAllUnlessErr(err, &envelope)
	if err != nil {
		return err
	}

	e.uuid = envelope.UUID
	e.destination = envelope.DestinationNodeUUID
	e.name = envelope.Name

	return nil
}

func (e *exit) MarshalJSON() ([]byte, error) {
	envelope := exitEnvelope{e.uuid, e.destination, e.name}
	return json.Marshal(envelope)
}
