package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/utils"
)

type expressionResponse struct {
	Result string   `json:"result"`
	Errors []string `json:"errors"`
}

type expressionRequest struct {
	Expression string          `json:"expression"`
	Context    json.RawMessage `json:"context"`
}

func handleExpression(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	expression := expressionRequest{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &expression); err != nil {
		return nil, err
	}

	if expression.Context == nil || expression.Expression == "" {
		return nil, fmt.Errorf("missing context or expression element")
	}

	context := utils.NewJSONFragment(expression.Context)

	// evaluate it
	result, err := excellent.EvaluateTemplateAsString(utils.NewDefaultEnvironment(), context, expression.Expression)
	if err != nil {
		return nil, err
	}

	return expressionResponse{Result: result}, err
}
