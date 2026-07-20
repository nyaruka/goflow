package contactql_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/contactql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEscapeValue(t *testing.T) {
	assert.Equal(t, `""`, contactql.EscapeValue(``))
	assert.Equal(t, `"bobby tables"`, contactql.EscapeValue(`bobby tables`))
	assert.Equal(t, `"\"\" OR (id = 1)"`, contactql.EscapeValue(`"" OR (id = 1)`))
	assert.Equal(t, `"\\\"foo"`, contactql.EscapeValue(`\"foo`))
}

// TestSimplifyLongChain checks that flattening a long chain of same-op conditions is linear. Parsing
// `a or b or c ...` produces a left leaning tree, and promoting the grand children used to copy them at
// every level, which is quadratic in the length of the chain. Queries are capped at MaxConditions so this
// exercises Simplify directly, keeping the algorithm honest rather than relying on that cap. The time
// bound is very loose compared to the ~0.2s this takes; it only exists to catch a quadratic regression.
func TestSimplifyLongChain(t *testing.T) {
	const n = 100000

	var node contactql.QueryNode = contactql.NewCondition(contactql.PropertyTypeAttribute, "name", contactql.OpEqual, "x")
	for range n {
		cond := contactql.NewCondition(contactql.PropertyTypeAttribute, "name", contactql.OpEqual, "x")
		node = contactql.NewBoolCombination(contactql.BoolOperatorOr, node, cond)
	}

	start := time.Now()
	simplified := node.Simplify()
	elapsed := time.Since(start)

	// the whole chain should collapse into a single OR node
	root, ok := simplified.(*contactql.BoolCombination)
	require.True(t, ok, "simplified root should be a bool combination")
	assert.Equal(t, contactql.BoolOperatorOr, root.Operator())
	assert.Len(t, root.Children(), n+1)

	assert.Less(t, elapsed, 10*time.Second, "flattening a chain of %d conditions took %s - has it become quadratic again?", n, elapsed)
}
