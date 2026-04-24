package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
)

// LLMService is an implementation of an LLM service for testing that echos the input.
type LLMService struct{}

func NewLLM() *LLMService {
	return &LLMService{}
}

var leetify = strings.NewReplacer(
	"a", "4", "A", "4",
	"b", "8", "B", "8",
	"e", "3", "E", "3",
	"g", "9", "G", "9",
	"i", "1", "I", "1",
	"l", "1", "L", "1",
	"o", "0", "O", "0",
	"s", "5", "S", "5",
	"t", "7", "T", "7",
).Replace

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
	} else if strings.HasPrefix(instructions, "Translate") { // "Translate..." leetifies the input; if "JSON" is mentioned, values of a string->[]string object
		if strings.Contains(instructions, "JSON") {
			obj := map[string][]string{}
			if err := json.Unmarshal([]byte(input), &obj); err != nil {
				return nil, fmt.Errorf("invalid JSON object input: %w", err)
			}
			for k, vs := range obj {
				for i, v := range vs {
					vs[i] = leetify(v)
				}
				obj[k] = vs
			}
			b, _ := json.Marshal(obj)
			output = string(b)
		} else {
			output = leetify(input)
		}
	} else {
		output = "You asked:\n\n" + instructions + "\n\n" + input
	}

	return &flows.LLMResponse{Output: output, TokensUsed: 123}, nil
}

var _ flows.LLMService = (*LLMService)(nil)
