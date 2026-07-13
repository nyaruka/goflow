package utils_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestExtractJSON(t *testing.T) {
	tcs := []struct {
		body []byte
		json []byte
	}{
		{[]byte(`{`), nil}, // invalid JSON
		{[]byte(`"x"`), []byte(`"x"`)},
		{[]byte(`{"foo": ["x"]}`), []byte(`{"foo": ["x"]}`)},
		{[]byte("\"a\x80\x81b\""), []byte(`"ab"`)},                     // invalid UTF-8 sequences stripped
		{[]byte("\u0000{\"foo\": 123\u0000}"), []byte(`{"foo": 123}`)}, // null chars stripped
		{[]byte(`"a\u0000b"`), []byte(`"ab"`)},                         // escaped null chars stripped
		{[]byte(`"01\02\03"`), nil},                                    // \0 not valid JSON escape
		{[]byte(`"01\\02\\03"`), []byte(`"01\\02\\03"`)},
	}

	for _, tc := range tcs {
		actual := utils.ExtractJSON(tc.body)
		assert.Equal(t, string(tc.json), string(actual), "extracted JSON mismatch for %s", string(tc.body))
	}

	asXValue := types.JSONToXValue([]byte(`{"foo": "01\\02\\03"}`))
	asXObject := asXValue.(*types.XObject)
	foo, _ := asXObject.Get("foo")
	assert.Equal(t, types.NewXText(`01\02\03`), foo)
	assert.Equal(t, `"01\\02\\03"`, string(jsonx.MustMarshal(foo)))
}
