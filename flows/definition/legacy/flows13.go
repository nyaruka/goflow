package legacy

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/legacy/expressions"
	"github.com/nyaruka/goflow/utils/uuids"
)

type migratedExit map[string]interface{}

func newExit(uuid uuids.UUID, destinationUUID uuids.UUID) migratedExit {
	d := map[string]interface{}{"uuid": uuid}
	if destinationUUID != "" {
		d["destination_uuid"] = destinationUUID
	}

	return migratedExit(d)
}

func (e migratedExit) UUID() uuids.UUID {
	return e["uuid"].(uuids.UUID)
}

type migratedNode map[string]interface{}

func newNode(uuid uuids.UUID, actions []flows.Action, router migratedRouter, exits []migratedExit) migratedNode {
	d := map[string]interface{}{
		"uuid":  uuid,
		"exits": exits,
	}
	if len(actions) > 0 {
		d["actions"] = actions
	}
	if router != nil {
		d["router"] = router
	}

	return migratedNode(d)
}

func (n migratedNode) UUID() uuids.UUID {
	return n["uuid"].(uuids.UUID)
}

type migratedLocalization map[envs.Language]map[uuids.UUID]map[string][]string

func (l migratedLocalization) addTranslation(lang envs.Language, itemUUID uuids.UUID, property string, translated []string) {
	_, found := l[lang]
	if !found {
		l[lang] = make(map[uuids.UUID]map[string][]string)
	}
	langTranslations := l[lang]

	_, found = langTranslations[itemUUID]
	if !found {
		langTranslations[itemUUID] = make(map[string][]string)
	}

	langTranslations[itemUUID][property] = translated
}

func (l migratedLocalization) addTranslationMap(baseLanguage envs.Language, mapped Translations, uuid uuids.UUID, property string) string {
	var inBaseLanguage string
	for language, item := range mapped {
		expression, _ := expressions.MigrateTemplate(item, nil)
		if language == baseLanguage {
			inBaseLanguage = expression
		} else if language != "base" {
			l.addTranslation(language, uuid, property, []string{expression})
		}
	}

	return inBaseLanguage
}

func (l migratedLocalization) addTranslationMultiMap(baseLanguage envs.Language, mapped map[envs.Language][]string, uuid uuids.UUID, property string) []string {
	var inBaseLanguage []string
	for language, items := range mapped {
		templates := make([]string, len(items))
		for i := range items {
			expression, _ := expressions.MigrateTemplate(items[i], nil)
			templates[i] = expression
		}
		if language != baseLanguage {
			l.addTranslation(language, uuid, property, templates)
		} else {
			inBaseLanguage = templates
		}
	}
	return inBaseLanguage
}

type migratedCategory map[string]interface{}

func newCategory(uuid uuids.UUID, name string, exitUUID uuids.UUID) migratedCategory {
	d := map[string]interface{}{"uuid": uuid}
	if name != "" {
		d["name"] = name
	}
	if exitUUID != "" {
		d["exit_uuid"] = exitUUID
	}

	return migratedCategory(d)
}

func (c migratedCategory) UUID() uuids.UUID {
	return c["uuid"].(uuids.UUID)
}

type migratedCase map[string]interface{}

func newCase(uuid uuids.UUID, type_ string, arguments []string, categoryUUID uuids.UUID) migratedCase {
	d := map[string]interface{}{
		"uuid":          uuid,
		"type":          type_,
		"category_uuid": categoryUUID,
	}
	if len(arguments) > 0 {
		d["arguments"] = arguments
	}

	return migratedCase(d)
}

func (c migratedCase) UUID() uuids.UUID {
	return c["uuid"].(uuids.UUID)
}

type migratedTimeout map[string]interface{}

func newTimeout(seconds int, categoryUUID uuids.UUID) migratedTimeout {
	return migratedTimeout(map[string]interface{}{"seconds": seconds, "category_uuid": categoryUUID})
}

type migratedHint map[string]interface{}

func newAudioHint() migratedHint {
	return migratedHint(map[string]interface{}{"type": "audio"})
}

func newImageHint() migratedHint {
	return migratedHint(map[string]interface{}{"type": "image"})
}

func newVideoHint() migratedHint {
	return migratedHint(map[string]interface{}{"type": "video"})
}

func newLocationHint() migratedHint {
	return migratedHint(map[string]interface{}{"type": "location"})
}

func newFixedDigitsHint(count int) migratedHint {
	return migratedHint(map[string]interface{}{"type": "digits", "count": count})
}

func newTerminatedDigitsHint(terminatedBy string) migratedHint {
	return migratedHint(map[string]interface{}{"type": "digits", "terminated_by": terminatedBy})
}

type migratedWait map[string]interface{}

func newMsgWait(timeout migratedTimeout, hint migratedHint) migratedWait {
	d := map[string]interface{}{"type": "msg"}
	if timeout != nil {
		d["timeout"] = timeout
	}
	if hint != nil {
		d["hint"] = hint
	}

	return migratedWait(d)
}

type migratedRouter map[string]interface{}

func newSwitchRouter(wait migratedWait, resultName string, categories []migratedCategory, operand string, cases []migratedCase, defaultCategory uuids.UUID) migratedRouter {
	d := map[string]interface{}{
		"type":                  "switch",
		"categories":            categories,
		"operand":               operand,
		"cases":                 cases,
		"default_category_uuid": defaultCategory,
	}
	if wait != nil {
		d["wait"] = wait
	}
	if resultName != "" {
		d["result_name"] = resultName
	}

	return migratedRouter(d)
}

func newRandomRouter(resultName string, categories []migratedCategory) migratedRouter {
	d := map[string]interface{}{
		"type":       "random",
		"categories": categories,
	}
	if resultName != "" {
		d["result_name"] = resultName
	}

	return migratedRouter(d)
}
