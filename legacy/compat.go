package legacy

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

// IsLegacyDefinition peeks at the given flow definition to determine if it is in legacy format
func IsLegacyDefinition(data json.RawMessage) bool {
	// any flow with root-level action_sets or rule_sets or flow_type property is considered to be in the new format
	frag1, _, _, _ := jsonparser.Get(data, "action_sets")
	frag2, _, _, _ := jsonparser.Get(data, "rule_sets")
	frag3, _, _, _ := jsonparser.Get(data, "flow_type")
	return frag1 != nil || frag2 != nil || frag3 != nil
}

// MigrateLegacyDefinition migrates a legacy definition
func MigrateLegacyDefinition(data json.RawMessage, baseMediaURL string) (json.RawMessage, error) {
	legacyFlow, err := ReadLegacyFlow(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read legacy flow")
	}

	flow, err := legacyFlow.Migrate(baseMediaURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to migrate legacy flow")
	}

	marshaled, err := utils.JSONMarshal(flow)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal migrated flow")
	}

	return marshaled, nil
}
