package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/nyaruka/goflow/flows"
)

// LLMService is an implementation of an LLM service for testing that echos the input.
type LLMService struct{}

func NewLLM() *LLMService {
	return &LLMService{}
}

func (s *LLMService) Response(ctx context.Context, instructions, input string, maxTokens int) (*flows.LLMResponse, error) {
	// If input is a JSON array of strings whose first element is an \error or
	// \return directive, apply the directive as if the input were just that
	// element. Lets callers that marshal their input as a JSON array still
	// exercise the directives without needing a wrapper.
	effective := input
	if len(input) > 0 && input[0] == '[' {
		var arr []string
		if err := json.Unmarshal([]byte(input), &arr); err == nil && len(arr) > 0 {
			if strings.HasPrefix(arr[0], "\\error ") || strings.HasPrefix(arr[0], "\\return ") {
				effective = arr[0]
			}
		}
	}

	var output string
	if strings.HasPrefix(effective, "\\error ") { // an input like "\error foo" will return the error "foo"
		return nil, errors.New(effective[7:])
	} else if strings.HasPrefix(effective, "\\return ") { // an input like "\return foo" will return "foo"
		output = effective[8:]
	} else if strings.HasPrefix(instructions, "Categorize") { // instructions like "Categorize... Category2, Category3]" will return "Category3"
		words := strings.Fields(instructions)
		output = strings.TrimSuffix(words[len(words)-1], "]")
	} else {
		output = "You asked:\n\n" + instructions + "\n\n" + input
	}

	return &flows.LLMResponse{Output: output, TokensUsed: 123}, nil
}

var _ flows.LLMService = (*LLMService)(nil)
