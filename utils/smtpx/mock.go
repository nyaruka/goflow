package smtpx

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// MockSender is a mocked sender for testing that just logs would-be commands
type MockSender struct {
	errs []error
	logs []string
}

// NewMockSender creates a new mock sender
func NewMockSender(errs ...error) *MockSender {
	return &MockSender{errs: errs}
}

// Logs returns the send logs
func (s *MockSender) Logs() []string {
	return s.logs
}

func (s *MockSender) Send(c *Client, m *Message) error {
	if len(s.errs) == 0 {
		panic(errors.Errorf("missing mock for send number %d", len(s.logs)))
	}

	err := s.errs[0]
	s.errs = s.errs[1:]

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
	return err
}
