package actions

import (
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"
	"github.com/nyaruka/goflow/flows/events"

	"github.com/pkg/errors"
)

func init() {
	registerType(TypeAddContactURN, func() flows.Action { return &AddContactURNAction{} })
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
//     "path": "@results.phone_number.value"
//   }
//
// @action add_contact_urn
type AddContactURNAction struct {
	baseAction
	universalAction

	Scheme string `json:"scheme" validate:"urnscheme"`
	Path   string `json:"path" validate:"required" engine:"evaluated"`
}

// NewAddContactURN creates a new add URN action
func NewAddContactURN(uuid flows.ActionUUID, scheme string, path string) *AddContactURNAction {
	return &AddContactURNAction{
		baseAction: newBaseAction(TypeAddContactURN, uuid),
		Scheme:     scheme,
		Path:       path,
	}
}

// Execute runs the labeling action
func (a *AddContactURNAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	// only generate event if run has a contact
	contact := run.Contact()
	if contact == nil {
		logEvent(events.NewErrorf("can't execute action in session without a contact"))
		return nil
	}

	evaluatedPath, err := run.EvaluateTemplate(a.Path)

	// if we received an error, log it although it might just be a non-expression like foo@bar.com
	if err != nil {
		logEvent(events.NewError(err))
	}

	evaluatedPath = strings.TrimSpace(evaluatedPath)
	if evaluatedPath == "" {
		logEvent(events.NewErrorf("can't add URN with empty path"))
		return nil
	}

	// if we don't have a valid URN, log error
	urn, err := urns.NewURNFromParts(a.Scheme, evaluatedPath, "", "")
	if err != nil {
		logEvent(events.NewError(errors.Wrapf(err, "unable to add URN '%s:%s'", a.Scheme, evaluatedPath)))
		return nil
	}

	a.applyModifier(run, modifiers.NewURNs([]urns.URN{urn}, modifiers.URNsAppend), logModifier, logEvent)
	return nil
}
