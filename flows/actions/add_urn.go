package actions

import (
	"fmt"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeAddURN is our type for add URN actions
const TypeAddURN string = "add_urn"

// AddURNAction can be used to add a URN to the current contact. An `add_urn` event
// will be created when this action is encountered. If there is no contact then this
// action will be ignored.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "add_urn",
//     "scheme": "tel",
//     "path": "@flow.phone_number"
//   }
// ```
//
// @action add_urn
type AddURNAction struct {
	BaseAction
	Scheme string `json:"scheme" validate:"urnscheme"`
	Path   string `json:"path" validate:"required"`
}

// Type returns the type of this action
func (a *AddURNAction) Type() string { return TypeAddURN }

// Validate validates our action is valid and has all the assets it needs
func (a *AddURNAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs the labeling action
func (a *AddURNAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	// only generate event if run has a contact
	contact := run.Contact()
	if contact == nil {
		return nil
	}

	evaluatedPath, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), a.Path, false)

	// if we received an error, log it
	if err != nil {
		log.Add(events.NewErrorEvent(err))
		return nil
	}

	urn := urns.NewURNFromParts(a.Scheme, evaluatedPath, "").Normalize("")

	// if we don't have a valid URN, log error
	if !urn.Validate() {
		log.Add(events.NewErrorEvent(fmt.Errorf("invalid URN: '%s'", string(urn))))
		return nil
	}

	log.Add(events.NewAddURNEvent(urn))
	return nil
}
