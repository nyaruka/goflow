package services_test

import (
	"context"
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test/services"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAirtimeService(t *testing.T) {
	ctx := context.Background()
	svc := services.NewAirtime("USD")

	uuid := flows.NewEventUUID()
	transfer, err := svc.Create(ctx, uuid, urns.URN("tel:+12025550100"), urns.URN("tel:+12025550101"), map[string]decimal.Decimal{
		"USD": decimal.RequireFromString("3"),
	}, func(*flows.HTTPLog) {})
	assert.NoError(t, err)
	assert.NotEmpty(t, transfer.ExternalID)
	assert.Equal(t, "USD", transfer.Currency)

	// Confirm is a no-op on this mock — returns nil without touching the transfer
	err = svc.Confirm(ctx, transfer, func(*flows.HTTPLog) {})
	assert.NoError(t, err)

	// recipients containing "666" cause Create to fail (used by other tests)
	_, err = svc.Create(ctx, uuid, urns.NilURN, urns.URN("tel:+1666"), map[string]decimal.Decimal{"USD": decimal.RequireFromString("3")}, func(*flows.HTTPLog) {})
	assert.EqualError(t, err, "invalid recipient number")
}
