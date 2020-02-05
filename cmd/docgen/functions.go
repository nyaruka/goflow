package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/nyaruka/goflow/utils/jsonx"
)

func init() {
	registerGenerator("function listing", generateFunctionListing)
}

type functionExample struct {
	Template string `json:"template"`
	Output   string `json:"output"`
}

type functionListing struct {
	Signature string             `json:"signature"`
	Summary   string             `json:"summary"`
	Detail    string             `json:"detail"`
	Examples  []*functionExample `json:"examples"`
}

func generateFunctionListing(baseDir string, outputDir string, items map[string][]*TaggedItem) error {
	funcItems := items["function"]
	listings := make([]*functionListing, len(funcItems))

	for i, funcItem := range funcItems {
		summary := funcItem.description[0]
		detail := strings.TrimSpace(strings.Join(funcItem.description[1:len(funcItem.description)-1], "\n"))

		examples := make([]*functionExample, len(funcItem.examples))
		for j := range funcItem.examples {
			parts := strings.Split(funcItem.examples[j], "→")
			examples[j] = &functionExample{Template: strings.TrimSpace(parts[0]), Output: strings.TrimSpace(parts[1])}
		}

		listings[i] = &functionListing{
			Signature: funcItem.tagValue + funcItem.tagExtra,
			Summary:   summary,
			Detail:    detail,
			Examples:  examples,
		}
	}

	data, err := jsonx.MarshalPretty(listings)
	if err != nil {
		return err
	}

	listingPath := path.Join(outputDir, "functions.json")

	if err := ioutil.WriteFile(listingPath, []byte(data), 0666); err != nil {
		return err
	}

	fmt.Printf(" > %d functions written to %s\n", len(listings), listingPath)

	// print table of function signatures and summaries
	//fmt.Printf("|Summary                                      |Signature                                                                                           |\n")
	//fmt.Printf("|---------------------------------------------|----------------------------------------------------------------------------------------------------|\n")
	//for _, fn := range listings {
	//	fmt.Printf("|%-45s|%-100s|\n", fn.Signature, fn.Summary)
	//}

	return nil
}
