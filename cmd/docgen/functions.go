package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/nyaruka/goflow/cmd/docgen/i18n"
	"github.com/nyaruka/goflow/utils"
)

// documentation extracted from the source code is in this language
const srcLanguage = "en_US"

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

	locales := i18n.NewLibrary(path.Join(baseDir, "locales"))
	pot := i18n.NewPOTemplate()

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

		pot.Add(summary)
		pot.Add(detail)
	}

	potPath := locales.POPath(srcLanguage, "functions")
	if err := pot.Save(potPath); err != nil {
		return err
	}

	fmt.Printf(" > localizable strings written to %s\n", potPath)

	for _, language := range locales.Languages() {
		locales.Activate(language, "functions")

		translated := make([]*functionListing, len(listings))

		for i, listing := range listings {
			translated[i] = &functionListing{
				Signature: listing.Signature,
				Summary:   i18n.GetText(listing.Summary),
				Detail:    i18n.GetText(listing.Detail),
				Examples:  listing.Examples,
			}
		}

		data, err := utils.JSONMarshalPretty(translated)
		if err != nil {
			return err
		}

		var listingPath string
		if language == srcLanguage {
			listingPath = path.Join(outputDir, "functions.json")
		} else {
			os.MkdirAll(path.Join(outputDir, language), 0755)
			listingPath = path.Join(outputDir, language, "functions.json")
		}

		if err := ioutil.WriteFile(listingPath, []byte(data), 0666); err != nil {
			return err
		}

		fmt.Printf(" > %d functions written to %s\n", len(listings), listingPath)
	}

	return nil
}
