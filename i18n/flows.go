package i18n

import (
	"errors"
	"fmt"
	"net/url"
	"sort"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/uuids"
)

// describes the location of a piece of extracted text
type textLocation struct {
	Flow     flows.Flow
	UUID     uuids.UUID
	Property string
	Index    int
}

type extractedText struct {
	Locations   []textLocation
	Base        string
	Translation string
	Unique      bool
}

// ExtractFromFlows extracts a PO file from a set of flows
func ExtractFromFlows(initialComment string, lang envs.Language, excludeArgs bool, sources ...flows.Flow) (*PO, error) {
	// check flows use same base language
	baseLanguage := envs.NilLanguage
	for _, flow := range sources {
		if baseLanguage == envs.NilLanguage {
			baseLanguage = flow.Language()
		} else if baseLanguage != flow.Language() {
			return nil, errors.New("flows use different base languages")
		}
	}

	extracted := extractFromFlows(lang, excludeArgs, sources)

	merged := mergeExtracted(extracted)

	return poFromExtracted(initialComment, lang, merged), nil
}

func extractFromFlows(lang envs.Language, excludeArgs bool, sources []flows.Flow) []*extractedText {
	extracted := make([]*extractedText, 0)

	for _, flow := range sources {
		var targetTranslation flows.Translation
		if lang != envs.NilLanguage {
			targetTranslation = flow.Localization().GetTranslation(lang)
		}

		for _, node := range flow.Nodes() {
			node.EnumerateLocalizables(func(uuid uuids.UUID, property string, texts []string) {
				if !excludeArgs || property != "arguments" {
					exts := extractFromProperty(flow, uuid, property, texts, targetTranslation)
					extracted = append(extracted, exts...)
				}
			})
		}
	}

	return extracted
}

func extractFromProperty(flow flows.Flow, uuid uuids.UUID, property string, texts []string, targetTranslation flows.Translation) []*extractedText {
	extracted := make([]*extractedText, 0)

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
			extracted = append(extracted, &extractedText{
				Locations: []textLocation{
					textLocation{
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

	return extracted
}

func mergeExtracted(extracted []*extractedText) []*extractedText {
	// organize extracted texts by their base text
	byBase := make(map[string][]*extractedText)
	for _, e := range extracted {
		byBase[e.Base] = append(byBase[e.Base], e)
	}

	// get the list of unique base text values and sort A-Z
	bases := make([]string, 0, len(byBase))
	for b := range byBase {
		bases = append(bases, b)
	}
	sort.Strings(bases)

	merged := make([]*extractedText, 0)

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
			locations := make([]textLocation, len(extractionsForBase))
			for i, ext := range extractionsForBase {
				locations[i] = ext.Locations[0]
			}

			merged = append(merged, &extractedText{
				Locations:   locations,
				Base:        base,
				Translation: singleTranslation,
				Unique:      true,
			})
		}
	}

	return merged
}

func poFromExtracted(initialComment string, lang envs.Language, extracted []*extractedText) *PO {
	pot := NewPO(initialComment, dates.Now(), lang.ToISO639_2(envs.NilCountry))

	for _, ext := range extracted {
		references := make([]string, len(ext.Locations))
		for i, loc := range ext.Locations {
			flowName := url.QueryEscape(loc.Flow.Name())
			references[i] = fmt.Sprintf("%s/%s/%s:%d", flowName, string(loc.UUID), loc.Property, loc.Index)
		}
		sort.Strings(references)

		context := ""
		if !ext.Unique {
			context = fmt.Sprintf("%s/%s:%d", string(ext.Locations[0].UUID), ext.Locations[0].Property, ext.Locations[0].Index)
		}

		entry := &Entry{
			Comment: Comment{
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
