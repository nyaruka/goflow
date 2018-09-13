package types

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// json serializable implementation of a flow asset
type flow struct {
	UUID_       assets.FlowUUID `json:"uuid" validate:"required,uuid4"`
	Name_       string          `json:"name"`
	Definition_ json.RawMessage
}

// UUID returns the UUID of the flow
func (f *flow) UUID() assets.FlowUUID { return f.UUID_ }

// Name returns the name of the flow
func (f *flow) Name() string { return f.Name_ }

func (f *flow) Definition() json.RawMessage { return f.Definition_ }

// ReadFlow reads a flow from the given JSON
func ReadFlow(data json.RawMessage) (assets.Flow, error) {
	f := &flow{Definition_: data}
	if err := utils.UnmarshalAndValidate(data, f); err != nil {
		return nil, fmt.Errorf("unable to read flow: %s", err)
	}
	return f, nil
}
