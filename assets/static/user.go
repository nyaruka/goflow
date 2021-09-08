package static

import (
	"github.com/nyaruka/goflow/assets"
)

// User is a JSON serializable implementation of a user asset
type User struct {
	Email_ string `json:"email" validate:"required"`
	Name_  string `json:"name"`
}

// NewUser creates a new user from the passed in email and name
func NewUser(email, name string) assets.User {
	return &User{Email_: email, Name_: name}
}

// Email returns the unique email address of the user
func (u *User) Email() string { return u.Email_ }

// Name returns the name of the user
func (u *User) Name() string { return u.Name_ }
