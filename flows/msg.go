package flows

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.RegisterValidatorAlias("msg_topic", "eq=event|eq=account|eq=purchase|eq=agent", func(validator.FieldError) string {
		return "is not a valid message topic"
	})
}

type UnsendableReason string

const (
	NilUnsendableReason           UnsendableReason = ""
	UnsendableReasonNoDestination UnsendableReason = "no_destination" // no sendable channel+URN pair
	UnsendableReasonContactStatus UnsendableReason = "contact_status" // contact is blocked or stopped or archived
)

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

	QuickReplies_     []string         `json:"quick_replies,omitempty"`
	Templating_       *MsgTemplating   `json:"templating,omitempty"`
	Topic_            MsgTopic         `json:"topic,omitempty"`
	Locale_           i18n.Locale      `json:"locale,omitempty"`
	UnsendableReason_ UnsendableReason `json:"unsendable_reason,omitempty"`
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
func NewMsgOut(urn urns.URN, channel *assets.ChannelReference, text string, attachments []utils.Attachment, quickReplies []string, templating *MsgTemplating, topic MsgTopic, locale i18n.Locale, reason UnsendableReason) *MsgOut {
	return &MsgOut{
		BaseMsg: BaseMsg{
			UUID_:        MsgUUID(uuids.New()),
			URN_:         urn,
			Channel_:     channel,
			Text_:        text,
			Attachments_: attachments,
		},
		QuickReplies_:     quickReplies,
		Templating_:       templating,
		Topic_:            topic,
		Locale_:           locale,
		UnsendableReason_: reason,
	}
}

// NewIVRMsgOut creates a new outgoing message for IVR
func NewIVRMsgOut(urn urns.URN, channel *assets.ChannelReference, text string, audioURL string, locale i18n.Locale) *MsgOut {
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
		Locale_:       locale,
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

// Locale returns the locale of this message (if any)
func (m *MsgOut) Locale() i18n.Locale { return m.Locale_ }

// UnsendableReason returns the reason this message can't be sent (if any)
func (m *MsgOut) UnsendableReason() UnsendableReason { return m.UnsendableReason_ }

type TemplatingParam struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TemplatingComponent struct {
	Type   string            `json:"type"`
	Name   string            `json:"name"`
	Params []TemplatingParam `json:"params"`
}

var templateVariableRegex = regexp.MustCompile(`{{(\d+)}}`)

// Preview returns the content of the given template component rendered using these templating params
func (tc *TemplatingComponent) Preview(c assets.TemplateComponent) string {
	content := c.Content()

	for i, p := range tc.Params {
		content = strings.ReplaceAll(content, fmt.Sprintf("{{%d}}", i+1), p.Value)
	}

	content = templateVariableRegex.ReplaceAllString(content, "")

	return content
}

// MsgTemplating represents any substituted message template that should be applied when sending this message
type MsgTemplating struct {
	Template_   *assets.TemplateReference `json:"template"`
	Namespace_  string                    `json:"namespace"`
	Components_ []*TemplatingComponent    `json:"components,omitempty"`
}

// NewMsgTemplating creates and returns a new msg template
func NewMsgTemplating(template *assets.TemplateReference, namespace string, components []*TemplatingComponent) *MsgTemplating {
	return &MsgTemplating{Template_: template, Namespace_: namespace, Components_: components}
}

// Template returns the template this msg template is for
func (t *MsgTemplating) Template() *assets.TemplateReference { return t.Template_ }

// Namespace returns the namespace that should be for the template
func (t *MsgTemplating) Namespace() string { return t.Namespace_ }

// Components returns the components that should be used for the templates
func (t *MsgTemplating) Components() []*TemplatingComponent { return t.Components_ }

// BroadcastTranslation is the broadcast content in a particular language
type BroadcastTranslation struct {
	Text         string             `json:"text"`
	Attachments  []utils.Attachment `json:"attachments,omitempty"`
	QuickReplies []string           `json:"quick_replies,omitempty"`
}

type BroadcastTranslations map[i18n.Language]*BroadcastTranslation

// ForContact is a utility to help callers select the translation for a contact
func (b BroadcastTranslations) ForContact(e envs.Environment, c *Contact, baseLanguage i18n.Language) (*BroadcastTranslation, i18n.Language) {
	// first try the contact language if it is valid
	if c.Language() != i18n.NilLanguage && slices.Contains(e.AllowedLanguages(), c.Language()) {
		t := b[c.Language()]
		if t != nil {
			return t, c.Language()
		}
	}

	// second try the default flow language
	t := b[e.DefaultLanguage()]
	if t != nil {
		return t, e.DefaultLanguage()
	}

	// finally return the base language
	return b[baseLanguage], baseLanguage
}

// Scan supports reading translation values from JSON in database
func (t *BroadcastTranslations) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &t)
}

func (t BroadcastTranslations) Value() (driver.Value, error) { return json.Marshal(t) }
