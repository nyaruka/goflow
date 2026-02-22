package cases

import (
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"

	"github.com/shopspring/decimal"
)

var altNumerals = map[rune]rune{
	// Arabic-Indic digits
	'٠': '0',
	'١': '1',
	'٢': '2',
	'٣': '3',
	'٤': '4',
	'٥': '5',
	'٦': '6',
	'٧': '7',
	'٨': '8',
	'٩': '9',

	// Eastern Arabic-Indic digits (Persian and Urdu)
	'۰': '0',
	'۱': '1',
	'۲': '2',
	'۳': '3',
	'۴': '4',
	'۵': '5',
	'۶': '6',
	'۷': '7',
	'۸': '8',
	'۹': '9',

	// Bengali
	'০': '0',
	'১': '1',
	'২': '2',
	'৩': '3',
	'৪': '4',
	'৫': '5',
	'৬': '6',
	'৭': '7',
	'৮': '8',
	'৯': '9',
}

func numeralMapper(r rune) rune {
	n, mapped := altNumerals[r]
	if mapped {
		return n
	}
	return r
}

// ParseDecimal parses a decimal from a string
func ParseDecimal(val string, format *envs.NumberFormat) (decimal.Decimal, error) {
	cleaned := strings.TrimSpace(val)

	// remove digit grouping symbol
	cleaned = strings.Replace(cleaned, format.DigitGroupingSymbol, "", -1)

	// replace non-period decimal symbols
	cleaned = strings.Replace(cleaned, format.DecimalSymbol, ".", -1)

	// replace non-Arabic (0-9) numerals with their equivalents
	cleaned = strings.Map(numeralMapper, cleaned)

	d, err := decimal.NewFromString(cleaned)
	if err != nil {
		return d, err
	}

	if err := types.CheckDecimalRange(d); err != nil {
		return decimal.Zero, err
	}

	return d, nil
}
