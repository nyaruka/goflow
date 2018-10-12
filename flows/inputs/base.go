package inputs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type readFunc func(session flows.Session, data json.RawMessage) (flows.Input, error)

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of input
func RegisterType(name string, f readFunc) {
	registeredTypes[name] = f
}

type baseInput struct {
	type_     string
	uuid      flows.InputUUID
	channel   *flows.Channel
	createdOn time.Time
}

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

// Resolve resolves the given key when this input is referenced in an expression
func (i *baseInput) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "type":
		return types.NewXText(i.type_)
	case "uuid":
		return types.NewXText(string(i.uuid))
	case "created_on":
		return types.NewXDateTime(i.createdOn)
	case "channel":
		return i.channel
	}

	return types.NewXResolveError(i, key)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseInputEnvelope struct {
	Type      string                   `json:"type" validate:"required"`
	UUID      flows.InputUUID          `json:"uuid"`
	Channel   *assets.ChannelReference `json:"channel,omitempty" validate:"omitempty,dive"`
	CreatedOn time.Time                `json:"created_on" validate:"required"`
}

// ReadInput reads an input from the given typed envelope
func ReadInput(session flows.Session, data json.RawMessage) (flows.Input, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: %s", typeName)
	}
	return f(session, data)
}

func (i *baseInput) unmarshal(session flows.Session, e *baseInputEnvelope) error {
	var err error

	i.type_ = e.Type
	i.uuid = e.UUID
	i.createdOn = e.CreatedOn

	if e.Channel != nil {
		i.channel, err = session.Assets().Channels().Get(e.Channel.UUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *baseInput) marshal(e *baseInputEnvelope) {
	e.Type = i.Type()
	e.UUID = i.UUID()
	e.CreatedOn = i.CreatedOn()

	if i.Channel() != nil {
		e.Channel = i.Channel().Reference()
	}
}
