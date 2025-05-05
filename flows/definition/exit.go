package definition

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type exit struct {
	uuid        flows.ExitUUID
	destination flows.NodeUUID
}

// NewExit creates a new exit
func NewExit(uuid flows.ExitUUID, destination flows.NodeUUID) flows.Exit {
	return &exit{uuid: uuid, destination: destination}
}

func (e *exit) UUID() flows.ExitUUID            { return e.uuid }
func (e *exit) DestinationUUID() flows.NodeUUID { return e.destination }

// LocalizationUUID gets the UUID which identifies this object for localization
func (e *exit) LocalizationUUID() uuids.UUID { return uuids.UUID(e.uuid) }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type exitEnvelope struct {
	UUID            flows.ExitUUID `json:"uuid"                       validate:"required,uuid"`
	DestinationUUID flows.NodeUUID `json:"destination_uuid,omitempty" validate:"omitempty,uuid"`
}

// UnmarshalJSON unmarshals a node exit from the given JSON
func (e *exit) UnmarshalJSON(data []byte) error {
	envelope := &exitEnvelope{}

	if err := utils.UnmarshalAndValidate(data, envelope); err != nil {
		return fmt.Errorf("unable to read exit: %w", err)
	}

	e.uuid = envelope.UUID
	e.destination = envelope.DestinationUUID
	return nil
}

// MarshalJSON marshals this node exit into JSON
func (e *exit) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&exitEnvelope{e.uuid, e.destination})
}
