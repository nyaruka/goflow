package dtone

import (
	"strings"

	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/shopspring/decimal"
)

// CSVStrings is a list of strings which can be automatically unmarshalled from a CSV list
type CSVStrings []string

// UnmarshalJSON unmarshals this list from a CSV string
func (l *CSVStrings) UnmarshalJSON(data []byte) error {
	var asString string
	if err := jsonx.Unmarshal(data, &asString); err != nil {
		return err
	}
	*l = strings.Split(asString, ",")
	return nil
}

// CSVDecimals is a list of decimals which can be automatically unmarshalled from a CSV list
type CSVDecimals []decimal.Decimal

// UnmarshalJSON unmarshals this list from a CSV string
func (l *CSVDecimals) UnmarshalJSON(data []byte) error {
	var asStrings CSVStrings
	if err := jsonx.Unmarshal(data, &asStrings); err != nil {
		return err
	}

	vals := make([]decimal.Decimal, len(asStrings))
	var err error
	for i := range asStrings {
		vals[i], err = decimal.NewFromString(asStrings[i])
		if err != nil {
			return err
		}
	}

	*l = vals
	return nil
}
