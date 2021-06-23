package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// User adds some functionality to user assets.
type User struct {
	assets.User
}

// NewUser returns a new user object from the given user asset
func NewUser(asset assets.User) *User {
	return &User{User: asset}
}

// Asset returns the underlying asset
func (u *User) Asset() assets.User { return u.User }

// Reference returns a reference to this user
func (u *User) Reference() *assets.UserReference {
	return assets.NewUserReference(u.Email(), u.Name())
}

// Context returns the properties available in expressions
//
//   __default__:text -> the name of the user
//   email:text -> the email address of the user
//   name:text -> the name of the user
//
// @context user
func (u *User) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"__default__": types.NewXText(u.Name()),
		"email":       types.NewXText(u.Email()),
		"name":        types.NewXText(u.Name()),
	}
}

// UserAssets provides access to all user assets
type UserAssets struct {
	all     []*User
	byEmail map[string]*User
}

// NewUserAssets creates a new set of user assets
func NewUserAssets(users []assets.User) *UserAssets {
	s := &UserAssets{
		all:     make([]*User, len(users)),
		byEmail: make(map[string]*User, len(users)),
	}
	for i, asset := range users {
		user := NewUser(asset)
		s.all[i] = user
		s.byEmail[user.Email()] = user
	}
	return s
}

// Get returns the user with the given email
func (s *UserAssets) Get(email string) *User {
	return s.byEmail[email]
}
