package refactor

import (
	"strings"

	"github.com/nyaruka/goflow/excellent"
)

// ContextRefRename returns a transformation function that renames context references
func ContextRefRename(from, to string) func(excellent.Expression) bool {
	return func(exp excellent.Expression) bool {
		changed := false
		exp.Visit(func(e excellent.Expression) {
			if ref, ok := e.(*excellent.ContextReference); ok && strings.EqualFold(ref.Name, from) {
				ref.Name = to
				changed = true
			}
		})
		return changed
	}
}
