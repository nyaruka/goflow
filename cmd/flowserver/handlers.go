package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/legacy"

	"github.com/pkg/errors"
)

const (
	maxRequestBytes int64 = 1048576
)

// Migrates a legacy flow to the new flow definition specification
//
//   {
//     "flow": {"uuid": "468621a8-32e6-4cd2-afc1-04416f7151f0", "action_sets": [], ...},
//     "include_ui": false
//   }
//
type migrateRequest struct {
	Flow          json.RawMessage `json:"flow"`
	CollapseExits *bool           `json:"collapse_exits"`
	IncludeUI     *bool           `json:"include_ui"`
}

func (s *FlowServer) handleMigrate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	migrate := migrateRequest{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxRequestBytes))
	if err != nil {
		return nil, err
	}

	if err := r.Body.Close(); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &migrate); err != nil {
		return nil, err
	}

	if migrate.Flow == nil {
		return nil, errors.Errorf("missing flow element")
	}

	legacyFlow, err := legacy.ReadLegacyFlow(migrate.Flow)
	if err != nil {
		return nil, err
	}

	collapseExits := migrate.CollapseExits == nil || *migrate.CollapseExits
	includeUI := migrate.IncludeUI == nil || *migrate.IncludeUI

	return legacyFlow.Migrate(collapseExits, includeUI)
}

// Returns the current version number
func (s *FlowServer) handleVersion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	response := map[string]string{
		"version": version,
	}
	return response, nil
}
