package inputs

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type readFunc func(flows.SessionAssets, []byte, assets.MissingCallback) (flows.Input, error)

var registeredTypes = map[string]readFunc{}

// registers a new type of input
func registerType(name string, f readFunc) {
	registeredTypes[name] = f
}

// base of all input types
type baseInput struct {
	type_     string
	uuid      flows.InputUUID
	channel   *flows.Channel
	createdOn time.Time
}

// creates a new base input
func newBaseInput(typeName string, uuid flows.InputUUID, channel *flows.Channel, createdOn time.Time) baseInput {
	return baseInput{
		type_:     typeName,
		uuid:      uuid,
		channel:   channel,
		createdOn: createdOn,
	}
}

// Type returns the type of this input
func (i *baseInput) Type() string { return i.type_ }

func (i *baseInput) UUID() flows.InputUUID   { return i.uuid }
func (i *baseInput) Channel() *flows.Channel { return i.channel }
func (i *baseInput) CreatedOn() time.Time    { return i.createdOn }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseEnvelope struct {
	Type      string                   `json:"type" validate:"required"`
	UUID      flows.InputUUID          `json:"uuid"`
	Channel   *assets.ChannelReference `json:"channel,omitempty" validate:"omitempty"`
	CreatedOn time.Time                `json:"created_on" validate:"required"`
}

// Read reads an input from the given typed envelope
func Read(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Input, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}

	return f(sa, data, missing)
}

func (i *baseInput) unmarshal(sessionAssets flows.SessionAssets, e *baseEnvelope, missing assets.MissingCallback) error {
	i.type_ = e.Type
	i.uuid = e.UUID
	i.createdOn = e.CreatedOn

	if e.Channel != nil {
		i.channel = sessionAssets.Channels().Get(e.Channel.UUID)
		if i.channel == nil {
			missing(e.Channel, nil)
			return nil
		}
	}
	return nil
}

func (i *baseInput) marshal(e *baseEnvelope) {
	e.Type = i.type_
	e.UUID = i.uuid
	e.CreatedOn = i.createdOn
	e.Channel = i.channel.Reference()
}
