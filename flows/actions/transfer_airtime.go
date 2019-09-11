package actions

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/pkg/errors"

	"github.com/shopspring/decimal"
)

func init() {
	registerType(TypeTransferAirtime, func() flows.Action { return &TransferAirtimeAction{} })
}

var statusCategories = map[flows.AirtimeTransferStatus]string{
	flows.AirtimeTransferStatusSuccess: "Success",
	flows.AirtimeTransferStatusFailed:  "Failure",
}

// TypeTransferAirtime is the type for the transfer airtime action
const TypeTransferAirtime string = "transfer_airtime"

// TransferAirtimeAction attempts to make an airtime transfer to the contact.
//
// An [event:email_created] event will be created for each email address.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "transfer_airtime",
//     "amounts": {"RWF": 500, "USD": 0.5},
//     "result_name": "reward_transfer"
//   }
//
// @action transfer_airtime
type TransferAirtimeAction struct {
	baseAction
	onlineAction

	Amounts    map[string]decimal.Decimal `json:"amounts"`
	ResultName string                     `json:"result_name,omitempty"`
}

// NewTransferAirtime creates a new airtime transfer action
func NewTransferAirtime(uuid flows.ActionUUID, amounts map[string]decimal.Decimal, resultName string) *TransferAirtimeAction {
	return &TransferAirtimeAction{
		baseAction: newBaseAction(TypeTransferAirtime, uuid),
		Amounts:    amounts,
		ResultName: resultName,
	}
}

// Execute creates the email events
func (a *TransferAirtimeAction) Execute(run flows.FlowRun, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	contact := run.Contact()
	if contact == nil {
		logEvent(events.NewError(errors.Errorf("can't execute action in session without a contact")))
		return nil
	}

	// check that our contact has a tel URN
	telURNs := contact.URNs().WithScheme(urns.TelScheme)
	if len(telURNs) == 0 {
		logEvent(events.NewError(errors.Errorf("can't transfer airtime to contact without a tel URN")))
		return nil
	}
	recipient := telURNs[0].URN()

	// if contact's preferred channel is tel, use that as the sender
	var sender urns.URN
	channel := contact.PreferredChannel()
	if channel != nil && channel.SupportsScheme(urns.TelScheme) {
		sender, _ = urns.Parse("tel:" + channel.Address())
	}

	airtimeSvc := run.Session().Engine().Services().Airtime(run.Session())

	transfer, err := airtimeSvc.Transfer(run.Session(), sender, recipient, a.Amounts)
	if err != nil {
		// an error without a transfer is considered a failure because we have nothing to route on
		if transfer == nil {
			a.fail(run, err, logEvent)
			return nil
		}
		logEvent(events.NewError(err))
	}

	logEvent(events.NewAirtimeTransferred(transfer))

	if a.ResultName != "" {
		value := transfer.ActualAmount.String()
		category := statusCategories[transfer.Status]

		a.saveResult(run, step, a.ResultName, value, category, "", "", nil, logEvent)
	}
	return nil
}
