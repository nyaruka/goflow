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
func EvaluateExpression(env utils.Environment, resolver types.Resolvable, template string) (interface{}, error) {
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

	visitor := NewVisitor(env, resolver)
	value := visitor.Visit(tree)

	err, isErr := value.(error)

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
func EvaluateTemplate(env utils.Environment, resolver types.Resolvable, template string, allowedTopLevels []string) (interface{}, error) {
	var buf bytes.Buffer
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we had one, then just return our string evaluation strategy
	if nextTT != EOF {
		return EvaluateTemplateAsString(env, resolver, template, false, allowedTopLevels)
	}

	switch tokenType {
	case IDENTIFIER:
		value := ResolveVariable(env, resolver, token)

		// didn't find it, our value is empty string
		if value == nil {
			value = ""
		}

		err, isErr := value.(error)

		// we got an error, return our raw value
		if isErr {
			buf.WriteString("@")
			buf.WriteString(token)
			return buf.String(), err
		}

		// found it, return that value
		return value, nil

	case EXPRESSION:
		value, err := EvaluateExpression(env, resolver, token)
		if err != nil {
			return buf.String(), err
		}

		return value, nil
	}

	// different type of token, return the string representation
	return EvaluateTemplateAsString(env, resolver, template, false, allowedTopLevels)
}

// EvaluateTemplateAsString evaluates the passed in template returning the string value of its execution
func EvaluateTemplateAsString(env utils.Environment, resolver types.Resolvable, template string, urlEncode bool, allowedTopLevels []string) (string, error) {
	var buf bytes.Buffer
	var errors TemplateErrors
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER:
			value := ResolveVariable(env, resolver, token)

			// didn't find it, our value is empty string
			if value == nil {
				value = ""
			}
			err, isErr := value.(error)

			// we got an error, return our raw variable
			if isErr {
				errors = append(errors, err)
			} else {
				strValue, _ := types.ToString(env, value)
				if urlEncode {
					strValue = url.QueryEscape(strValue)
				}

				buf.WriteString(strValue)
			}
		case EXPRESSION:
			value, err := EvaluateExpression(env, resolver, token)
			if err != nil {
				errors = append(errors, err)
			} else {
				strValue, _ := types.ToString(env, value)
				if urlEncode {
					strValue = url.QueryEscape(strValue)
				}

				buf.WriteString(strValue)
			}

		}
	}

	if len(errors) > 0 {
		return buf.String(), errors
	}
	return buf.String(), nil
}

// ResolveVariable will resolve the passed in string variable given in dot notation and return
// the value as defined by the VariableResolver passed in.
//
// Example syntaxes:
//      foo.bar.0  - 0th element of bar slice within foo, could also be "0" key in bar map within foo
//      foo.bar[0] - same as above
func ResolveVariable(env utils.Environment, variable interface{}, key string) interface{} {
	var err error

	err, isErr := variable.(error)
	if isErr {
		return err
	}

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

		if types.IsNil(variable) {
			return fmt.Errorf("can't resolve key '%s' of nil", key)
		}

		// is our key numeric?
		index, err := strconv.Atoi(key)
		if err == nil {
			indexable, isIndexable := variable.(types.Indexable)
			if isIndexable {
				if index >= indexable.Length() || index < -indexable.Length() {
					return fmt.Errorf("index %d out of range for %d items", index, indexable.Length())
				}
				if index < 0 {
					index += indexable.Length()
				}
				variable = indexable.Index(index)
				continue
			}
		}

		resolver, isResolver := variable.(types.Resolvable)

		// look it up in our resolver
		if isResolver {
			variable = resolver.Resolve(key)

			err, isErr := variable.(error)
			if isErr {
				return err
			}

		} else {
			return fmt.Errorf("can't resolve key '%s' of type %s", key, reflect.TypeOf(variable))
		}
	}

	// check what we are returning is a type that expressions understand
	_, _, err = types.ToXAtom(env, variable)
	if err != nil {
		_, isAtomizable := variable.(types.Atomizable)
		if !isAtomizable {
			panic(fmt.Sprintf("key '%s' of resolved to usupported type %s", key, reflect.TypeOf(variable)))
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
