package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestValidateHTTPMethod(t *testing.T) {
	for _, method := range []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"} {
		assert.NoError(t, utils.Validate(&struct {
			Method string `validate:"http_method"`
		}{Method: method}))
	}

	assert.EqualError(t, utils.Validate(&struct {
		Method string `validate:"http_method"`
	}{Method: "xxxx"}), "field 'Method' is not a valid HTTP method")
}
