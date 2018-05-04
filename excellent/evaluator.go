package excellent

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/excellent/gen"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// EvaluateExpression evalutes the passed in template, returning the raw value it evaluates to
func EvaluateExpression(env utils.Environment, context types.XValue, template string) (types.XValue, error) {
	errors := NewErrorListener()

	input := antlr.NewInputStream(template)
	lexer := gen.NewExcellent2Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := gen.NewExcellent2Parser(stream)
	p.AddErrorListener(errors)
	tree := p.Parse()

	// if we ran into errors parsing, bail
	if errors.HasErrors() {
		return nil, fmt.Errorf(errors.Errors())
	}

	visitor := NewVisitor(env, context)
	value := toXValue(visitor.Visit(tree))

	err, isErr := value.(types.XError)

	// did our evaluation result in an error? return that
	if isErr {
		return nil, err
	}

	// all is good, return our value
	return value, nil
}

// EvaluateTemplate tries to evaluate the passed in template into an object, this only works if the template
// is a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In cases
// which are not a single identifier or expression, we return the stringified value
func EvaluateTemplate(env utils.Environment, context types.XValue, template string, allowedTopLevels []string) (types.XValue, error) {
	var buf bytes.Buffer
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we had one, then just return our string evaluation strategy
	if nextTT != EOF {
		asStr, err := EvaluateTemplateAsString(env, context, template, false, allowedTopLevels)
		return types.NewXText(asStr), err
	}

	switch tokenType {
	case IDENTIFIER:
		value := ResolveValue(env, context, token)
		err, isErr := value.(error)

		// we got an error, return our raw value
		if isErr {
			buf.WriteString("@")
			buf.WriteString(token)
			return types.NewXText(buf.String()), err
		}

		// found it, return that value
		return value, nil

	case EXPRESSION:
		value, err := EvaluateExpression(env, context, token)
		if err != nil {
			return types.NewXText(buf.String()), err
		}

		return value, nil
	}

	// different type of token, return the string representation
	asStr, err := EvaluateTemplateAsString(env, context, template, false, allowedTopLevels)
	return types.NewXText(asStr), err
}

// EvaluateTemplateAsString evaluates the passed in template returning the string value of its execution
func EvaluateTemplateAsString(env utils.Environment, context types.XValue, template string, urlEncode bool, allowedTopLevels []string) (string, error) {
	var buf bytes.Buffer
	var errors TemplateErrors
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER:
			value := ResolveValue(env, context, token)

			// didn't find it, our value is empty string
			if value == nil {
				value = types.XTextEmpty
			}

			// we got an error, return our raw variable
			if types.IsXError(value) {
				errors = append(errors, value.(types.XError))
			} else {
				strValue, _ := types.ToXText(value)
				if urlEncode {
					strValue = types.NewXText(url.QueryEscape(strValue.Native()))
				}

				buf.WriteString(strValue.Native())
			}
		case EXPRESSION:
			value, err := EvaluateExpression(env, context, token)

			if err != nil {
				errors = append(errors, err)
			} else {
				strValue, _ := types.ToXText(value)
				if urlEncode {
					strValue = types.NewXText(url.QueryEscape(strValue.Native()))
				}

				buf.WriteString(strValue.Native())
			}
		}
	}

	if len(errors) > 0 {
		return buf.String(), errors
	}
	return buf.String(), nil
}

// ResolveValue will resolve the passed in string variable given in dot notation and return
// the value as defined by the Resolvable passed in.
//
// Example syntaxes:
//      foo.bar.0  - 0th element of bar slice within foo, could also be "0" key in bar map within foo
//      foo.bar[0] - same as above
func ResolveValue(env utils.Environment, variable types.XValue, key string) types.XValue {
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
			return types.NewXErrorf("can't resolve key '%s' of nil", key)
		}

		// is our key numeric?
		index, err := strconv.Atoi(key)
		if err == nil {
			indexable, isIndexable := variable.(types.XIndexable)
			if isIndexable {
				if index >= indexable.Length() || index < -indexable.Length() {
					return types.NewXErrorf("index %d out of range for %d items", index, indexable.Length())
				}
				if index < 0 {
					index += indexable.Length()
				}
				variable = indexable.Index(index)
				continue
			}
		}

		resolver, isResolver := variable.(types.XResolvable)

		// look it up in our resolver
		if isResolver {
			variable = resolver.Resolve(key)

			if types.IsXError(variable) {
				return variable
			}

		} else {
			return types.NewXErrorf("can't resolve key '%s' of type %s", key, reflect.TypeOf(variable))
		}
	}

	return variable
}

// popNextVariable pops the next variable off our string:
//     foo.bar.baz -> "foo", "bar.baz"
//     foo[0].bar -> "foo", "[0].baz"
//     foo.0.bar -> "foo", "0.baz"
//     [0].bar -> "0", "bar"
//     foo["my key"] -> "foo", "my key"
func popNextVariable(input string) (string, string) {
	var keyStart = 0
	var keyEnd = -1
	var restStart = -1

	for i, c := range input {
		if i == 0 && c == '[' {
			keyStart++
		} else if c == '[' {
			keyEnd = i
			restStart = i
			break
		} else if c == ']' {
			keyEnd = i
			restStart = i + 1
			break
		} else if c == '.' {
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
