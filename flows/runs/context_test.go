package runs_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sessionAssets = `{
    "channels": [
        {
            "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
            "name": "My Android Phone",
            "address": "+12345671111",
            "schemes": ["tel"],
            "roles": ["send", "receive"]
        }
    ],
    "flows": [
        {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "No Related Runs",
            "spec_version": "12.0",
            "language": "eng",
            "type": "messaging",
            "revision": 123,
            "nodes": [
                {
                    "uuid": "3dcccbb4-d29c-41dd-a01f-16d814c9ab82",
                    "wait": {
                        "type": "msg",
                        "timeout": 600
                    },
                    "router": {
                        "type": "switch",
                        "default_category_uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
                        "operand": "@input"
                    },
                    "exits": [
                        {
                            "uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
                            "name": "All Responses",
                            "destination_uuid": null
                        }
                    ]
                }
            ]
		}
	]
}`

var sessionTrigger = `{
    "type": "manual",
    "triggered_on": "2017-12-31T11:31:15.035757258-02:00",
    "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "No Related Runs"},
    "contact": {
        "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
        "id": 1234567,
        "name": "Ryan Lewis",
        "language": "eng",
        "timezone": "America/Guayaquil",
        "created_on": "2018-06-20T11:40:30.123456789-00:00",
        "urns": [ "tel:+12065551212"]
    },
    "environment": {
        "date_format": "YYYY-MM-DD",
        "default_language": "eng",
        "allowed_languages": [
            "eng", 
            "spa"
        ],
        "redaction_policy": "none",
        "time_format": "hh:mm",
        "timezone": "America/Guayaquil"
    }
}`

func TestRelatedRunContext(t *testing.T) {
	// create a run with no parent or child
	session, err := test.CreateSession([]byte(sessionAssets), "")
	require.NoError(t, err)

	trigger, err := triggers.ReadTrigger(session.Assets(), []byte(sessionTrigger), assets.IgnoreMissing)
	require.NoError(t, err)

	_, err = session.Start(trigger)
	require.NoError(t, err)

	run := session.Runs()[0]

	// check that trying to resolve parent is an error
	val, err := run.EvaluateTemplateValue(`@parent.contact`)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXErrorf("null has no property 'contact'"), val)

	// check that trying to resolve child is an error
	val, err = run.EvaluateTemplateValue(`@child.contact`)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXErrorf("null has no property 'contact'"), val)
}
