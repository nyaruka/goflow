package engine_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/engine"
	"github.com/stretchr/testify/assert"
)

func TestEmptyServices(t *testing.T) {
	// default engine configation provides no services for anything
	eng := engine.NewBuilder().Build()

	webhookSvc, err := eng.Services().Webhook(nil)
	assert.EqualError(t, err, "no webhook service factory configured")
	assert.Nil(t, webhookSvc)

	modelSvc, err := eng.Services().Model(nil)
	assert.EqualError(t, err, "no model service factory configured")
	assert.Nil(t, modelSvc)

	airtimeSvc, err := eng.Services().Airtime(nil)
	assert.EqualError(t, err, "no airtime service factory configured")
	assert.Nil(t, airtimeSvc)
}
