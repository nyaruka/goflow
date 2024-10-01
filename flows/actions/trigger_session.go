package actions

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// max number of times a session can trigger another session without there being input from the contact
const maxAncestorsSinceInput = 5

func init() {
	registerType(TypeTriggerSession, func() flows.Action { return &TriggerSessionAction{} })
}

// TypeTriggerSession is the type for the trigger session action
const TypeTriggerSession string = "trigger_session"

// TriggerSessionAction can be used to trigger sessions for another contact. A [event:session_triggered] event will be
// created and it's the responsibility of the caller to act on that by initiating a new session with the flow engine.
// The contact can be specified via a reference or as a URN. In the latter case the contact will be created if they
// don't exist.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "trigger_session",
//	  "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Registration"},
//	  "contact": {"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a", "name": "Bob"},
//	  "interrupt": true
//	}
//
// @action start_session
type TriggerSessionAction struct {
	baseAction
	onlineAction

	Flow      *assets.FlowReference   `json:"flow"               validate:"required"`
	Contact   *flows.ContactReference `json:"contact,omitempty"`
	URN       string                  `json:"urn,omitempty"      engine:"evaluated"`
	Interrupt bool                    `json:"interrupt"`
}

// NewTriggerSession creates a new trigger session action
func NewTriggerSession(uuid flows.ActionUUID, flow *assets.FlowReference, contact *flows.ContactReference, urn string, interrupt bool) *TriggerSessionAction {
	return &TriggerSessionAction{
		baseAction: newBaseAction(TypeTriggerSession, uuid),
		Flow:       flow,
		Contact:    contact,
		URN:        urn,
		Interrupt:  interrupt,
	}
}

// Validate validates our action is valid
func (a *TriggerSessionAction) Validate() error {
	if (a.Contact != nil && a.URN != "") || (a.Contact == nil && a.URN == "") {
		return fmt.Errorf("must specify either contact or urn")
	}
	return nil
}

// Execute runs our action
func (a *TriggerSessionAction) Execute(run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	urn := a.resolveURN(run, logEvent)

	if urn == urns.NilURN && a.Contact == nil {
		return nil
	}

	// check that flow exists - error event if not
	flow, err := run.Session().Assets().Flows().Get(a.Flow.UUID)
	if err != nil {
		logEvent(events.NewDependencyError(a.Flow))
		return nil
	}

	// loop footgun prevention
	ref := run.Session().History()
	if ref.AncestorsSinceInput >= maxAncestorsSinceInput {
		logEvent(events.NewErrorf("too many sessions have been spawned since the last time input was received"))
		return nil
	}

	runSnapshot, err := jsonx.Marshal(run.Snapshot())
	if err != nil {
		return err
	}

	history := flows.NewChildHistory(run.Session())

	logEvent(events.NewSessionTriggered(flow.Reference(false), a.Contact, urn, a.Interrupt, runSnapshot, history))
	return nil
}

func (a *TriggerSessionAction) resolveURN(run flows.Run, logEvent flows.EventCallback) urns.URN {
	if a.URN == "" {
		return urns.NilURN
	}

	// otherwise this is a variable reference so evaluate it
	evaluatedURN, ok := run.EvaluateTemplate(a.URN, logEvent)
	if !ok {
		return urns.NilURN
	}

	// if we have a valid URN now, return it
	urn := urns.URN(evaluatedURN)
	if urn.Validate() == nil {
		return urn.Normalize()
	}

	// otherwise try to parse as phone number
	parsedTel := utils.ParsePhoneNumber(evaluatedURN, run.Session().MergedEnvironment().DefaultCountry())
	if parsedTel != "" {
		urn, _ := urns.New(urns.Phone, parsedTel)
		return urn.Normalize()
	}

	return urns.NilURN
}
