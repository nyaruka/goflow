package smtpx

import (
	"errors"
	"fmt"
	"strings"
)

// MockSender is a mocked sender for testing that just logs would-be commands
type MockSender struct {
	err  string
	logs []string
}

// NewMockSender creates a new mock sender
func NewMockSender(err string) *MockSender {
	return &MockSender{err: err}
}

// Logs returns the send logs
func (s *MockSender) Logs() []string {
	return s.logs
}

func (s *MockSender) Send(c *Client, m *Message) error {
	if s.err != "" {
		return errors.New(s.err)
	}

	b := &strings.Builder{}
	b.WriteString("HELO localhost\n")
	b.WriteString(fmt.Sprintf("MAIL FROM:<%s>\n", c.from))
	for _, r := range m.recipients {
		b.WriteString(fmt.Sprintf("RCPT TO:<%s>\n", r))
	}
	b.WriteString("DATA\n")
	b.WriteString(fmt.Sprintf("%s\n", m.text))
	b.WriteString(".\n")
	b.WriteString("QUIT\n")

	s.logs = append(s.logs, b.String())
	return nil
}
