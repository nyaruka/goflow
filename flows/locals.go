package flows

import (
	"encoding/json"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Locals is a map of local variables for a run
type Locals struct {
	vals map[string]types.XValue
}

// NewLocals creates a new empty locals instance
func NewLocals() *Locals {
	return &Locals{make(map[string]types.XValue)}
}

func (l *Locals) Get(key string) types.XValue {
	return l.vals[key]
}

func (l *Locals) Set(key string, value types.XValue) {
	l.vals[key] = value
}

func (l *Locals) IsZero() bool {
	return len(l.vals) == 0
}

// Context returns the properties available in expressions
func (l *Locals) Context(env envs.Environment) map[string]types.XValue {
	return l.vals
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func (l *Locals) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.vals)
}

func (l *Locals) UnmarshalJSON(data []byte) error {
	obj, err := types.ReadXObject(data)
	if err != nil {
		return err
	}

	l.vals = make(map[string]types.XValue)
	for _, p := range obj.Properties() {
		l.vals[p], _ = obj.Get(p)
	}
	return nil
}
