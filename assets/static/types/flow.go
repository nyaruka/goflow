package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"

	"github.com/buger/jsonparser"
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

	// if there's a metadata field, this is a legacy definition, so read the UUID/name from there
	legacyMetadata, _, _, _ := jsonparser.Get(data, "metadata")
	if legacyMetadata != nil {
		return json.Unmarshal(legacyMetadata, (*alias)(f))
	}

	return json.Unmarshal(data, (*alias)(f))
}
