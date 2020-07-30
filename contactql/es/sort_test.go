package es_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/contactql/es"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/stretchr/testify/assert"
)

func TestElasticSort(t *testing.T) {
	resolver := newMockResolver()

	tcs := []struct {
		Label   string
		Sort    string
		Elastic string
		Error   error
	}{
		{"empty", "", `{"id":{"order":"desc"}}`, nil},
		{"descending created_on", "-created_on", `{"created_on":{"order":"desc"}}`, nil},
		{"ascending name", "name", `{"name.keyword":{"order":"asc"}}`, nil},
		{"descending language", "-language", `{"language":{"order":"desc"}}`, nil},
		{"descending numeric", "-AGE", `{"fields.number":{"nested":{"filter":{"term":{"fields.field":"6b6a43fa-a26d-4017-bede-328bcdd5c93b"}},"path":"fields"},"order":"desc"}}`, nil},
		{"ascending text", "color", `{"fields.text":{"nested":{"filter":{"term":{"fields.field":"ecc7b13b-c698-4f46-8a90-24a8fab6fe34"}},"path":"fields"},"order":"asc"}}`, nil},
		{"descending date", "-dob", `{"fields.datetime":{"nested":{"filter":{"term":{"fields.field":"cbd3fc0e-9b74-4207-a8c7-248082bb4572"}},"path":"fields"},"order":"desc"}}`, nil},
		{"descending state", "-state", `{"fields.state_keyword":{"nested":{"filter":{"term":{"fields.field":"67663ad1-3abc-42dd-a162-09df2dea66ec"}},"path":"fields"},"order":"desc"}}`, nil},
		{"ascending district", "district", `{"fields.district_keyword":{"nested":{"filter":{"term":{"fields.field":"54c72635-d747-4e45-883c-099d57dd998e"}},"path":"fields"},"order":"asc"}}`, nil},
		{"ascending ward", "ward", `{"fields.ward_keyword":{"nested":{"filter":{"term":{"fields.field":"fde8f740-c337-421b-8abb-83b954897c80"}},"path":"fields"},"order":"asc"}}`, nil},

		{"unknown field", "foo", "", fmt.Errorf("unable to find field with name: foo")},
	}

	for _, tc := range tcs {
		sort, err := es.ToElasticFieldSort(resolver, tc.Sort)

		if err != nil {
			assert.Equal(t, tc.Error.Error(), err.Error())
			continue
		}

		src, _ := sort.Source()
		encoded, _ := jsonx.Marshal(src)
		assert.Equal(t, tc.Elastic, string(encoded))
	}
}
