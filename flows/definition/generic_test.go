package definition_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/definition"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type unmarshalable struct{}

func (u unmarshalable) MarshalJSON() ([]byte, error) { return nil, errors.New("doh!") }

func TestReadFlowFromGeneric(t *testing.T) {
	// try to read a generic that can't even be marshaled
	_, err := definition.ReadFlowFromGeneric(map[string]interface{}{
		"foo": unmarshalable{},
	})
	assert.Error(t, err)

	// try to read a generic that isn't a valid flow
	_, err = definition.ReadFlowFromGeneric(map[string]interface{}{
		"foo": "bar",
	})
	assert.Error(t, err)

	// try same thing with the Must version of this function
	assert.Panics(t, func() {
		definition.MustReadFlowFromGeneric(map[string]interface{}{
			"foo": "bar",
		})
	})

	// read a valid generic flow
	flow, err := definition.ReadFlowFromGeneric(map[string]interface{}{
		"uuid":         "786750b5-b3c7-4ccf-869b-69c9aeb8d891",
		"name":         "Empty Flow",
		"spec_version": "13.0",
		"language":     "eng",
		"type":         "messaging",
		"nodes":        []interface{}{},
	})
	assert.NoError(t, err)
	assert.Equal(t, "Empty Flow", flow.Name())
}
