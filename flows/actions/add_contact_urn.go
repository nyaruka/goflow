package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeAddContactURN, func() flows.Action { return &AddContactURN{} })
}

// TypeAddContactURN is our type for the add URN action
const TypeAddContactURN string = "add_contact_urn"

// AddContactURN can be used to add a URN to the current contact. A [event:contact_urns_changed] event
// will be created when this action is encountered.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "add_contact_urn",
//	  "scheme": "tel",
//	  "path": "@results.phone_number.value"
//	}
//
// @action add_contact_urn
type AddContactURN struct {
	baseAction
	universalAction

	Scheme string `json:"scheme" validate:"urnscheme"`
	Path   string `json:"path" validate:"required" engine:"evaluated"`
}

// NewAddContactURN creates a new add URN action
func NewAddContactURN(uuid flows.ActionUUID, scheme string, path string) *AddContactURN {
	return &AddContactURN{
		baseAction: newBaseAction(TypeAddContactURN, uuid),
		Scheme:     scheme,
		Path:       path,
	}
}

// Execute runs the labeling action
func (a *AddContactURN) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	evaluatedPath, _ := run.EvaluateTemplate(a.Path, log)
	evaluatedPath = strings.TrimSpace(evaluatedPath)
	if evaluatedPath == "" {
		log(events.NewError("can't add URN with empty path"))
		return nil
	}

	// create URN - modifier will take care of validating it
	urn := urns.URN(fmt.Sprintf("%s:%s", a.Scheme, evaluatedPath))

	_, err := a.applyModifier(ctx, run, modifiers.NewURNs([]urns.URN{urn}, modifiers.URNsAppend), log)
	return err
}
