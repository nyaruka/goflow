package definition

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type exit struct {
	uuid        flows.ExitUUID
	destination flows.NodeUUID
	name        string
}

// NewExit creates a new exit
func NewExit(uuid flows.ExitUUID, destination flows.NodeUUID, name string) flows.Exit {
	return &exit{uuid: uuid, destination: destination, name: name}
}

func (e *exit) UUID() flows.ExitUUID                { return e.uuid }
func (e *exit) DestinationNodeUUID() flows.NodeUUID { return e.destination }
func (e *exit) Name() string                        { return e.name }

type exitEnvelope struct {
	UUID                flows.ExitUUID `json:"uuid"                               validate:"required,uuid4"`
	DestinationNodeUUID flows.NodeUUID `json:"destination_node_uuid,omitempty"    validate:"omitempty,uuid4"`
	Name                string         `json:"name,omitempty"`
}

// UnmarshalJSON unmarshals a node exit from the given JSON
func (e *exit) UnmarshalJSON(data []byte) error {
	envelope := &exitEnvelope{}
	err := utils.UnmarshalAndValidate(data, envelope)
	if err != nil {
		return errors.Wrap(err, "unable to read exit")
	}

	e.uuid = envelope.UUID
	e.destination = envelope.DestinationNodeUUID
	e.name = envelope.Name

	return nil
}

// MarshalJSON marshals this node exit into JSON
func (e *exit) MarshalJSON() ([]byte, error) {
	envelope := &exitEnvelope{e.uuid, e.destination, e.name}
	return json.Marshal(envelope)
}
