package migrations

import (
	"github.com/Masterminds/semver"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/excellent/tools"
)

func init() {
	registerMigration(semver.MustParse("13.4.0"), Migrate13_4)
	registerMigration(semver.MustParse("13.3.0"), Migrate13_3)
	registerMigration(semver.MustParse("13.2.0"), Migrate13_2)
	registerMigration(semver.MustParse("13.1.0"), Migrate13_1)
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
					bodyCompUUID := uuids.New()
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
		refactored, _ := tools.RefactorTemplate(s, []string{"webhook"}, tools.ContextRefRename("webhook", "webhook.json"))
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
					templating["uuid"] = uuids.New()
				}
			}
		}
	}
	return f, nil
}
