package mobile

import (
	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// ReadLegacyOrNewFlow reads either a legacy or new flow
func ReadLegacyOrNewFlow(definition string) (string, error) {
	flow, err := legacy.ReadLegacyOrNewFlow([]byte(definition))
	if err != nil {
		return "", errors.Wrap(err, "unable to read flow")
	}

	marshaled, err := utils.JSONMarshal(flow)
	if err != nil {
		return "", errors.Wrap(err, "unable to marshal flow")
	}

	return string(marshaled), nil
}
