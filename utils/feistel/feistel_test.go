package feistel_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils/feistel"
	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecode(t *testing.T) {
	keys := []int64{0xA3B1C, 0xD2E3F, 0x1A2B3, 0xC0FFEE}
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
		actual, err := feistel.Encode(tc.id, keys)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual, "encoding mismatch for id %d", tc.id)

		decoded, err := feistel.Decode(tc.expected, keys)
		assert.NoError(t, err)
		assert.Equal(t, tc.id, decoded, "decoding mismatch for %s", tc.expected)
	}

	_, err := feistel.Decode("E2E6MXXX", keys) // too long
	assert.EqualError(t, err, "code must be 6 or 7 characters")

	_, err = feistel.Decode("E2E6M", keys) // too short
	assert.EqualError(t, err, "code must be 6 or 7 characters")

	_, err = feistel.Decode("E2E6M0", keys) // invalid char 0
	assert.EqualError(t, err, "invalid character in code")
}
