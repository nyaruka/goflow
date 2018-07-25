package client_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/extensions/transferto"

	"github.com/stretchr/testify/assert"
)

func TestCSVList(t *testing.T) {
	s := &struct {
		List1 transferto.CSVList `json:"list1"`
		List2 transferto.CSVList `json:"list2"`
	}{}
	err := json.Unmarshal([]byte(`{"list1":"foo","list2":"foo,bar"}`), s)
	assert.NoError(t, err)
	assert.Equal(t, transferto.CSVList{"foo"}, s.List1)
	assert.Equal(t, transferto.CSVList{"foo", "bar"}, s.List2)
}
