package excellent

import (
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// EvaluateExpression evalutes the passed in Excellent expression, returning the typed value it evaluates to,
// which might be an error, e.g. "2 / 3" or "contact.fields.age"
func EvaluateExpression(env utils.Environment, context types.XValue, expression string) types.XValue {
	visitor := newEvaluationVisitor(env, context)

	return toXValue(VisitExpression(expression, visitor))
}

// EvaluateTemplate tries to evaluate the passed in template into an object, this only works if the template
// is a single identifier or expression, ie: "@contact" or "@(first(contact.urns))". In cases
// which are not a single identifier or expression, we return the stringified value
func EvaluateTemplate(env utils.Environment, context types.XValue, template string, allowedTopLevels []string) (types.XValue, error) {
	template = strings.TrimSpace(template)
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)

	// parse our first token
	tokenType, token := scanner.Scan()

	// try to scan to our next token
	nextTT, _ := scanner.Scan()

	// if we only have an identifier or an expression, evaluate it on its own
	if nextTT == EOF {
		switch tokenType {
		case IDENTIFIER, EXPRESSION:
			return EvaluateExpression(env, context, token), nil
		}
	}

	// otherwise fallback to full template evaluation
	asStr, err := EvaluateTemplateAsString(env, context, template, allowedTopLevels)
	return types.NewXText(asStr), err
}

// EvaluateTemplateAsString evaluates the passed in template returning the string value of its execution
func EvaluateTemplateAsString(env utils.Environment, context types.XValue, template string, allowedTopLevels []string) (string, error) {
	buf := &strings.Builder{}

	err := VisitTemplate(template, allowedTopLevels, func(tokenType XTokenType, token string) error {
		switch tokenType {
		case BODY:
			buf.WriteString(token)
		case IDENTIFIER, EXPRESSION:
			value := EvaluateExpression(env, context, token)

			// if we got an error, return that
			if types.IsXError(value) {
				return value.(error)
			}

			// if not, stringify value and append to the output
			strValue, _ := types.ToXText(env, value)
			buf.WriteString(strValue.Native())
		}
		return nil
	})

	return buf.String(), err
}

// VisitTemplate scans the given template and calls the callback for each token encountered
func VisitTemplate(template string, allowedTopLevels []string, callback func(XTokenType, string) error) error {
	scanner := NewXScanner(strings.NewReader(template), allowedTopLevels)
	errors := NewTemplateErrors()

	for tokenType, token := scanner.Scan(); tokenType != EOF; tokenType, token = scanner.Scan() {
		err := callback(tokenType, token)
		if err != nil {
			var repr string
			if tokenType == IDENTIFIER {
				repr = "@" + token
			} else {
				repr = "@(" + token + ")"
			}

			errors.Add(repr, err.Error())
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// lookup an index on the given value
func lookupIndex(env utils.Environment, value types.XValue, index types.XNumber) types.XValue {
	indexable, isIndexable := value.(types.XIndexable)

	if !isIndexable || utils.IsNil(indexable) {
		return types.NewXErrorf("%s is not indexable", value.Describe())
	}

	indexAsInt, xerr := types.ToInteger(env, index)
	if xerr != nil {
		return xerr
	}

	if indexAsInt >= indexable.Length() || indexAsInt < -indexable.Length() {
		return types.NewXErrorf("index %d out of range for %d items", indexAsInt, indexable.Length())
	}
	if indexAsInt < 0 {
		indexAsInt += indexable.Length()
	}
	return indexable.Index(indexAsInt)
}

// lookup a named property on the given value
func lookupProperty(env utils.Environment, variable types.XValue, key string) types.XValue {
	resolver, isResolver := variable.(types.XResolvable)

	if !isResolver || utils.IsNil(resolver) {
		return types.NewXErrorf("%s has no property '%s'", types.Describe(variable), key)
	}

	return resolver.Resolve(env, key)
}
