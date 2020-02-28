package dtone_test

import (
	"testing"

	"github.com/nyaruka/goflow/services/airtime/dtone"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCSVStrings(t *testing.T) {
	s := &struct {
		List1 dtone.CSVStrings `json:"list1"`
		List2 dtone.CSVStrings `json:"list2"`
	}{}
	err := jsonx.Unmarshal([]byte(`{"list1":"foo","list2":"foo,bar"}`), s)
	assert.NoError(t, err)
	assert.Equal(t, dtone.CSVStrings{"foo"}, s.List1)
	assert.Equal(t, dtone.CSVStrings{"foo", "bar"}, s.List2)

	// try with invalid JSON
	err = jsonx.Unmarshal([]byte(`{"list1":true}`), s)
	assert.Error(t, err)
}

func TestCSVDecimals(t *testing.T) {
	s := &struct {
		List1 dtone.CSVDecimals `json:"list1"`
		List2 dtone.CSVDecimals `json:"list2"`
	}{}
	err := jsonx.Unmarshal([]byte(`{"list1":"12.34","list2":"12.34,56.78"}`), s)
	assert.NoError(t, err)
	assert.Equal(t, dtone.CSVDecimals{decimal.RequireFromString("12.34")}, s.List1)
	assert.Equal(t, dtone.CSVDecimals{decimal.RequireFromString("12.34"), decimal.RequireFromString("56.78")}, s.List2)

	// try with invalid JSON
	err = jsonx.Unmarshal([]byte(`{"list1":true}`), s)
	assert.Error(t, err)
}
