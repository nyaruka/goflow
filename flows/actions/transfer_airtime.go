package actions

import (
	"context"
	"errors"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
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
// An [event:airtime_transferred] event will be created if the airtime could be sent.
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
func (a *TransferAirtime) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	transfer, err := a.transfer(ctx, run, logEvent)
	if err != nil {
		logEvent(events.NewError(err.Error()))
	}

	if transfer != nil {
		run.Locals().Set(TransferAirtimeOutputLocal, transfer.ExternalID)
	} else {
		run.Locals().Set(TransferAirtimeOutputLocal, "")
	}

	return nil
}

func (a *TransferAirtime) transfer(ctx context.Context, run flows.Run, logEvent flows.EventCallback) (*flows.AirtimeTransfer, error) {
	// fail if we don't have a contact
	contact := run.Contact()

	// fail if the contact doesn't have a phone URN or whatsap URN
	telURNs := contact.URNs().WithScheme(urns.Phone.Prefix, urns.WhatsApp.Prefix)
	if len(telURNs) == 0 {
		return nil, errors.New("can't transfer airtime to contact without a phone number")
	}

	recipient := telURNs[0].URN()

	// if contact's preferred channel is a phone number, use that as the sender
	var sender urns.URN
	channel := contact.PreferredChannel()
	if channel != nil && channel.SupportsScheme(recipient.Scheme()) {
		sender, _ = urns.Parse(recipient.Scheme() + ":" + channel.Address())
	}

	svc, err := run.Session().Engine().Services().Airtime(run.Session().Assets())
	if err != nil {
		return nil, err
	}

	httpLogger := &flows.HTTPLogger{}

	transfer, err := svc.Transfer(ctx, sender, recipient, a.Amounts, httpLogger.Log)
	if transfer != nil { // can be non-nil for failed transfer
		logEvent(events.NewAirtimeTransferred(transfer, httpLogger.Logs))
	}
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (a *TransferAirtime) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	local(TransferAirtimeOutputLocal)
}
