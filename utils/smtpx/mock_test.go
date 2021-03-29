package smtpx_test

import (
	"errors"
	"testing"

	"github.com/nyaruka/goflow/utils/smtpx"

	"github.com/stretchr/testify/assert"
)

func TestMockSender(t *testing.T) {
	defer smtpx.SetSender(smtpx.DefaultSender)

	// a sender which succeeds
	sender := smtpx.NewMockSender(nil, nil)
	smtpx.SetSender(sender)

	c := smtpx.NewClient("mail.temba.io", 255, "leah", "pass123", "updates@temba.io")

	msg1 := smtpx.NewMessage([]string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Updates", "Have a great week", "<p>Have a great week</p>")
	msg2 := smtpx.NewMessage([]string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Updates", "Have a great weekend", "")

	err := smtpx.Send(c, msg1, nil)
	assert.NoError(t, err)
	err = smtpx.Send(c, msg2, nil)
	assert.NoError(t, err)

	assert.Equal(t, []string{
		"HELO localhost\nMAIL FROM:<updates@temba.io>\nRCPT TO:<bob@nyaruka.com>\nRCPT TO:<jim@nyaruka.com>\nDATA\nHave a great week\n.\nQUIT\n",
		"HELO localhost\nMAIL FROM:<updates@temba.io>\nRCPT TO:<bob@nyaruka.com>\nRCPT TO:<jim@nyaruka.com>\nDATA\nHave a great weekend\n.\nQUIT\n",
	}, sender.Logs())

	// a sender which errors
	sender = smtpx.NewMockSender(errors.New("oops can't send"), errors.New("421 Service not available, closing transmission channel"))
	smtpx.SetSender(sender)

	err = smtpx.Send(c, msg1, nil)
	assert.EqualError(t, err, "oops can't send")
	assert.Equal(t, 1, len(sender.Logs()))

	err = smtpx.Send(c, msg2, nil)
	assert.EqualError(t, err, "421 Service not available, closing transmission channel")
	assert.Equal(t, 2, len(sender.Logs()))

	// we panic if we run out of mocks
	assert.PanicsWithError(t, "missing mock for send number 2", func() { smtpx.Send(c, msg2, nil) })
}
