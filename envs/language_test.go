package envs_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
)

func TestLanguage(t *testing.T) {
	lang, err := envs.ParseLanguage("ENG")
	assert.NoError(t, err)
	assert.Equal(t, envs.Language("eng"), lang)

	_, err = envs.ParseLanguage("base")
	assert.EqualError(t, err, "iso-639-3 codes must be 3 characters, got: base")

	_, err = envs.ParseLanguage("xzx")
	assert.EqualError(t, err, "unrecognized language code: xzx")

	v, err := envs.Language("eng").Value()
	assert.NoError(t, err)
	assert.Equal(t, "eng", v)

	v, err = envs.NilLanguage.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)

	var l envs.Language
	assert.NoError(t, l.Scan("eng"))
	assert.Equal(t, envs.Language("eng"), l)

	assert.NoError(t, l.Scan(nil))
	assert.Equal(t, envs.NilLanguage, l)
}
