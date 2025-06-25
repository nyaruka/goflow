package actions

import (
	"context"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// max number of times a session can trigger another session without there being input from the contact
const maxAncestorsSinceInput = 5

func init() {
	registerType(TypeStartSession, func() flows.Action { return &StartSession{} })
}

// TypeStartSession is the type for the start session action
const TypeStartSession string = "start_session"

// StartSession can be used to trigger sessions for other contacts and groups. A [event:session_triggered] event
// will be created and it's the responsibility of the caller to act on that by initiating a new session with the flow engine.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "start_session",
//	  "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Registration"},
//	  "groups": [
//	    {"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a", "name": "Customers"}
//	  ],
//	  "exclusions": {"in_a_flow": true}
//	}
//
// @action start_session
type StartSession struct {
	baseAction
	onlineAction
	otherContactsAction

	Flow          *assets.FlowReference `json:"flow" validate:"required"`
	Exclusions    events.Exclusions     `json:"exclusions"`
	CreateContact bool                  `json:"create_contact,omitempty"`
}

// NewStartSession creates a new start session action
func NewStartSession(uuid flows.ActionUUID, flow *assets.FlowReference, groups []*assets.GroupReference, contacts []*flows.ContactReference, contactQuery string, urns []urns.URN, legacyVars []string, createContact bool) *StartSession {
	return &StartSession{
		baseAction: newBaseAction(TypeStartSession, uuid),
		otherContactsAction: otherContactsAction{
			Groups:       groups,
			Contacts:     contacts,
			ContactQuery: contactQuery,
			URNs:         urns,
			LegacyVars:   legacyVars,
		},
		Flow:          flow,
		CreateContact: createContact,
	}
}

// Execute runs our action
func (a *StartSession) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	groupRefs, contactRefs, contactQuery, urnList, err := a.resolveRecipients(run, logEvent)
	if err != nil {
		return err
	}

	// check that flow exists - error event if not
	flow, err := run.Session().Assets().Flows().Get(a.Flow.UUID)
	if err != nil {
		logEvent(events.NewDependencyError(a.Flow))
		return nil
	}

	// batch footgun prevention
	if run.Session().BatchStart() && (len(groupRefs) > 0 || contactQuery != "") {
		logEvent(events.NewError("can't start new sessions for groups or queries during batch starts"))
		return nil
	}

	// loop footgun prevention
	ref := run.Session().History()
	if ref.AncestorsSinceInput >= maxAncestorsSinceInput {
		logEvent(events.NewError("too many sessions have been spawned since the last time input was received"))
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

	logEvent(events.NewSessionTriggered(flow.Reference(false), groupRefs, contactRefs, contactQuery, a.Exclusions, a.CreateContact, urnList, runSnapshot, history))
	return nil
}

func (a *StartSession) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	a.otherContactsAction.Inspect(dependency, local, result)

	dependency(a.Flow)
}
