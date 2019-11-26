package smtpx_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils/smtpx"

	"github.com/stretchr/testify/assert"
)

func TestMockSender(t *testing.T) {
	defer smtpx.SetSender(smtpx.DefaultSender)

	sender := smtpx.NewMockSender()

	smtpx.SetSender(sender)

	err := smtpx.Send("mail.temba.io", 255, "leah", "pass123", "updates@temba.io", []string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Updates", "Have a great week")

	assert.NoError(t, err)
	assert.Equal(t, []string{"HELO localhost\nMAIL FROM:<updates@temba.io>\nRCPT TO:<bob@nyaruka.com>\nRCPT TO:<jim@nyaruka.com>\nDATA\nHave a great week\n.\nQUIT\n"}, sender.Logs())
}
