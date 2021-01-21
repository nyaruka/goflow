package flows

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.RegisterValidatorAlias("dial_status", "eq=answered|eq=no_answer|eq=busy|eq=failed", func(validator.FieldError) string {
		return "is not a valid dial status"
	})
}

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
	Status   DialStatus `json:"status" validate:"required,dial_status"`
	Duration int        `json:"duration"`
}

// NewDial creates a new dial
func NewDial(status DialStatus, duration int) *Dial {
	return &Dial{Status: status, Duration: duration}
}

// Context for dial resumes additionally exposes the dial object
func (d *Dial) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"status":   types.NewXText(string(d.Status)),
		"duration": types.NewXNumberFromInt(d.Duration),
	}
}
