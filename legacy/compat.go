package legacy

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

// IsLegacyDefinition peeks at the given flow definition to determine if it is in legacy format
func IsLegacyDefinition(data json.RawMessage) bool {
	// any flow with a root-level flow_type property is considered to be in legacy format
	_, _, _, err := jsonparser.Get(data, "flow_type")
	return err == nil
}

// MigrateLegacyDefinition migrates a legacy definition
func MigrateLegacyDefinition(data json.RawMessage, baseMediaURL string) (json.RawMessage, error) {
	legacyFlow, err := ReadLegacyFlow(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read legacy flow")
	}

	flow, err := legacyFlow.Migrate(true, baseMediaURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to migrate legacy flow")
	}

	marshaled, err := utils.JSONMarshal(flow)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal migrated flow")
	}

	return marshaled, nil
}
