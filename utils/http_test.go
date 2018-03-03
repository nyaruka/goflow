package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateHTTPMethod(t *testing.T) {
	for _, method := range []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"} {
		assert.NoError(t, Validate(&struct {
			Method string `validate:"http_method"`
		}{Method: method}))
	}

	assert.EqualError(t, Validate(&struct {
		Method string `validate:"http_method"`
	}{Method: "xxxx"}), "field 'Method' is not a valid HTTP method")
}
