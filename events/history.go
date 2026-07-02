package events

import "github.com/nyaruka/gocommon/uuids"

// SessionUUID is the UUID of a session
type SessionUUID uuids.UUID

// NewSessionUUID generates a new UUID for a session
func NewSessionUUID() SessionUUID { return SessionUUID(uuids.NewV7()) }

// SessionHistory provides information about the sessions that caused this session
type SessionHistory struct {
	ParentUUID          SessionUUID `json:"parent_uuid"`
	Ancestors           int         `json:"ancestors"`
	AncestorsSinceInput int         `json:"ancestors_since_input"`
}

// Advance moves history forward to a new parent
func (h *SessionHistory) Advance(newParent SessionUUID, receivedInput bool) *SessionHistory {
	ancestorsSinceinput := 0
	if !receivedInput {
		ancestorsSinceinput = h.AncestorsSinceInput + 1
	}

	return &SessionHistory{
		ParentUUID:          newParent,
		Ancestors:           h.Ancestors + 1,
		AncestorsSinceInput: ancestorsSinceinput,
	}
}

// EmptyHistory is used for a session which has no history
var EmptyHistory = &SessionHistory{}
