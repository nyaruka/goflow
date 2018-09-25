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
	Login    string `json:"login"`
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

// AllowedFlowTypes returns the flow types which this action is allowed to occur in
func (a *TransferAirtimeAction) AllowedFlowTypes() []flows.FlowType {
	return []flows.FlowType{flows.FlowTypeMessaging, flows.FlowTypeVoice}
}

// Execute runs this action
func (a *TransferAirtimeAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	contact := run.Contact()
	if contact == nil {
		log.Add(events.NewErrorEvent(fmt.Errorf("can't execute action in session without a contact")))
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

	currency, amount, err := attemptTransfer(run.Contact().PreferredChannel(), config, a.Amounts, telURNs[0].Path(), run.Session().HTTPClient())

	if err != nil {
		log.Add(events.NewErrorEvent(err))
		log.Add(NewFailedAirtimeTransferredEvent())
		return nil
	}

	log.Add(NewAirtimeTransferredEvent(currency, amount))
	return nil
}

// attempts to make the transfer, returning the actual amount transfered or an error
func attemptTransfer(channel *flows.Channel, config *transferToConfig, amounts map[string]decimal.Decimal, recipient string, httpClient *utils.HTTPClient) (string, decimal.Decimal, error) {
	cl := client.NewTransferToClient(config.Login, config.APIToken, httpClient)

	info, err := cl.MSISDNInfo(recipient, config.Currency, "1")
	if err != nil {
		return "", decimal.Zero, err
	}

	// look up the amount to send in this currency
	amount, hasAmount := amounts[info.DestinationCurrency]
	if !hasAmount {
		return "", decimal.Zero, fmt.Errorf("no amount configured for transfers in %s", info.DestinationCurrency)
	}

	if info.OpenRange {
		// TODO add support for open-range topups once we can find numbers to test this with
		// see https://shop.transferto.com/shop/v3/doc/TransferTo_API_OR.pdf
		return "", decimal.Zero, fmt.Errorf("transferto account is configured for open-range which is not yet supported")
	}

	// find the product closest to our desired amount
	var useProduct string
	useAmount := decimal.Zero
	for p, product := range info.ProductList {
		price := info.LocalInfoValueList[p]
		if price.GreaterThan(useAmount) && price.LessThanOrEqual(amount) {
			useProduct = product
			useAmount = price
		}
	}

	reservedID, err := cl.ReserveID()
	if err != nil {
		return "", decimal.Zero, err
	}

	var fromMSISDN string
	if channel != nil {
		fromMSISDN = channel.Address()
	}

	topup, err := cl.Topup(reservedID, fromMSISDN, recipient, useProduct, "")
	if err != nil {
		return "", decimal.Zero, err
	}

	return topup.DestinationCurrency, topup.ActualProductSent, nil
}
