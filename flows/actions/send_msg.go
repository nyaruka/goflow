package actions

import (
	"fmt"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeSendMsg, func() flows.Action { return &SendMsgAction{} })
}

type msgDestination struct {
	urn     urns.URN
	channel *flows.Channel
}

// TypeSendMsg is the type for the send message action
const TypeSendMsg string = "send_msg"

// SendMsgAction can be used to reply to the current contact in a flow. The text field may contain templates. The action
// will attempt to find pairs of URNs and channels which can be used for sending. If it can't find such a pair, it will
// create a message without a channel or URN.
//
// A [event:msg_created] event will be created with the evaluated text.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "send_msg",
//     "text": "Hi @contact.name, are you ready to complete today's survey?",
//     "attachments": [],
//     "all_urns": false
//   }
//
// @action send_msg
type SendMsgAction struct {
	BaseAction
	universalAction

	Text         string   `json:"text"`
	Attachments  []string `json:"attachments,omitempty"`
	QuickReplies []string `json:"quick_replies,omitempty"`
	AllURNs      bool     `json:"all_urns,omitempty"`
}

// NewSendMsgAction creates a new send msg action
func NewSendMsgAction(uuid flows.ActionUUID, text string, attachments []string, quickReplies []string, allURNs bool) *SendMsgAction {
	return &SendMsgAction{
		BaseAction:   NewBaseAction(TypeSendMsg, uuid),
		Text:         text,
		Attachments:  attachments,
		QuickReplies: quickReplies,
		AllURNs:      allURNs,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *SendMsgAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute runs this action
func (a *SendMsgAction) Execute(run flows.FlowRun, step flows.Step) error {
	if run.Contact() == nil {
		a.logError(run, step, fmt.Errorf("can't execute action in session without a contact"))
		return nil
	}

	evaluatedText, evaluatedAttachments, evaluatedQuickReplies := a.evaluateMessage(run, step, nil, a.Text, a.Attachments, a.QuickReplies)

	channels := run.Session().Assets().Channels()
	destinations := []msgDestination{}

	for _, u := range run.Contact().URNs() {
		channel := channels.GetForURN(u, assets.ChannelRoleSend)
		if channel != nil {
			destinations = append(destinations, msgDestination{urn: u.URN, channel: channel})

			// if we're not sending to all URNs we just need the first sendable URN
			if !a.AllURNs {
				break
			}
		}
	}

	// create a new message for each URN+channel destination
	for _, dest := range destinations {
		var channelRef *assets.ChannelReference
		if dest.channel != nil {
			channelRef = assets.NewChannelReference(dest.channel.UUID(), dest.channel.Name())
		}

		msg := flows.NewMsgOut(dest.urn, channelRef, evaluatedText, evaluatedAttachments, evaluatedQuickReplies)
		a.log(run, step, events.NewMsgCreatedEvent(msg))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgOut(urns.NilURN, nil, evaluatedText, evaluatedAttachments, evaluatedQuickReplies)
		a.log(run, step, events.NewMsgCreatedEvent(msg))
	}

	return nil
}
