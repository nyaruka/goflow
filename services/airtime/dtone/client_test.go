package dtone_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/services/airtime/dtone"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

var lookupNumberResponse = `[
	{
		"country": {
		"iso_code": "ECU",
		"name": "Ecuador",
		"regions": null
		},
		"id": 1596,
		"identified": true,
		"name": "Claro Ecuador",
		"regions": null
	},
	{
		"country": {
		"iso_code": "ECU",
		"name": "Ecuador",
		"regions": null
		},
		"id": 1597,
		"identified": false,
		"name": "CNT Ecuador",
		"regions": null
	}
]`

var productsResponse = `[
	{
		"availability_zones": [
			"INTERNATIONAL"
		],
		"benefits": [
			{
			"additional_information": null,
			"amount": {
				"base": 3,
				"promotion_bonus": 0,
				"total_excluding_tax": 3
			},
			"type": "CREDITS",
			"unit": "USD",
			"unit_type": "CURRENCY"
			}
		],
		"description": "",
		"destination": {
			"amount": 3,
			"unit": "USD",
			"unit_type": "CURRENCY"
		},
		"id": 6035,
		"name": "3 USD",
		"operator": {
			"country": {
			"iso_code": "ECU",
			"name": "Ecuador",
			"regions": null
			},
			"id": 1596,
			"name": "Claro Ecuador",
			"regions": null
		},
		"prices": {
			"retail": {
			"amount": 4,
			"fee": 0,
			"unit": "USD",
			"unit_type": "CURRENCY"
			},
			"wholesale": {
			"amount": 3.6,
			"fee": 0,
			"unit": "USD",
			"unit_type": "CURRENCY"
			}
		},
		"promotions": null,
		"rates": {
			"base": 0.833333333333333,
			"retail": 0.75,
			"wholesale": 0.833333333333333
		},
		"regions": null,
		"required_beneficiary_fields": null,
		"required_credit_party_identifier_fields": [
			[
			"mobile_number"
			]
		],
		"required_debit_party_identifier_fields": null,
		"required_sender_fields": null,
		"service": {
			"id": 1,
			"name": "Mobile"
		},
		"source": {
			"amount": 3.6,
			"unit": "USD",
			"unit_type": "CURRENCY"
		},
		"type": "FIXED_VALUE_RECHARGE",
		"validity": null
		},
		{
		"availability_zones": [
			"INTERNATIONAL"
		],
		"benefits": [
			{
			"additional_information": null,
			"amount": {
				"base": 6,
				"promotion_bonus": 0,
				"total_excluding_tax": 6
			},
			"type": "CREDITS",
			"unit": "USD",
			"unit_type": "CURRENCY"
			}
		],
		"description": "",
		"destination": {
			"amount": 6,
			"unit": "USD",
			"unit_type": "CURRENCY"
		},
		"id": 6036,
		"name": "6 USD",
		"operator": {
			"country": {
			"iso_code": "ECU",
			"name": "Ecuador",
			"regions": null
			},
			"id": 1596,
			"name": "Claro Ecuador",
			"regions": null
		},
		"prices": {
			"retail": {
			"amount": 7,
			"fee": 0,
			"unit": "USD",
			"unit_type": "CURRENCY"
			},
			"wholesale": {
			"amount": 6.3,
			"fee": 0,
			"unit": "USD",
			"unit_type": "CURRENCY"
			}
		},
		"promotions": null,
		"rates": {
			"base": 0.952380952380952,
			"retail": 0.857142857142857,
			"wholesale": 0.952380952380952
		},
		"regions": null,
		"required_beneficiary_fields": null,
		"required_credit_party_identifier_fields": [
			[
				"mobile_number"
			]
		],
		"required_debit_party_identifier_fields": null,
		"required_sender_fields": null,
		"service": {
			"id": 1,
			"name": "Mobile"
		},
		"source": {
			"amount": 6.3,
			"unit": "USD",
			"unit_type": "CURRENCY"
		},
		"type": "FIXED_VALUE_RECHARGE",
		"validity": null
	}
]`

var transactionConfirmedResponse = `{
	"benefits": [
		{
			"additional_information": null,
			"amount": {
				"base": 3,
				"promotion_bonus": 0,
				"total_excluding_tax": 3
			},
			"type": "CREDITS",
			"unit": "USD",
			"unit_type": "CURRENCY"
		}
	],
	"confirmation_date": "2021-03-24T20:05:06.111631000Z",
	"confirmation_expiration_date": "2021-03-24T21:05:05.883561000Z",
	"creation_date": "2021-03-24T20:05:05.883561000Z",
	"credit_party_identifier": {
		"mobile_number": "+593979123456"
	},
	"external_id": "EX12345",
	"id": 2237512891,
	"prices": {
		"retail": {
			"amount": 4,
			"fee": 0,
			"unit": "USD",
			"unit_type": "CURRENCY"
		},
		"wholesale": {
			"amount": 3.6,
			"fee": 0,
			"unit": "USD",
			"unit_type": "CURRENCY"
		}
	},
	"product": {
		"description": "",
		"id": 6035,
		"name": "3 USD",
		"operator": {
			"country": {
				"iso_code": "ECU",
				"name": "Ecuador",
				"regions": null
			},
			"id": 1596,
			"name": "Claro Ecuador",
			"regions": null
		},
		"regions": null,
		"service": {
			"id": 1,
			"name": "Mobile"
		},
		"type": "FIXED_VALUE_RECHARGE"
	},
	"promotions": null,
	"rates": {
		"base": 0.833333333333333,
		"retail": 0.75,
		"wholesale": 0.833333333333333
	},
	"status": {
		"class": {
			"id": 2,
			"message": "CONFIRMED"
		},
		"id": 20000,
		"message": "CONFIRMED"
	}
}`

