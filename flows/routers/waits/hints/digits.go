package hints

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeDigits, func() flows.Hint { return &Digits{} })
}

// TypeDigits is the type of our digits hint
const TypeDigits string = "digits"

// Digits requests a message containing one or more digits
type Digits struct {
	baseHint

	Count        *int   `json:"count,omitempty"`
	TerminatedBy string `json:"terminated_by,omitempty"`
}

// NewFixedDigits creates a new digits hint for a fixed count of digits
func NewFixedDigits(count int) *Digits {
	return &Digits{
		baseHint: newBaseHint(TypeDigits),
		Count:    &count,
	}
}

// NewTerminatedDigits creates a new digits hint for a sequence of digits terminated by the given key
func NewTerminatedDigits(terminatedBy string) *Digits {
	return &Digits{
		baseHint:     newBaseHint(TypeDigits),
		TerminatedBy: terminatedBy,
	}
}
