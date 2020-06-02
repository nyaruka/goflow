package dtone_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/services/airtime/dtone"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)
	defer dates.SetNowSource(dates.DefaultNowSource)

	mocks := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://airtime-api.dtone.com/cgi-bin/shop/topup": {
			httpx.NewMockResponse(200, nil, "info_txt=pong\r\n"),                  // successful ping
			httpx.NewMockResponse(400, nil, "error_code=1\r\nerror_txt=Oops\r\n"), // unsuccessful ping
			httpx.NewMockResponse(200, nil, withCRLF(msisdnResponse)),             // successful msdninfo query
			httpx.NewMockResponse(200, nil, "xxx=yyy\r\n"),                        // unexpected response to msdninfo query
			httpx.NewMockResponse(200, nil, withCRLF(reserveResponse)),            // successful reserve ID request
			httpx.NewMockResponse(200, nil, "xxx=yyy\r\n"),                        // unexpected response to reserve ID request
			httpx.NewMockResponse(200, nil, withCRLF(topupResponse)),              // successful topup request
			httpx.NewMockResponse(200, nil, "xxx=yyy\r\n"),                        // unexpected response to topup request
			httpx.MockConnectionError,                                             // timeout
		},
	})

	httpx.SetRequestor(mocks)
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 9, 15, 25, 30, 123456789, time.UTC)))

	cl := dtone.NewClient(http.DefaultClient, nil, "joe", "1234567")

	// test ping action
	trace, err := cl.Ping()
	assert.NoError(t, err)
	test.AssertSnapshot(t, "ping_request", string(trace.RequestTrace))

	// test when ping returns error
	_, err = cl.Ping()
	assert.EqualError(t, err, "DTOne API request failed: Oops (1)")

	// test MSISDN info query
	info, trace, err := cl.MSISDNInfo("+593970000001", "USD", "1")
	assert.NoError(t, err)
	test.AssertSnapshot(t, "msisdn_request", string(trace.RequestTrace))
	assert.Equal(t, "Ecuador", info.Country)

	// test MSISDN info query when response is wrong format
	info, trace, err = cl.MSISDNInfo("+593970000001", "USD", "1")
	assert.EqualError(t, err, "DTOne API request failed: field 'destination_currency' is required, field 'product_list' is required, field 'local_info_value_list' is required")
	assert.Nil(t, info)

	// test reserve ID action
	reservedID, trace, err := cl.ReserveID()
	assert.NoError(t, err)
	test.AssertSnapshot(t, "reserve_id_request", string(trace.RequestTrace))
	assert.Equal(t, 123456789, reservedID.ReservedID)

	// test reserve ID action when response is wrong format
	reservedID, _, err = cl.ReserveID()
	assert.EqualError(t, err, "DTOne API request failed: field 'reserved_id' is required")
	assert.Nil(t, reservedID)

	// test topup action
	topup, trace, err := cl.Topup(123455, "593999000001", "593999000002", "1", "2")
	assert.NoError(t, err)
	test.AssertSnapshot(t, "topup_request", string(trace.RequestTrace))
	assert.Equal(t, decimal.RequireFromString("1"), topup.ActualProductSent)

	// test topup action when response is wrong format
	topup, _, err = cl.Topup(123455, "593999000001", "593999000002", "1", "2")
	assert.EqualError(t, err, "DTOne API request failed: field 'destination_currency' is required")
	assert.Nil(t, topup)

	// test timeout still gives us a trace
	trace, err = cl.Ping()
	assert.EqualError(t, err, "unable to connect to server")
	assert.Equal(t, "POST /cgi-bin/shop/topup HTTP/1.1\r\nHost: airtime-api.dtone.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 76\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept-Encoding: gzip\r\n\r\naction=ping&key=1570634754123&login=joe&md5=9b97b9694adede6840de4e8056245f6d", string(trace.RequestTrace))
	assert.Equal(t, "", string(trace.ResponseTrace))

	assert.False(t, mocks.HasUnused())
}
