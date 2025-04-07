package services

import (
	"context"
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
	var output string
	if strings.HasPrefix(input, "\\error ") { // an input like "\error foo" will return the error "foo"
		return nil, errors.New(input[7:])
	} else if strings.HasPrefix(input, "\\return ") { // an input like "\return foo" will return "foo"
		output = input[8:]
	} else if strings.HasPrefix(instructions, "Categorize") { // instructions like "Categorize... Category2, Category3" will return "Category3"
		words := strings.Fields(instructions)
		output = words[len(words)-1]
	} else {
		output = "You asked:\n\n" + instructions + "\n\n" + input
	}

	return &flows.LLMResponse{Output: output, TokensUsed: 123}, nil
}

var _ flows.LLMService = (*LLMService)(nil)
