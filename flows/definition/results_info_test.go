package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultInfosFromNames(t *testing.T) {
	assert.Equal(t, []resultInfo{}, resultInfosFromNames(nil))
	assert.Equal(t, []resultInfo{
		{Name: "Response 1", Key: "response_1"},
		{Name: "Favorite Beer", Key: "favorite_beer"},
	}, resultInfosFromNames([]string{"Response 1", "Response-1", "Favorite Beer"}))
}
