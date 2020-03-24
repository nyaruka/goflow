package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/buger/jsonparser"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/gettext"
	"github.com/nyaruka/goflow/utils/uuids"

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

// TextLocation describes the location of a piece of extracted text
type TextLocation struct {
	Flow     flows.Flow
	UUID     uuids.UUID
	Property string
	Index    int
}

type ExtractedText struct {
	Locations   []TextLocation
	Base        string
	Translation string
	Unique      bool
}

func (t *ExtractedText) OnlyFromArguments() bool {
	for _, loc := range t.Locations {
		if loc.Property != "arguments" {
			return false
		}
	}
	return true
}

func FlowXGetText(lang envs.Language, excludeArgs bool, paths []string, writer io.Writer) error {
	sources, err := loadFlows(paths)
	if err != nil {
		return err
	}

	extracted, err := extractText(lang, excludeArgs, sources)
	if err != nil {
		return err
	}

	merged := mergeExtracted(extracted)
	pot := createPOT(lang, merged)
	pot.Write(writer)

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

func extractText(lang envs.Language, excludeArgs bool, sources []flows.Flow) ([]*ExtractedText, error) {
	baseLanguage := envs.NilLanguage
	extracted := make([]*ExtractedText, 0)
	for _, flow := range sources {
		if baseLanguage == envs.NilLanguage {
			baseLanguage = flow.Language()
		} else if baseLanguage != flow.Language() {
			return nil, errors.New("flows use different base languages")
		}

		baseTranslation := flow.ExtractBaseTranslation()
		var targetTranslation flows.Translation
		if lang != envs.NilLanguage {
			targetTranslation = flow.Localization().GetTranslation(lang)
		}

		baseTranslation.Enumerate(func(uuid uuids.UUID, property string, texts []string) {
			// look up target translation if we have one
			targets := make([]string, len(texts))
			if targetTranslation != nil {
				translation := targetTranslation.GetTextArray(uuid, property)
				if translation != nil {
					for t := range targets {
						if t < len(translation) {
							targets[t] = translation[t]
						}
					}
				}
			}

			for t, text := range texts {
				if text != "" {
					extracted = append(extracted, &ExtractedText{
						Locations: []TextLocation{
							TextLocation{
								Flow:     flow,
								UUID:     uuid,
								Property: property,
								Index:    t,
							},
						},
						Base:        text,
						Translation: targets[t],
						Unique:      false,
					})
				}
			}
		})
	}

	return extracted, nil
}

func mergeExtracted(extracted []*ExtractedText) []*ExtractedText {
	// organize extracted texts by their base text
	byBase := make(map[string][]*ExtractedText)
	for _, e := range extracted {
		byBase[e.Base] = append(byBase[e.Base], e)
	}

	// get the list of unique base text values and sort A-Z
	bases := make([]string, 0, len(byBase))
	for b := range byBase {
		bases = append(bases, b)
	}
	sort.Strings(bases)

	merged := make([]*ExtractedText, 0)

	for _, base := range bases {
		extractionsForBase := byBase[base]

		differingTranslations := false
		singleTranslation := extractionsForBase[0].Translation
		for _, ext := range extractionsForBase[1:] {
			if singleTranslation == "" {
				singleTranslation = ext.Translation
			}
			if ext.Translation != "" && ext.Translation != singleTranslation {
				differingTranslations = true
			}
		}

		if differingTranslations {
			// we have differing translations, keep extractions for each location separate
			for _, e := range extractionsForBase {
				merged = append(merged, e)
			}
		} else {
			// all translations were the same, create merged extraction for all locations
			locations := make([]TextLocation, len(extractionsForBase))
			for i, ext := range extractionsForBase {
				locations[i] = ext.Locations[0]
			}

			merged = append(merged, &ExtractedText{
				Locations:   locations,
				Base:        base,
				Translation: singleTranslation,
				Unique:      true,
			})
		}
	}

	return merged
}

func createPOT(lang envs.Language, extracted []*ExtractedText) *gettext.PO {
	pot := gettext.NewPO("Generated by flowxgettext", dates.Now(), lang.ToISO639_2(envs.NilCountry))

	for _, ext := range extracted {
		context := ""
		if !ext.Unique {
			context = fmt.Sprintf("%s/%s:%d", string(ext.Locations[0].UUID), ext.Locations[0].Property, ext.Locations[0].Index)
		}

		comment := ""
		if ext.OnlyFromArguments() {
			comment = "only test arguments"
		}

		references := make([]string, len(ext.Locations))
		for i, loc := range ext.Locations {
			references[i] = fmt.Sprintf("%s/%s/%s:%d", loc.Flow.UUID(), string(loc.UUID), loc.Property, loc.Index)
		}

		entry := &gettext.Entry{
			Comment: gettext.Comment{
				Extracted:  comment,
				References: references,
			},
			MsgContext: context,
			MsgID:      ext.Base,
			MsgStr:     ext.Translation,
		}

		pot.AddEntry(entry)
	}

	return pot
}
