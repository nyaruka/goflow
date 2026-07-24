package services

import (
	"context"

	"github.com/nyaruka/goflow/flows"
)

// Email is an implementation of an email service for testing which just fakes sending the email
type Email struct{}

func NewEmail() *Email {
	return &Email{}
}

func (s *Email) Send(ctx context.Context, addresses []string, subject, body string) error {
	return nil
}

var _ flows.EmailService = (*Email)(nil)
