package dtone_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/services/airtime/dtone"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)
	defer dates.SetNowSource(dates.DefaultNowSource)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"https://airtime-api.dtone.com/cgi-bin/shop/topup": []*httpx.MockResponse{
			httpx.NewMockResponse(200, "info_txt=pong\r\n"),
			httpx.NewMockResponse(400, "error_code=1\r\nerror_txt=Oops\r\n"),
			httpx.NewMockResponse(200, "country=Rwanda\r\n"),
			httpx.NewMockResponse(200, "reserved_id=123456789\r\n"),
		},
	}))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 9, 15, 25, 30, 123456789, time.UTC)))

	cl := dtone.NewClient(http.DefaultClient, "joe", "1234567")

	// test ping action
	trace, err := cl.Ping()
	assert.NoError(t, err)
	assert.Equal(t, "POST /cgi-bin/shop/topup HTTP/1.1\r\nHost: airtime-api.dtone.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 76\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept-Encoding: gzip\r\n\r\naction=ping&key=1570634730123&login=joe&md5=dbbf7ba35d298f2772712124ef9aab4f", string(trace.RequestTrace))

	// test when ping returns error
	_, err = cl.Ping()
	assert.EqualError(t, err, "transferto API call returned an error: Oops (1)")

	// test MSISDN info query
	info, trace, err := cl.MSISDNInfo("+593970000001", "USD", "1")
	assert.NoError(t, err)
	assert.Equal(t, "POST /cgi-bin/shop/topup HTTP/1.1\r\nHost: airtime-api.dtone.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 155\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept-Encoding: gzip\r\n\r\naction=msisdn_info&currency=USD&delivered_amount_info=1&destination_msisdn=%2B593970000001&key=1570634736123&login=joe&md5=922e0934bef1f4c5fa7ed60b20ba5085", string(trace.RequestTrace))
	assert.Equal(t, "Rwanda", info.Country)

	// test reserve ID action
	reservedID, err := cl.ReserveID()
	assert.NoError(t, err)
	assert.Equal(t, 123456789, reservedID)
}
