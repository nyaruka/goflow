package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/services/airtime/dtone"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const usage = `usage: transferairtime [flags] <destnumber> <amount> <currency>`

var verbose bool

func main() {
	var dtoneKey, dtoneSecret string
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&dtoneKey, "dtone.key", "", "API key for DTOne service")
	flags.StringVar(&dtoneSecret, "dtone.secret", "", "API secret for DTOne service")
	flags.BoolVar(&verbose, "v", false, "enable verbose logging")
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) != 3 {
		fmt.Println(usage)
		flags.PrintDefaults()
		os.Exit(1)
	}

	destination, err := urns.NewTelURNForCountry(args[0], "")
	if err != nil {
		fmt.Printf("%s isn't a valid phone number\n", args[0])
		os.Exit(1)
	}

	amount, err := decimal.NewFromString(args[1])
	if err != nil {
		fmt.Printf("%s isn't a valid amount\n", args[1])
		os.Exit(1)
	}

	if dtoneKey == "" || dtoneSecret == "" {
		fmt.Println("no airtime service credentials provided")
		os.Exit(1)
	}

	httpx.SetDebug(verbose)

	svcFactory := func(flows.Session) (flows.AirtimeService, error) {
		return dtone.NewService(http.DefaultClient, nil, dtoneKey, dtoneSecret), nil
	}

	if err := transferAirtime(destination, amount, args[2], svcFactory); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const assetsTemplate = `
{
	"flows": [
		{
            "uuid": "2374f60d-7412-442c-9177-585967afa972",
            "name": "Airtime",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": [
				{
					"uuid": "f365a3fb-cc94-4f93-9263-f3dc9d94c0d3",
					"actions": [
						{
							"uuid": "ad1b1c41-553a-413d-9619-813cdc578933",
							"type": "transfer_airtime",
							"amounts": {"%s": %s},
							"result_name": "Transfer"
						}
					],
					"exits": [
						{
							"uuid": "d3add939-33e3-43f9-bb27-0e7100830c8a"
						}
					]
				}
			]
        }
	],
	"channels": []
}
`

func transferAirtime(destination urns.URN, amount decimal.Decimal, currency string, svcFactory engine.AirtimeServiceFactory) error {
	// create a flow to do make the transfer
	assetsJSON := fmt.Sprintf(assetsTemplate, currency, amount.String())
	source, err := static.NewSource([]byte(assetsJSON))
	if err != nil {
		return err
	}

	env := envs.NewBuilder().Build()

	sa, err := engine.NewSessionAssets(env, source, nil)
	if err != nil {
		return errors.Wrap(err, "error parsing assets")
	}

	eng := engine.NewBuilder().WithAirtimeServiceFactory(svcFactory).Build()
	contact := flows.NewEmptyContact(sa, "", "", nil)
	contact.AddURN(destination, nil)

	trigger := triggers.NewBuilder(env, assets.NewFlowReference(assets.FlowUUID("2374f60d-7412-442c-9177-585967afa972"), "Airtime"), contact).Manual().Build()

	_, sprint, err := eng.NewSession(sa, trigger)
	if err != nil {
		return err
	}

	for _, event := range sprint.Events() {
		marshaled, _ := jsonx.Marshal(event)
		fmt.Println(string(marshaled))
	}

	return nil
}
