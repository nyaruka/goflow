package actions

import (
	"context"
	"errors"
	"github.com/nyaruka/goflow/core"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/events"
	"github.com/nyaruka/goflow/flows"
	"github.com/shopspring/decimal"
)

func init() {
	registerType(TypeTransferAirtime, func() flows.Action { return &TransferAirtime{} })
}

const (
	// TypeTransferAirtime is the type for the transfer airtime action
	TypeTransferAirtime string = "transfer_airtime"

	TransferAirtimeOutputLocal = "_new_transfer"
)

// TransferAirtime attempts to make an airtime transfer to the contact.
//
// An [event:airtime_created] event will be created if the airtime transfer could be initiated.
// The action sets a `_new_transfer` local to the UUID of the airtime_created event when the
// transfer is initiated, and to an empty string otherwise.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "transfer_airtime",
//	  "amounts": {"RWF": 500, "USD": 0.5}
//	}
//
// @action transfer_airtime
type TransferAirtime struct {
	baseAction
	onlineAction

	Amounts map[string]decimal.Decimal `json:"amounts" validate:"required"`
}

// NewTransferAirtime creates a new airtime transfer action
func NewTransferAirtime(uuid flows.ActionUUID, amounts map[string]decimal.Decimal) *TransferAirtime {
	return &TransferAirtime{
		baseAction: newBaseAction(TypeTransferAirtime, uuid),
		Amounts:    amounts,
	}
}

// Execute executes the transfer action
func (a *TransferAirtime) Execute(ctx context.Context, run flows.Run, step flows.Step, log events.EventLogger) error {
	airtime, err := a.transfer(ctx, run, log)
	if err != nil {
		log(events.NewRawError(err))
	}

	if airtime != nil {
		run.Locals().Set(TransferAirtimeOutputLocal, string(airtime.UUID()))
	} else {
		run.Locals().Set(TransferAirtimeOutputLocal, "")
	}

	return nil
}

func (a *TransferAirtime) transfer(ctx context.Context, run flows.Run, log events.EventLogger) (*events.AirtimeCreated, error) {
	// fail if we don't have a contact
	contact := run.Contact()

	// fail if the contact doesn't have a phone URN or whatsap URN
	telURNs := contact.URNs().WithScheme(urns.Phone.Prefix, urns.WhatsApp.Prefix)
	if len(telURNs) == 0 {
		return nil, errors.New("can't transfer airtime to contact without a phone number")
	}

	recipient := telURNs[0]

	// if contact's preferred channel is a phone number, use that as the sender
	var sender urns.URN
	channel := contact.PreferredChannel()
	if channel != nil && channel.SupportsScheme(recipient.Scheme) {
		sender, _ = urns.Parse(recipient.Scheme + ":" + channel.Address())
	}

	svc, err := run.Session().Engine().Services().Airtime(run.Session().Assets())
	if err != nil {
		return nil, err
	}

	httpLogger := &core.HTTPLogger{}

	// pre-allocate the event UUID so the service can pass it through to the provider as a stable reference
	transferUUID := events.NewEventUUID()

	transfer, err := svc.Create(ctx, transferUUID, sender, recipient.Identity(), a.Amounts, httpLogger.Log)
	if err != nil {
		return nil, err
	}

	evt := events.NewAirtimeCreated(transferUUID, transfer, httpLogger.Logs)
	log(evt)

	return evt, nil
}

func (a *TransferAirtime) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	local(TransferAirtimeOutputLocal)
}
