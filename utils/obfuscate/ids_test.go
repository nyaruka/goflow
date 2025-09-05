package obfuscate_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils/obfuscate"
	"github.com/stretchr/testify/assert"
)

func TestIDs(t *testing.T) {
	key := [4]int64{0xA3B1C, 0xD2E3F, 0x1A2B3, 0xC0FFEE}
	tcs := []struct {
		id       int64
		expected string
	}{
		{1, "E2E6MX"},
		{2, "3MWB69"},
		{3, "Q3GP9G"},
		{4, "U6Y6T5"},
		{5, "SJPWLU"},
		{12345, "A6YWQL"},
		{999_999_999, "KNGEUX"},
		{1_073_741_823, "NVQ26R"},
		{1_073_741_824, "GQENS3N"},
		{1_073_741_825, "NTA3479"},
		{1_073_741_826, "KYEKD42"},
		{9_999_999_999, "P7U8B2J"},
	}
	for _, tc := range tcs {
		actual, err := obfuscate.EncodeID(tc.id, key)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual, "encoding mismatch for id %d", tc.id)

		assert.True(t, obfuscate.WasID(tc.expected), "WasID should be true for %s", tc.expected)

		decoded, err := obfuscate.DecodeID(tc.expected, key)
		assert.NoError(t, err)
		assert.Equal(t, tc.id, decoded, "decoding mismatch for %s", tc.expected)
	}

	// decode error cases
	for _, id := range []string{
		"E2E6MXXX", // too long
		"E2E6M",    // too short
		"E2E6M0",   // invalid char 0
	} {
		assert.False(t, obfuscate.WasID(id), "WasID should be false for %s", id)

		_, err := obfuscate.DecodeID(id, key)
		assert.EqualError(t, err, "code is not a valid obfuscated ID")
	}
}
