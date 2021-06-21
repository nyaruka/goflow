package flows

import (
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

type User struct {
	email string
	name  string
}

func NewUser(email, name string) *User {
	return &User{email: email, name: name}
}

func (u *User) Email() string { return u.email }
func (u *User) Name() string  { return u.name }

// Context returns the properties available in expressions
//
//   email:text -> the email address of the user
//   name:text -> the name of the user
//
// @context user
func (u *User) Context(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"email": types.NewXText(u.email),
		"name":  types.NewXText(u.name),
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type userEnvelope struct {
	Email string `json:"email" validate:"required"`
	Name  string `json:"name"`
}

// UmarshalJSON unmarshals this object from JSON
func (u *User) UnmarshalJSON(data []byte) error {
	// can be read from email string
	if data[0] == '"' {
		var email string
		if err := jsonx.Unmarshal(data, &email); err != nil {
			return err
		}
		u.email = email
	} else {
		e := &userEnvelope{}
		if err := utils.UnmarshalAndValidate(data, e); err != nil {
			return err
		}
		u.email = e.Email
		u.name = e.Name
	}

	return nil
}

// MarshalJSON marshals this object into JSON
func (u *User) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&userEnvelope{Email: u.email, Name: u.name})
}
