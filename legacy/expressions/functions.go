package expressions

import (
	"fmt"
	"strings"
)

// migrates a function call in an legacy expression
type callMigrator func(funcName string, params []interface{}) (string, error)

// leaves a function call as is
func asIs() callMigrator {
	return func(funcName string, params []interface{}) (string, error) {
		return renderCall(funcName, params)
	}
}

// migrates a function call as a simple rename of the function
func asRename(newName string) callMigrator {
	return func(funcName string, params []interface{}) (string, error) {
		return renderCall(newName, params)
	}
}

// migrates a function call using a template
func asTemplate(template string) callMigrator {
	return func(funcName string, params []interface{}) (string, error) {
		numParamPlaceholders := strings.Count(template, "%s") + strings.Count(template, "%v")
		if numParamPlaceholders > len(params) {
			return "", fmt.Errorf("expecting %d params whilst migrating call to %s but got %d", numParamPlaceholders, funcName, len(params))
		}
		renderedParams, err := renderParams(params)
		if err != nil {
			return "", err
		}

		paramsAsInterfaces := make([]interface{}, len(renderedParams))
		for p := range renderedParams {
			paramsAsInterfaces[p] = renderedParams[p]
		}

		return fmt.Sprintf(template, paramsAsInterfaces...), nil
	}
}

// migrates a function call using a template based on the number of provided params
func asParamTemplates(newName string, templates map[int]string) callMigrator {
	return func(funcName string, params []interface{}) (string, error) {
		paramsTemplate, hasTemplate := templates[len(params)]
		if !hasTemplate {
			return "", fmt.Errorf("don't know how to migrate call to %s with %d parameters", funcName, len(params))
		}

		template := fmt.Sprintf("%s(%s)", newName, paramsTemplate)

		return asTemplate(template)(funcName, params)
	}
}

// migrates a function call by joining its parameters with the given delimiter
func asJoin(delimiter string) callMigrator {
	return func(funcName string, params []interface{}) (string, error) {
		renderedParams, err := renderParams(params)
		if err != nil {
			return "", err
		}
		return strings.Join(renderedParams, delimiter), nil
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
	"date":              asTemplate(`datetime("%s-%s-%s")`),
	"datedif":           asRename(`datetime_diff`),
	"datevalue":         asRename(`datetime`),
	"day":               asTemplate(`format_datetime(%s, "D")`),
	"days":              asTemplate(`datetime_diff(%s, %s, "D")`),
	"edate":             asTemplate(`datetime_add(%s, %s, "M")`),
	"exp":               asTemplate(`2.718281828459045 ^ %s`),
	"false":             asTemplate(`false`), // becomes just a keyword
	"field":             asParamTemplates(`field`, map[int]string{2: `%s, %s - 1`, 3: `%s, %s - 1, %s`}),
	"first_word":        asTemplate(`word(%s, 0)`),
	"fixed":             asParamTemplates(`format_number`, map[int]string{1: `%s`, 2: `%s, %s`, 3: `%s, %s, %s`}),
	"format_date":       asRename(`format_datetime`),
	"format_location":   asIs(),
	"hour":              asTemplate(`format_datetime(%s, "h")`),
	"if":                asIs(),
	"int":               asRename(`round_down`),
	"left":              asIs(),
	"len":               asRename(`length`),
	"lower":             asIs(),
	"max":               asIs(),
	"min":               asIs(),
	"minute":            asTemplate(`format_datetime(%s, "m")`),
	"mod":               asIs(),
	"month":             asTemplate(`format_datetime(%s, "M")`),
	"now":               asIs(),
	"or":                asIs(),
	"percent":           asIs(),
	"power":             asTemplate(`%s ^ %s`),
	"proper":            asRename(`title`),
	"rand":              asIs(),
	"randbetween":       asRename(`rand_between`),
	"read_digits":       asRename(`read_chars`),
	"regex_group":       asIs(),
	"remove_first_word": asIs(),
	"rept":              asRename(`repeat`),
	"right":             asIs(),
	"round":             asIs(),
	"rounddown":         asRename(`round_down`),
	"roundup":           asRename(`round_up`),
	"second":            asTemplate(`format_datetime(%s, "s")`),
	"substitute":        asRename(`replace`),
	"sum":               asJoin(` + `),
	"time":              asTemplate(`time(%s %s %s)`), // special case format, we sum these parts into seconds for datetime_add
	"timevalue":         asRename(`parse_datetime`),
	"today":             asIs(),
	"true":              asTemplate(`true`), // becomes just a keyword
	"trunc":             asRename(`round_down`),
	"unichar":           asRename(`char`),
	"unicode":           asRename(`code`),
	"upper":             asIs(),
	"weekday":           asIs(),
	"word_count":        asIs(),
	"word_slice":        asParamTemplates("word_slice", map[int]string{2: `%s, %s - 1`, 3: `%s, %s - 1, %s - 1`, 4: `%s, %s - 1, %s - 1`}),
	"word":              asTemplate(`word(%s, %s - 1)`),
	"year":              asTemplate(`format_datetime(%s, "YYYY")`),
}

func migrateFunctionCall(funcName string, params []interface{}) (string, error) {
	migrator, hasMigrator := callMigrators[funcName]
	if hasMigrator {
		return migrator(funcName, params)
	}

	// if we don't recognize this function, return it as it is
	return renderCall(funcName, params)
}

func renderParams(params []interface{}) ([]string, error) {
	rendered := make([]string, len(params))
	var err error

	for p := range params {
		rendered[p], err = toString(params[p])
		if err != nil {
			return nil, err
		}
	}

	return rendered, nil
}

// renders a function call from the given function name and parameters
func renderCall(funcName string, params []interface{}) (string, error) {
	renderedParams, err := renderParams(params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s)", funcName, strings.Join(renderedParams, ", ")), nil
}
