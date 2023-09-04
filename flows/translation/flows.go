package translation

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/po"
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

func getBaseLanguage(set []flows.Flow) i18n.Language {
	if len(set) == 0 {
		return i18n.NilLanguage
	}
	baseLanguage := set[0].Language()
	for _, flow := range set[1:] {
		if baseLanguage != flow.Language() {
			return i18n.NilLanguage
		}
	}
	return baseLanguage
}

// ExtractFromFlows extracts a PO file from a set of flows
func ExtractFromFlows(initialComment string, translationsLanguage i18n.Language, excludeProperties []string, sources ...flows.Flow) (*po.PO, error) {
	// check all flows have same base language
	baseLanguage := getBaseLanguage(sources)
	if baseLanguage == i18n.NilLanguage {
		return nil, errors.New("can't extract from flows with differing base languages")
	} else if translationsLanguage == baseLanguage {
		translationsLanguage = i18n.NilLanguage // we'll create a POT in the base language (i.e. no translations)
	}

	extracted := findLocalizedText(translationsLanguage, excludeProperties, sources)

	merged := mergeExtracted(extracted)

	return poFromExtracted(sources, initialComment, translationsLanguage, merged), nil
}

func findLocalizedText(translationsLanguage i18n.Language, excludeProperties []string, sources []flows.Flow) []*localizedText {
	exclude := utils.Set(excludeProperties)
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

func extractFromProperty(translationsLanguage i18n.Language, flow flows.Flow, uuid uuids.UUID, property string, texts []string) []*localizedText {
	extracted := make([]*localizedText, 0)

	// look up target translation if we have a translation language
	targets := make([]string, len(texts))
	if translationsLanguage != i18n.NilLanguage {
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

func poFromExtracted(sources []flows.Flow, initialComment string, lang i18n.Language, extracted []*localizedText) *po.PO {
	flowUUIDs := make([]string, len(sources))
	for i, f := range sources {
		flowUUIDs[i] = string(f.UUID())
	}

	header := po.NewHeader(initialComment, dates.Now(), i18n.NewLocale(lang, i18n.NilCountry))
	header.Custom["Source-Flows"] = strings.Join(flowUUIDs, "; ")
	header.Custom["Language-3"] = string(lang)
	p := po.NewPO(header)

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

		entry := &po.Entry{
			Comment: po.Comment{
				References: references,
			},
			MsgContext: context,
			MsgID:      ext.Base,
			MsgStr:     ext.Translation,
		}

		p.AddEntry(entry)
	}

	return p
}

// ImportIntoFlows imports translations from the given PO into the given flows
func ImportIntoFlows(p *po.PO, translationsLanguage i18n.Language, excludeProperties []string, targets ...flows.Flow) error {
	baseLanguage := getBaseLanguage(targets)
	if baseLanguage == i18n.NilLanguage {
		return errors.New("can't import into flows with differing base languages")
	} else if translationsLanguage == baseLanguage {
		return errors.New("can't import as the flow base language")
	}

	updates := CalculateFlowUpdates(p, translationsLanguage, excludeProperties, targets...)

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
func CalculateFlowUpdates(p *po.PO, translationsLanguage i18n.Language, excludeProperties []string, targets ...flows.Flow) []*TranslationUpdate {
	localized := findLocalizedText(translationsLanguage, excludeProperties, targets)
	localizedByContext := make(map[string][]*localizedText)
	localizedByMsgID := make(map[string][]*localizedText)

	for _, lt := range localized {
		context := lt.Locations[0].MsgContext()
		localizedByContext[context] = append(localizedByContext[context], lt)
		localizedByMsgID[lt.Base] = append(localizedByMsgID[lt.Base], lt)
	}

	updates := make([]*TranslationUpdate, 0)
	addUpdate := func(lt *localizedText, e *po.Entry) {
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
	for _, entry := range p.Entries {
		if entry.MsgContext == "" {
			for _, lt := range localizedByMsgID[entry.MsgID] {
				addUpdate(lt, entry)
			}
		}
	}

	// create more specific context based updates
	for _, entry := range p.Entries {
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

func applyUpdates(updates []*TranslationUpdate, translationsLanguage i18n.Language) {
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
