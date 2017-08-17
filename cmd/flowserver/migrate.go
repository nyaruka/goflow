package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/flows"

	"github.com/nyaruka/goflow/flows/definition"
)

type migrateRequest struct {
	Flows        []json.RawMessage          `json:"flows"`
	FieldMapping map[string]flows.FieldUUID `json:"field_mapping"`
}

func handleMigrate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	migrate := migrateRequest{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}

	if err := r.Body.Close(); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &migrate); err != nil {
		return nil, err
	}

	if migrate.Flows == nil {
		return nil, fmt.Errorf("missing flows element")
	}

	flows, err := definition.ReadLegacyFlows(migrate.Flows, migrate.FieldMapping)
	if err != nil {
		return nil, err
	}

	return flows, err
}
