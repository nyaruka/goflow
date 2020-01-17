package migrations

import (
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/Masterminds/semver"
)

func init() {
	registerMigration(semver.MustParse("13.1.0"), Migrate13_1)
	registerMigration(semver.MustParse("13.2.0"), Migrate13_2)
}

// Migrate13_1 adds UUID to send_msg templating
func Migrate13_1(f Flow) (Flow, error) {
	for _, node := range f.Nodes() {
		for _, action := range node.Actions() {
			if action.Type() == "send_msg" {
				templating, _ := action["templating"].(map[string]interface{})
				if templating != nil {
					templating["uuid"] = uuids.New()
				}
			}
		}
	}
	return f, nil
}

// Migrate13_2 adds response_as_extra to all webhook actions
func Migrate13_2(f Flow) (Flow, error) {
	for _, node := range f.Nodes() {
		for _, action := range node.Actions() {
			if action.Type() == "call_webhook" {
				action["response_as_extra"] = true
			}
		}
	}
	return f, nil
}
