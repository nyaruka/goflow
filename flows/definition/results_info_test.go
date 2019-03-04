package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultsInfo(t *testing.T) {
	r := newResultsInfo([]string{"Age", "Response 1", "Response-1"})

	assert.Equal(t, resultsInfo(map[string][]string{"age": {"Age"}, "response_1": {"Response 1", "Response-1"}}), r)
}
