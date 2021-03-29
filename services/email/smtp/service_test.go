package smtp_test

import (
	"testing"

	"github.com/nyaruka/goflow/services/email/smtp"
	"github.com/nyaruka/goflow/utils/smtpx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	defer smtpx.SetSender(smtpx.DefaultSender)

	sender := smtpx.NewMockSender(nil, nil)
	smtpx.SetSender(sender)

	// try with invalid URL
	_, err := smtp.NewService("xyz", nil)
	assert.EqualError(t, err, "connection URL must use SMTP scheme")

	svc, err := smtp.NewService("smtp://leah:pass123@temba.io:255?from=updates@temba.io", nil)
	require.NoError(t, err)

	err = svc.Send(nil, []string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Updates", "Have a great week")

	assert.NoError(t, err)
	assert.Equal(t, []string{"HELO localhost\nMAIL FROM:<updates@temba.io>\nRCPT TO:<bob@nyaruka.com>\nRCPT TO:<jim@nyaruka.com>\nDATA\nHave a great week\n.\nQUIT\n"}, sender.Logs())

	// if body is blank, we'll use a placeholder
	err = svc.Send(nil, []string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Updates", " ")

	assert.NoError(t, err)
	assert.Equal(t, "HELO localhost\nMAIL FROM:<updates@temba.io>\nRCPT TO:<bob@nyaruka.com>\nRCPT TO:<jim@nyaruka.com>\nDATA\n(empty body)\n.\nQUIT\n", sender.Logs()[1])
}
