package static

import (
	"github.com/nyaruka/goflow/assets"
)

// User is a JSON serializable implementation of a user asset
type User struct {
	UUID_  assets.UserUUID `json:"uuid"  validate:"required"`
	Name_  string          `json:"name"`
	Email_ string          `json:"email" validate:"required"`
}

// NewUser creates a new user from the passed in email and name
func NewUser(uuid assets.UserUUID, name, email string) assets.User {
	return &User{UUID_: uuid, Name_: name, Email_: email}
}

// UUID returns the UUID of the user
func (u *User) UUID() assets.UserUUID { return u.UUID_ }

// Email returns the unique email address of the user
func (u *User) Email() string { return u.Email_ }

// Name returns the name of the user
func (u *User) Name() string { return u.Name_ }
