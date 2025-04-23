package migrations

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/refactor"
)

func init() {
	registerMigration(semver.MustParse("14.2.0"), Migrate14_2)
	registerMigration(semver.MustParse("14.1.0"), Migrate14_1)
	registerMigration(semver.MustParse("14.0.0"), Migrate14_0)
	registerMigration(semver.MustParse("13.6.1"), Migrate13_6_1)
	registerMigration(semver.MustParse("13.6.0"), Migrate13_6)
	registerMigration(semver.MustParse("13.5.0"), Migrate13_5)
	registerMigration(semver.MustParse("13.4.0"), Migrate13_4)
	registerMigration(semver.MustParse("13.3.0"), Migrate13_3)
	registerMigration(semver.MustParse("13.2.0"), Migrate13_2)
	registerMigration(semver.MustParse("13.1.0"), Migrate13_1)
}

// Migrate14_2 changes body to note on open ticket actions and cleans up invalid localization languages.
//
// @version 14_2 "14.2"
func Migrate14_2(f Flow, cfg *Config) (Flow, error) {
	for _, node := range f.Nodes() {
		for _, action := range node.Actions() {
			if action.Type() == "open_ticket" {
				body, _ := action["body"].(string)
				if body != "" {
					action["note"] = body
					delete(action, "body")
				}
			}
		}
	}

	if localization := f.Localization(); localization != nil {
		for _, lang := range localization.Languages() {
			if len(lang) != 3 {
				delete(localization, string(lang))
			}
		}
	}

	return f, nil
}

// Migrate14_1 changes webhook nodes to split on @webhook instead of the action's result.
//
// @version 14_1 "14.1"
func Migrate14_1(f Flow, cfg *Config) (Flow, error) {
	webhookActions := []string{"call_webhook", "call_resthook"}
	maxQuickReplies := 10

	// replace any @results.* operands in webhook nodes with @webhook.status
	for _, node := range f.Nodes() {
		actions := node.Actions()
		router := node.Router()

		// ignore if this isn't a webhook or resthook split
		if len(actions) != 1 || !slices.Contains(webhookActions, actions[0].Type()) || router == nil || router.Type() != "switch" {
			continue
		}

		operand, _ := router["operand"].(string)
		cases, _ := router["cases"].([]any)

		// ignore if it already isn't splitting on a result
		if !strings.HasPrefix(operand, "@results.") || len(cases) == 0 {
			continue
		}

		case0, _ := cases[0].(map[string]any)
		case0["type"] = "has_number_between"
		case0["arguments"] = []any{"200", "299"}

		router["operand"] = "@webhook.status"
		router["cases"] = []any{case0}
	}

	// trim any quick replies to a max of 10
	for _, node := range f.Nodes() {
		for _, action := range node.Actions() {
			if action.Type() == "send_msg" || action.Type() == "send_broadcast" {
				quickReplies, ok := action["quick_replies"].([]any)
				if ok && len(quickReplies) > maxQuickReplies {
					action["quick_replies"] = quickReplies[:maxQuickReplies]
				}
			}
		}
	}

	return f, nil
}

// Migrate14_0 fixes invalid expires values and categories with missing names.
// Note that this is a major version change because of other additions to the flow spec that don't require migration.
//
// @version 14_0 "14.0"
func Migrate14_0(f Flow, cfg *Config) (Flow, error) {
	maxExpires := map[string]int{
		"messaging": 20160, // two weeks
		"voice":     15,
	}

	expires, ok := f["expire_after_minutes"]
	if ok {
		expiresNum, ok := expires.(json.Number)
		if ok {
			expiresInt, err := expiresNum.Int64()
			if err == nil {
				f["expire_after_minutes"] = json.Number(fmt.Sprint(min(int(expiresInt), maxExpires[f.Type()])))
			}
		}
	}

	for _, node := range f.Nodes() {
		router := node.Router()
		if router != nil {
			categories, _ := router["categories"].([]any)
			for _, cat := range categories {
				category, _ := cat.(map[string]any)
				if category != nil {
					name, _ := category["name"].(string)
					if name == "" {
						category["name"] = "Match"
					}
				}
			}
		}
	}

	return f, nil
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

// Migrate13_6 ensures that names of results and categories respect definition limits.
//
// @version 13_6 "13.6"
func Migrate13_6(f Flow, cfg *Config) (Flow, error) {
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

// Migrate13_5 converts the `templating` object in [action:send_msg] actions to use a merged list of variables.
//
// @version 13_5 "13.5"
func Migrate13_5(f Flow, cfg *Config) (Flow, error) {
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

// Migrate13_4 converts the `templating` object in [action:send_msg] actions to use a list of components.
//
// @version 13_4 "13.4"
func Migrate13_4(f Flow, cfg *Config) (Flow, error) {
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

// Migrate13_3 refactors template expressions that reference @webhook to use @webhook.json
//
// @version 13_3 "13.3"
func Migrate13_3(f Flow, cfg *Config) (Flow, error) {
	RewriteTemplates(f, GetTemplateCatalog(semver.MustParse("13.2.0")), func(s string) string {
		// some optimizations here...
		//   1. we can parse templates as if @(...) and @webhook are only valid top-levels
		//   2. we can treat adding .json as a lookup to webhook as a simple renaming of webhook to webhook.json
		refactored, _ := refactor.Template(s, []string{"webhook"}, refactor.ContextRefRename("webhook", "webhook.json"))
		return refactored
	})
	return f, nil
}

// Migrate13_2 replaces `base` as a flow language with `und` which indicates text with undetermined language
// in the ISO-639-3 standard.
//
// @version 13_2 "13.2"
func Migrate13_2(f Flow, cfg *Config) (Flow, error) {
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

// Migrate13_1 adds a `uuid` property to templating objects in [action:send_msg] actions.
//
// @version 13_1 "13.1"
func Migrate13_1(f Flow, cfg *Config) (Flow, error) {
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
