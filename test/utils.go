package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/utils"

	diff "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NormalizeJSON re-formats the given JSON
func NormalizeJSON(data json.RawMessage) ([]byte, error) {
	var asMap map[string]interface{}
	if err := json.Unmarshal(data, &asMap); err != nil {
		return nil, err
	}

	return utils.JSONMarshalPretty(asMap)
}

// AssertEqualJSON checks two JSON strings for equality
func AssertEqualJSON(t *testing.T, expected json.RawMessage, actual json.RawMessage, msg string, msgArgs ...interface{}) bool {
	expectedNormalized, err := NormalizeJSON(expected)
	require.NoError(t, err)

	actualNormalized, err := NormalizeJSON(actual)
	require.NoError(t, err)

	differ := diff.New()
	diffs := differ.DiffMain(string(expectedNormalized), string(actualNormalized), false)

	if len(diffs) != 1 || diffs[0].Type != diff.DiffEqual {
		message := fmt.Sprintf(msg, msgArgs...)
		assert.Fail(t, message, differ.DiffPrettyText(diffs))
		return false
	}
	return true
}
