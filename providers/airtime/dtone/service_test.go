package dtone_test

import (
	"testing"

	"github.com/nyaruka/goflow/providers/airtime/dtone"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	provider := dtone.NewProvider("login", "token", "RWF")

	assert.NotNil(t, provider)
}
