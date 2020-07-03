package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/buger/jsonparser"
	diff "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertXEqual is equivalent to assert.Equal for two XValue instances
func AssertXEqual(t *testing.T, expected types.XValue, actual types.XValue, msgAndArgs ...interface{}) bool {
	if !types.Equals(expected, actual) {
		return assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s", expected, actual), msgAndArgs...)
	}
	return true
}

// NormalizeJSON re-formats the given JSON
func NormalizeJSON(data json.RawMessage) ([]byte, error) {
	var asGeneric interface{}
	if err := jsonx.Unmarshal(data, &asGeneric); err != nil {
		return nil, err
	}
	return jsonx.MarshalPretty(asGeneric)
}

// AssertEqualJSON checks two JSON strings for equality
func AssertEqualJSON(t *testing.T, expected json.RawMessage, actual json.RawMessage, msg string, msgArgs ...interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}

	message := fmt.Sprintf(msg, msgArgs...)

	expectedNormalized, err := NormalizeJSON(expected)
	require.NoError(t, err, "%s: unable to normalize expected JSON: %s", message, string(expected))

	actualNormalized, err := NormalizeJSON(actual)
	require.NoError(t, err, "%s: unable to normalize actual JSON: %s", message, string(actual))

	differ := diff.New()
	diffs := differ.DiffMain(string(expectedNormalized), string(actualNormalized), false)

	if len(diffs) != 1 || diffs[0].Type != diff.DiffEqual {
		assert.Fail(t, message, differ.DiffPrettyText(diffs))
		return false
	}
	return true
}

// JSONReplace replaces a node in JSON
func JSONReplace(data json.RawMessage, path []string, value json.RawMessage) json.RawMessage {
	newData, err := jsonparser.Set(data, value, path...)
	if err != nil {
		panic("unable to replace JSON")
	}
	return newData
}

// JSONDelete deletes a node in JSON
func JSONDelete(data json.RawMessage, path []string) json.RawMessage {
	return jsonparser.Delete(data, path...)
}
