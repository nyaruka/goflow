package actions

import (
	"fmt"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeAddContactURN, func() flows.Action { return &AddContactURNAction{} })
}

// TypeAddContactURN is our type for the add URN action
const TypeAddContactURN string = "add_contact_urn"

// AddContactURNAction can be used to add a URN to the current contact. A [event:contact_urn_added] event
// will be created when this action is encountered. If there is no contact then this
// action will be ignored.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_contact_urn",
//     "scheme": "tel",
//     "path": "@run.results.phone_number"
//   }
//
// @action add_contact_urn
type AddContactURNAction struct {
	BaseAction
	universalAction

	Scheme string `json:"scheme" validate:"urnscheme"`
	Path   string `json:"path" validate:"required"`
}

// Type returns the type of this action
func (a *AddContactURNAction) Type() string { return TypeAddContactURN }

// Validate validates our action is valid and has all the assets it needs
func (a *AddContactURNAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs the labeling action
func (a *AddContactURNAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	// only generate event if run has a contact
	contact := run.Contact()
	if contact == nil {
		log.Add(events.NewFatalErrorEvent(fmt.Errorf("can't execute action in session without a contact")))
		return nil
	}

	evaluatedPath, err := run.EvaluateTemplateAsString(a.Path, false)

	// if we received an error, log it although it might just be a non-expression like foo@bar.com
	if err != nil {
		log.Add(events.NewErrorEvent(err))
	}

	// if we don't have a valid URN, log error
	urn, err := urns.NewURNFromParts(a.Scheme, evaluatedPath, "", "")
	if err != nil {
		log.Add(events.NewErrorEvent(fmt.Errorf("unable to add URN '%s:%s': %s", a.Scheme, evaluatedPath, err.Error())))
		return nil
	}

	if !run.Contact().HasURN(urn) {
		log.Add(events.NewURNAddedEvent(urn))
	}

	return nil
}
