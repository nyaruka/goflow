package client_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/extensions/transferto/client"

	"github.com/stretchr/testify/assert"
)

func TestCSVList(t *testing.T) {
	s := &struct {
		List1 client.CSVStrings `json:"list1"`
		List2 client.CSVStrings `json:"list2"`
	}{}
	err := json.Unmarshal([]byte(`{"list1":"foo","list2":"foo,bar"}`), s)
	assert.NoError(t, err)
	assert.Equal(t, client.CSVStrings{"foo"}, s.List1)
	assert.Equal(t, client.CSVStrings{"foo", "bar"}, s.List2)
}
