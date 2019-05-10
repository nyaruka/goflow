package test_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestAssertEqualJSON(t *testing.T) {
	assert.True(t, test.AssertEqualJSON(t, json.RawMessage(`{"foo":1,"bar":2}`), json.RawMessage(`{"bar": 2, "foo": 1}`), "doh!"))
}

func TestJSONReplace(t *testing.T) {
	assert.Equal(t, json.RawMessage(`{"foo":"x","bar":2}`), test.JSONReplace(json.RawMessage(`{"foo":1,"bar":2}`), []string{"foo"}, json.RawMessage(`"x"`)))
}
