package translation

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/i18n"
)

// describes the location of a piece of extracted text
type textLocation struct {
	Flow     flows.Flow
	UUID     uuids.UUID
	Property string
	Index    int
}

// String returns a full string representation of this location for use in PO reference comments
func (l *textLocation) String() string {
	flowName := url.QueryEscape(l.Flow.Name())

	return fmt.Sprintf("%s/%s/%s:%d", flowName, string(l.UUID), l.Property, l.Index)
}

func (l *textLocation) MsgContext() string {
	return fmt.Sprintf("%s/%s:%d", string(l.UUID), l.Property, l.Index)
}

type localizedText struct {
	Locations   []textLocation
	Base        string
	Translation string
	Unique      bool
}

func getBaseLanguage(set []flows.Flow) envs.Language {
	if len(set) == 0 {
		return envs.NilLanguage
	}
	baseLanguage := set[0].Language()
	for _, flow := range set[1:] {
		if baseLanguage != flow.Language() {
			return envs.NilLanguage
		}
	}
	return baseLanguage
}

// ExtractFromFlows extracts a PO file from a set of flows
func ExtractFromFlows(initialComment string, translationsLanguage envs.Language, excludeProperties []string, sources ...flows.Flow) (*i18n.PO, error) {
	// check all flows have same base language
	baseLanguage := getBaseLanguage(sources)
	if baseLanguage == envs.NilLanguage {
		return nil, errors.New("can't extract from flows with differing base languages")
	} else if translationsLanguage == baseLanguage {
		translationsLanguage = envs.NilLanguage // we'll create a POT in the base language (i.e. no translations)
	}

	extracted := findLocalizedText(translationsLanguage, excludeProperties, sources)

	merged := mergeExtracted(extracted)

	return poFromExtracted(sources, initialComment, translationsLanguage, merged), nil
}

func findLocalizedText(translationsLanguage envs.Language, excludeProperties []string, sources []flows.Flow) []*localizedText {
	exclude := utils.StringSet(excludeProperties)
	extracted := make([]*localizedText, 0)

	for _, flow := range sources {
		for _, node := range flow.Nodes() {
			node.EnumerateLocalizables(func(uuid uuids.UUID, property string, texts []string, write func([]string)) {
				if !exclude[property] {
					exts := extractFromProperty(translationsLanguage, flow, uuid, property, texts)
					extracted = append(extracted, exts...)
				}
			})
		}
	}

	return extracted
}

