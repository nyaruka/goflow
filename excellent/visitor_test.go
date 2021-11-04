package excellent

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/stretchr/testify/assert"
)

func TestExpressionsToString(t *testing.T) {
	assert.Equal(t, `("")`, (&Parentheses{exp: &TextLiteral{val: types.XTextEmpty}}).String())
	assert.Equal(t, `""`, (&TextLiteral{val: types.XTextEmpty}).String())
	assert.Equal(t, `"abc"`, (&TextLiteral{val: types.NewXText("abc")}).String())
	assert.Equal(t, `"don't say \"hello\""`, (&TextLiteral{val: types.NewXText(`don't say "hello"`)}).String())
	assert.Equal(t, `123.5`, (&NumberLiteral{val: types.RequireXNumberFromString(`123.5`)}).String())
	assert.Equal(t, `123`, (&NumberLiteral{val: types.RequireXNumberFromString(`123.0`)}).String())
	assert.Equal(t, `true`, (&BooleanLiteral{val: types.XBooleanTrue}).String())
	assert.Equal(t, `false`, (&BooleanLiteral{val: types.XBooleanFalse}).String())
	assert.Equal(t, `null`, (&NullLiteral{}).String())
}
