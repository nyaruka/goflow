package services

import (
	"context"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
)

// LLMService is an implementation of an LLM service for testing that echos the input.
type LLMService struct{}

func NewLLM() *LLMService {
	return &LLMService{}
}

func (s *LLMService) Response(ctx context.Context, env envs.Environment, instructions, input string) (string, error) {
	return "You asked:\n\n" + instructions + "\n\n" + input, nil
}

var _ flows.LLMService = (*LLMService)(nil)
