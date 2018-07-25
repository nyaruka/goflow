package transferto

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/extensions/transferto/client"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

func init() {
	actions.RegisterType(TypeTransferAirtime, func() flows.Action { return &TransferAirtimeAction{} })
}

type transferToConfig struct {
	APIToken string `json:"api_token"`
	Account  string `json:"account"`
	Currency string `json:"currency"`
	Disabled bool   `json:"disabled"`
}

// TypeTransferAirtime is the type constant for our airtime action
var TypeTransferAirtime = "transfer_airtime"

// TransferAirtimeAction attempts to make a TransferTo airtime transfer to the contact
type TransferAirtimeAction struct {
	actions.BaseAction

	Amounts map[string]decimal.Decimal `json:"amounts"`
}

// Type returns the type of this router
func (a *TransferAirtimeAction) Type() string { return TypeTransferAirtime }

// Validate validates our action is valid and has all the assets it needs
func (a *TransferAirtimeAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *TransferAirtimeAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	contact := run.Contact()
	if contact == nil {
		log.Add(events.NewFatalErrorEvent(fmt.Errorf("can't execute action in session without a contact")))
		return nil
	}

	// log error and return if we don't have a configuration
	rawConfig := run.Session().Environment().Extension("transferto")
	if rawConfig == nil {
		log.Add(events.NewErrorEvent(fmt.Errorf("missing transferto configuration")))
		log.Add(NewFailedAirtimeTransferredEvent())
		return nil
	}

	config := &transferToConfig{}
	if err := json.Unmarshal(rawConfig, config); err != nil {
		return fmt.Errorf("unable to read config: %s", err)
	}

	// if airtime transferred are disabled, return a mock event
	if config.Disabled {
		log.Add(NewAirtimeTransferredEvent(config.Currency, decimal.RequireFromString("1")))
		return nil
	}

	// check that our contact has a tel URN
	telURNs := contact.URNs().WithScheme(urns.TelScheme)
	if len(telURNs) == 0 {
		log.Add(events.NewErrorEvent(fmt.Errorf("can't transfer airtime to contact without a tel URN")))
		log.Add(NewFailedAirtimeTransferredEvent())
		return nil
	}

	amount, err := attemptTransfer(run.Session(), run.Contact().PreferredChannel(), config, a.Amounts, telURNs[0].Path())

	if err != nil {
		log.Add(events.NewErrorEvent(err))
		log.Add(NewFailedAirtimeTransferredEvent())
		return nil
	}

	log.Add(NewAirtimeTransferredEvent(config.Currency, amount))
	return nil
}

// attempts to make the transfer, returning the amount transfered or an error
func attemptTransfer(session flows.Session, channel flows.Channel, config *transferToConfig, amounts map[string]decimal.Decimal, recipient string) (decimal.Decimal, error) {
	cl := client.NewTransferToClient(config.Account, config.APIToken, session.HTTPClient())

	info, err := cl.MSISDNInfo(recipient, config.Currency, "1")
	if err != nil {
		return decimal.Zero, err
	}

	countryCode := utils.CountryCodeFromName(info.Country)
	amount, hasAmount := amounts[countryCode]
	if !hasAmount {
		return decimal.Zero, fmt.Errorf("no amount configured for transfers to %s (%s)", info.Country, countryCode)
	}

	// find the product closest to our desired amount
	var useProduct string
	for p, product := range info.ProductList {
		price := info.LocalInfoValueList[p]
		if price <= amount {
			useProduct = product
		} else {
			break
		}
	}

	// TODO complicated product/skuid stuff

	reservedID, err := cl.ReserveID()
	if err != nil {
		return decimal.Zero, err
	}

	var fromMSISDN string
	if channel != nil {
		fromMSISDN = channel.Address()
	}

	_, err = cl.Topup(reservedID, fromMSISDN, recipient, config.Currency, "", "")
	if err != nil {
		return decimal.Zero, err
	}

	return amount, nil
}
