package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/pkg/errors"
)

const usage = `usage: flowinfo <inspect|localization> <path.json>`

func main() {
	if len(os.Args) != 3 || (os.Args[1] != "inspect" && os.Args[1] != "localization") {
		fmt.Println(usage)
		os.Exit(1)
	}

	if err := flowInfo(os.Args[1], os.Args[2]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func flowInfo(action, path string) error {
	definitionJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "error reading flow definition '%s'", path)
	}

	flow, err := definition.ReadFlow(definitionJSON, &migrations.Config{BaseMediaURL: "http://temba.io"})
	if err != nil {
		return err
	}

	if action == "inspect" {
		inspectionJSON, _ := jsonx.MarshalPretty(flow.Inspect(nil))
		fmt.Println(string(inspectionJSON))
	} else if action == "localization" {
		langs := flow.Localization().Languages()

		headers := []string{string(flow.Language())}
		for _, l := range langs {
			headers = append(headers, string(l))
		}

		writer := csv.NewWriter(os.Stdout)
		writer.Write(headers)

		base := flow.ExtractBaseTranslation()
		base.Enumerate(func(uuid uuids.UUID, property string, texts []string) {
			if len(texts) > 0 && texts[0] != "" {
				writeTranslationRows(writer, flow, uuid, property, texts, langs)
			}
		})

		writer.Flush()
	}

	return nil
}

func writeTranslationRows(writer *csv.Writer, flow flows.Flow, uuid uuids.UUID, property string, texts []string, langs []envs.Language) {
	rows := make([][]string, len(texts))

	for t, text := range texts {
		rows[t] = make([]string, len(langs)+1)
		rows[t][0] = text
	}

	for l, lang := range langs {
		translation := flow.Localization().GetTranslation(lang)
		translated := translation.GetTextArray(uuid, property)

		for t, text := range translated {
			if t < len(texts) {
				rows[t][l+1] = text
			}
		}
	}

	for _, row := range rows {
		writer.Write(row)
	}
}
