package definition

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testReadingInvalidFlow(t *testing.T, file string, expectedErr string) {
	var err error
	var assetsJSON json.RawMessage
	assetsJSON, err = ioutil.ReadFile(file)
	assert.NoError(t, err)

	_, err = ReadFlow(assetsJSON)
	assert.EqualError(t, err, expectedErr)
}

func TestReadFlow(t *testing.T) {
	testReadingInvalidFlow(t,
		"testdata/flow_with_duplicate_node_uuid.json",
		"duplicate node uuid: a58be63b-907d-4a1a-856b-0bb5579d7507",
	)
	testReadingInvalidFlow(t,
		"testdata/flow_with_invalid_default_exit.json",
		"router is invalid on node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: default exit 0680b01f-ba0b-48f4-a688-d2f963130126 is not a valid exit",
	)
	testReadingInvalidFlow(t,
		"testdata/flow_with_invalid_case_exit.json",
		"router is invalid on node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: case exit 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid exit",
	)
	testReadingInvalidFlow(t,
		"testdata/flow_with_invalid_case_exit.json",
		"router is invalid on node[uuid=a58be63b-907d-4a1a-856b-0bb5579d7507]: case exit 37d8813f-1402-4ad2-9cc2-e9054a96525b is not a valid exit",
	)
}
