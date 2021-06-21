package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	u := flows.NewUser("bob@nyaruka.com", "Bob")
	assert.Equal(t, "bob@nyaruka.com", u.Email())
	assert.Equal(t, "Bob", u.Name())

	env := envs.NewBuilder().Build()
	assert.Equal(t, map[string]types.XValue{
		"__default__": types.NewXText("bob@nyaruka.com"),
		"email":       types.NewXText("bob@nyaruka.com"),
		"name":        types.NewXText("Bob"),
	}, u.Context(env))

	marshaled, err := jsonx.Marshal(u)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"email":"bob@nyaruka.com","name":"Bob"}`), marshaled)

	// unmarshal from JSON object
	u2 := &flows.User{}
	err = jsonx.Unmarshal([]byte(`{"email":"jim@nyaruka.com","name":"Jim"}`), u2)
	assert.NoError(t, err)
	assert.Equal(t, flows.NewUser("jim@nyaruka.com", "Jim"), u2)

	// error if email missing
	u3 := &flows.User{}
	err = jsonx.Unmarshal([]byte(`{"email":"","name":"Jim"}`), u3)
	assert.EqualError(t, err, "field 'email' is required")

	// can also unmarshal from string as email (triggers will look like this on old sessions)
	u4 := &flows.User{}
	err = jsonx.Unmarshal([]byte(`"eve@nyaruka.com"`), u4)
	assert.NoError(t, err)
	assert.Equal(t, flows.NewUser("eve@nyaruka.com", ""), u4)
}
