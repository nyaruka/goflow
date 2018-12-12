package mobile

import (
	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// MigrateLegacyFlow migrates a legacy flow definitin
func MigrateLegacyFlow(definition string) (string, error) {
	legacyFlow, err := legacy.ReadLegacyFlow([]byte(definition))
	if err != nil {
		return "", errors.Wrap(err, "unable to read legacy flow")
	}

	flow, err := legacyFlow.Migrate(false, false)
	if err != nil {
		return "", errors.Wrap(err, "unable to migrate legacy flow")
	}

	marshaled, err := utils.JSONMarshal(flow)
	if err != nil {
		return "", errors.Wrap(err, "unable to marshal migrated flow")
	}

	return string(marshaled), nil
}
