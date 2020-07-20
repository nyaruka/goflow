package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UpdateSnapshots indicates whether tests should update snapshots
var UpdateSnapshots bool

func init() {
	flag.BoolVar(&UpdateSnapshots, "update", false, "whether to update test snapshots")
}

// AssertSnapshot checks that the file contains the expected text.
// However it creates the file if -update was set or file doesn't exist.
func AssertSnapshot(t *testing.T, name, expected string) {
	path := fmt.Sprintf("testdata/%s_%s.snap", t.Name(), name)
	_, err := os.Stat(path)

	if UpdateSnapshots || os.IsNotExist(err) {
		err := ioutil.WriteFile(path, []byte(expected), 0666)
		require.NoError(t, err, "error writing snapshot file %s", path)
	} else {
		data, err := ioutil.ReadFile(path)
		require.NoError(t, err, "error reading snapshot file %s", path)

		assert.Equal(t, string(data), expected)
	}
}
