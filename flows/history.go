package flows

// SessionHistory provides information about the sessions that caused this session
type SessionHistory struct {
	ParentUUID          SessionUUID `json:"parent_uuid"`
	Ancestors           int         `json:"ancestors"`
	AncestorsSinceInput int         `json:"ancestors_since_input"`
}

// EmptyHistory is used for a session which has no history
var EmptyHistory = SessionHistory{}

// NewChildHistory creates a new history for a child of the given session
func NewChildHistory(parent Session) SessionHistory {
	h := parent.History()
	ancestors := h.Ancestors + 1
	ancestorsSinceinput := 0

	if !sessionReceivedInput(parent) {
		ancestorsSinceinput = h.AncestorsSinceInput + 1
	}

	return SessionHistory{
		ParentUUID:          parent.UUID(),
		Ancestors:           ancestors,
		AncestorsSinceInput: ancestorsSinceinput,
	}
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
