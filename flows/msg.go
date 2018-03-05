package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/satori/go.uuid"
)

// BaseMsg represents a incoming or outgoing message with the session contact
type BaseMsg struct {
	UUID_        MsgUUID           `json:"uuid" validate:"required,uuid4"`
	URN_         urns.URN          `json:"urn" validate:"required"`
	Channel_     *ChannelReference `json:"channel,omitempty"`
	Text_        string            `json:"text"`
	Attachments_ []Attachment      `json:"attachments,omitempty"`
}

// MsgIn represents a incoming message from the session contact
type MsgIn struct {
	BaseMsg
}

// MsgOut represents a outgoing message to the session contact
type MsgOut struct {
	BaseMsg
	QuickReplies_ []string `json:"quick_replies,omitempty"`
}

// NewMsgIn creates a new incoming message
func NewMsgIn(uuid MsgUUID, urn urns.URN, channel Channel, text string, attachments []Attachment) *MsgIn {
	var channelRef *ChannelReference
	if channel != nil {
		channelRef = channel.Reference()
	}

	return &MsgIn{
		BaseMsg: BaseMsg{
			UUID_:        uuid,
			URN_:         urn,
			Channel_:     channelRef,
			Text_:        text,
			Attachments_: attachments,
		},
	}
}

// NewMsgOut creates a new outgoing message
func NewMsgOut(urn urns.URN, channel Channel, text string, attachments []Attachment, quickReplies []string) *MsgOut {
	var channelRef *ChannelReference
	if channel != nil {
		channelRef = channel.Reference()
	}

	return &MsgOut{
		BaseMsg: BaseMsg{
			UUID_:        MsgUUID(uuid.NewV4().String()),
			URN_:         urn,
			Channel_:     channelRef,
			Text_:        text,
			Attachments_: attachments,
		},
		QuickReplies_: quickReplies,
	}
}

func (m *BaseMsg) UUID() MsgUUID              { return m.UUID_ }
func (m *BaseMsg) URN() urns.URN              { return m.URN_ }
func (m *BaseMsg) Channel() *ChannelReference { return m.Channel_ }
func (m *BaseMsg) Text() string               { return m.Text_ }
func (m *BaseMsg) Attachments() []Attachment  { return m.Attachments_ }

func (m *MsgOut) QuickReplies() []string { return m.QuickReplies_ }
