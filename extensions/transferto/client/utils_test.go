package client_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/extensions/transferto/client"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCSVStrings(t *testing.T) {
	s := &struct {
		List1 client.CSVStrings `json:"list1"`
		List2 client.CSVStrings `json:"list2"`
	}{}
	err := json.Unmarshal([]byte(`{"list1":"foo","list2":"foo,bar"}`), s)
	assert.NoError(t, err)
	assert.Equal(t, client.CSVStrings{"foo"}, s.List1)
	assert.Equal(t, client.CSVStrings{"foo", "bar"}, s.List2)

	// try with invalid JSON
	err = json.Unmarshal([]byte(`{,`), s)
	assert.Error(t, err)
}

func TestCSVDecimals(t *testing.T) {
	s := &struct {
		List1 client.CSVDecimals `json:"list1"`
		List2 client.CSVDecimals `json:"list2"`
	}{}
	err := json.Unmarshal([]byte(`{"list1":"12.34","list2":"12.34,56.78"}`), s)
	assert.NoError(t, err)
	assert.Equal(t, client.CSVDecimals{decimal.RequireFromString("12.34")}, s.List1)
	assert.Equal(t, client.CSVDecimals{decimal.RequireFromString("12.34"), decimal.RequireFromString("56.78")}, s.List2)

	// try with invalid JSON
	err = json.Unmarshal([]byte(`{,`), s)
	assert.Error(t, err)
}
