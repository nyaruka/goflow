package flow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func readJSON(file string) ([]byte, error) {
	fmt.Printf("Parsing: %s\n", file)

	dir, exists := os.LookupEnv("GOPATH")
	fmt.Println(dir)
	if !exists {
		fmt.Printf("Missing")
	}

	flowJSON, err := ioutil.ReadFile(dir + "/test_flows/migrate/" + file)
	raw := json.RawMessage(flowJSON)

	flows, err := readLegacyFlows(raw)
	asJSON, err := json.MarshalIndent(flows, "", "  ")

	fmt.Printf("\njson:\n%s\n", asJSON)

	return asJSON, err
}

func TestSimpleMigration(t *testing.T) {
	readJSON("meningitis.json")
}
