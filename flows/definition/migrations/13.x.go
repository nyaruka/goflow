package migrations

import (
	"strings"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/refactor"
)

func init() {
	registerMigration(semver.MustParse("13.6.1"), Migrate13_6_1)
	registerMigration(semver.MustParse("13.6.0"), Migrate13_6_0)
	registerMigration(semver.MustParse("13.5.0"), Migrate13_5_0)
	registerMigration(semver.MustParse("13.4.0"), Migrate13_4_0)
	registerMigration(semver.MustParse("13.3.0"), Migrate13_3_0)
	registerMigration(semver.MustParse("13.2.0"), Migrate13_2_0)
	registerMigration(semver.MustParse("13.1.0"), Migrate13_1_0)
}

// Migrate13_6_1 fixes result lookups that need to be truncated.
//
// @version 13_6_1 "13.6.1"
func Migrate13_6_1(f Flow, cfg *Config) (Flow, error) {
	const maxResultRef = 64

	RewriteTemplates(f, GetTemplateCatalog(semver.MustParse("13.6.0")), func(s string) string {
		// refactor any @result.* or @(...) template to find result lookups that need to be truncated
		refactored, _ := refactor.Template(s, []string{"results"}, func(exp excellent.Expression) bool {
			changed := false

			exp.Visit(func(e excellent.Expression) {
				switch typed := e.(type) {
				case *excellent.DotLookup:
					if asRef, isRef := typed.Container.(*excellent.ContextReference); isRef {
						if asRef.Name == "results" {
							old := typed.Lookup
							typed.Lookup = stringsx.Truncate(old, maxResultRef)
							if typed.Lookup != old {
								changed = true
							}
						}
					}
				}
			})

			return changed
		})

		return refactored
	})
	return f, nil
}

// Migrate13_6_0 ensures that names of results and categories respect definition limits.
//
// @version 13_6_0 "13.6.0"
func Migrate13_6_0(f Flow, cfg *Config) (Flow, error) {
	const maxResultName = 64
	const maxCategoryName = 36

	truncate := func(s string, max int) string {
		return strings.TrimSpace(stringsx.Truncate(s, max)) // so we don't leave trailing spaces
	}

	for _, node := range f.Nodes() {
		for _, action := range node.Actions() {
			if action.Type() == "set_run_result" {
				name, _ := action["name"].(string)
				category, _ := action["category"].(string)

				if len(name) > maxResultName {
					action["name"] = truncate(name, maxResultName)
				}
				if len(category) > maxCategoryName {
					action["category"] = truncate(category, maxCategoryName)
				}
			}
		}

		router := node.Router()
		if router != nil {
			resultName, _ := router["result_name"].(string)
			categories, _ := router["categories"].([]any)

			if len(resultName) > maxResultName {
				router["result_name"] = truncate(resultName, maxResultName)
			}

			for _, cat := range categories {
				category, _ := cat.(map[string]any)
				if category != nil {
					name, _ := category["name"].(string)

					if len(name) > maxCategoryName {
						category["name"] = truncate(name, maxCategoryName)
					}
				}
			}
		}
	}
	return f, nil
}

