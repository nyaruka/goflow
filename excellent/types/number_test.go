package types_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestXNumberMarshaling(t *testing.T) {
	// unmarshal with quotes
	var num types.XNumber
	err := json.Unmarshal([]byte(`"23.45"`), &num)
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("23.45"), num)

	// unmarshal without quotes
	err = json.Unmarshal([]byte(`34.56`), &num)
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("34.56"), num)

	// marshal (doesn't use quotes)
	data, err := json.Marshal(types.RequireXNumberFromString("23.45"))
	assert.NoError(t, err)
	assert.Equal(t, []byte(`23.45`), data)
}
