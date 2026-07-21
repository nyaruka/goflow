package cases_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/routers/cases"

	"github.com/stretchr/testify/assert"
)

func TestParseNumber(t *testing.T) {
	num := types.RequireXNumberFromString
	spaFormat := &envs.NumberFormat{DecimalSymbol: ",", DigitGroupingSymbol: "."}

	parseTests := []struct {
		input    string
		expected *types.XNumber
		format   *envs.NumberFormat
	}{
		{"1", num("1"), envs.DefaultNumberFormat},
		{"1234", num("1234"), envs.DefaultNumberFormat},
		{"1,234.567", num("1234.567"), envs.DefaultNumberFormat},
		{"1.234,567", num("1234.567"), spaFormat},
		{".1234", num("0.1234"), envs.DefaultNumberFormat},
		{" .1234 ", num("0.1234"), envs.DefaultNumberFormat},
		{"100.00", num("100.00"), envs.DefaultNumberFormat},

		// Eastern Arabic
		{"١", num("1"), envs.DefaultNumberFormat},
		{"١٢٣٤", num("1234"), envs.DefaultNumberFormat},
		{"١,٢٣٤.٥٦٧", num("1234.567"), envs.DefaultNumberFormat},
		{"١.٢٣٤,٥٦٧", num("1234.567"), spaFormat},
		{"٠.٨٩", num("0.89"), envs.DefaultNumberFormat},
		{".١٢٣٤", num("0.1234"), envs.DefaultNumberFormat},

		// Bengali
		{"১,২৩৪.৫৬৭", num("1234.567"), envs.DefaultNumberFormat},
		{"০.৮৯", num("0.89"), envs.DefaultNumberFormat},
	}

	for _, test := range parseTests {
		val, err := cases.ParseNumber(test.input, test.format)

		assert.NoError(t, err)
		assert.Equal(t, test.expected, val, "parse number failed for input '%s'", test.input)
	}

	// test that oversized numbers are rejected
	_, err := cases.ParseNumber("1234567890123456789012345678901234567", envs.DefaultNumberFormat)
	assert.EqualError(t, err, "number has too many digits")

	// test that scientific notation is rejected
	_, err = cases.ParseNumber("1e10", envs.DefaultNumberFormat)
	assert.EqualError(t, err, "not a valid number format")
}
