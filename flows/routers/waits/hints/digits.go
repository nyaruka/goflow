package hints

import (
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeDigits, func() flows.Hint { return &DigitsHint{} })
}

// TypeDigits is the type of our digits hint
const TypeDigits string = "digits"

// DigitsHint requests a message containing one or more digits
type DigitsHint struct {
	baseHint

	Count        *int   `json:"count,omitempty"`
	TerminatedBy string `json:"terminated_by,omitempty"`
}

// NewFixedDigitsHint creates a new digits hint for a fixed count of digits
func NewFixedDigitsHint(count int) *DigitsHint {
	return &DigitsHint{
		baseHint: newBaseHint(TypeDigits),
		Count:    &count,
	}
}

// NewTerminatedDigitsHint creates a new digits hint for a sequence of digits terminated by the given key
func NewTerminatedDigitsHint(terminatedBy string) *DigitsHint {
	return &DigitsHint{
		baseHint:     newBaseHint(TypeDigits),
		TerminatedBy: terminatedBy,
	}
}
