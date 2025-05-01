package test_test

import (
	"testing"

	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

func TestAssertEqualJSON(t *testing.T) {
	assert.True(t, test.AssertEqualJSON(t, []byte(`{"foo":1,"bar":2}`), []byte(`{"bar": 2, "foo": 1}`), "doh!"))
}

func TestJSONReplace(t *testing.T) {
	assert.Equal(t, []byte(`{"foo":"x","bar":2}`), test.JSONReplace([]byte(`{"foo":1,"bar":2}`), []string{"foo"}, []byte(`"x"`)))
}
