package refactor

import (
	"strings"

	"github.com/nyaruka/goflow/excellent"
)

// Template refactors the passed in template
func Template(template string, allowedTopLevels []string, tx func(excellent.Expression) bool) (string, error) {
	buf := &strings.Builder{}

	err := excellent.VisitTemplate(template, allowedTopLevels, false, func(tokenType excellent.XTokenType, token string) error {
		switch tokenType {
		case excellent.BODY:
			buf.WriteString(token)
		case excellent.IDENTIFIER, excellent.EXPRESSION:
			refactored, err := expression(token, tx)

			// if we got an error, return that, and rewrite original expression
			if err != nil {
				buf.WriteString(wrapExpression(tokenType, token))
				return err
			}

			// if not, append refactored expresion to the output
			buf.WriteString(wrapExpression(tokenType, refactored))
		}
		return nil
	})

	return buf.String(), err
}

// refactors the passed in expression by applying a transformation function
func expression(expression string, tx func(excellent.Expression) bool) (string, error) {
	parsed, err := excellent.Parse(expression, nil)
	if err != nil {
		return "", err
	}

	// if transformer actually changes anything, return reformatted expression
	if tx(parsed) {
		return parsed.String(), nil
	}

	// otherwise keep original
	return expression, nil
}

func wrapExpression(tokenType excellent.XTokenType, token string) string {
	if tokenType == excellent.IDENTIFIER {
		return "@" + token
	}
	return "@(" + token + ")"
}
