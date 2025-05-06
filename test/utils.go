package test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/random"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/excellent/types"
	diff "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockUniverse mocks time and random generators for testing
func MockUniverse() func() {
	now := dates.NewSequentialNow(time.Date(2025, 5, 4, 12, 30, 45, 123456789, time.UTC), time.Second)

	dates.SetNowFunc(now)
	uuids.SetGenerator(uuids.NewSeededGenerator(123456, now))
	random.SetGenerator(random.NewSeededGenerator(123456))

	return func() {
		dates.SetNowFunc(time.Now)
		uuids.SetGenerator(uuids.DefaultGenerator)
		random.SetGenerator(random.DefaultGenerator)
	}
}

// AssertXEqual is equivalent to assert.Equal for two XValue instances
func AssertXEqual(t *testing.T, expected types.XValue, actual types.XValue, msgAndArgs ...any) bool {
	if !types.Equals(expected, actual) {
		return assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s", expected, actual), msgAndArgs...)
	}
	return true
}

// NormalizeJSON re-formats the given JSON
func NormalizeJSON(data []byte) ([]byte, error) {
	var asGeneric any
	if err := jsonx.Unmarshal(data, &asGeneric); err != nil {
		return nil, err
	}
	return jsonx.MarshalPretty(asGeneric)
}

// AssertEqualJSON checks two JSON strings for equality
func AssertEqualJSON(t *testing.T, expected []byte, actual []byte, msgAndArgs ...any) bool {
	if expected == nil && actual == nil {
		return true
	}

	message := fmtMsgAndArgs(msgAndArgs)

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
func JSONReplace(data []byte, path []string, value []byte) []byte {
	newData, err := jsonparser.Set(data, value, path...)
	if err != nil {
		panic("unable to replace JSON")
	}
	return newData
}

// JSONDelete deletes a node in JSON
func JSONDelete(data []byte, path []string) []byte {
	return jsonparser.Delete(data, path...)
}

func fmtMsgAndArgs(msgAndArgs []any) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

// RootError returns the root cause of an error by following the unwrap chain
func RootError(err error) error {
	for {
		cause := errors.Unwrap(err)
		if cause == nil {
			return err
		}
		err = cause
	}
}
