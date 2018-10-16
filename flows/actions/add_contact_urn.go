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
// will be created when this action is encountered.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_contact_urn",
//     "scheme": "tel",
//     "path": "@results.phone_number"
//   }
//
// @action add_contact_urn
type AddContactURNAction struct {
	BaseAction
	universalAction

	Scheme string `json:"scheme" validate:"urnscheme"`
	Path   string `json:"path" validate:"required"`
}

// NewAddContactURNAction creates a new add URN action
func NewAddContactURNAction(uuid flows.ActionUUID, scheme string, path string) *AddContactURNAction {
	return &AddContactURNAction{
		BaseAction: NewBaseAction(TypeAddContactURN, uuid),
		Scheme:     scheme,
		Path:       path,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *AddContactURNAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute runs the labeling action
func (a *AddContactURNAction) Execute(run flows.FlowRun, step flows.Step) error {
	// only generate event if run has a contact
	contact := run.Contact()
	if contact == nil {
		a.logError(run, step, fmt.Errorf("can't execute action in session without a contact"))
		return nil
	}

	evaluatedPath, err := run.EvaluateTemplateAsString(a.Path, false)

	// if we received an error, log it although it might just be a non-expression like foo@bar.com
	if err != nil {
		a.logError(run, step, err)
	}

	// if we don't have a valid URN, log error
	urn, err := urns.NewURNFromParts(a.Scheme, evaluatedPath, "", "")
	if err != nil {
		a.logError(run, step, fmt.Errorf("unable to add URN '%s:%s': %s", a.Scheme, evaluatedPath, err.Error()))
		return nil
	}

	if !run.Contact().HasURN(urn) {
		run.Contact().AddURN(urn)
		a.log(run, step, events.NewURNAddedEvent(urn))

		a.reevaluateDynamicGroups(run, step)
	}

	return nil
}
