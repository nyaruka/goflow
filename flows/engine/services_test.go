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

	classificationSvc, err := eng.Services().Classification(nil, nil)
	assert.EqualError(t, err, "no classification service factory configured")
	assert.Nil(t, classificationSvc)

	airtimeSvc, err := eng.Services().Airtime(nil)
	assert.EqualError(t, err, "no airtime service factory configured")
	assert.Nil(t, airtimeSvc)
}
