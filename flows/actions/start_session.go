package actions

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeStartSession, func() flows.Action { return &StartSessionAction{} })
}

// TypeStartSession is the type for the start session action
const TypeStartSession string = "start_session"

// StartSessionAction can be used to trigger sessions for other contacts and groups. A [event:session_triggered] event
// will be created and it's the responsibility of the caller to act on that by initiating a new session with the flow engine.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "start_session",
//     "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Registration"},
//     "groups": [
//       {"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a", "name": "Customers"}
//     ]
//   }
//
// @action start_session
type StartSessionAction struct {
	BaseAction
	onlineAction

	URNs          []urns.URN                `json:"urns,omitempty"`
	Contacts      []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	Groups        []*assets.GroupReference  `json:"groups,omitempty" validate:"dive"`
	LegacyVars    []string                  `json:"legacy_vars,omitempty"`
	Flow          *assets.FlowReference     `json:"flow" validate:"required"`
	CreateContact bool                      `json:"create_contact,omitempty"`
}

// NewStartSessionAction creates a new start session action
func NewStartSessionAction(uuid flows.ActionUUID, urns []urns.URN, contacts []*flows.ContactReference, groups []*assets.GroupReference, legacyVars []string, flow *assets.FlowReference, createContact bool) *StartSessionAction {
	return &StartSessionAction{
		BaseAction:    NewBaseAction(TypeStartSession, uuid),
		URNs:          urns,
		Contacts:      contacts,
		Groups:        groups,
		LegacyVars:    legacyVars,
		Flow:          flow,
		CreateContact: createContact,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *StartSessionAction) Validate(assets flows.SessionAssets) error {
	// check we have the flow
	if _, err := assets.Flows().Get(a.Flow.UUID); err != nil {
		return err
	}
	// check we have all groups
	return a.validateGroups(assets, a.Groups)
}

// Execute runs our action
func (a *StartSessionAction) Execute(run flows.FlowRun, step flows.Step) error {
	urnList, contactRefs, groupRefs, err := a.resolveContactsAndGroups(run, step, a.URNs, a.Contacts, a.Groups, a.LegacyVars)
	if err != nil {
		return err
	}

	runSnapshot, err := json.Marshal(run.Snapshot())
	if err != nil {
		return err
	}

	a.log(run, step, events.NewSessionTriggeredEvent(a.Flow, urnList, contactRefs, groupRefs, a.CreateContact, runSnapshot))
	return nil
}
