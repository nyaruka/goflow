package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestCountryCodeFromName(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{"Rwanda", "RW"},
		{"United States of America", "US"},
		{"United States", "US"},
		{"United Kingdom", "GB"},
		{"Ivory Coast", "CI"},
		{"Democratic Republic of the Congo", "CD"},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.code, utils.CountryCodeFromName(tc.name))
	}
}
