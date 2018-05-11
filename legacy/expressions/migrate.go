package expressions

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/legacy/gen"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var datePrefixes = []string{
	"today()",
	"yesterday()",
	"now()",
	"time",
	"timevalue",
}

// MigrateTemplate will take a legacy expression and translate it to the new syntax
func MigrateTemplate(template string, extraAs ExtraVarsMapping) (string, error) {
	migrationVarMapper := newMigrationVarMapper(extraAs)

	return migrateLegacyTemplateAsString(migrationVarMapper, template)
}

func migrateLegacyTemplateAsString(resolver types.XValue, template string) (string, error) {
	var buf bytes.Buffer
	scanner := excellent.NewXScanner(strings.NewReader(template), ContextTopLevels)
	errors := excellent.NewTemplateErrors()

	for tokenType, token := scanner.Scan(); tokenType != excellent.EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case excellent.BODY:
			buf.WriteString(token)
		case excellent.IDENTIFIER:
			value := excellent.ResolveValue(nil, resolver, token)
			if value == nil {
				errors.Add(fmt.Sprintf("@%s", token), "unable to map")
				buf.WriteString("@")
				buf.WriteString(token)
			} else {
				strValue, _ := toString(value)

				if strValue == token {
					buf.WriteString("@" + token)
				} else {
					// if expression has been changed, then it might need to be wrapped in @(...)
					buf.WriteString(wrapRawExpression(strValue))
				}
			}

		case excellent.EXPRESSION:
			value, err := migrateExpression(nil, resolver, token)
			buf.WriteString("@(")
			if err != nil {
				errors.Add(fmt.Sprintf("@(%s)", token), err.Error())
				buf.WriteString(token)
			} else {
				strValue, err := toString(value)
				if err != nil {
					buf.WriteString(token)
					errors.Add(fmt.Sprintf("@(%s)", token), err.Error())
				} else {
					buf.WriteString(strValue)
				}
			}
			buf.WriteString(")")
		}
	}

	if errors.HasErrors() {
		return buf.String(), errors
	}
	return buf.String(), nil
}

func toString(params interface{}) (string, error) {
	switch typed := params.(type) {
	case types.XValue:
		str, xerr := types.ToXText(typed)
		return str.Native(), xerr
	case string:
		return typed, nil

	case []interface{}:
		strArr := make([]string, len(typed))
		for i := range strArr {
			str, err := toString(typed[i])
			if err != nil {
				return "", err
			}
			strArr[i] = str
		}
		return strings.Join(strArr, ", "), nil
	default:
		panic(fmt.Sprintf("can't toString a %s", typed))
	}
}

// migrates an old expression into a new format expression
func migrateExpression(env utils.Environment, resolver types.XValue, expression string) (interface{}, error) {
	errListener := excellent.NewErrorListener(expression)

	input := antlr.NewInputStream(expression)
	lexer := gen.NewExcellent1Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellent1Parser(stream)
	p.AddErrorListener(errListener)

	// speed up parsing
	p.SetErrorHandler(antlr.NewBailErrorStrategy())
	p.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)
	// TODO: add second stage - https://github.com/antlr/antlr4/issues/192

	// this is still super slow, awaiting fixes in golang antlr
	// leaving this debug in until then
	// start := time.Now()
	tree := p.Parse()
	// timeTrack(start, "Parsing")

	// if we ran into errors parsing, return the first one
	if len(errListener.Errors()) > 0 {
		return nil, errListener.Errors()[0]
	}

	visitor := newLegacyVisitor(env, resolver)
	value := visitor.Visit(tree)
	err, isErr := value.(error)

	// did our evaluation result in an error? return that
	if isErr {
		return nil, err
	}

	// all is good, return our value
	return value, nil
}

func isDate(operand string) bool {
	for i := range datePrefixes {
		if strings.HasPrefix(operand, datePrefixes[i]) {
			return true
		}
	}
	return false
}

func wrapRawExpression(raw string) string {
	for _, topLevel := range runs.RunContextTopLevels {
		if strings.HasPrefix(raw, topLevel+".") || raw == topLevel {
			return "@" + raw
		}
	}
	return fmt.Sprintf("@(%s)", raw)
}

// convertTimeToSeconds takes a old TIME(0,2,5) like expression
// and returns the numeric value in seconds
func convertTimeToSeconds(operand string) (string, bool) {
	converted := false
	if strings.HasPrefix(operand, "time(") {
		var hours, minutes, seconds int
		parsed, _ := fmt.Sscanf(operand, "time(%d %d %d)", &hours, &minutes, &seconds)
		if parsed == 3 {
			operand = strconv.Itoa(seconds + (minutes * 60) + (hours * 3600))
			converted = true
		}
	}
	return operand, converted
}
