package definition

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type exit struct {
	uuid        core.ExitUUID
	destination core.NodeUUID
}

// NewExit creates a new exit
func NewExit(uuid core.ExitUUID, destination core.NodeUUID) flows.Exit {
	return &exit{uuid: uuid, destination: destination}
}

func (e *exit) UUID() core.ExitUUID            { return e.uuid }
func (e *exit) DestinationUUID() core.NodeUUID { return e.destination }

// LocalizationUUID gets the UUID which identifies this object for localization
func (e *exit) LocalizationUUID() uuids.UUID { return uuids.UUID(e.uuid) }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type exitEnvelope struct {
	UUID            core.ExitUUID `json:"uuid"                       validate:"required,uuid"`
	DestinationUUID core.NodeUUID `json:"destination_uuid,omitempty" validate:"omitempty,uuid"`
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
