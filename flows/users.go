package flows

import (
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
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
	if u == nil {
		return nil
	}
	return assets.NewUserReference(u.UUID(), u.Name())
}

// Format returns a friendly string version of this user depending on what fields are set
func (u *User) Format() string {
	// if user has a name set, use that
	if u.Name() != "" {
		return u.Name()
	}

	// otherwise use email
	return u.Email()
}

// Context returns the properties available in expressions
//
//	__default__:text -> the name or email
//	email:text -> the email address of the user
//	name:text -> the name of the user
//	first_name:text -> the first name of the user
//
// @context user
func (u *User) Context(env envs.Environment) map[string]types.XValue {
	firstName := types.XTextEmpty

	names := utils.TokenizeString(u.Name())
	if len(names) >= 1 {
		firstName = types.NewXText(names[0])
	}

	return map[string]types.XValue{
		"__default__": types.NewXText(u.Format()),
		"email":       types.NewXText(u.Email()),
		"name":        types.NewXText(u.Name()),
		"first_name":  firstName,
	}
}

// UserAssets provides access to all user assets
type UserAssets struct {
	all    []*User
	byUUID map[assets.UserUUID]*User
}

// NewUserAssets creates a new set of user assets
func NewUserAssets(users []assets.User) *UserAssets {
	s := &UserAssets{
		all:    make([]*User, len(users)),
		byUUID: make(map[assets.UserUUID]*User, len(users)),
	}
	for i, asset := range users {
		user := NewUser(asset)
		s.all[i] = user
		s.byUUID[user.UUID()] = user
	}
	return s
}

// Get returns the user with the given UUID
func (s *UserAssets) Get(uuid assets.UserUUID) *User {
	return s.byUUID[uuid]
}

// FindByEmail looks for a user with the given email (case-insensitive)
func (s *UserAssets) FindByEmail(email string) *User {
	email = strings.ToLower(email)
	for _, user := range s.all {
		if strings.ToLower(user.Email()) == email {
			return user
		}
	}
	return nil
}
