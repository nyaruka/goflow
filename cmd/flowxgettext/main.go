package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/buger/jsonparser"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/i18n"

	"github.com/pkg/errors"
)

const usage = `usage: flowxgettext [flags] <flowfile>...`

func main() {
	var excludeArgs bool
	var lang string
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&lang, "lang", "", "translation language to extract")
	flags.BoolVar(&excludeArgs, "exclude-args", false, "whether to exclude localized router arguments")
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) == 0 {
		fmt.Println(usage)
		flags.PrintDefaults()
		os.Exit(1)
	}
	if err := FlowXGetText(envs.Language(lang), excludeArgs, args, os.Stdout); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func FlowXGetText(lang envs.Language, excludeArgs bool, paths []string, writer io.Writer) error {
	sources, err := loadFlows(paths)
	if err != nil {
		return err
	}

	po, err := i18n.ExtractFromFlows(lang, excludeArgs, sources...)

	po.Write(writer)

	return nil
}

// loads all the flows in the given file paths which may be asset files or single flow definitions
func loadFlows(paths []string) ([]flows.Flow, error) {
	flows := make([]flows.Flow, 0)
	for _, path := range paths {
		fileJSON, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading flow file '%s'", path)
		}

		var flowDefs []json.RawMessage

		flowsSection, _, _, err := jsonparser.Get(fileJSON, "flows")
		if err == nil {
			// file is a set of assets with a flow section
			jsonparser.ArrayEach(flowsSection, func(flowJSON []byte, dataType jsonparser.ValueType, offset int, err error) {
				flowDefs = append(flowDefs, flowJSON)
			})
		} else {
			// file is a single flow definition
			flowDefs = append(flowDefs, fileJSON)
		}

		for _, flowDef := range flowDefs {
			flow, err := definition.ReadFlow(flowDef, &migrations.Config{BaseMediaURL: "http://temba.io"})
			if err != nil {
				return nil, errors.Wrapf(err, "error reading flow '%s'", path)
			}
			flows = append(flows, flow)
		}
	}

	return flows, nil
}
