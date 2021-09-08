package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/utils"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {
	utils.RegisterStructValidator(UserReferenceValidation, UserReference{})
}

// User is an person who can trigger flows or be assigned tickets etc.
//
//   {
//     "email": "bob@nyaruka.com",
//     "name": "Bob"
//   }
//
// @asset user
type User interface {
	Email() string
	Name() string
}

// UserReference is used to reference a user
type UserReference struct {
	Email      string `json:"email,omitempty" validate:"omitempty,email"`
	Name       string `json:"name,omitempty"`
	EmailMatch string `json:"email_match,omitempty" engine:"evaluated"`
}

// NewUserReference creates a new user reference with the given key and name
func NewUserReference(email, name string) *UserReference {
	return &UserReference{Email: email, Name: name}
}

// NewVariableUserReference creates a new user reference from the given templatized email match
func NewVariableUserReference(emailMatch string) *UserReference {
	return &UserReference{EmailMatch: emailMatch}
}

// Type returns the name of the asset type
func (r *UserReference) Type() string {
	return "user"
}

// Identity returns the unique identity of the asset
func (r *UserReference) Identity() string {
	return r.Email
}

// Variable returns whether this a variable (vs concrete) reference
func (r *UserReference) Variable() bool {
	return r.Identity() == ""
}

func (r *UserReference) String() string {
	return fmt.Sprintf("%s[email=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

// UmarshalJSON unmarshals this object from JSON
func (r *UserReference) UnmarshalJSON(data []byte) error {
	// can be read from email string
	if data[0] == '"' {
		var email string
		if err := jsonx.Unmarshal(data, &email); err != nil {
			return err
		}
		r.Email = email
		r.Name = ""
		return nil
	}

	// or a JSON object with email/name properties
	var raw map[string]string
	if err := jsonx.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Email = raw["email"]
	r.Name = raw["name"]
	r.EmailMatch = raw["email_match"]
	return nil
}

var _ Reference = (*UserReference)(nil)

//------------------------------------------------------------------------------------------
// Validation
//------------------------------------------------------------------------------------------

// UserReferenceValidation validates that the given user reference is either a concrete
// reference or an email matcher
func UserReferenceValidation(sl validator.StructLevel) {
	ref := sl.Current().Interface().(UserReference)
	if neitherOrBoth(string(ref.Email), ref.EmailMatch) {
		sl.ReportError(ref.Email, "email", "Email", "mutually_exclusive", "email_match")
		sl.ReportError(ref.EmailMatch, "email_match", "EmailMatch", "mutually_exclusive", "email")
	}
}
