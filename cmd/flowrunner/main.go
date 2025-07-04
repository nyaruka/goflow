package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/utils"
)

const contactJSON = `{
	"uuid": "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3",
	"id": 1234567,
	"name": "Ben Haggerty",
	"status": "active",
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

const usage = `usage: flowrunner [flags] <assets.json> [flow_uuid]`

func main() {
	var initialMsg, contactLang, witToken string
	var printRepro bool
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&initialMsg, "msg", "", "initial message to trigger session with")
	flags.StringVar(&contactLang, "lang", "eng", "initial language of the contact")
	flags.StringVar(&witToken, "wit.token", "", "access token for wit.ai")
	flags.BoolVar(&printRepro, "repro", false, "print repro afterwards")
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if !(len(args) == 1 || len(args) == 2) {
		fmt.Println(usage)
		flags.PrintDefaults()
		os.Exit(1)
	}

	assetsPath := args[0]
	var flowUUID assets.FlowUUID
	if len(args) == 2 {
		flowUUID = assets.FlowUUID(args[1])
	}

	engine := createEngine(witToken)

	repro, err := RunFlow(engine, assetsPath, flowUUID, initialMsg, i18n.Language(contactLang), os.Stdin, os.Stdout)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if printRepro {
		fmt.Println("---------------------------------------")
		marshaledRepro, _ := jsonx.MarshalPretty(repro)
		fmt.Println(string(marshaledRepro))
	}
}

func createEngine(witToken string) flows.Engine {
	builder := engine.NewBuilder().
		WithWebhookServiceFactory(webhooks.NewServiceFactory(http.DefaultClient, nil, nil, map[string]string{"User-Agent": "goflow-runner"}, 10000))

	if witToken != "" {
		builder.WithClassificationServiceFactory(func(classifier *flows.Classifier) (flows.ClassificationService, error) {
			if classifier.Type() == "wit" {
				return wit.NewService(http.DefaultClient, nil, classifier, witToken), nil
			}
			return nil, errors.New("only classifiers of type wit supported")
		})
	}

	return builder.Build()
}

// RunFlow steps through a flow
func RunFlow(eng flows.Engine, assetsPath string, flowUUID assets.FlowUUID, initialMsg string, contactLang i18n.Language, in io.Reader, out io.Writer) (*Repro, error) {
	ctx := context.Background()

	assetsJSON, err := os.ReadFile(assetsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading assets file '%s': %w", assetsPath, err)
	}

	// if user didn't provide a flow UUID, look for the UUID of the first flow
	if flowUUID == "" {
		uuidBytes, _, _, err := jsonparser.Get(assetsJSON, "flows", "[0]", "uuid")
		if err != nil {
			return nil, errors.New("no flows found in assets file")
		}
		flowUUID = assets.FlowUUID(uuidBytes)
	}

	source, err := static.NewSource(assetsJSON)
	if err != nil {
		return nil, err
	}

	sa, err := engine.NewSessionAssets(envs.NewBuilder().Build(), source, nil)
	if err != nil {
		return nil, fmt.Errorf("error parsing assets: %w", err)
	}

	flow, err := sa.Flows().Get(assets.FlowUUID(flowUUID))
	if err != nil {
		return nil, err
	}

	contact, err := flows.ReadContact(sa, []byte(contactJSON), assets.PanicOnMissing)
	if err != nil {
		return nil, err
	}
	contact.SetLanguage(contactLang)

	// create our environment
	la, _ := time.LoadLocation("America/Los_Angeles")
	env := envs.NewBuilder().WithTimezone(la).WithAllowedLanguages(flow.Language(), contact.Language()).Build()

	repro := &Repro{}
	var call *flows.Call

	if initialMsg != "" {
		msg := events.NewMsgReceived(createMessage(contact, initialMsg))
		repro.Trigger = triggers.NewBuilder(flow.Reference(false)).Msg(msg).Build()

		printEvents([]flows.Event{msg}, out)
	} else {
		tb := triggers.NewBuilder(flow.Reference(false)).Manual()

		// if we're starting a voice flow we need a call
		if flow.Type() == flows.FlowTypeVoice {
			channel := sa.Channels().GetForURN(flows.NewContactURN(urns.URN("tel:+12065551212"), nil), assets.ChannelRoleCall)
			call = flows.NewCall("01978a2f-ad9a-7f2e-ad44-6e7547078cec", channel, urns.URN("tel:+12065551212"))
		}

		repro.Trigger = tb.Build()
	}
	fmt.Fprintf(out, "Starting flow '%s'....\n---------------------------------------\n", flow.Name())

	// start our session
	session, sprint, err := eng.NewSession(ctx, sa, env, contact, repro.Trigger, call)
	if err != nil {
		return nil, err
	}

	printEvents(sprint.Events(), out)
	scanner := bufio.NewScanner(in)

	for session.Status() == flows.SessionStatusWaiting {

		// ask for input
		fmt.Fprintf(out, "> ")
		scanner.Scan()

		text := scanner.Text()
		var resume flows.Resume

		// create our resume
		if text == "/timeout" {
			resume = resumes.NewWaitTimeout(events.NewWaitTimedOut())
		} else if strings.HasPrefix(text, "/dial") {
			status := flows.DialStatus(strings.TrimSpace(text[5:]))
			resume = resumes.NewDial(events.NewDialEnded(flows.NewDial(status, 10)))
		} else {
			msg := events.NewMsgReceived(createMessage(contact, scanner.Text()))
			resume = resumes.NewMsg(msg)

			printEvents([]flows.Event{msg}, out)
		}

		repro.Resumes = append(repro.Resumes, resume)

		sprint, err := session.Resume(ctx, resume)
		if err != nil {
			return nil, err
		}

		printEvents(sprint.Events(), out)
	}

	return repro, nil
}

func createMessage(contact *flows.Contact, text string) *flows.MsgIn {
	return flows.NewMsgIn(contact.URNs()[0].URN(), nil, text, []utils.Attachment{}, "")
}

func printEvents(log []flows.Event, out io.Writer) {
	for _, event := range log {
		PrintEvent(event, out)
		fmt.Fprintln(out)
	}
}

// PrintEvent prints out the given event to the given writer
func PrintEvent(event flows.Event, out io.Writer) {
	var msg string
	switch typed := event.(type) {
	case *events.BroadcastCreated:
		text := typed.Translations[typed.BaseLanguage].Text
		msg = fmt.Sprintf("🔉 broadcasted '%s' to ...", text)
	case *events.ContactFieldChanged:
		var action string
		if typed.Value != nil {
			action = fmt.Sprintf("changed to '%s'", typed.Value.Text.Native())
		} else {
			action = "cleared"
		}
		msg = fmt.Sprintf("✏️ field '%s' %s", typed.Field.Key, action)
	case *events.ContactGroupsChanged:
		msgs := make([]string, 0)
		if len(typed.GroupsAdded) > 0 {
			groups := make([]string, len(typed.GroupsAdded))
			for i, group := range typed.GroupsAdded {
				groups[i] = fmt.Sprintf("'%s'", group.Name)
			}
			msgs = append(msgs, "added to "+strings.Join(groups, ", "))
		}
		if len(typed.GroupsRemoved) > 0 {
			groups := make([]string, len(typed.GroupsRemoved))
			for i, group := range typed.GroupsRemoved {
				groups[i] = fmt.Sprintf("'%s'", group.Name)
			}
			msgs = append(msgs, "removed from "+strings.Join(groups, ", "))
		}
		msg = fmt.Sprintf("👪 %s", strings.Join(msgs, ", "))
	case *events.ContactLanguageChanged:
		msg = fmt.Sprintf("🌐 language changed to '%s'", typed.Language)
	case *events.ContactNameChanged:
		msg = fmt.Sprintf("📛 name changed to '%s'", typed.Name)
	case *events.ContactRefreshed:
		msg = "👤 contact refreshed on resume"
	case *events.ContactTimezoneChanged:
		msg = fmt.Sprintf("🕑 timezone changed to '%s'", typed.Timezone)
	case *events.DialEnded:
		msg = fmt.Sprintf("☎️ dial ended with '%s'", typed.Dial.Status)
	case *events.DialWait:
		msg = "⏳ waiting for dial (type /dial <answered|no_answer|busy|failed>)..."
	case *events.EmailSent:
		msg = fmt.Sprintf("✉️ email sent with subject '%s'", typed.Subject)
	case *events.EnvironmentRefreshed:
		msg = "⚙️ environment refreshed on resume"
	case *events.Error:
		msg = fmt.Sprintf("⚠️ %s", typed.Text)
	case *events.Failure:
		msg = fmt.Sprintf("🛑 %s", typed.Text)
	case *events.FlowEntered:
		msg = fmt.Sprintf("↪️ entered flow '%s'", typed.Flow.Name)
	case *events.InputLabelsAdded:
		labels := make([]string, len(typed.Labels))
		for i, label := range typed.Labels {
			labels[i] = fmt.Sprintf("'%s'", label.Name)
		}
		msg = fmt.Sprintf("🏷️ labeled with %s", strings.Join(labels, ", "))
	case *events.IVRCreated:
		msg = fmt.Sprintf("📞 IVR created \"%s\"", typed.Msg.Text())
	case *events.MsgCreated:
		msg = fmt.Sprintf("💬 message created \"%s\"", typed.Msg.Text())
	case *events.MsgReceived:
		msg = fmt.Sprintf("📥 message received \"%s\"", typed.Msg.Text())
	case *events.MsgWait:
		if typed.TimeoutSeconds != nil {
			msg = fmt.Sprintf("⏳ waiting for message (%d sec timeout, type /timeout to simulate)...", *typed.TimeoutSeconds)
		} else {
			msg = "⏳ waiting for message..."
		}
	case *events.RunExpired:
		msg = "📆 exiting due to expiration"
	case *events.RunResultChanged:
		msg = fmt.Sprintf("📈 run result '%s' changed to '%s' with category '%s'", typed.Name, typed.Value, typed.Category)
	case *events.ServiceCalled:
		switch typed.Service {
		case "classifier":
			msg = fmt.Sprintf("👁️‍🗨️ NLU classifier '%s' called", typed.Classifier.Name)
		}
	case *events.SessionTriggered:
		msg = fmt.Sprintf("🏁 session triggered for '%s'", typed.Flow.Name)
	case *events.TicketOpened:
		msg = fmt.Sprintf("🎟️ ticket opened with topic \"%s\"", typed.Ticket.Topic.Name)
	case *events.WaitTimedOut:
		msg = "⏲️ resuming due to wait timeout"
	case *events.WebhookCalled:
		url := stringsx.TruncateEllipsis(typed.URL, 50)
		msg = fmt.Sprintf("☁️ called %s", url)
	default:
		msg = fmt.Sprintf("❓ %s event", typed.Type())
	}

	fmt.Fprint(out, msg)
}

// Repro describes the trigger and resumes needed to reproduce this session
type Repro struct {
	Trigger flows.Trigger  `json:"trigger"`
	Resumes []flows.Resume `json:"resumes,omitempty"`
}
