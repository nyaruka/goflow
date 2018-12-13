package actions

import (
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"

	"github.com/pkg/errors"
)

func init() {
	RegisterType(TypeAddContactURN, func() flows.Action { return &AddContactURNAction{} })
}

// TypeAddContactURN is our type for the add URN action
const TypeAddContactURN string = "add_contact_urn"

// AddContactURNAction can be used to add a URN to the current contact. A [event:contact_urns_changed] event
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
		a.logError(run, step, errors.Errorf("can't execute action in session without a contact"))
		return nil
	}

	evaluatedPath, err := run.EvaluateTemplateAsString(a.Path)

	// if we received an error, log it although it might just be a non-expression like foo@bar.com
	if err != nil {
		a.logError(run, step, err)
	}

	evaluatedPath = strings.TrimSpace(evaluatedPath)
	if evaluatedPath == "" {
		a.logError(run, step, errors.Errorf("can't add URN with empty path"))
		return nil
	}

	// if we don't have a valid URN, log error
	urn, err := urns.NewURNFromParts(a.Scheme, evaluatedPath, "", "")
	if err != nil {
		a.logError(run, step, errors.Wrapf(err, "unable to add URN '%s:%s'", a.Scheme, evaluatedPath))
		return nil
	}

	a.applyModifier(run, step, modifiers.NewURNModifier(urn, modifiers.URNAppend))
	return nil
}
