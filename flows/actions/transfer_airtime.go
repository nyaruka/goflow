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
	registerType(TypeTransferAirtime, func() flows.Action { return &TransferAirtimeAction{} })
}

var transferCategories = []string{CategorySuccess, CategoryFailure}

// TypeTransferAirtime is the type for the transfer airtime action
const TypeTransferAirtime string = "transfer_airtime"

// TransferAirtimeAction attempts to make an airtime transfer to the contact.
//
// An [event:airtime_transferred] event will be created if the airtime could be sent.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "transfer_airtime",
//	  "amounts": {"RWF": 500, "USD": 0.5},
//	  "result_name": "Reward Transfer"
//	}
//
// @action transfer_airtime
type TransferAirtimeAction struct {
	baseAction
	onlineAction

	Amounts    map[string]decimal.Decimal `json:"amounts"     validate:"required"`
	ResultName string                     `json:"result_name" validate:"required,result_name"`
}

// NewTransferAirtime creates a new airtime transfer action
func NewTransferAirtime(uuid flows.ActionUUID, amounts map[string]decimal.Decimal, resultName string) *TransferAirtimeAction {
	return &TransferAirtimeAction{
		baseAction: newBaseAction(TypeTransferAirtime, uuid),
		Amounts:    amounts,
		ResultName: resultName,
	}
}

// Execute executes the transfer action
func (a *TransferAirtimeAction) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	transfer, err := a.transfer(ctx, run, logEvent)
	if err != nil {
		logEvent(events.NewError(err.Error()))

		a.saveFailure(run, step, logEvent)
	} else {
		a.saveSuccess(run, step, transfer, logEvent)
	}

	return nil
}

func (a *TransferAirtimeAction) transfer(ctx context.Context, run flows.Run, logEvent flows.EventCallback) (*flows.AirtimeTransfer, error) {
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
	if transfer != nil {
		logEvent(events.NewAirtimeTransferred(transfer, httpLogger.Logs))
	}

	return transfer, err
}

func (a *TransferAirtimeAction) saveSuccess(run flows.Run, step flows.Step, transfer *flows.AirtimeTransfer, logEvent flows.EventCallback) {
	a.saveResult(run, step, a.ResultName, transfer.ExternalID, CategorySuccess, "", "", nil, logEvent)
}

func (a *TransferAirtimeAction) saveFailure(run flows.Run, step flows.Step, logEvent flows.EventCallback) {
	a.saveResult(run, step, a.ResultName, "", CategoryFailure, "", "", nil, logEvent)
}

func (a *TransferAirtimeAction) Inspect(result func(*flows.ResultInfo), dependency func(assets.Reference)) {
	if a.ResultName != "" {
		result(flows.NewResultInfo(a.ResultName, transferCategories))
	}
}
