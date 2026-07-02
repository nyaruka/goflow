package core

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.RegisterValidatorAlias("contact_status", "eq=active|eq=blocked|eq=stopped|eq=archived", func(validator.FieldError) string {
		return "is not a valid contact status"
	})
}

// ContactUUID is the UUID of a contact
type ContactUUID uuids.UUID

// NewContactUUID generates a new UUID for a contact
func NewContactUUID() ContactUUID { return ContactUUID(uuids.NewV4()) }

// ContactStatus is status in which a contact is in
type ContactStatus string

const (
	// ContactStatusActive is the contact status of active
	ContactStatusActive ContactStatus = "active"

	// ContactStatusBlocked is the contact status of blocked
	ContactStatusBlocked ContactStatus = "blocked"

	// ContactStatusStopped is the contact status of stopped
	ContactStatusStopped ContactStatus = "stopped"

	// ContactStatusArchived is the contact status of archived
	ContactStatusArchived ContactStatus = "archived"
)

// ContactReference is used to reference a contact
type ContactReference struct {
	UUID ContactUUID `json:"uuid" validate:"required,uuid"`
	Name string      `json:"name"`
}

// NewContactReference creates a new contact reference with the given UUID and name
func NewContactReference(uuid ContactUUID, name string) *ContactReference {
	return &ContactReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *ContactReference) Type() string {
	return "contact"
}

// Identity returns the unique identity of the asset
func (r *ContactReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *ContactReference) Variable() bool {
	return r.Identity() == ""
}

func (r *ContactReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ assets.Reference = (*ContactReference)(nil)
