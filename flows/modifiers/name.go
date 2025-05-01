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
	registerType(TypeName, readNameModifier)
}

// TypeName is the type of our name modifier
const TypeName string = "name"

// NameModifier modifies the name of a contact
type NameModifier struct {
	baseModifier

	Name string `json:"name"`
}

// NewName creates a new name modifier
func NewName(name string) *NameModifier {
	return &NameModifier{
		baseModifier: newBaseModifier(TypeName),
		Name:         name,
	}
}

// Apply applies this modification to the given contact
func (m *NameModifier) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	if contact.Name() != m.Name {
		// truncate value if necessary
		name := stringsx.Truncate(m.Name, eng.Options().MaxFieldChars)

		contact.SetName(name)
		log(events.NewContactNameChanged(name))
		return true
	}
	return false
}

var _ flows.Modifier = (*NameModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readNameModifier(assets flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &NameModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
