package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeAddContactURN, func() flows.Action { return &AddContactURN{} })
}

const (
	// TypeAddContactURN is our type for the add URN action
	TypeAddContactURN string = "add_contact_urn"

	// AddURNOutputLocal receives the identity of the URN if it was added or if it was already present
	AddURNOutputLocal = "_has_urn"
)

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
		log(events.NewError("Can't add URN with empty path", ""))
		return nil
	}

	// create URN - modifier will take care of validating it
	urn := urns.URN(fmt.Sprintf("%s:%s", a.Scheme, evaluatedPath))
	urn = urn.Normalize()

	_, err := a.applyModifier(ctx, run, modifiers.NewURNs([]urns.URN{urn}, modifiers.URNsAppend), log)
	if run.Contact().HasURN(urn) {
		run.Locals().Set(AddURNOutputLocal, string(urn))
	} else {
		run.Locals().Set(AddURNOutputLocal, "")
	}

	return err
}

func (a *AddContactURN) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	local(AddURNOutputLocal)
}
