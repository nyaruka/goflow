package expressions

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/excellent"
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
func MigrateTemplate(template string, extraAs ExtraVarsMapping, defaultToSelf bool) (string, error) {
	migrationVarMapper := newMigrationVarMapper(extraAs)

	return migrateLegacyTemplateAsString(migrationVarMapper, template, defaultToSelf)
}

func migrateLegacyTemplateAsString(resolver Resolvable, template string, defaultToSelf bool) (string, error) {
	var buf bytes.Buffer
	scanner := excellent.NewXScanner(strings.NewReader(template), ContextTopLevels)
	scanner.SetUnescapeBody(false)
	errors := excellent.NewTemplateErrors()

	for tokenType, token := scanner.Scan(); tokenType != excellent.EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case excellent.BODY:
			buf.WriteString(token)
		case excellent.IDENTIFIER:
			value := resolveLookup(nil, resolver, token)
			if value == nil {
				errors.Add(fmt.Sprintf("@%s", token), "unable to migrate variable")
				buf.WriteString("@")
				buf.WriteString(token)
			} else {
				strValue, _ := toString(value)

				var errorAs string
				if defaultToSelf {
					errorAs = "@" + token
				}

				// expression might need to be wrapped in @(...) or call to @(default(...))
				buf.WriteString(wrapRawExpression(strValue, errorAs))
			}

		case excellent.EXPRESSION:
			value, err := migrateExpression(nil, resolver, token)
			if err != nil {
				errors.Add(fmt.Sprintf("@(%s)", token), err.Error())
				buf.WriteString("@(")
				buf.WriteString(token)
				buf.WriteString(")")
			} else {
				strValue, _ := toString(value)

				var errorAs string
				if defaultToSelf {
					errorAs = "@(" + token + ")"
				}

				// expression might need to be wrapped in @(...) or call to @(default(...))
				buf.WriteString(wrapRawExpression(strValue, errorAs))
			}
		}
	}

	if errors.HasErrors() {
		return buf.String(), errors
	}
	return buf.String(), nil
}

func toString(params interface{}) (string, error) {
	switch typed := params.(type) {
	case error:
		return "", typed
	case string:
		return typed, nil
	case Resolvable:
		return typed.String(), nil
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
		panic(fmt.Sprintf("can't toString a %T %s", typed, typed))
	}
}

// migrates an old expression into a new format expression
func migrateExpression(env utils.Environment, resolver interface{}, expression string) (interface{}, error) {
	errListener := excellent.NewErrorListener(expression)

	input := antlr.NewInputStream(expression)
	lexer := gen.NewExcellent1Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellent1Parser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)

	// speed up parsing
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

func resolveLookup(env utils.Environment, variable interface{}, key string) interface{} {
	// self referencing
	if key == "" {
		return variable
	}

	// strip leading '.'
	if key[0] == '.' {
		key = key[1:]
	}

	rest := key
	for rest != "" {
		key, rest = popNextVariable(rest)

		if utils.IsNil(variable) {
			return fmt.Errorf("%s has no property '%s'", variable, key)
		}

		resolver, isResolver := variable.(Resolvable)

		// look it up in our resolver
		if isResolver {
			variable = resolver.Resolve(key)

			_, isErr := variable.(error)
			if isErr {
				return variable
			}

		} else {
			return fmt.Errorf("%s has no property '%s'", variable, key)
		}
	}

	return variable
}

// popNextVariable pops the next variable off our string:
//     foo.bar.baz -> "foo", "bar.baz"
//     foo.0.bar -> "foo", "0.baz"
func popNextVariable(input string) (string, string) {
	var keyStart = 0
	var keyEnd = -1
	var restStart = -1

	for i, c := range input {
		if c == '.' {
			keyEnd = i
			restStart = i + 1
			break
		}
	}

	if keyEnd == -1 {
		return input, ""
	}

	key := strings.Trim(input[keyStart:keyEnd], "\"")
	rest := input[restStart:]
	return key, rest
}

func isDate(operand string) bool {
	for i := range datePrefixes {
		if strings.HasPrefix(operand, datePrefixes[i]) {
			return true
		}
	}
	return false
}

var identifierRegex = regexp.MustCompile(`^\pL+[\pL\pN_.]*$`)

func isValidIdentifier(expression string) bool {
	if !identifierRegex.MatchString(expression) {
		return false
	}

	for _, topLevel := range runs.RunContextTopLevels {
		if strings.HasPrefix(expression, topLevel+".") || expression == topLevel {
			return true
		}
	}

	return false
}

// takes a raw expression and wraps it for inclusion in a template, e.g. now() -> @(now())
func wrapRawExpression(expression string, errorAs string) string {
	if errorAs != "" {
		expression = fmt.Sprintf(`if(is_error(%s), %s, %s)`, expression, strconv.Quote(errorAs), expression)
	}

	if !isValidIdentifier(expression) {
		expression = "(" + expression + ")"
	}

	return "@" + expression
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
