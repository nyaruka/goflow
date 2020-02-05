package utils_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	_ "github.com/nyaruka/goflow/envs"
	_ "github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

type BaseObject struct {
	Foo string `json:"foo" validate:"required"`
}

type SubObject struct {
	UUID      string `json:"uuid" validate:"uuid"`
	UUID4     string `json:"uuid4" validate:"uuid4"`
	SomeValue int    `json:"some_value" validate:"eq=2|eq=3"`
}

type TestObject struct {
	BaseObject
	Bar        SubObject `json:"bar" validate:"required"`
	Things     []string  `json:"things" validate:"min=1,max=3,dive,http_method"`
	DateFormat string    `json:"date_format" validate:"date_format"`
	TimeFormat string    `json:"time_format" validate:"time_format"`
	Topic      string    `json:"topic" validate:"msg_topic"`
}

func TestValidate(t *testing.T) {
	// test with valid object
	errs := utils.Validate(&TestObject{
		BaseObject: BaseObject{Foo: "hello"},
		Bar: SubObject{
			UUID:      "ffffffff-ffff-ffff-bf1a-4186adc14195",
			UUID4:     "f0a26027-9ae9-422a-bf1a-4186adc14195",
			SomeValue: 2,
		},
		Things:     []string{"GET", "POST", "PATCH"},
		DateFormat: "DD-MM-YYYY",
		TimeFormat: "hh:mm:ss",
		Topic:      "account",
	})
	assert.Nil(t, errs)

	// test with invalid object
	errs = utils.Validate(&TestObject{
		BaseObject: BaseObject{Foo: ""},
		Bar: SubObject{
			UUID:      "12345abcdefe",
			UUID4:     "ffffffff-ffff-ffff-bf1a-4186adc14195",
			SomeValue: 0,
		},
		Things:     nil,
		DateFormat: "hh:mm",
		TimeFormat: "DD-MM",
		Topic:      "beer",
	})
	assert.NotNil(t, errs)

	// check the individual error messages
	msgs := strings.Split(errs.Error(), ", ")
	assert.Equal(t, []string{
		`field 'foo' is required`,
		`field 'bar.uuid' must be a valid UUID`,
		`field 'bar.uuid4' must be a valid UUID4`,
		`field 'bar.some_value' failed tag 'eq=2|eq=3'`,
		`field 'things' must have a minimum of 1 items`,
		`field 'date_format' is not a valid date format`,
		`field 'time_format' is not a valid time format`,
		`field 'topic' is not a valid message topic`,
	}, msgs)

	// test with another invalid object
	errs = utils.Validate(&TestObject{
		BaseObject: BaseObject{Foo: "hello"},
		Bar: SubObject{
			UUID:      "ffffffff-ffff-ffff-bf1a-4186adc14195",
			UUID4:     "f0a26027-9ae9-422a-bf1a-4186adc14195",
			SomeValue: 2,
		},
		Things: []string{"UGHHH"},
		Topic:  "football",
	})
	assert.NotNil(t, errs)

	// check the individual error messages
	msgs = strings.Split(errs.Error(), ", ")
	assert.Equal(t, []string{
		`field 'things[0]' is not a valid HTTP method`,
		`field 'topic' is not a valid message topic`,
	}, msgs)
}

func TestUnmarshalAndValidate(t *testing.T) {
	o := &BaseObject{}
	err := utils.UnmarshalAndValidate([]byte(`{}`), o)

	assert.EqualError(t, err, "field 'foo' is required")

	err = utils.UnmarshalAndValidate([]byte(`{"foo": "123"}`), o)

	assert.NoError(t, err)
	assert.Equal(t, "123", o.Foo)

	err = utils.UnmarshalAndValidateWithLimit(ioutil.NopCloser(bytes.NewReader([]byte(`{"foo": "abc"}`))), o, 100)

	assert.NoError(t, err)
	assert.Equal(t, "abc", o.Foo)

	err = utils.UnmarshalAndValidateWithLimit(ioutil.NopCloser(bytes.NewReader([]byte(`{"foo": "abc"}`))), o, 5)

	assert.EqualError(t, err, "unexpected end of JSON input")
}
