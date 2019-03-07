package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultsInfo(t *testing.T) {
	r := newResultsInfo([]string{"Age", "Response 1", "Response-1", "Response 1"})

	assert.Equal(t, resultsInfo(map[string]*resultInfo{
		"age":        {Names: []string{"Age"}},
		"response_1": {Names: []string{"Response 1", "Response-1"}},
	}), r)
}
