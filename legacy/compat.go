package legacy

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// IsLegacyDefinition peeks at the given flow definition to determine if it is in legacy format
func IsLegacyDefinition(data json.RawMessage) bool {
	header := &flowHeader{}
	if err := utils.UnmarshalAndValidate(data, header); err != nil {
		return false
	}

	// any flow definition with a metadata section is handled as a legacy definition
	return header.Metadata != nil
}

// MigrateLegacyDefinition migrates a legacy definition
func MigrateLegacyDefinition(data json.RawMessage) (json.RawMessage, error) {
	legacyFlow, err := ReadLegacyFlow(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read legacy flow")
	}

	flow, err := legacyFlow.Migrate(true, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to migrate legacy flow")
	}

	marshaled, err := utils.JSONMarshal(flow)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal migrated flow")
	}

	return marshaled, nil
}
