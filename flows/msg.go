package flows

import (
	"fmt"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/uuids"
)

func init() {
	utils.Validator.RegisterAlias("msg_topic", "eq=event|eq=account|eq=purchase|eq=agent")
}

// MsgTopic is the topic, as required by some channel types
type MsgTopic string

// possible msg topic values
const (
	NilMsgTopic      MsgTopic = ""
	MsgTopicEvent    MsgTopic = "event"
	MsgTopicAccount  MsgTopic = "account"
	MsgTopicPurchase MsgTopic = "purchase"
	MsgTopicAgent    MsgTopic = "agent"
)

// BaseMsg represents a incoming or outgoing message with the session contact
type BaseMsg struct {
	UUID_        MsgUUID                  `json:"uuid"`
	ID_          MsgID                    `json:"id,omitempty"`
	URN_         urns.URN                 `json:"urn,omitempty" validate:"omitempty,urn"`
	Channel_     *assets.ChannelReference `json:"channel,omitempty"`
	Text_        string                   `json:"text"`
	Attachments_ []utils.Attachment       `json:"attachments,omitempty"`
}

// MsgIn represents a incoming message from the session contact
type MsgIn struct {
	BaseMsg

	ExternalID_ string `json:"external_id,omitempty"`
}

// MsgOut represents a outgoing message to the session contact
type MsgOut struct {
	BaseMsg

	QuickReplies_ []string       `json:"quick_replies,omitempty"`
	Templating_   *MsgTemplating `json:"templating,omitempty"`
	Topic_        MsgTopic       `json:"topic,omitempty"`
	TextLanguage  envs.Language  `json:"text_language,omitempty"`
}

// NewMsgIn creates a new incoming message
func NewMsgIn(uuid MsgUUID, urn urns.URN, channel *assets.ChannelReference, text string, attachments []utils.Attachment) *MsgIn {
	return &MsgIn{
		BaseMsg: BaseMsg{
			UUID_:        uuid,
			URN_:         urn,
			Channel_:     channel,
			Text_:        text,
			Attachments_: attachments,
		},
	}
}

// NewMsgOut creates a new outgoing message
func NewMsgOut(urn urns.URN, channel *assets.ChannelReference, text string, attachments []utils.Attachment, quickReplies []string, templating *MsgTemplating, topic MsgTopic) *MsgOut {
	return &MsgOut{
		BaseMsg: BaseMsg{
			UUID_:        MsgUUID(uuids.New()),
			URN_:         urn,
			Channel_:     channel,
			Text_:        text,
			Attachments_: attachments,
		},
		QuickReplies_: quickReplies,
		Templating_:   templating,
		Topic_:        topic,
	}
}

// NewIVRMsgOut creates a new outgoing message for IVR
func NewIVRMsgOut(urn urns.URN, channel *assets.ChannelReference, text string, textLanguage envs.Language, audioURL string) *MsgOut {
	var attachments []utils.Attachment
	if audioURL != "" {
		attachments = []utils.Attachment{utils.Attachment(fmt.Sprintf("audio:%s", audioURL))}
	}

	return &MsgOut{
		BaseMsg: BaseMsg{
			UUID_:        MsgUUID(uuids.New()),
			URN_:         urn,
			Channel_:     channel,
			Text_:        text,
			Attachments_: attachments,
		},
		QuickReplies_: nil,
		Templating_:   nil,
		Topic_:        NilMsgTopic,
		TextLanguage:  textLanguage,
	}
}

// UUID returns the UUID of this message
func (m *BaseMsg) UUID() MsgUUID { return m.UUID_ }

// ID returns the internal ID of this message
func (m *BaseMsg) ID() MsgID { return m.ID_ }

// SetID sets the internal ID of this message
func (m *BaseMsg) SetID(id MsgID) { m.ID_ = id }

// URN returns the URN of this message
func (m *BaseMsg) URN() urns.URN { return m.URN_ }

// SetURN returns the URN of this message
func (m *BaseMsg) SetURN(urn urns.URN) { m.URN_ = urn }

// Channel returns the channel of this message
func (m *BaseMsg) Channel() *assets.ChannelReference { return m.Channel_ }

// Text returns the text of this message
func (m *BaseMsg) Text() string { return m.Text_ }

// Attachments returns the attachments of this message
func (m *BaseMsg) Attachments() []utils.Attachment { return m.Attachments_ }

// ExternalID returns the optional external ID of this incoming message
func (m *MsgIn) ExternalID() string { return m.ExternalID_ }

// SetExternalID sets the external ID of this message
func (m *MsgIn) SetExternalID(id string) { m.ExternalID_ = id }

// QuickReplies returns the quick replies of this outgoing message
func (m *MsgOut) QuickReplies() []string { return m.QuickReplies_ }

// Templating returns the templating to use to send this message (if any)
func (m *MsgOut) Templating() *MsgTemplating { return m.Templating_ }

// Topic returns the topic to use to send this message (if any)
func (m *MsgOut) Topic() MsgTopic { return m.Topic_ }

// MsgTemplating represents any substituted message template that should be applied when sending this message
type MsgTemplating struct {
	Template_  *assets.TemplateReference `json:"template"`
	Language_  envs.Language             `json:"language"`
	Variables_ []string                  `json:"variables,omitempty"`
}

// Template returns the template this msg template is for
func (t MsgTemplating) Template() *assets.TemplateReference { return t.Template_ }

// Language returns the language that should be used for the template
func (t MsgTemplating) Language() envs.Language { return t.Language_ }

// Variables returns the variables that should be substituted in the template
func (t MsgTemplating) Variables() []string { return t.Variables_ }

// NewMsgTemplating creates and returns a new msg template
func NewMsgTemplating(template *assets.TemplateReference, language envs.Language, variables []string) *MsgTemplating {
	return &MsgTemplating{
		Template_:  template,
		Language_:  language,
		Variables_: variables,
	}
}
