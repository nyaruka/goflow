package flows

// DialStatus is the type for different dial statuses
type DialStatus string

// possible dial status values
const (
	DialStatusAnswered DialStatus = "answered"
	DialStatusNoAnswer DialStatus = "no_answer"
	DialStatusBusy     DialStatus = "busy"
	DialStatusFailed   DialStatus = "failed"
)

// Dial represents a dialed call or attempt to dial a phone number
type Dial struct {
	Status   DialStatus `json:"status" validate:"required"`
	Duration int        `json:"duration"`
}
