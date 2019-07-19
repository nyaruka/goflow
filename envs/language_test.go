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
}
