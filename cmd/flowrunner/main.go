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

	"github.com/pkg/errors"
)

var contactJSON = `{
	"uuid": "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3",
	"id": 1234567,
	"name": "Ben Haggerty",
	"created_on": "2018-01-01T12:00:00.000000000-00:00",
	"fields": {},
	"timezone": "America/Guayaquil",
	"urns": [
		"tel:+12065551212",
		"facebook:1122334455667788",
		"mailto:ben@macklemore"
	]
}
`

var usage = `usage: flowrunner [flags] <assets.json> <flow_uuid>`

func main() {
	var initialMsg, contactLang string
	var printRepro bool
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&initialMsg, "msg", "", "initial message to trigger session with")
	flags.StringVar(&contactLang, "lang", "eng", "initial language of the contact")
	flags.BoolVar(&printRepro, "repro", false, "print repro afterwards")
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) != 2 {
		fmt.Println(usage)
		flags.PrintDefaults()
		os.Exit(1)
	}

	assetsPath := args[0]
	flowUUID := assets.FlowUUID(args[1])

	repro, err := RunFlow(assetsPath, flowUUID, initialMsg, utils.Language(contactLang), os.Stdin, os.Stdout)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if printRepro {
		fmt.Println("---------------------------------------")
		marshaledRepro, _ := utils.JSONMarshalPretty(repro)
		fmt.Println(string(marshaledRepro))
	}
}

// RunFlow steps through a flow
func RunFlow(assetsPath string, flowUUID assets.FlowUUID, initialMsg string, contactLang utils.Language, in io.Reader, out io.Writer) (*Repro, error) {
	source, err := static.LoadSource(assetsPath)
	if err != nil {
		return nil, err
	}

	sessionAssets, err := engine.NewSessionAssets(source)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing assets")
	}

	contact, err := flows.ReadContact(sessionAssets, json.RawMessage(contactJSON), assets.PanicOnMissing)
	if err != nil {
		return nil, err
	}
	contact.SetLanguage(contactLang)

	flow, err := sessionAssets.Flows().Get(flowUUID)
	if err != nil {
		return nil, err
	}

	// create our environment
	la, _ := time.LoadLocation("America/Los_Angeles")
	languages := []utils.Language{flow.Language(), contact.Language()}
	env := utils.NewEnvironmentBuilder().WithTimezone(la).WithAllowedLanguages(languages).Build()

	repro := &Repro{}

	if initialMsg != "" {
		msg := createMessage(contact, initialMsg)
		repro.Trigger = triggers.NewMsgTrigger(env, flow.Reference(), contact, msg, nil)
	} else {
		repro.Trigger = triggers.NewManualTrigger(env, flow.Reference(), contact, nil)
	}
	fmt.Fprintf(out, "Starting flow '%s'....\n---------------------------------------\n", flow.Name())

	eng := engine.NewBuilder().WithDefaultUserAgent("goflow-flowrunner").Build()
	session := eng.NewSession(sessionAssets)

	// start our session
	sprint, err := session.Start(repro.Trigger)
	if err != nil {
		return nil, err
	}

	printEvents(sprint.Events(), out)
	scanner := bufio.NewScanner(in)

	for session.Wait() != nil {

		// ask for input
		fmt.Fprintf(out, "> ")
		scanner.Scan()

		// create our resume
		msg := createMessage(contact, scanner.Text())
		resume := resumes.NewMsgResume(nil, nil, msg)
		repro.Resumes = append(repro.Resumes, resume)

		sprint, err := session.Resume(resume)
		if err != nil {
			return nil, err
		}

		printEvents(sprint.Events(), out)
	}

	return repro, nil
}

func createMessage(contact *flows.Contact, text string) *flows.MsgIn {
	return flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), contact.URNs()[0].URN(), nil, text, []flows.Attachment{})
}

func printEvents(log []flows.Event, out io.Writer) {
	for _, event := range log {
		var msg string
		switch typed := event.(type) {
		case *events.ContactNameChangedEvent:
			msg = fmt.Sprintf("📛 name changed to %s", typed.Name)
		case *events.ContactLanguageChangedEvent:
			msg = fmt.Sprintf("🌐 language changed to %s", typed.Language)
		case *events.ContactTimezoneChangedEvent:
			msg = fmt.Sprintf("🕑 timezone changed to %s", typed.Timezone)
		case *events.ErrorEvent:
			msg = fmt.Sprintf("⚠️ %s", typed.Text)
		case *events.MsgCreatedEvent:
			msg = fmt.Sprintf("💬 \"%s\"", typed.Msg.Text())
		case *events.MsgReceivedEvent:
			msg = fmt.Sprintf("📥 received message '%s'", typed.Msg.Text())
		case *events.MsgWaitEvent:
			msg = fmt.Sprintf("⏳ waiting for message....")
		case *events.RunResultChangedEvent:
			msg = fmt.Sprintf("📈 run result '%s' changed to '%s'", typed.Name, typed.Value)
		default:
			msg = fmt.Sprintf("❓ %s event", typed.Type())
		}

		fmt.Fprintln(out, msg)
	}
}

// Repro describes the trigger and resumes needed to reproduce this session
type Repro struct {
	Trigger flows.Trigger  `json:"trigger"`
	Resumes []flows.Resume `json:"resumes"`
}
