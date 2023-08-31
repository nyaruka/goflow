package migrations

import (
	"github.com/Masterminds/semver"
	"github.com/nyaruka/gocommon/uuids"
)

func init() {
	registerMigration(semver.MustParse("13.2.0"), Migrate13_2)
	registerMigration(semver.MustParse("13.1.0"), Migrate13_1)
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
