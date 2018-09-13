package flows

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// BaseMsg represents a incoming or outgoing message with the session contact
type BaseMsg struct {
	UUID_        MsgUUID                  `json:"uuid"`
	ID_          MsgID                    `json:"id,omitempty"`
	URN_         urns.URN                 `json:"urn" validate:"omitempty,urn"`
	Channel_     *assets.ChannelReference `json:"channel,omitempty"`
	Text_        string                   `json:"text"`
	Attachments_ []Attachment             `json:"attachments,omitempty"`
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
func NewMsgIn(uuid MsgUUID, id MsgID, urn urns.URN, channel *Channel, text string, attachments []Attachment) *MsgIn {
	var channelRef *assets.ChannelReference
	if channel != nil {
		channelRef = channel.Reference()
	}

	return &MsgIn{
		BaseMsg: BaseMsg{
			UUID_:        uuid,
			ID_:          id,
			URN_:         urn,
			Channel_:     channelRef,
			Text_:        text,
			Attachments_: attachments,
		},
	}
}

// NewMsgOut creates a new outgoing message
func NewMsgOut(urn urns.URN, channel *Channel, text string, attachments []Attachment, quickReplies []string) *MsgOut {
	var channelRef *assets.ChannelReference
	if channel != nil {
		channelRef = channel.Reference()
	}

	return &MsgOut{
		BaseMsg: BaseMsg{
			UUID_:        MsgUUID(utils.NewUUID()),
			URN_:         urn,
			Channel_:     channelRef,
			Text_:        text,
			Attachments_: attachments,
		},
		QuickReplies_: quickReplies,
	}
}

// UUID returns the UUID of this message
func (m *BaseMsg) UUID() MsgUUID { return m.UUID_ }

// ID returns the ID of this message
func (m *BaseMsg) ID() MsgID { return m.ID_ }

// URN returns the URN of this message
func (m *BaseMsg) URN() urns.URN { return m.URN_ }

// Channel returns the channel of this message
func (m *BaseMsg) Channel() *assets.ChannelReference { return m.Channel_ }

// Text returns the text of this message
func (m *BaseMsg) Text() string { return m.Text_ }

// Attachments returns the attachments of this message
func (m *BaseMsg) Attachments() []Attachment { return m.Attachments_ }

// QuickReplies returns the quick replies of this outgoing message
func (m *MsgOut) QuickReplies() []string { return m.QuickReplies_ }
