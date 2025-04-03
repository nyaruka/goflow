package flows

import (
	"encoding/json"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

var localNamePattern = regexp.MustCompile(`^[a-z_][a-z0-9_]{0,63}$`)

func init() {
	utils.RegisterValidatorTag("local_name",
		func(fl validator.FieldLevel) bool {
			return localNamePattern.MatchString(fl.Field().String())
		},
		func(validator.FieldError) string { return "is not a valid local variable name" },
	)
}

// Locals is a map of local variables for a run
type Locals struct {
	vals map[string]string
}

// NewLocals creates a new empty locals instance
func NewLocals() *Locals {
	return &Locals{make(map[string]string)}
}

func (l *Locals) Get(key string) string {
	return l.vals[key]
}

func (l *Locals) Set(key string, value string) {
	l.vals[key] = value
}

func (l *Locals) IsZero() bool {
	return len(l.vals) == 0
}

// Context returns the properties available in expressions
func (l *Locals) Context(env envs.Environment) map[string]types.XValue {
	vals := make(map[string]types.XValue)
	for k, v := range l.vals {
		vals[k] = types.NewXText(v)
	}
	return vals
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func (l *Locals) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.vals)
}

func (l *Locals) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.vals)
}
