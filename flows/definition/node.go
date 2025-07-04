package definition

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/utils"
)

type node struct {
	uuid    flows.NodeUUID
	actions []flows.Action
	router  flows.Router
	exits   []flows.Exit
}

// NewNode creates a new flow node
func NewNode(uuid flows.NodeUUID, actions []flows.Action, router flows.Router, exits []flows.Exit) flows.Node {
	return &node{
		uuid:    uuid,
		actions: actions,
		router:  router,
		exits:   exits,
	}
}

func (n *node) UUID() flows.NodeUUID    { return n.uuid }
func (n *node) Actions() []flows.Action { return n.actions }
func (n *node) Router() flows.Router    { return n.router }
func (n *node) Exits() []flows.Exit     { return n.exits }

func (n *node) Validate(flow flows.Flow, seenUUIDs map[uuids.UUID]bool) error {
	if len(n.actions) > flows.MaxActionsPerNode {
		return fmt.Errorf("node can't have more than %d actions (has %d)", flows.MaxActionsPerNode, len(n.actions))
	}
	if len(n.exits) > flows.MaxExitsPerNode {
		return fmt.Errorf("node can't have more than %d exits (has %d)", flows.MaxExitsPerNode, len(n.exits))
	}

	// validate all the node's actions
	for _, action := range n.Actions() {

		// check that this action is valid for this flow type
		if !flow.Type().Allows(action) {
			return fmt.Errorf("action type '%s' is not allowed in a flow of type '%s'", action.Type(), flow.Type())
		}

		uuidAlreadySeen := seenUUIDs[uuids.UUID(action.UUID())]
		if uuidAlreadySeen {
			return fmt.Errorf("action UUID %s isn't unique", action.UUID())
		}
		seenUUIDs[uuids.UUID(action.UUID())] = true

		if err := action.Validate(); err != nil {
			return fmt.Errorf("invalid action[uuid=%s, type=%s]: %w", action.UUID(), action.Type(), err)
		}
	}

	// check the router if there is one
	if n.Router() != nil {
		if err := n.Router().Validate(flow, n.Exits()); err != nil {
			return fmt.Errorf("invalid router: %w", err)
		}
	}

	// check every exit has a unique UUID and valid destination
	for _, exit := range n.Exits() {
		uuidAlreadySeen := seenUUIDs[uuids.UUID(exit.UUID())]
		if uuidAlreadySeen {
			return fmt.Errorf("exit UUID %s isn't unique", exit.UUID())
		}
		seenUUIDs[uuids.UUID(exit.UUID())] = true

		if exit.DestinationUUID() != "" && flow.GetNode(exit.DestinationUUID()) == nil {
			return fmt.Errorf("destination %s of exit[uuid=%s] isn't a known node", exit.DestinationUUID(), exit.UUID())
		}
	}

	return nil
}

// Inspect reports on the results, dependencies etc used by this node
func (n *node) Inspect(dependency func(flows.Action, flows.Router, assets.Reference), local func(flows.Action, flows.Router, string), result func(flows.Action, flows.Router, *flows.ResultInfo)) {
	for _, action := range n.actions {
		action.Inspect(func(r assets.Reference) {
			dependency(action, nil, r)
		}, func(l string) {
			local(action, nil, l)
		}, func(r *flows.ResultInfo) {
			result(action, nil, r)
		})
	}

	if n.router != nil {
		n.router.Inspect(func(r *flows.ResultInfo) {
			result(nil, n.router, r)
		}, func(r assets.Reference) {
			dependency(nil, n.router, r)
		})
	}
}

// EnumerateTemplates enumerates all expressions on this object
func (n *node) EnumerateTemplates(localization flows.Localization, include func(flows.Action, flows.Router, i18n.Language, string)) {
	for _, action := range n.actions {
		inspect.Templates(action, localization, func(l i18n.Language, t string) {
			include(action, nil, l, t)
		})
	}

	if n.router != nil {
		n.router.EnumerateTemplates(localization, func(l i18n.Language, t string) {
			include(nil, n.router, l, t)
		})
	}
}

// EnumerateLocalizables enumerates all localizable text on this object
func (n *node) EnumerateLocalizables(include func(uuids.UUID, string, []string, func([]string))) {
	for _, action := range n.actions {
		inspect.LocalizableText(action, include)
	}

	if n.router != nil {
		n.router.EnumerateLocalizables(include)
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type nodeEnvelope struct {
	UUID    flows.NodeUUID    `json:"uuid"               validate:"required,uuid"`
	Actions []json.RawMessage `json:"actions,omitempty"  validate:"dive,required"`
	Router  json.RawMessage   `json:"router,omitempty"`
	Exits   []*exit           `json:"exits"              validate:"required,min=1,dive,required"`
}

// UnmarshalJSON unmarshals a flow node from the given JSON
func (n *node) UnmarshalJSON(data []byte) error {
	e := &nodeEnvelope{}
	err := utils.UnmarshalAndValidate(data, e)
	if err != nil {
		return fmt.Errorf("unable to read node: %w", err)
	}

	n.uuid = e.UUID

	// instantiate the right kind of router
	if e.Router != nil {
		n.router, err = routers.Read(e.Router)
		if err != nil {
			return fmt.Errorf("unable to read router: %w", err)
		}
	}

	// and the right kind of actions
	n.actions = make([]flows.Action, len(e.Actions))
	for i := range e.Actions {
		n.actions[i], err = actions.Read(e.Actions[i])
		if err != nil {
			return fmt.Errorf("unable to read action: %w", err)
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

	e.Actions = make([]json.RawMessage, len(n.actions))
	for i := range n.actions {
		e.Actions[i], err = jsonx.Marshal(n.actions[i])
		if err != nil {
			return nil, err
		}
	}

	if n.router != nil {
		e.Router, err = jsonx.Marshal(n.router)
		if err != nil {
			return nil, err
		}
	}

	e.Exits = make([]*exit, len(n.exits))
	for i := range n.exits {
		e.Exits[i] = n.exits[i].(*exit)
	}

	return jsonx.Marshal(e)
}
