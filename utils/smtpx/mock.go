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

func NewMockSender(err string) *MockSender {
	return &MockSender{err: err}
}

func (s *MockSender) Logs() []string {
	return s.logs
}

func (s *MockSender) Send(host string, port int, username, password, from string, recipients []string, subject, body string) error {
	if s.err != "" {
		return errors.New(s.err)
	}

	b := &strings.Builder{}
	b.WriteString("HELO localhost\n")
	b.WriteString(fmt.Sprintf("MAIL FROM:<%s>\n", from))
	for _, r := range recipients {
		b.WriteString(fmt.Sprintf("RCPT TO:<%s>\n", r))
	}
	b.WriteString("DATA\n")
	b.WriteString(fmt.Sprintf("%s\n", body))
	b.WriteString(".\n")
	b.WriteString("QUIT\n")

	s.logs = append(s.logs, b.String())
	return nil
}
