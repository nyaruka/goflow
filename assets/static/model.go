package static

import (
	"github.com/nyaruka/goflow/assets"
)

// Model is a JSON serializable implementation of a model asset
type Model struct {
	UUID_  assets.ModelUUID   `json:"uuid"  validate:"required,uuid"`
	Name_  string             `json:"name"`
	Type_  string             `json:"type"`
	Roles_ []assets.ModelRole `json:"roles" validate:"min=1,dive,eq=editing|eq=engine"`
}

// NewModel creates a new model
func NewModel(uuid assets.ModelUUID, name string, type_ string, roles []assets.ModelRole) assets.Model {
	return &Model{
		UUID_:  uuid,
		Name_:  name,
		Type_:  type_,
		Roles_: roles,
	}
}

// UUID returns the UUID of this model
func (m *Model) UUID() assets.ModelUUID { return m.UUID_ }

// Name returns the name of this model
func (m *Model) Name() string { return m.Name_ }

// Type returns the type of this model
func (m *Model) Type() string { return m.Type_ }

// Roles returns the roles of this model
func (m *Model) Roles() []assets.ModelRole { return m.Roles_ }
