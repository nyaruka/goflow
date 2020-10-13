package cases

import (
	"strings"

	"github.com/nyaruka/goflow/envs"

	"github.com/shopspring/decimal"
)

var altNumerals = map[rune]rune{
	// Eastern Arabic
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

	return decimal.NewFromString(cleaned)
}
