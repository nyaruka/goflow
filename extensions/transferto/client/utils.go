package client

import (
	"encoding/json"
	"strings"

	"github.com/shopspring/decimal"
)

// CSVStringList is a list of strings which can be automatically unmarshalled from a CSV list
type CSVStringList []string

// UnmarshalJSON unmarshals this list from a CSV string
func (l *CSVList) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}
	*l = strings.Split(asString, ",")
	return nil
}

// CSVDecimalList is a list of decimals which can be automatically unmarshalled from a CSV list
type CSVDecimalList []decimal.Decimal

// UnmarshalJSON unmarshals this list from a CSV string
func (l *CSVDecimalList) UnmarshalJSON(data []byte) error {
	var asStrings CSVStringList
	if err := json.Unmarshal(data, &asStrings); err != nil {
		return err
	}

	vals := make([]decimal.Decimal, len(asStrings))
	for v := range asStrings {
		vals[v], err = decimal.NewFromString(asStrings[v])
		if err != nil {
			return err
		}
	}

	*l = vals
	return nil
}
