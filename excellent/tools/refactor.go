package tools

import (
	"strings"

	"github.com/nyaruka/goflow/excellent"
)

// RefactorTemplate refactors the passed in template
func RefactorTemplate(template string, allowedTopLevels []string) (string, error) {
	buf := &strings.Builder{}

	err := excellent.VisitTemplate(template, allowedTopLevels, func(tokenType excellent.XTokenType, token string) error {
		switch tokenType {
		case excellent.BODY:
			buf.WriteString(token)
		case excellent.IDENTIFIER, excellent.EXPRESSION:
			refactored, err := refactorExpression(token)

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

// RefactorTemplate refactors the passed in template
func refactorExpression(expression string) (string, error) {
	parsed, err := excellent.Parse(expression, nil)
	if err != nil {
		return "", err
	}

	return parsed.String(), nil
}

func wrapExpression(tokenType excellent.XTokenType, token string) string {
	if tokenType == excellent.IDENTIFIER {
		return "@" + token
	}
	return "@(" + token + ")"
}
