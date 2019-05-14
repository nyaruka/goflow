package test

import (
	"encoding/json"
	"fmt"
	"testing"

	xtest "github.com/nyaruka/goflow/excellent/test"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/buger/jsonparser"
	diff "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertXEqual(t *testing.T, expected types.XValue, actual types.XValue, msgAndArgs ...interface{}) bool {
	return xtest.AssertEqual(t, expected, actual, msgAndArgs...)
}

// NormalizeJSON re-formats the given JSON
func NormalizeJSON(data json.RawMessage) ([]byte, error) {
	var asGeneric interface{}
	if err := json.Unmarshal(data, &asGeneric); err != nil {
		return nil, err
	}
	return utils.JSONMarshalPretty(asGeneric)
}

// AssertEqualJSON checks two JSON strings for equality
func AssertEqualJSON(t *testing.T, expected json.RawMessage, actual json.RawMessage, msg string, msgArgs ...interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}

	expectedNormalized, err := NormalizeJSON(expected)
	require.NoError(t, err, "unable to normalize expected JSON: %s", string(expected))

	actualNormalized, err := NormalizeJSON(actual)
	require.NoError(t, err, "unable to normalize actual JSON: %s", string(actual))

	differ := diff.New()
	diffs := differ.DiffMain(string(expectedNormalized), string(actualNormalized), false)

	if len(diffs) != 1 || diffs[0].Type != diff.DiffEqual {
		message := fmt.Sprintf(msg, msgArgs...)
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
