package dtone_test

import (
	"testing"

	"github.com/nyaruka/goflow/providers/airtime/dtone"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	svc := dtone.NewService("login", "token", "RWF")

	assert.NotNil(t, svc)
}
