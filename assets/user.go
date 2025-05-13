package assets

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.RegisterStructValidator(UserReferenceValidation, UserReference{})
}

// UserUUID is the UUID of a user
type UserUUID uuids.UUID

// User is an person who can trigger flows or be assigned tickets etc.
//
//	{
//	  "uuid": "aefbc3b2-2f36-4a26-aa54-5fa20f761f99",
//	  "name": "Bob",
//	  "email": "bob@nyaruka.com"
//	}
//
// @asset user
type User interface {
	UUID() UserUUID
	Name() string
	Email() string
}

// UserReference is used to reference a user
type UserReference struct {
	UUID       UserUUID `json:"uuid,omitempty" validate:"omitempty,uuid"`
	Name       string   `json:"name,omitempty"`
	EmailMatch string   `json:"email,omitempty" engine:"evaluated"` // TODO should really be email_match in JSON
}

// NewUserReference creates a new user reference with the given UUID and name
func NewUserReference(uuid UserUUID, name string) *UserReference {
	return &UserReference{UUID: uuid, Name: name}
}

// NewVariableUserReference creates a new user reference from the given templatized email match
func NewVariableUserReference(emailMatch string) *UserReference {
	return &UserReference{EmailMatch: emailMatch}
}

// Type returns the name of the asset type
func (r *UserReference) Type() string {
	return "user"
}

// GenericUUID returns the untyped UUID
func (r *UserReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *UserReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *UserReference) Variable() bool {
	return r.Identity() == ""
}

func (r *UserReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ Reference = (*UserReference)(nil)

//------------------------------------------------------------------------------------------
// Validation
//------------------------------------------------------------------------------------------

// UserReferenceValidation validates that the given user reference is either a concrete
// reference or an email matcher
func UserReferenceValidation(sl validator.StructLevel) {
	ref := sl.Current().Interface().(UserReference)
	if neitherOrBoth(string(ref.UUID), ref.EmailMatch) {
		sl.ReportError(ref.UUID, "uuid", "UUID", "mutually_exclusive", "email")
		sl.ReportError(ref.EmailMatch, "email", "EmailMatch", "mutually_exclusive", "uuid")
	}
}
