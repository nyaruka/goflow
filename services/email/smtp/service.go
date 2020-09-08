package smtp

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/smtpx"
)

type service struct {
	smtpClient *smtpx.Client
}

// NewService creates a new SMTP email service
func NewService(smtpURL string) (flows.EmailService, error) {
	c, err := smtpx.NewClientFromURL(smtpURL)
	if err != nil {
		return nil, err
	}

	return &service{smtpClient: c}, nil
}

func (s *service) Send(session flows.Session, addresses []string, subject, body string) error {
	m := smtpx.NewMessage(addresses, subject, body, "")
	return smtpx.Send(s.smtpClient, m)
}
