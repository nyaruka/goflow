package main

// go install github.com/nyaruka/goflow/cmd/flowmigrate
// cat legacy_flow.json | flowmigrate
// cat legacy_export.json | jq '.flows[0]' | flowmigrate

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"
)

func main() {
	var includeUI, collapseExits bool
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.BoolVar(&includeUI, "collapse-exits", false, "Collapse ruleset exits with same category")
	flags.BoolVar(&includeUI, "include-ui", false, "Include UI configuration")
	flags.Parse(os.Args[1:])

	reader := bufio.NewReader(os.Stdin)

	output, err := Migrate(reader, collapseExits, includeUI)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
	}
}

// Migrate reads a legacy flow definition as JSON and migrates it
func Migrate(reader io.Reader, collapseExits, includeUI bool) ([]byte, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	flow, err := legacy.ReadLegacyFlow(data)
	if err != nil {
		return nil, err
	}

	migrated, err := flow.Migrate(collapseExits, includeUI)
	if err != nil {
		return nil, err
	}

	return utils.JSONMarshalPretty(migrated)
}
