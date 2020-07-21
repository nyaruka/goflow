package flows

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

// NewChildHistory creates a new history for a child of the given session
func NewChildHistory(parent Session) *SessionHistory {
	return parent.History().Advance(parent.UUID(), sessionReceivedInput(parent))
}

// looks through a session's events to see if it received input
func sessionReceivedInput(s Session) bool {
	for _, r := range s.Runs() {
		if r.ReceivedInput() {
			return true
		}
	}
	return false
}
