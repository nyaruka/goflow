package transferto

import (
	"fmt"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
)

func init() {
	actions.RegisterType(TypeTransferAirtime, func() flows.Action { return &TransferAirtimeAction{} })
}

type transferToConfig struct {
	APIToken string `mapstructure:"transferto.api_token"`
	Account  string `mapstructure:"transferto.account"`
	Currency string `mapstructure:"transferto.currency"`
	Disabled bool   `mapstructure:"transferto.disabled"`
}

// TypeTransferAirtime is the type constant for our airtime action
var TypeTransferAirtime = "transfer_airtime"

// TransferAirtimeAction attempts to make a TransferTo airtime transfer to the contact
type TransferAirtimeAction struct {
	actions.BaseAction
}

// Type returns the type of this router
func (a *TransferAirtimeAction) Type() string { return TypeTransferAirtime }

// Validate validates our action is valid and has all the assets it needs
func (a *TransferAirtimeAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *TransferAirtimeAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	config := &transferToConfig{}
	if err := run.Session().EngineConfig().ReadInto(config); err != nil {
		return fmt.Errorf("unable to read transferto config: %s", err)
	}

	// TODO

	log.Add(NewAirtimeTransferedEvent("RWF", 100, "success"))
	return nil
}
