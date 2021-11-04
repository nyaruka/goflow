package tools

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/functions"
)

// FindContextRefsInTemplate audits context references in the given template. Note that the case of
// the found references is preserved as these may be significant, e.g. ["X"] vs ["x"] in JSON
func FindContextRefsInTemplate(template string, allowedTopLevels []string, callback func([]string)) error {
	// wrap callback to exclude function references
	wrapped := func(p []string) {
		if functions.Lookup(p[0]) == nil {
			callback(p)
		}
	}

	return excellent.VisitTemplate(template, allowedTopLevels, func(tokenType excellent.XTokenType, token string) error {
		switch tokenType {
		case excellent.IDENTIFIER, excellent.EXPRESSION:
			excellent.Parse(token, wrapped)
		}
		return nil
	})
}
