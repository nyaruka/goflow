package expressions

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/legacy/gen"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// ContextTopLevels are the allowed top-level identifiers in legacy expressions, i.e. @contact.bar is valid but @foo.bar isn't
var ContextTopLevels = []string{"channel", "child", "contact", "date", "extra", "flow", "parent", "step"}

var functionReturnTypes = map[string]string{
	"abs":                 "number",
	"datetime_add":        "datetime",
	"datetime_from_parts": "datetime",
	"datetime":            "datetime",
	"date":                "date",
	"format_date":         "date",
	"max":                 "number",
	"mean":                "number",
	"min":                 "number",
	"mod":                 "number",
	"now":                 "datetime",
	"sum":                 "number",
	"rand":                "number",
	"round":               "number",
	"round_down":          "number",
	"round_up":            "number",
	"time":                "time",
	"time_from_parts":     "time",
	"today":               "date",
}

// MigrateOptions are options for how expressions are migrated
type MigrateOptions struct {
	DefaultToSelf bool
	URLEncode     bool
	RawDates      bool
}

var defaultOptions = &MigrateOptions{DefaultToSelf: false, URLEncode: false, RawDates: false}

// MigrateTemplate will take a legacy expression and translate it to the new syntax
func MigrateTemplate(template string, options *MigrateOptions) (string, error) {
	if options == nil {
		options = defaultOptions
	}

	return migrateLegacyTemplateAsString(template, options)
}

func migrateLegacyTemplateAsString(template string, options *MigrateOptions) (string, error) {
	var buf bytes.Buffer
	scanner := excellent.NewXScanner(strings.NewReader(template), ContextTopLevels)
	scanner.SetUnescapeBody(false)
	errors := excellent.NewTemplateErrors()

	for tokenType, token := scanner.Scan(); tokenType != excellent.EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case excellent.BODY:
			buf.WriteString(token)
		case excellent.IDENTIFIER:
			value := MigrateContextReference(token, options.RawDates)

			var errorAs string
			if options.DefaultToSelf {
				errorAs = "@" + token
			}

			// optionally wrap expression so that it is URL encoded or defaults to itself on error
			buf.WriteString(wrapRawExpression(value, errorAs, options.URLEncode))

		case excellent.EXPRESSION:
			// special case of @("") which was a common workaround for the editor requiring a
			// non-empty string, but is no longer needed and can be replaced by an empty string
			if token == `""` {
				continue
			}

			value, err := migrateExpression(nil, token, options)
			if err != nil {
				errors.Add(fmt.Sprintf("@(%s)", token), err.Error())
				buf.WriteString("@(")
				buf.WriteString(token)
				buf.WriteString(")")
			} else {
				var errorAs string
				if options.DefaultToSelf {
					errorAs = "@(" + token + ")"
				}

				// optionally wrap expression so that it is URL encoded or defaults to itself on error
				buf.WriteString(wrapRawExpression(value, errorAs, options.URLEncode))
			}
		}
	}

	if errors.HasErrors() {
		return buf.String(), errors
	}
	return buf.String(), nil
}

// migrates an old expression into a new format expression
func migrateExpression(env envs.Environment, expression string, options *MigrateOptions) (string, error) {
	errListener := excellent.NewErrorListener(expression)

	input := antlr.NewInputStream(expression)
	lexer := gen.NewExcellent1Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellent1Parser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)

	// speed up parsing
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)

	tree := p.Parse()

	// if we ran into errors parsing, return the first one
	if len(errListener.Errors()) > 0 {
		return "", errListener.Errors()[0]
	}

	visitor := newLegacyVisitor(env, options)
	value := visitor.Visit(tree)
	err, isErr := value.(error)

	// did our evaluation result in an error? return that
	if isErr {
		return "", err
	}

	// all is good, return our value
	return value.(string), nil
}

var functionCallRegex = regexp.MustCompile(`^(\w+)\(`)

func inferType(operand string) string {
	// if we have an integer literal, we're a number
	_, numErr := strconv.Atoi(operand)
	if numErr == nil {
		return "number"
	}

	// if this looks like a function call, lookup its return type
	matches := functionCallRegex.FindStringSubmatch(operand)
	if matches != nil {
		return functionReturnTypes[matches[1]]
	}
	return ""
}

var identifierRegex = regexp.MustCompile(`^\pL+[\pL\pN_.]*$`)

func isValidIdentifier(expression string) bool {
	if !identifierRegex.MatchString(expression) {
		return false
	}

	for _, topLevel := range flows.RunContextTopLevels {
		if strings.HasPrefix(expression, topLevel+".") || expression == topLevel {
			return true
		}
	}

	return false
}

// takes a raw expression and wraps it for inclusion in a template, e.g. now() -> @(now())
func wrapRawExpression(expression string, errorAs string, urlEncode bool) string {
	if errorAs != "" {
		expression = fmt.Sprintf(`if(is_error(%s), %s, %s)`, expression, strconv.Quote(errorAs), expression)
	}

	if urlEncode {
		expression = wrap(expression, "url_encode")
	}

	if !isValidIdentifier(expression) {
		expression = "(" + expression + ")"
	}

	return "@" + expression
}

func wrap(expression, funcName string) string {
	return fmt.Sprintf("%s(%s)", funcName, expression)
}

// MigrateStringLiteral migrates a string literal (legacy expressions use Excel "" escaping)
func MigrateStringLiteral(s string) string {
	// strip surrounding quotes
	s = s[1 : len(s)-1]

	// replace any escaped quotes
	s = strings.Replace(s, `""`, `\"`, -1)

	// re-quote
	return `"` + s + `"`
}
