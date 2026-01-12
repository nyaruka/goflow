package modifiers

import (
	"context"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeName, readName)
}

// TypeName is the type of our name modifier
const TypeName string = "name"

// Name modifies the name of a contact
type Name struct {
	baseModifier

	name string
}

// NewName creates a new name modifier
func NewName(name string) *Name {
	return &Name{
		baseModifier: newBaseModifier(TypeName),
		name:         name,
	}
}

// Apply applies this modification to the given contact
func (m *Name) Apply(ctx context.Context, eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventLogger) (bool, error) {
	if contact.Name() != m.name {
		// truncate value if necessary
		name := stringsx.Truncate(m.name, eng.Options().MaxFieldChars)

		contact.SetName(name)
		log(events.NewContactNameChanged(name))
		return true, nil
	}
	return false, nil
}

var _ flows.Modifier = (*Name)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type nameEnvelope struct {
	utils.TypedEnvelope

	Name string `json:"name"`
}

func readName(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &nameEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewName(e.Name), nil
}

func (m *Name) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&nameEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Name:          m.name,
	})
}
