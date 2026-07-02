package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nyaruka/goflow/events"
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

func translate(s string) (string, error) {
	if s == "error" {
		return "", errors.New("simulated LLM error")
	}
	if s == "untranslatable" {
		return "<CANT>", nil
	}
	return leetify(s), nil
}

func (s *LLMService) Response(ctx context.Context, instructions, input string, maxTokens int) (*events.LLMResponse, error) {
	var output string
	if strings.HasPrefix(input, "\\error ") { // an input like "\error foo" will return the error "foo"
		return nil, errors.New(input[7:])
	} else if strings.HasPrefix(input, "\\return ") { // an input like "\return foo" will return "foo"
		output = input[8:]
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
					tv, err := translate(v)
					if err != nil {
						return nil, err
					}
					vs[i] = tv
				}
				obj[k] = vs
			}
			b, _ := json.Marshal(obj)
			output = string(b)
		} else {
			tv, err := translate(input)
			if err != nil {
				return nil, err
			}
			output = tv
		}
	} else {
		output = "You asked:\n\n" + instructions + "\n\n" + input
	}

	return &events.LLMResponse{Output: output, TokensInput: 45, TokensOutput: 78}, nil
}

var _ flows.LLMService = (*LLMService)(nil)
