package dtone_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/providers/airtime/dtone"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	provider := dtone.NewProvider(http.DefaultClient, "login", "token", "RWF")

	assert.NotNil(t, provider)
}
