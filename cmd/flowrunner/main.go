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
	source, err := static.LoadStaticSource(assetsPath)
	if err != nil {
		return nil, err
	}

	assets, err := engine.NewSessionAssets(source)
	if err != nil {
		return nil, fmt.Errorf("error parsing assets: %s", err)
	}

	httpClient := utils.NewHTTPClient("goflow-flowrunner")
	session := engine.NewSession(assets, engine.NewDefaultConfig(), httpClient)

	contact, err := flows.ReadContact(session.Assets(), json.RawMessage(contactJSON), true)
	if err != nil {
		return nil, err
	}
	contact.SetLanguage(contactLang)

	flow, err := session.Assets().Flows().Get(flowUUID)
	if err != nil {
		return nil, err
	}

	// create our environment
	la, _ := time.LoadLocation("America/Los_Angeles")
	languages := []utils.Language{flow.Language(), contact.Language()}
	env := utils.NewEnvironment(utils.DateFormatYearMonthDay, utils.TimeFormatHourMinute, la, utils.NilLanguage, languages, utils.DefaultNumberFormat, utils.RedactionPolicyNone)

	repro := &Repro{}

	if initialMsg != "" {
		msg := createMessage(contact, initialMsg)
		repro.Trigger = triggers.NewMsgTrigger(env, contact, flow.Reference(), msg, nil, time.Now())
	} else {
		repro.Trigger = triggers.NewManualTrigger(env, contact, flow.Reference(), nil, time.Now())
	}
	fmt.Fprintf(out, "Starting flow '%s'....\n---------------------------------------\n", flow.Name())

	// start our session
	if err := session.Start(repro.Trigger); err != nil {
		return nil, err
	}

	printEvents(session, out)
	scanner := bufio.NewScanner(in)

	for session.Wait() != nil {

		// ask for input
		fmt.Fprintf(out, "> ")
		scanner.Scan()

		// create our resume
		msg := createMessage(contact, scanner.Text())
		resume := resumes.NewMsgResume(nil, nil, msg)
		repro.Resumes = append(repro.Resumes, resume)

		if err := session.Resume(resume); err != nil {
			return nil, err
		}

		printEvents(session, out)
	}

	return repro, nil
}

func createMessage(contact *flows.Contact, text string) *flows.MsgIn {
	return flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), flows.NilMsgID, contact.URNs()[0].URN, nil, text, []flows.Attachment{})
}

func printEvents(session flows.Session, out io.Writer) {
	for _, event := range session.Events() {
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

type Repro struct {
	Trigger flows.Trigger  `json:"trigger"`
	Resumes []flows.Resume `json:"resumes"`
}

func (r *Repro) MarshalJSON() ([]byte, error) {
	envelope := &struct {
		Trigger *utils.TypedEnvelope   `json:"trigger"`
		Resumes []*utils.TypedEnvelope `json:"resumes"`
	}{}

	envelope.Trigger, _ = utils.EnvelopeFromTyped(r.Trigger)
	envelope.Resumes = make([]*utils.TypedEnvelope, len(r.Resumes))
	for i := range r.Resumes {
		envelope.Resumes[i], _ = utils.EnvelopeFromTyped(r.Resumes[i])
	}

	return json.Marshal(envelope)
}
