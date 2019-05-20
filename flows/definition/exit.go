package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
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
func (e *exit) LocalizationUUID() utils.UUID { return utils.UUID(e.uuid) }

func (e *exit) Inspect(inspect func(flows.Inspectable)) {
	inspect(e)
}

// EnumerateTemplates enumerates all expressions on this object
func (e *exit) EnumerateTemplates(include flows.TemplateIncluder) {}

// EnumerateDependencies enumerates all dependencies on this object
func (e *exit) EnumerateDependencies(localization flows.Localization, include func(assets.Reference)) {
}

// EnumerateResults enumerates all potential results on this object
func (e *exit) EnumerateResults(include func(*flows.ResultSpec)) {}

// EnumerateElementUUIDs enumerates all element UUIDs on this object
func (e *exit) EnumerateElementUUIDs(include func(*utils.UUID)) {
	include((*utils.UUID)(&e.uuid))
	include((*utils.UUID)(&e.destination))
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type exitEnvelope struct {
	UUID            flows.ExitUUID `json:"uuid"                       validate:"required,uuid4"`
	DestinationUUID flows.NodeUUID `json:"destination_uuid,omitempty" validate:"omitempty,uuid4"`
}

// UnmarshalJSON unmarshals a node exit from the given JSON
func (e *exit) UnmarshalJSON(data []byte) error {
	envelope := &exitEnvelope{}

	if err := utils.UnmarshalAndValidate(data, envelope); err != nil {
		return errors.Wrap(err, "unable to read exit")
	}

	e.uuid = envelope.UUID
	e.destination = envelope.DestinationUUID
	return nil
}

// MarshalJSON marshals this node exit into JSON
func (e *exit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&exitEnvelope{e.uuid, e.destination})
}
