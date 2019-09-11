package dtone_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nyaruka/goflow/providers/dtone"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(testAPIHandler))
	defer ts1.Close()

	cl := dtone.NewClient("joe", "1234567", http.DefaultClient)
	cl.SetAPIURL(ts1.URL)

	// test ping action
	err := cl.Ping()
	assert.NoError(t, err)

	// test MSISDN info query
	info, err := cl.MSISDNInfo("+593970000001", "USD", "1")
	assert.NoError(t, err)
	assert.Equal(t, "Rwanda", info.Country)

	// test reserve ID action
	reservedID, err := cl.ReserveID()
	assert.NoError(t, err)
	assert.Equal(t, 123456789, reservedID)

	// start a new test server which always returns errors
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "error_code=1\r\nerror_txt=Oops\r\n") }))
	defer ts2.Close()
	cl.SetAPIURL(ts2.URL)

	err = cl.Ping()
	assert.EqualError(t, err, "transferto API call returned an error: Oops (1)")
}

func testAPIHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	action := r.PostFormValue("action")
	switch action {
	case "ping":
		fmt.Fprint(w, "info_txt=pong\r\n")
	case "msisdn_info":
		fmt.Fprint(w, "country=Rwanda\r\n")
	case "reserve_id":
		fmt.Fprint(w, "reserved_id=123456789\r\n")
	default:
		fmt.Fprint(w, "error_code=6\r\nerror_txt=Unknown action\r\n")
	}
}
