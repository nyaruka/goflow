package services

import (
	"context"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
)

// LLMService is an implementation of an LLM service for testing.
type LLMService struct {
	Responses map[string]string
}

func NewLLM(responses map[string]string) *LLMService {
	return &LLMService{Responses: responses}
}

func (s *LLMService) Response(ctx context.Context, env envs.Environment, instructions, input string) (string, error) {
	output := s.Responses[input]
	if output == "" {
		output = "I'm sorry I dont't understand"
	}

	return s.Responses[input], nil
}

var _ flows.LLMService = (*LLMService)(nil)
