package services

import (
	"context"
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
	if strings.HasPrefix(input, "\\return ") { // an input like "\return foo" will return "foo"
		output = input[8:]
	} else {
		output = "You asked:\n\n" + instructions + "\n\n" + input
	}

	return &flows.LLMResponse{Output: output, TokensUsed: 123}, nil
}

var _ flows.LLMService = (*LLMService)(nil)
