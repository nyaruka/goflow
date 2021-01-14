package cases_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/routers/cases"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestParseDecimal(t *testing.T) {
	dec := decimal.RequireFromString
	spaFormat := &envs.NumberFormat{DecimalSymbol: ",", DigitGroupingSymbol: "."}

	parseTests := []struct {
		input    string
		expected decimal.Decimal
		format   *envs.NumberFormat
	}{
		{"1", dec("1"), envs.DefaultNumberFormat},
		{"1234", dec("1234"), envs.DefaultNumberFormat},
		{"1,234.567", dec("1234.567"), envs.DefaultNumberFormat},
		{"1.234,567", dec("1234.567"), spaFormat},
		{".1234", dec("0.1234"), envs.DefaultNumberFormat},
		{" .1234 ", dec("0.1234"), envs.DefaultNumberFormat},
		{"100.00", dec("100.00"), envs.DefaultNumberFormat},

		// Eastern Arabic
		{"١", dec("1"), envs.DefaultNumberFormat},
		{"١٢٣٤", dec("1234"), envs.DefaultNumberFormat},
		{"١,٢٣٤.٥٦٧", dec("1234.567"), envs.DefaultNumberFormat},
		{"١.٢٣٤,٥٦٧", dec("1234.567"), spaFormat},
		{"٠.٨٩", dec("0.89"), envs.DefaultNumberFormat},
		{".١٢٣٤", dec("0.1234"), envs.DefaultNumberFormat},

		// Bengali
		{"১,২৩৪.৫৬৭", dec("1234.567"), envs.DefaultNumberFormat},
		{"০.৮৯", dec("0.89"), envs.DefaultNumberFormat},
	}

	for _, test := range parseTests {
		val, err := cases.ParseDecimal(test.input, test.format)

		assert.NoError(t, err)
		assert.Equal(t, test.expected, val, "parse decimal failed for input '%s'", test.input)
	}
}
