package static_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets/static"
	"github.com/stretchr/testify/assert"
)

var assetsJSON = `{
	"flows": [
		{
            "uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
            "name": "Empty",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
	],
	"fields": [
        {"uuid": "d66a7823-eada-40e5-9a3a-57239d4690bf", "key": "gender", "name": "Gender", "type": "text"},
        {"uuid": "f1b5aea6-6586-41c7-9020-1a6326cc6565", "key": "age", "name": "Age", "type": "number"}
    ],
	"groups": [
		{
			"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
			"name": "Survey Audience"
		}
	],
	"labels": [
		{
			"uuid": "18644b27-fb7f-40e1-b8f4-4ea8999129ef",
			"name": "Spam"
		}
	],
	"resthooks": [
		{
			"slug": "new-registration",
			"subscribers": [
				"http://temba.io/"
			]
		}
	]
}`

func TestSource(t *testing.T) {
	src := static.NewEmptySource()
	channels, err := src.Channels()
	assert.NoError(t, err)
	assert.Len(t, channels, 0)

	_, err = static.NewSource([]byte(`{`))
	assert.EqualError(t, err, "unable to read assets: unexpected end of JSON input")

	src, err = static.NewSource(json.RawMessage(assetsJSON))
	assert.NoError(t, err)

	channels, err = src.Channels()
	assert.NoError(t, err)
	assert.Len(t, channels, 0)

	classifiers, err := src.Classifiers()
	assert.NoError(t, err)
	assert.Len(t, classifiers, 0)

	fields, err := src.Fields()
	assert.NoError(t, err)
	assert.Len(t, fields, 2)

	flow, err := src.FlowByUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	assert.NoError(t, err)
	assert.Equal(t, "Empty", flow.Name())

	flow, err = src.FlowByName("Empty")
	assert.NoError(t, err)
	assert.Equal(t, "Empty", flow.Name())

	globals, err := src.Globals()
	assert.NoError(t, err)
	assert.Len(t, globals, 0)

	groups, err := src.Groups()
	assert.NoError(t, err)
	assert.Len(t, groups, 1)

	labels, err := src.Labels()
	assert.NoError(t, err)
	assert.Len(t, labels, 1)

	locations, err := src.Locations()
	assert.NoError(t, err)
	assert.Len(t, locations, 0)

	resthooks, err := src.Resthooks()
	assert.NoError(t, err)
	assert.Len(t, resthooks, 1)

	templates, err := src.Templates()
	assert.NoError(t, err)
	assert.Len(t, templates, 0)

	ticketers, err := src.Ticketers()
	assert.NoError(t, err)
	assert.Len(t, ticketers, 0)

	topics, err := src.Topics()
	assert.NoError(t, err)
	assert.Len(t, topics, 0)

	users, err := src.Users()
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}