// Migrate13_5_0 converts the `templating` object in [action:send_msg] actions to use a merged list of variables.
//
// @version 13_5_0 "13.5.0"
func Migrate13_5_0(f Flow, cfg *Config) (Flow, error) {
	localization := f.Localization()

	for _, node := range f.Nodes() {
		for _, action := range node.Actions() {
			if action.Type() == "send_msg" {
				templating, _ := action["templating"].(map[string]any)

				if templating != nil {
					variables := make([]string, 0, 5)
					localizedVariables := make(map[i18n.Language][]string)

					// the languages for which any component has param translations
					localizedLangs := make(map[i18n.Language]bool)

					components, _ := templating["components"].([]any)
					for i := range components {
						comp, _ := components[i].(map[string]any)
						compUUID := GetObjectUUID(comp)
						compParams, _ := comp["params"].([]any)
						compParamsAsStrings := make([]string, len(compParams))
						for j := range compParams {
							p, _ := compParams[j].(string)
							variables = append(variables, p)
							compParamsAsStrings[j] = p
						}

						if localization != nil {
							for _, lang := range localization.Languages() {
								langTrans := localization.GetLanguageTranslation(lang)
								if langTrans != nil {
									params := langTrans.GetTranslation(compUUID, "params")
									if params != nil {
										localizedVariables[lang] = append(localizedVariables[lang], params...)
										langTrans.DeleteTranslation(compUUID, "params")
										localizedLangs[lang] = true
									} else {
										// maybe this component's params aren't translated but others are
										localizedVariables[lang] = append(localizedVariables[lang], compParamsAsStrings...)
									}
								}
							}
						}
					}

					action["template"] = templating["template"]
					action["template_variables"] = variables
					delete(action, "templating")

					if localization != nil {
						for lang, langVariables := range localizedVariables {
							if localizedLangs[lang] {
								langTrans := localization.GetLanguageTranslation(lang)
								langTrans.SetTranslation(action.UUID(), "template_variables", langVariables)
							}
						}
					}
				}
			}
		}
	}
	return f, nil
}

// Migrate13_4_0 converts the `templating` object in [action:send_msg] actions to use a list of components.
//
// @version 13_4_0 "13.4.0"
func Migrate13_4_0(f Flow, cfg *Config) (Flow, error) {
	localization := f.Localization()

	for _, node := range f.Nodes() {
		for _, action := range node.Actions() {
			if action.Type() == "send_msg" {
				templating, _ := action["templating"].(map[string]any)
				if templating != nil {
					templatingUUID := GetObjectUUID(templating)
					bodyCompUUID := uuids.NewV4()
					variables, _ := templating["variables"].([]any)
					if variables == nil {
						variables = []any{}
					}
					templating["components"] = []map[string]any{
						{"uuid": bodyCompUUID, "name": "body", "params": variables},
					}

					if localization != nil {
						for _, lang := range localization.Languages() {
							langTrans := localization.GetLanguageTranslation(lang)
							if langTrans != nil {
								vars := langTrans.GetTranslation(templatingUUID, "variables")
								if vars != nil {
									langTrans.SetTranslation(bodyCompUUID, "params", vars)
									langTrans.DeleteTranslation(templatingUUID, "variables")
								}
							}
						}
					}

					delete(templating, "uuid")
					delete(templating, "variables")
				}
			}
		}
	}
	return f, nil
}

// Migrate13_3_0 refactors template expressions that reference @webhook to use @webhook.json
//
// @version 13_3_0 "13.3.0"
func Migrate13_3_0(f Flow, cfg *Config) (Flow, error) {
	RewriteTemplates(f, GetTemplateCatalog(semver.MustParse("13.2.0")), func(s string) string {
		// some optimizations here...
		//   1. we can parse templates as if @(...) and @webhook are only valid top-levels
		//   2. we can treat adding .json as a lookup to webhook as a simple renaming of webhook to webhook.json
		refactored, _ := refactor.Template(s, []string{"webhook"}, refactor.ContextRefRename("webhook", "webhook.json"))
		return refactored
	})
	return f, nil
}

// Migrate13_2_0 replaces `base` as a flow language with `und` which indicates text with undetermined language
// in the ISO-639-3 standard.
//
// @version 13_2_0 "13.2.0"
func Migrate13_2_0(f Flow, cfg *Config) (Flow, error) {
	language, _ := f["language"].(string)
	localization := f.Localization()

	// if we don't have a valid language, replace it
	if len(language) != 3 {
		f["language"] = "und"
		if localization != nil {
			delete(localization, "und")
		}
	}

	return f, nil
}

// Migrate13_1_0 adds a `uuid` property to templating objects in [action:send_msg] actions.
//
// @version 13_1_0 "13.1.0"
func Migrate13_1_0(f Flow, cfg *Config) (Flow, error) {
	for _, node := range f.Nodes() {
		for _, action := range node.Actions() {
			if action.Type() == "send_msg" {
				templating, _ := action["templating"].(map[string]any)
				if templating != nil {
					templating["uuid"] = uuids.NewV4()
				}
			}
		}
	}
	return f, nil
}
