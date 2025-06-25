package modifiers

import (
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

	Name string `json:"name"`
}

// NewName creates a new name modifier
func NewName(name string) *Name {
	return &Name{
		baseModifier: newBaseModifier(TypeName),
		Name:         name,
	}
}

// Apply applies this modification to the given contact
func (m *Name) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	if contact.Name() != m.Name {
		// truncate value if necessary
		name := stringsx.Truncate(m.Name, eng.Options().MaxFieldChars)

		contact.SetName(name)
		log(events.NewContactNameChanged(name))
		return true
	}
	return false
}

var _ flows.Modifier = (*Name)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readName(assets flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &Name{}
	return m, utils.UnmarshalAndValidate(data, m)
}
