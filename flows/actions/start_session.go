package actions

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/urns"
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
	Groups        []*flows.GroupReference   `json:"groups,omitempty" validate:"dive"`
	LegacyVars    []string                  `json:"legacy_vars,omitempty"`
	Flow          *flows.FlowReference      `json:"flow" validate:"required"`
	CreateContact bool                      `json:"create_contact,omitempty"`
}

// Type returns the type of this action
func (a *StartSessionAction) Type() string { return TypeStartSession }

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
func (a *StartSessionAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	urnList, contactRefs, groupRefs, err := a.resolveContactsAndGroups(run, step, a.URNs, a.Contacts, a.Groups, a.LegacyVars, log)
	if err != nil {
		return err
	}

	runSnapshot, err := json.Marshal(run.Snapshot())
	if err != nil {
		return err
	}

	log.Add(events.NewSessionTriggeredEvent(a.Flow, urnList, contactRefs, groupRefs, a.CreateContact, runSnapshot))
	return nil
}
