package types_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/stretchr/testify/assert"
)

func TestXText(t *testing.T) {
	env := envs.NewBuilder().Build()

	// test equality
	assert.True(t, types.NewXText("abc").Equals(types.NewXText("abc")))
	assert.False(t, types.NewXText("abc").Equals(types.NewXText("def")))

	// test comparison
	assert.Equal(t, 0, types.NewXText("abc").Compare(types.NewXText("abc")))
	assert.Equal(t, 1, types.NewXText("def").Compare(types.NewXText("abc")))
	assert.Equal(t, -1, types.NewXText("abc").Compare(types.NewXText("def")))

	// test length
	assert.Equal(t, 0, types.NewXText("").Length())
	assert.Equal(t, 3, types.NewXText("abc").Length())
	assert.Equal(t, 2, types.NewXText("世界").Length())
	assert.Equal(t, 1, types.NewXText("😁").Length())

	assert.Equal(t, `abc`, types.NewXText("abc").Render())
	assert.Equal(t, `abc`, types.NewXText("abc").Format(env))
	assert.Equal(t, `XText("abc")`, types.NewXText("abc").String())

	// unmarshal
	var val types.XText
	err := jsonx.Unmarshal([]byte(`"hello"`), &val)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXText("hello"), val)
}
