package flows

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/utils"
)

type UnsendableReason string

const (
	// max length of a message attachment (type:url)
	MaxAttachmentLength = 2048

	// max length of a quick reply
	MaxQuickReplyLength = 64

	NilUnsendableReason            UnsendableReason = ""
	UnsendableReasonNoDestination  UnsendableReason = "no_destination"  // no sendable channel+URN pair
	UnsendableReasonContactBlocked UnsendableReason = "contact_blocked" // contact is blocked
	UnsendableReasonContactStopped UnsendableReason = "contact_stopped" // contact is stopped
	UnsendableReasonContactArchived UnsendableReason = "contact_archived" // contact is archived
)

// BaseMsg represents a incoming or outgoing message with the session contact
type BaseMsg struct {
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

	QuickReplies_     []QuickReply     `json:"quick_replies,omitempty"`
	Templating_       *MsgTemplating   `json:"templating,omitempty"`
	Locale_           i18n.Locale      `json:"locale,omitempty"`
	UnsendableReason_ UnsendableReason `json:"unsendable_reason,omitempty"`
}

// NewMsgIn creates a new incoming message
func NewMsgIn(urn urns.URN, channel *assets.ChannelReference, text string, attachments []utils.Attachment, externalID string) *MsgIn {
	return &MsgIn{
		BaseMsg: BaseMsg{
			URN_:         urn,
			Channel_:     channel,
			Text_:        text,
			Attachments_: attachments,
		},
		ExternalID_: externalID,
	}
}

// NewMsgOut creates a new outgoing message
func NewMsgOut(urn urns.URN, channel *assets.ChannelReference, content *MsgContent, templating *MsgTemplating, locale i18n.Locale, reason UnsendableReason) *MsgOut {
	return &MsgOut{
		BaseMsg: BaseMsg{
			URN_:         urn,
			Channel_:     channel,
			Text_:        content.Text,
			Attachments_: content.Attachments,
		},
		QuickReplies_:     content.QuickReplies,
		Templating_:       templating,
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
			URN_:         urn,
			Channel_:     channel,
			Text_:        text,
			Attachments_: attachments,
		},
		QuickReplies_: nil,
		Templating_:   nil,
		Locale_:       locale,
	}
}

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

// QuickReplies returns the quick replies of this outgoing message
func (m *MsgOut) QuickReplies() []QuickReply { return m.QuickReplies_ }

// Templating returns the templating to use to send this message (if any)
func (m *MsgOut) Templating() *MsgTemplating { return m.Templating_ }

// Locale returns the locale of this message (if any)
func (m *MsgOut) Locale() i18n.Locale { return m.Locale_ }

// UnsendableReason returns the reason this message can't be sent (if any)
func (m *MsgOut) UnsendableReason() UnsendableReason { return m.UnsendableReason_ }

type TemplatingVariable struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type TemplatingComponent struct {
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Variables map[string]int `json:"variables"`
}

// MsgTemplating represents any substituted message template that should be applied when sending this message
type MsgTemplating struct {
	Template   *assets.TemplateReference `json:"template"`
	Components []*TemplatingComponent    `json:"components,omitempty"`
	Variables  []*TemplatingVariable     `json:"variables,omitempty"`
}

// NewMsgTemplating creates and returns a new msg template
func NewMsgTemplating(template *assets.TemplateReference, components []*TemplatingComponent, variables []*TemplatingVariable) *MsgTemplating {
	return &MsgTemplating{Template: template, Components: components, Variables: variables}
}

type QuickReply struct {
	Text  string `json:"text"`
	Extra string `json:"extra,omitempty"`
}

// MarshalText marshals a quick reply into a text representation using a new line if extra is present
func (q QuickReply) MarshalText() (text []byte, err error) {
	vs := []string{q.Text}
	if q.Extra != "" {
		vs = append(vs, q.Extra)
	}
	return []byte(strings.Join(vs, "\n")), nil
}

func (q *QuickReply) UnmarshalText(text []byte) error {
	vs := strings.SplitN(string(text), "\n", 2)
	q.Text = vs[0]
	if len(vs) > 1 {
		q.Extra = vs[1]
	}
	return nil
}

func (q QuickReply) MarshalJSON() ([]byte, error) {
	// alias our type so we don't end up here again
	type alias QuickReply

	// we need to provide a MarshalJSON or the json package uses our MarshalText
	return jsonx.Marshal((alias)(q))
}

func (q *QuickReply) UnmarshalJSON(d []byte) error {
	// if we just have a string we unmarshal it into the text field
	if len(d) > 2 && d[0] == '"' && d[len(d)-1] == '"' {
		return jsonx.Unmarshal(d, &q.Text)
	}

	// alias our type so we don't end up here again
	type alias QuickReply

	return jsonx.Unmarshal(d, (*alias)(q))
}

// MsgContent is message content in a particular language
type MsgContent struct {
	Text         string             `json:"text"`
	Attachments  []utils.Attachment `json:"attachments,omitempty"`
	QuickReplies []QuickReply       `json:"quick_replies,omitempty"`
}

func (c *MsgContent) Empty() bool {
	return c.Text == "" && len(c.Attachments) == 0 && len(c.QuickReplies) == 0
}

type BroadcastUUID uuids.UUID

func NewBroadcastUUID() BroadcastUUID { return BroadcastUUID(uuids.NewV7()) }

type BroadcastTranslations map[i18n.Language]*MsgContent

// ForContact is a utility to help callers get the message content for a contact
func (b BroadcastTranslations) ForContact(e envs.Environment, c *Contact, baseLanguage i18n.Language) (*MsgContent, i18n.Locale) {
	// get the set of languages to merge translations from
	languages := make([]i18n.Language, 0, 3)

	// highest priority is the contact language if it is valid
	if c.Language() != i18n.NilLanguage && slices.Contains(e.AllowedLanguages(), c.Language()) {
		languages = append(languages, c.Language())
	}

	// then the default workspace language, then the base language
	languages = append(languages, e.DefaultLanguage(), baseLanguage)

	content := &MsgContent{}
	language := i18n.NilLanguage
	country := e.DefaultCountry()
	if c.Country() != i18n.NilCountry {
		country = c.Country()
	}

	for _, lang := range languages {
		trans := b[lang]
		if trans != nil {
			if content.Text == "" && trans.Text != "" {
				content.Text = trans.Text
				language = lang
			}
			if len(content.Attachments) == 0 && len(trans.Attachments) > 0 {
				content.Attachments = trans.Attachments
			}
			if len(content.QuickReplies) == 0 && len(trans.QuickReplies) > 0 {
				content.QuickReplies = trans.QuickReplies
			}
		}
	}

	return content, i18n.NewLocale(language, country)
}
