package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestXDateMarshaling(t *testing.T) {
	var date types.XDate
	err := json.Unmarshal([]byte(`"2018-04-09T17:01:30Z"`), &date)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXDate(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), date)

	// marshal
	data, err := json.Marshal(types.NewXDate(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)))
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"2018-04-09T17:01:30Z"`), data)
}
