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
		{`@(upper("abc"))`, []scannedToken{{excellent.EXPRESSION_PARENTHESES, `upper("abc")`}}},
		{` @(upper("abc")) `, []scannedToken{{excellent.BODY, " "}, {excellent.EXPRESSION_PARENTHESES, `upper("abc")`}, {excellent.BODY, " "}}},
		{`@(`, []scannedToken{{excellent.BODY, `@(`}}},
		{`@(")`, []scannedToken{{excellent.BODY, `@(")`}}},
		{`@(")")`, []scannedToken{{excellent.EXPRESSION_PARENTHESES, `")"`}}},
		{`@("zz\"zz")`, []scannedToken{{excellent.EXPRESSION_PARENTHESES, `"zz\"zz"`}}},
		{`{{upper("abc")}}`, []scannedToken{{excellent.EXPRESSION_BRACES, `upper("abc")`}}},
		{`Hi {{ contact }}!`, []scannedToken{{excellent.BODY, "Hi "}, {excellent.EXPRESSION_BRACES, " contact "}, {excellent.BODY, "!"}}},
		{`{{(((1)))}}`, []scannedToken{{excellent.EXPRESSION_BRACES, `(((1)))`}}},
		{`{{"}}"}}`, []scannedToken{{excellent.EXPRESSION_BRACES, `"}}"`}}},
		{`{{"zz\"}}"}}`, []scannedToken{{excellent.EXPRESSION_BRACES, `"zz\"}}"`}}},
		{`{{x}}}`, []scannedToken{{excellent.EXPRESSION_BRACES, `x`}, {excellent.BODY, "}"}}},
		{`{{a}b}}`, []scannedToken{{excellent.EXPRESSION_BRACES, `a}b`}}},
		{`{{`, []scannedToken{{excellent.BODY, `{{`}}},
		{`{{"}}`, []scannedToken{{excellent.BODY, `{{"}}`}}},
		{`{x} {`, []scannedToken{{excellent.BODY, `{x} {`}}},
		{`a{{b}}c{{d}}`, []scannedToken{{excellent.BODY, "a"}, {excellent.EXPRESSION_BRACES, "b"}, {excellent.BODY, "c"}, {excellent.EXPRESSION_BRACES, "d"}}},
		{`@contact{{x}}@(y)`, []scannedToken{{excellent.IDENTIFIER, "contact"}, {excellent.EXPRESSION_BRACES, "x"}, {excellent.EXPRESSION_PARENTHESES, "y"}}},
		{`\{{x}}`, []scannedToken{{excellent.BODY, `{{x}}`}}},
		{`a \{{x}} {{y}}`, []scannedToken{{excellent.BODY, `a {{x}} `}, {excellent.EXPRESSION_BRACES, "y"}}},
		{`\\{{x}}`, []scannedToken{{excellent.BODY, `\`}, {excellent.EXPRESSION_BRACES, "x"}}},
		{`\{x} \\x \`, []scannedToken{{excellent.BODY, `\{x} \\x \`}}},
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

func TestScannerWithoutBraces(t *testing.T) {
	scanner := excellent.NewXScanner(strings.NewReader(`{{x}} \{{y}} @(z)`), nil)
	scanner.SetBraces(false)

	tokens := make([]scannedToken, 0)
	for tokenType, value := scanner.Scan(); tokenType != excellent.EOF; tokenType, value = scanner.Scan() {
		tokens = append(tokens, scannedToken{tokenType, value})
	}

	assert.Equal(t, []scannedToken{{excellent.BODY, `{{x}} \{{y}} `}, {excellent.EXPRESSION_PARENTHESES, "z"}}, tokens)
}

func TestScannerWithoutUnescaping(t *testing.T) {
	tests := []struct {
		input  string
		tokens []scannedToken
	}{
		{`@@x`, []scannedToken{{excellent.BODY, "@@x"}}},
		{`\{{x}}`, []scannedToken{{excellent.BODY, `\{{x}}`}}},
		{`\\{{x}}`, []scannedToken{{excellent.BODY, `\\`}, {excellent.EXPRESSION_BRACES, "x"}}},
	}

	for _, test := range tests {
		scanner := excellent.NewXScanner(strings.NewReader(test.input), nil)
		scanner.SetUnescapeBody(false)

		tokens := make([]scannedToken, 0)
		for tokenType, value := scanner.Scan(); tokenType != excellent.EOF; tokenType, value = scanner.Scan() {
			tokens = append(tokens, scannedToken{tokenType, value})
		}

		assert.Equal(t, test.tokens, tokens, "scan failed for input %s", test.input)
	}
}
