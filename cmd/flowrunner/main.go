package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	_ "github.com/nyaruka/goflow/extensions/transferto"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
)

var contactJSON = `{
	"uuid": "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3",
	"id": 1234567,
	"name": "Ben Haggerty",
	"created_on": "2018-01-01T12:00:00.000000000-00:00",
	"fields": {
		"first_name": {
			"text": "Ben"
		}
	},
	"language": "eng",
	"timezone": "America/Guayaquil",
	"urns": [
		"tel:+12065551212",
		"facebook:1122334455667788",
		"mailto:ben@macklemore"
	]
}
`

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Printf("\nUsage: flowrunner <assets.json> <flow_uuid>\n\n")
		os.Exit(1)
	}

	assetsPath := flag.Args()[0]
	flowUUID := assets.FlowUUID(flag.Args()[1])

	if err := RunFlow(assetsPath, flowUUID, os.Stdin, os.Stdout); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// RunFlow steps through a flow
func RunFlow(assetsPath string, flowUUID assets.FlowUUID, in io.Reader, out io.Writer) error {
	source, err := static.LoadStaticSource(assetsPath)
	if err != nil {
		return err
	}

	// create our environment
	la, _ := time.LoadLocation("America/Los_Angeles")
	env := utils.NewEnvironment(utils.DateFormatYearMonthDay, utils.TimeFormatHourMinute, la, utils.NilLanguage, nil, utils.DefaultNumberFormat, utils.RedactionPolicyNone)

	assets, err := engine.NewSessionAssets(source)
	if err != nil {
		return fmt.Errorf("error parsing assets: %s", err)
	}

	httpClient := utils.NewHTTPClient("goflow-flowrunner")
	session := engine.NewSession(assets, engine.NewDefaultConfig(), httpClient)

	contact, err := flows.ReadContact(session.Assets(), json.RawMessage(contactJSON), true)
	if err != nil {
		return err
	}
	flow, err := session.Assets().Flows().Get(flowUUID)
	if err != nil {
		return err
	}

	trigger := triggers.NewManualTrigger(env, contact, flow.Reference(), nil, time.Now())
	fmt.Fprintf(out, "Starting flow '%s'....\n", flow.Name())

	// start our session
	if err := session.Start(trigger); err != nil {
		return err
	}

	printEvents(session, out)
	scanner := bufio.NewScanner(in)

	for session.Wait() != nil {

		// ask for input
		fmt.Fprintf(out, "> ")
		scanner.Scan()

		// create our resume
		msg := flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), flows.NilMsgID, contact.URNs()[0].URN, nil, scanner.Text(), []flows.Attachment{})
		resume := resumes.NewMsgResume(nil, nil, msg)

		if err := session.Resume(resume); err != nil {
			return err
		}

		printEvents(session, out)
	}
	return nil
}

func printEvents(session flows.Session, out io.Writer) {
	// print any msg_created events
	for _, event := range session.Events() {
		if event.Type() == events.TypeMsgCreated {
			fmt.Fprintf(out, "💬 %s\n", event.(*events.MsgCreatedEvent).Msg.Text())
		}
	}
}
