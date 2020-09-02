package actions

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// max number of times a session can trigger another session without there being input from the contact
const maxAncestorsSinceInput = 5

func init() {
	registerType(TypeStartSession, func() flows.Action { return &StartSessionAction{} })
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
	baseAction
	onlineAction
	otherContactsAction

	Flow          *assets.FlowReference `json:"flow" validate:"required"`
	CreateContact bool                  `json:"create_contact,omitempty"`
}

// NewStartSession creates a new start session action
func NewStartSession(uuid flows.ActionUUID, flow *assets.FlowReference, urns []urns.URN, contacts []*flows.ContactReference, groups []*assets.GroupReference, legacyVars []string, createContact bool) *StartSessionAction {
	return &StartSessionAction{
		baseAction: newBaseAction(TypeStartSession, uuid),
		otherContactsAction: otherContactsAction{
			URNs:       urns,
			Contacts:   contacts,
			Groups:     groups,
			LegacyVars: legacyVars,
		},
		Flow:          flow,
		CreateContact: createContact,
	}
}

// Execute runs our action
func (a *StartSessionAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	groupRefs, contactRefs, contactQuery, urnList, err := a.resolveRecipients(run, logEvent)
	if err != nil {
		return err
	}

	// batch footgun prevention
	if run.Session().BatchStart() && (len(groupRefs) > 0 || contactQuery != "") {
		logEvent(events.NewErrorf("can't start new sessions for groups or queries during batch starts"))
		return nil
	}

	// loop footgun prevention
	ref := run.Session().History()
	if ref.AncestorsSinceInput >= maxAncestorsSinceInput {
		logEvent(events.NewErrorf("too many sessions have been spawned since the last time input was received"))
		return nil
	}

	// if we don't have any recipients, noop
	if !(len(urnList) > 0 || len(groupRefs) > 0 || len(contactRefs) > 0 || a.ContactQuery != "" || a.CreateContact) {
		return nil
	}

	runSnapshot, err := jsonx.Marshal(run.Snapshot())
	if err != nil {
		return err
	}

	history := flows.NewChildHistory(run.Session())

	logEvent(events.NewSessionTriggered(a.Flow, groupRefs, contactRefs, contactQuery, a.CreateContact, urnList, runSnapshot, history))
	return nil
}
