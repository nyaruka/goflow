package excellent_test

import (
	"strings"
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/stretchr/testify/assert"
)

type scannedToken struct {
	tokenType excellent.XTokenType
	value     string
}

func TestScanner(t *testing.T) {
	tests := []struct {
		input  string
		tokens []scannedToken
	}{
		{`@contact`, []scannedToken{{excellent.IDENTIFIER, "contact"}}},
		{`Hi @contact how are you?`, []scannedToken{{excellent.BODY, "Hi "}, {excellent.IDENTIFIER, "contact"}, {excellent.BODY, " how are you?"}}},
		{`@contact...?`, []scannedToken{{excellent.IDENTIFIER, "contact"}, {excellent.BODY, "...?"}}},
		{`My Twitter is @bob`, []scannedToken{{excellent.BODY, "My Twitter is "}, {excellent.BODY, "@bob"}}},
		{`@(upper("abc"))`, []scannedToken{{excellent.EXPRESSION, `upper("abc")`}}},
		{` @(upper("abc")) `, []scannedToken{{excellent.BODY, " "}, {excellent.EXPRESSION, `upper("abc")`}, {excellent.BODY, " "}}},
		{`@(`, []scannedToken{{excellent.BODY, `@(`}}},
		{`@(")`, []scannedToken{{excellent.BODY, `@(")`}}},
		{`@(")")`, []scannedToken{{excellent.EXPRESSION, `")"`}}},
		{`@("zz\"zz")`, []scannedToken{{excellent.EXPRESSION, `"zz\"zz"`}}},
	}

	for _, test := range tests {
		scanner := excellent.NewXScanner(strings.NewReader(test.input), []string{"contact", "flow"})

		tokens := make([]scannedToken, 0)
		for tokenType, value := scanner.Scan(); tokenType != excellent.EOF; tokenType, value = scanner.Scan() {
			tokens = append(tokens, scannedToken{tokenType, value})
		}

		assert.Equal(t, test.tokens, tokens, "scan failed for input %s", test.input)
	}
}