var transactionRejectedResponse = `{
	"benefits": [
		{
			"additional_information": null,
			"amount": {
				"base": 3,
				"promotion_bonus": 0,
				"total_excluding_tax": 3
			},
			"type": "CREDITS",
			"unit": "USD",
			"unit_type": "CURRENCY"
		}
	],
	"confirmation_date": "2021-03-24T20:05:06.111631000Z",
	"confirmation_expiration_date": "2021-03-24T21:05:05.883561000Z",
	"creation_date": "2021-03-24T20:05:05.883561000Z",
	"credit_party_identifier": {
		"mobile_number": "+593979123456"
	},
	"external_id": "EX12345",
	"id": 2237512891,
	"prices": {
		"retail": {
			"amount": 4,
			"fee": 0,
			"unit": "USD",
			"unit_type": "CURRENCY"
		},
		"wholesale": {
			"amount": 3.6,
			"fee": 0,
			"unit": "USD",
			"unit_type": "CURRENCY"
		}
	},
	"product": {
		"description": "",
		"id": 6035,
		"name": "3 USD",
		"operator": {
			"country": {
				"iso_code": "ECU",
				"name": "Ecuador",
				"regions": null
			},
			"id": 1596,
			"name": "Claro Ecuador",
			"regions": null
		},
		"regions": null,
		"service": {
			"id": 1,
			"name": "Mobile"
		},
		"type": "FIXED_VALUE_RECHARGE"
	},
	"promotions": null,
	"rates": {
		"base": 0.833333333333333,
		"retail": 0.75,
		"wholesale": 0.833333333333333
	},
	"status": {
		"class": {
			"id": 3,
			"message": "REJECTED"
		},
		"id": 30400,
		"message": "REJECTED-OPERATOR-CURRENTLY-UNAVAILABLE"
	}
}`

func TestClient(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)
	defer dates.SetNowSource(dates.DefaultNowSource)

	mocks := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://dvs-api.dtone.com/v1/lookup/mobile-number/+593979123456": {
			httpx.NewMockResponse(200, nil, lookupNumberResponse), // successful mobile number lookup
			httpx.MockConnectionError,                             // timeout
		},
		"https://dvs-api.dtone.com/v1/products?type=FIXED_VALUE_RECHARGE&operator_id=1596&per_page=100": {
			httpx.NewMockResponse(200, nil, productsResponse),
		},
		"https://dvs-api.dtone.com/v1/sync/transactions": {
			httpx.NewMockResponse(200, nil, transactionRejectedResponse),
		},
	})

	httpx.SetRequestor(mocks)
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 9, 15, 25, 30, 123456789, time.UTC)))

	cl := dtone.NewClient(http.DefaultClient, nil, "key123", "sesame")

	// test lookup mobile number
	operators, trace, err := cl.LookupMobileNumber("+593979123456")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(operators))
	assert.Equal(t, 1596, operators[0].ID)
	assert.Equal(t, "Claro Ecuador", operators[0].Name)
	test.AssertSnapshot(t, "lookup_mobile_number", string(trace.RequestTrace))

	// test with error
	operators, _, err = cl.LookupMobileNumber("+593979123456")
	assert.EqualError(t, err, "unable to connect to server")
	assert.Nil(t, operators)

	// fetch products for that operator
	products, trace, err := cl.Products("FIXED_VALUE_RECHARGE", 1596)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, 6035, products[0].ID)
	assert.Equal(t, "3 USD", products[0].Name)
	assert.Equal(t, 6036, products[1].ID)
	assert.Equal(t, "6 USD", products[1].Name)
	test.AssertSnapshot(t, "products", string(trace.RequestTrace))

	// create a synchronous transaction
	tx, trace, err := cl.TransactionSync("EX12345", 6035, "+593979123456")
	assert.NoError(t, err)
	assert.Equal(t, int64(2237512891), tx.ID)
	assert.Equal(t, "EX12345", tx.ExternalID)
	assert.Equal(t, "REJECTED-OPERATOR-CURRENTLY-UNAVAILABLE", tx.Status.Message)
	test.AssertSnapshot(t, "transaction_sync", string(trace.RequestTrace))
}
