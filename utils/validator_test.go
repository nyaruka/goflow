package utils_test

import (
	"strings"
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

type SubObject struct {
	UUID      string `json:"uuid" validate:"uuid"`
	UUID4     string `json:"uuid4" validate:"uuid4"`
	SomeValue int    `json:"some_value" validate:"eq=2|eq=3"`
}

type TestObject struct {
	Foo    string    `json:"foo" validate:"required"`
	Bar    SubObject `json:"bar" validate:"required"`
	Things []string  `json:"things" validate:"min=1,max=3,dive,http_method"`
}

func TestValidate(t *testing.T) {
	// test with valid object
	errs := utils.ValidateAs(&TestObject{
		Foo: "hello",
		Bar: SubObject{
			UUID:      "ffffffff-ffff-ffff-bf1a-4186adc14195",
			UUID4:     "f0a26027-9ae9-422a-bf1a-4186adc14195",
			SomeValue: 2,
		},
		Things: []string{"GET", "POST", "PATCH"},
	}, "")
	assert.Nil(t, errs)

	// test with invalid object
	errs = utils.ValidateAs(&TestObject{
		Foo: "",
		Bar: SubObject{
			UUID:      "12345abcdefe",
			UUID4:     "ffffffff-ffff-ffff-bf1a-4186adc14195",
			SomeValue: 0,
		},
		Things: nil,
	}, "")
	assert.NotNil(t, errs)

	// check the individual error messages
	msgs := strings.Split(errs.Error(), ", ")
	assert.Equal(t, []string{
		`field 'foo' is required`,
		`field 'bar.uuid' must be a valid UUID`,
		`field 'bar.uuid4' must be a valid UUID4`,
		`field 'bar.some_value' failed tag 'eq=2|eq=3'`,
		`field 'things' must have a minimum of 1 items`,
	}, msgs)

	// test with another invalid object and an explicit object path
	errs = utils.ValidateAs(&TestObject{
		Foo: "hello",
		Bar: SubObject{
			UUID:      "ffffffff-ffff-ffff-bf1a-4186adc14195",
			UUID4:     "f0a26027-9ae9-422a-bf1a-4186adc14195",
			SomeValue: 2,
		},
		Things: []string{"UGHHH"},
	}, "blob.thing.test_object")
	assert.NotNil(t, errs)

	// check the individual error messages
	msgs = strings.Split(errs.Error(), ", ")
	assert.Equal(t, []string{
		`field 'blob.thing.test_object.things[0]' is not a valid HTTP method`,
	}, msgs)
}
