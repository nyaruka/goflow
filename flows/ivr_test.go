package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	d := flows.NewDial(flows.DialStatusNoAnswer, 5)

	// test marshalling
	marshalled, err := jsonx.Marshal(d)
	assert.NoError(t, err)
	assert.Equal(t, `{"status":"no_answer","duration":5}`, string(marshalled))

	// and unmarsalling
	d2 := &flows.Dial{}
	err = jsonx.Unmarshal(marshalled, d2)
	assert.NoError(t, err)
	assert.Equal(t, flows.DialStatusNoAnswer, d2.Status)
	assert.Equal(t, 5, d2.Duration)

	// test status validation
	err = utils.UnmarshalAndValidate([]byte(`{"status":"broken","duration":5}`), d2)
	assert.EqualError(t, err, "field 'status' is not a valid dial status")

	// test context
	assert.Equal(t, map[string]types.XValue{
		"status":   types.NewXText("no_answer"),
		"duration": types.NewXNumberFromInt(5),
	}, d.Context(envs.NewBuilder().Build()))
}
