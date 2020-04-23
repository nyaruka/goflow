package dtone_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/airtime/dtone"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const msisdnResponse = `country=Ecuador
countryid=727
operator=Movistar Ecuador
operatorid=1472
connection_status=100
destination_msisdn=593999000001
destination_currency=USD
product_list=1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,33,40,44,50,55,60,70
retail_price_list=1.30,2.60,3.90,5.20,6.50,7.80,9.00,10.40,11.70,13.00,14.30,15.60,16.90,18.20,19.50,20.80,22.10,23.40,24.70,26.00,27.30,28.60,29.90,31.20,32.50,33.80,35.00,36.30,37.60,38.90,40.60,51.90,54.10,64.80,67.60,73.80,86.10
wholesale_price_list=0.99,1.98,2.97,3.96,4.95,5.82,6.79,7.76,8.73,9.70,10.67,11.64,12.61,13.58,14.55,15.52,16.49,17.46,18.43,19.40,20.37,21.34,22.31,23.28,24.25,25.22,26.19,27.16,28.13,29.10,32.45,38.80,43.26,48.50,54.07,58.98,68.81
local_info_value_list=1.00,2.00,3.00,4.00,5.00,6.00,7.00,8.00,9.00,10.00,11.00,12.00,13.00,14.00,15.00,16.00,17.00,18.00,19.00,20.00,21.00,22.00,23.00,24.00,25.00,26.00,27.00,28.00,29.00,30.00,33.00,40.00,44.00,50.00,55.00,60.00,70.00
local_info_amount_list=1.00,2.00,3.00,4.00,5.00,6.00,7.00,8.00,9.00,10.00,11.00,12.00,13.00,14.00,15.00,16.00,17.00,18.00,19.00,20.00,21.00,22.00,23.00,24.00,25.00,26.00,27.00,28.00,29.00,30.00,33.00,40.00,44.00,50.00,55.00,60.00,70.00
local_info_currency=USD
authentication_key=4433322221111
error_code=0
error_txt=Transaction successful
`

const reserveResponse = `reserved_id=123456789
authentication_key=4433322221111
error_code=0
error_txt=Transaction successful
`

const topupResponse = `transactionid=837765537
msisdn=a friend
destination_msisdn=593999000001
country=Ecuador
countryid=727
operator=Movistar Ecuador
operatorid=1472
reference_operator=
originating_currency=USD
destination_currency=USD
product_requested=1
actual_product_sent=1
wholesale_price=0.99
retail_price=1.30
balance=27.03
sms_sent=yes
sms=
cid1=
cid2=
cid3=
authentication_key=4433322221111
error_code=0
error_txt=Transaction successful
`

var withCRLF = func(s string) string { return strings.Replace(s, "\n", "\r\n", -1) }

func TestServiceWithSuccessfulTopup(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	mocks := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://airtime-api.dtone.com/cgi-bin/shop/topup": {
			httpx.NewMockResponse(200, nil, withCRLF(msisdnResponse)),
			httpx.NewMockResponse(200, nil, withCRLF(reserveResponse)),
			httpx.NewMockResponse(200, nil, withCRLF(topupResponse)),
		},
	})

	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))
	httpx.SetRequestor(mocks)
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 9, 15, 25, 30, 123456789, time.UTC)))

	svc := dtone.NewService(http.DefaultClient, nil, "login", "token", "USD")

	httpLogger := &flows.HTTPLogger{}

	transfer, err := svc.Transfer(
		session,
		urns.URN("tel:+593979099111"),
		urns.URN("tel:+593979099111"),
		map[string]decimal.Decimal{"USD": decimal.RequireFromString("1.5")},
		httpLogger.Log,
	)
	assert.NoError(t, err)
	assert.Equal(t, &flows.AirtimeTransfer{
		Sender:        urns.URN("tel:+593979099111"),
		Recipient:     urns.URN("tel:+593979099111"),
		Currency:      "USD",
		DesiredAmount: decimal.RequireFromString("1.5"),
		ActualAmount:  decimal.RequireFromString("1"), // closest product
	}, transfer)

	assert.Equal(t, 3, len(httpLogger.Logs))

	assert.False(t, mocks.HasUnused())
}

func TestServiceFailedTransfers(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	mocks := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://airtime-api.dtone.com/cgi-bin/shop/topup": {
			httpx.NewMockResponse(200, nil, withCRLF(msisdnResponse)),
			httpx.NewMockResponse(200, nil, withCRLF(msisdnResponse)),
		},
	})

	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))
	httpx.SetRequestor(mocks)
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 9, 15, 25, 30, 123456789, time.UTC)))

	svc := dtone.NewService(http.DefaultClient, nil, "login", "token", "USD")

	httpLogger := &flows.HTTPLogger{}

	// try when currency not configured
	transfer, err := svc.Transfer(
		session,
		urns.URN("tel:+593979099111"),
		urns.URN("tel:+593979099222"),
		map[string]decimal.Decimal{"RWF": decimal.RequireFromString("1000")},
		httpLogger.Log,
	)
	assert.EqualError(t, err, "no amount configured for transfers in USD")
	assert.NotNil(t, transfer)
	assert.Equal(t, urns.URN("tel:+593979099111"), transfer.Sender)
	assert.Equal(t, urns.URN("tel:+593979099222"), transfer.Recipient)
	assert.Equal(t, "USD", transfer.Currency)
	assert.Equal(t, decimal.Zero, transfer.DesiredAmount)
	assert.Equal(t, decimal.Zero, transfer.ActualAmount)
	assert.Equal(t, 1, len(httpLogger.Logs))

	// try when amount is smaller than minimum in currency
	transfer, err = svc.Transfer(
		session,
		urns.URN("tel:+593979099111"),
		urns.URN("tel:+593979099222"),
		map[string]decimal.Decimal{"USD": decimal.RequireFromString("0.10")},
		httpLogger.Log,
	)
	assert.EqualError(t, err, "amount requested is smaller than the minimum topup of 1 USD")
	assert.NotNil(t, transfer)
	assert.Equal(t, urns.URN("tel:+593979099111"), transfer.Sender)
	assert.Equal(t, urns.URN("tel:+593979099222"), transfer.Recipient)
	assert.Equal(t, "USD", transfer.Currency)
	assert.Equal(t, decimal.RequireFromString("0.10"), transfer.DesiredAmount)
	assert.Equal(t, decimal.Zero, transfer.ActualAmount)
	assert.Equal(t, 2, len(httpLogger.Logs))

	assert.False(t, mocks.HasUnused())
}
