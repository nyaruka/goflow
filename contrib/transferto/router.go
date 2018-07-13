package transferto

import (
	"fmt"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers"
)

func init() {
	routers.RegisterType(TypeAirtimeRouter, func() flows.Router { return &AirtimeRouter{} })
}

type transferToConfig struct {
	APIToken string `mapstructure:"transferto.api_token"`
	Account  string `mapstructure:"transferto.account"`
	Currency string `mapstructure:"transferto.currency"`
	Disabled bool   `mapstructure:"transferto.disabled"`
}

// TypeAirtimeRouter is the type constant for our airtime router
var TypeAirtimeRouter = "airtime"

// AirtimeRouter attempts to make a TransferTo airtime transfer to the contact. If it succeeds it will take
// the first exit, otherwise the second exit.
type AirtimeRouter struct {
	routers.BaseRouter
}

// Type returns the type of this router
func (r *AirtimeRouter) Type() string { return TypeAirtimeRouter }

// Validate validates the arguments for this router
func (r *AirtimeRouter) Validate(exits []flows.Exit) error {
	if len(exits) != 2 {
		return fmt.Errorf("airtime router requires 2 exits, have %d", len(exits))
	}
	return nil
}

func (r *AirtimeRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (*string, flows.Route, error) {
	config := &transferToConfig{}
	if err := run.Session().EngineConfig().ReadInto(config); err != nil {
		return nil, flows.NoRoute, fmt.Errorf("unable to read transferto config: %s", err)
	}

	return nil, flows.NoRoute, nil
}
