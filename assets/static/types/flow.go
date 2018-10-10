package types

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// Flow is a JSON serializable implementation of a flow asset
type Flow struct {
	UUID_       assets.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name_       string          `json:"name"`
	Definition_ json.RawMessage
}

// UUID returns the UUID of the flow
func (f *Flow) UUID() assets.FlowUUID { return f.UUID_ }

// Name returns the name of the flow
func (f *Flow) Name() string { return f.Name_ }

func (f *Flow) Definition() json.RawMessage { return f.Definition_ }

func (f *Flow) UnmarshalJSON(data []byte) error {
	f.Definition_ = data

	// alias our type so we don't end up here again
	type alias Flow
	return json.Unmarshal(data, (*alias)(f))
}

// ReadFlow reads a flow from the given JSON
func ReadFlow(data json.RawMessage) (assets.Flow, error) {
	f := &Flow{Definition_: data}
	if err := utils.UnmarshalAndValidate(data, f); err != nil {
		return nil, fmt.Errorf("unable to read flow: %s", err)
	}
	return f, nil
}
