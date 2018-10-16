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

func (n *node) Validate(assets flows.SessionAssets, context *flows.ValidationContext, flow flows.Flow, seenUUIDs map[utils.UUID]bool) error {
	// validate all the node's actions
	for _, action := range n.Actions() {

		// check that this action is valid for this flow type
		isValidInType := false
		for _, allowedType := range action.AllowedFlowTypes() {
			if flow.Type() == allowedType {
				isValidInType = true
				break
			}
		}
		if !isValidInType {
			return fmt.Errorf("action type '%s' is not allowed in a flow of type '%s'", action.Type(), flow.Type())
		}

		uuidAlreadySeen := seenUUIDs[utils.UUID(action.UUID())]
		if uuidAlreadySeen {
			return fmt.Errorf("action UUID %s isn't unique", action.UUID())
		}
		seenUUIDs[utils.UUID(action.UUID())] = true

		if err := action.Validate(assets, context); err != nil {
			return fmt.Errorf("validation failed for action[uuid=%s, type=%s]: %v", action.UUID(), action.Type(), err)
		}
	}

	// check the router if there is one
	if n.Router() != nil {
		if err := n.Router().Validate(n.Exits()); err != nil {
			return fmt.Errorf("validation failed for router: %s", err)
		}
	}

	// check we have at least one exit
	if len(n.Exits()) < 1 {
		return fmt.Errorf("nodes must have at least one exit")
	}

	// check every exit has a unique UUID and valid destination
	for _, exit := range n.Exits() {
		uuidAlreadySeen := seenUUIDs[utils.UUID(exit.UUID())]
		if uuidAlreadySeen {
			return fmt.Errorf("exit UUID %s isn't unique", exit.UUID())
		}
		seenUUIDs[utils.UUID(exit.UUID())] = true

		if exit.DestinationNodeUUID() != "" && flow.GetNode(exit.DestinationNodeUUID()) == nil {
			return fmt.Errorf("destination %s of exit[uuid=%s] isn't a known node", exit.DestinationNodeUUID(), exit.UUID())
		}
	}

	return nil
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