func extractFromProperty(translationsLanguage envs.Language, flow flows.Flow, uuid uuids.UUID, property string, texts []string) []*localizedText {
	extracted := make([]*localizedText, 0)

	// look up target translation if we have a translation language
	targets := make([]string, len(texts))
	if translationsLanguage != envs.NilLanguage {
		translation := flow.Localization().GetItemTranslation(translationsLanguage, uuid, property)
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
			extracted = append(extracted, &localizedText{
				Locations: []textLocation{
					{
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

func mergeExtracted(extracted []*localizedText) []*localizedText {
	// organize extracted texts by their base text
	byBase := make(map[string][]*localizedText)
	for _, e := range extracted {
		byBase[e.Base] = append(byBase[e.Base], e)
	}

	// get the list of unique base text values and sort A-Z
	bases := make([]string, 0, len(byBase))
	for b := range byBase {
		bases = append(bases, b)
	}
	sort.Strings(bases)

	merged := make([]*localizedText, 0)

	for _, base := range bases {
		extractionsForBase := byBase[base]

		majorityTranslation := majorityTranslation(extractionsForBase)

		// all extractions with majority translation or no translation get merged into a new context-less extraction
		mergedLocations := make([]textLocation, 0)

		for _, ext := range extractionsForBase {
			if ext.Translation == majorityTranslation || ext.Translation == "" {
				mergedLocations = append(mergedLocations, ext.Locations[0])
			} else {
				merged = append(merged, ext)
			}
		}

		merged = append(merged, &localizedText{
			Locations:   mergedLocations,
			Base:        base,
			Translation: majorityTranslation,
			Unique:      true,
		})
	}

	return merged
}

// finds the majority non-empty translation
func majorityTranslation(extracted []*localizedText) string {
	counts := make(map[string]int)
	for _, e := range extracted {
		if e.Translation != "" {
			counts[e.Translation]++
		}
	}
	max := 0
	majority := ""
	for _, e := range extracted {
		if counts[e.Translation] > max {
			majority = e.Translation
			max = counts[e.Translation]
		}
	}
	return majority
}

func poFromExtracted(sources []flows.Flow, initialComment string, lang envs.Language, extracted []*localizedText) *i18n.PO {
	flowUUIDs := make([]string, len(sources))
	for i, f := range sources {
		flowUUIDs[i] = string(f.UUID())
	}

	header := i18n.NewPOHeader(initialComment, dates.Now(), envs.NewLocale(lang, envs.NilCountry).ToBCP47())
	header.Custom["Source-Flows"] = strings.Join(flowUUIDs, "; ")
	header.Custom["Language-3"] = string(lang)
	po := i18n.NewPO(header)

	for _, ext := range extracted {
		references := make([]string, len(ext.Locations))
		for i, loc := range ext.Locations {
			references[i] = loc.String()
		}
		sort.Strings(references)

		context := ""
		if !ext.Unique {
			context = ext.Locations[0].MsgContext()
		}

		entry := &i18n.POEntry{
			Comment: i18n.POComment{
				References: references,
			},
			MsgContext: context,
			MsgID:      ext.Base,
			MsgStr:     ext.Translation,
		}

		po.AddEntry(entry)
	}

	return po
}

// ImportIntoFlows imports translations from the given PO into the given flows
func ImportIntoFlows(po *i18n.PO, translationsLanguage envs.Language, targets ...flows.Flow) error {
	baseLanguage := getBaseLanguage(targets)
	if baseLanguage == envs.NilLanguage {
		return errors.New("can't import into flows with differing base languages")
	} else if translationsLanguage == baseLanguage {
		return errors.New("can't import as the flow base language")
	}

	updates := CalculateFlowUpdates(po, translationsLanguage, targets...)

	applyUpdates(updates, translationsLanguage)

	return nil
}

// TranslationUpdate describs a change to be made to a flow translation
type TranslationUpdate struct {
	textLocation
	Base string
	Old  string
	New  string
}

func (u *TranslationUpdate) String() string {
	return fmt.Sprintf("%s %s -> %s", u.textLocation.String(), strconv.Quote(u.Old), strconv.Quote(u.New))
}

// CalculateFlowUpdates calculates what updates should be made to translations in the given flows
func CalculateFlowUpdates(po *i18n.PO, translationsLanguage envs.Language, targets ...flows.Flow) []*TranslationUpdate {
	localized := findLocalizedText(translationsLanguage, nil, targets)
	localizedByContext := make(map[string][]*localizedText)
	localizedByMsgID := make(map[string][]*localizedText)

	for _, lt := range localized {
		context := lt.Locations[0].MsgContext()
		localizedByContext[context] = append(localizedByContext[context], lt)
		localizedByMsgID[lt.Base] = append(localizedByMsgID[lt.Base], lt)
	}

	updates := make([]*TranslationUpdate, 0)
	addUpdate := func(lt *localizedText, e *i18n.POEntry) {
		// only update if translation has actually changed
		if lt.Translation != e.MsgStr {
			updates = append(updates, &TranslationUpdate{
				textLocation: lt.Locations[0],
				Base:         lt.Base,
				Old:          lt.Translation,
				New:          e.MsgStr,
			})
		}
	}

	// create all context-less updates
	for _, entry := range po.Entries {
		if entry.MsgContext == "" {
			for _, lt := range localizedByMsgID[entry.MsgID] {
				addUpdate(lt, entry)
			}
		}
	}

	// create more specific context based updates
	for _, entry := range po.Entries {
		if entry.MsgContext != "" {
			for _, lt := range localizedByContext[entry.MsgContext] {
				// only update if base text is still the same
				if lt.Base == entry.MsgID {
					addUpdate(lt, entry)
				}
			}
		}
	}

	// de-duplicate by location
	locationsSeen := make(map[string]int)
	deduped := make([]*TranslationUpdate, 0)
	for i, update := range updates {
		locationStr := update.textLocation.MsgContext()

		locIndex, existing := locationsSeen[locationStr]
		if existing {
			deduped[locIndex] = update
		} else {
			deduped = append(deduped, update)
			locationsSeen[locationStr] = i
		}
	}

	return deduped
}

func applyUpdates(updates []*TranslationUpdate, translationsLanguage envs.Language) {
	for _, update := range updates {
		localization := update.textLocation.Flow.Localization()
		texts := localization.GetItemTranslation(translationsLanguage, update.UUID, update.Property)

		// grow existing array if necessary
		for len(texts) < (update.Index + 1) {
			texts = append(texts, "")
		}

		texts[update.Index] = update.New

		localization.SetItemTranslation(translationsLanguage, update.UUID, update.Property, texts)
	}
}
