package routers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/pkg/errors"
)

type Category struct {
	uuid     flows.CategoryUUID
	name     string
	exitUUID flows.ExitUUID
}

// NewCategory creates a new category
func NewCategory(uuid flows.CategoryUUID, name string, exit flows.ExitUUID) *Category {
	return &Category{uuid: uuid, name: name, exitUUID: exit}
}

func (c *Category) UUID() flows.CategoryUUID { return c.uuid }
func (c *Category) Name() string             { return c.name }
func (c *Category) ExitUUID() flows.ExitUUID { return c.exitUUID }

// LocalizationUUID gets the UUID which identifies this object for localization
func (c *Category) LocalizationUUID() uuids.UUID { return uuids.UUID(c.uuid) }

var _ flows.Category = (*Category)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type categoryEnvelope struct {
	UUID     flows.CategoryUUID `json:"uuid"                validate:"required,uuid4"`
	Name     string             `json:"name,omitempty"`
	ExitUUID flows.ExitUUID     `json:"exit_uuid,omitempty" validate:"required,uuid4"`
}

// ReadCategory unmarshals a router category from the given JSON
func ReadCategory(data []byte) (flows.Category, error) {
	e := &categoryEnvelope{}

	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, errors.Wrap(err, "unable to read category")
	}

	return NewCategory(e.UUID, e.Name, e.ExitUUID), nil
}

// MarshalJSON marshals this node category into JSON
func (c *Category) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&categoryEnvelope{
		c.uuid,
		c.name,
		c.exitUUID,
	})
}
