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

	flowJSON, err := ioutil.ReadFile("testdata/legacy/" + file)
	raw := json.RawMessage(flowJSON)

	flows, err := ReadLegacyFlows(raw)
	if err != nil {
		fmt.Println(err)
	}

	asJSON, err := json.MarshalIndent(flows, "", "  ")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\njson:\n%s\n", asJSON)

	return asJSON, err
}

func TestSimpleMigration(t *testing.T) {
	readJSON("lots_of_action.json")
}
