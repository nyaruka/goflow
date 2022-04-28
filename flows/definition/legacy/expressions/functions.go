package expressions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// migrates a parameter value in an legacy expression
type paramMigrator func(param string) string

// migrates a function call in an legacy expression
type callMigrator func(funcName string, params []string) (string, error)

// leaves a function call as is
func asIs() callMigrator {
	return func(funcName string, params []string) (string, error) {
		return renderCall(funcName, params)
	}
}

// migrates a function call as a simple rename of the function
func asRename(newName string) callMigrator {
	return func(funcName string, params []string) (string, error) {
		return renderCall(newName, params)
	}
}

// migrates a function call using a template
func asTemplate(template string) callMigrator {
	return func(funcName string, params []string) (string, error) {
		numParamPlaceholders := strings.Count(template, "%s") + strings.Count(template, "%v")
		if numParamPlaceholders > len(params) {
			return "", errors.Errorf("expecting %d params whilst migrating call to %s but got %d", numParamPlaceholders, funcName, len(params))
		}

		paramsAsInterfaces := make([]interface{}, len(params))
		for i := range params {
			paramsAsInterfaces[i] = params[i]
		}

		return fmt.Sprintf(template, paramsAsInterfaces...), nil
	}
}

// migrates a function call by joining its parameters with the given delimiter
func asJoin(delimiter string) callMigrator {
	return func(funcName string, params []string) (string, error) {
		return strings.Join(params, delimiter), nil
	}
}

// migrates a function call using migrators for each parameter
func asParamMigrators(newName string, paramMigrators ...paramMigrator) callMigrator {
	return asParamMigratorsWithDefaults(newName, nil, paramMigrators...)
}

// migrates a function call using migrators for each parameter and also defaults for params not provided
func asParamMigratorsWithDefaults(newName string, defaults []string, paramMigrators ...paramMigrator) callMigrator {
	return func(funcName string, oldParams []string) (string, error) {
		if len(oldParams) > len(paramMigrators) {
			return "", errors.Errorf("don't know how to migrate call to %s with %d parameters", funcName, len(oldParams))
		}

		newParams := make([]string, utils.Max(len(oldParams), len(defaults)))

		for i := range newParams {
			var param string
			if i < len(oldParams) {
				param = oldParams[i]
			} else {
				param = defaults[i]
			}
			newParams[i] = paramMigrators[i](param)
		}

		return renderCall(newName, newParams)
	}
}

// migrates a parameter as is
func paramAsIs() paramMigrator {
	return func(param string) string { return param }
}

// migrates a parameter to it decremented by one
func paramDecremented() paramMigrator {
	return func(param string) string {
		// if param is a number literal then we can do the decrementing now
		asInt, err := strconv.Atoi(param)
		if err == nil {
			return strconv.Itoa(asInt - 1)
		}

		// if not return a decrementing expression
		return fmt.Sprintf("%s - 1", param)
	}
}

// migrates the by_spaces param used by several string tokenizing functions
func paramBySpaces() paramMigrator {
	return func(param string) string {
		if strings.TrimSpace(strings.ToLower(param)) == "true" {
			return `" \t"`
		}
		return `NULL`
	}
}

var callMigrators = map[string]callMigrator{
	"abs":               asIs(),
	"and":               asIs(),
	"average":           asRename(`mean`),
	"char":              asIs(),
	"clean":             asIs(),
	"code":              asIs(),
	"concatenate":       asJoin(` & `),
	"date":              asRename(`date_from_parts`),
	"datedif":           asRename(`datetime_diff`),
	"datevalue":         asRename(`date`),
	"day":               asTemplate(`format_date(%s, "D")`),
	"days":              asTemplate(`datetime_diff(%[2]s, %[1]s, "D")`),
	"edate":             asTemplate(`datetime_add(%s, %s, "M")`),
	"epoch":             asIs(),
	"exp":               asTemplate(`2.718281828459045 ^ %s`),
	"false":             asTemplate(`false`), // becomes just a keyword
	"field":             asParamMigrators(`field`, paramAsIs(), paramDecremented(), paramAsIs()),
	"first_word":        asTemplate(`word(%s, 0)`),
	"fixed":             asParamMigratorsWithDefaults(`format_number`, []string{"", "2"}, paramAsIs(), paramAsIs(), paramAsIs()),
	"format_date":       asRename(`format_datetime`),
	"format_location":   asIs(),
	"hour":              asTemplate(`format_datetime(%s, "tt")`),
	"if":                asIs(),
	"int":               asRename(`round_down`),
	"left":              asTemplate(`text_slice(%[1]s, 0, %[2]s)`),
	"len":               asRename(`text_length`),
	"lower":             asIs(),
	"max":               asIs(),
	"min":               asIs(),
	"minute":            asTemplate(`format_datetime(%s, "m")`),
	"mod":               asIs(),
	"month":             asTemplate(`format_date(%s, "M")`),
	"now":               asIs(),
	"or":                asIs(),
	"percent":           asIs(),
	"power":             asTemplate(`%s ^ %s`),
	"proper":            asRename(`title`),
	"rand":              asIs(),
	"randbetween":       asRename(`rand_between`),
	"read_digits":       asRename(`read_chars`),
	"regex_group":       asRename(`regex_match`),
	"remove_first_word": asIs(),
	"rept":              asRename(`repeat`),
	"right":             asTemplate(`text_slice(%[1]s, -%[2]s)`),
	"round":             asIs(),
	"rounddown":         asRename(`round_down`),
	"roundup":           asRename(`round_up`),
	"second":            asTemplate(`format_datetime(%s, "s")`),
	"substitute":        asRename(`replace`),
	"sum":               asJoin(` + `),
	"time":              asTemplate(`time_from_parts(%s, %s, %s)`),
	"timevalue":         asTemplate(`time(%s)`),
	"today":             asIs(),
	"true":              asTemplate(`true`), // becomes just a keyword
	"trunc":             asRename(`round_down`),
	"unichar":           asRename(`char`),
	"unicode":           asRename(`code`),
	"upper":             asIs(),
	"weekday":           asTemplate(`weekday(%s) + 1`),
	"word_count":        asParamMigrators(`word_count`, paramAsIs(), paramBySpaces()),
	"word_slice":        asParamMigrators(`word_slice`, paramAsIs(), paramDecremented(), paramDecremented(), paramBySpaces()),
	"word":              asParamMigrators(`word`, paramAsIs(), paramDecremented(), paramBySpaces()),
	"year":              asTemplate(`format_date(%s, "YYYY")`),
}

func migrateFunctionCall(funcName string, params []string) (string, error) {
	migrator, hasMigrator := callMigrators[funcName]
	if hasMigrator {
		return migrator(funcName, params)
	}

	// if we don't recognize this function, return it as it is
	return renderCall(funcName, params)
}

// renders a function call from the given function name and parameters
func renderCall(funcName string, params []string) (string, error) {
	return fmt.Sprintf("%s(%s)", funcName, strings.Join(params, ", ")), nil
}
