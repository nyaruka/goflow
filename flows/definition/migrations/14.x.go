package migrations

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/Masterminds/semver"
)

func init() {
	registerMigration(semver.MustParse("14.3.0"), Migrate14_3_0)
	registerMigration(semver.MustParse("14.2.0"), Migrate14_2_0)
	registerMigration(semver.MustParse("14.1.0"), Migrate14_1_0)
	registerMigration(semver.MustParse("14.0.0"), Migrate14_0_0)
}

// Migrate14_3_0 changes airtime and ticket nodes to split on their output local instead of the action's result.
//
// @version 14_3_0 "14.3.0"
func Migrate14_3_0(f Flow, cfg *Config) (Flow, error) {
	actionTypes := map[string]string{"open_ticket": "_new_ticket", "transfer_airtime": "_new_transfer"}

	// replace any @results.* operands
	for _, node := range f.Nodes() {
		actions := node.Actions()
		router := node.Router()

		// ignore if this isn't a airtime or ticket split
		if len(actions) != 1 || actionTypes[actions[0].Type()] == "" || router == nil || router.Type() != "switch" {
			continue
		}

		action := actions[0]
		operand, _ := router["operand"].(string)
		cases, _ := router["cases"].([]any)

		// ignore if it already isn't splitting on a result
		if !strings.HasPrefix(operand, "@results.") || len(cases) == 0 {
			continue
		}

		case0, _ := cases[0].(map[string]any)
		case0["type"] = "has_text"
		case0["arguments"] = []any{}

		router["operand"] = "@locals." + actionTypes[action.Type()]
		router["cases"] = []any{case0}

		// move result to router
		router["result_name"] = action["result_name"]
		delete(action, "result_name")
	}

	return f, nil
}

// Migrate14_2_0 changes body to note on open ticket actions and cleans up invalid localization languages.
//
// @version 14_2_0 "14.2.0"
func Migrate14_2_0(f Flow, cfg *Config) (Flow, error) {
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

// Migrate14_1_0 changes webhook nodes to split on @webhook instead of the action's result.
//
// @version 14_1_0 "14.1.0"
func Migrate14_1_0(f Flow, cfg *Config) (Flow, error) {
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

// Migrate14_0_0 fixes invalid expires values and categories with missing names.
// Note that this is a major version change because of other additions to the flow spec that don't require migration.
//
// @version 14_0_0 "14.0.0"
func Migrate14_0_0(f Flow, cfg *Config) (Flow, error) {
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
