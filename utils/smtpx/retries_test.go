package smtpx_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils/smtpx"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSendWithRetries(t *testing.T) {
	msg := smtpx.NewMessage([]string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Updates", "Have a great weekend", "")
	c := smtpx.NewClient("mail.temba.io", 255, "leah", "pass123", "updates@temba.io")

	// a sender which errors
	sender := smtpx.NewMockSender(
		errors.New("535 5.7.8 Username and Password not accepted"),            // a non-retriable 5xx error
		errors.New("oops can't send"),                                         // a non-retriable error with no code
		errors.New("421 Service not available, closing transmission channel"), // a retriable error
		errors.New("432 4.7.12 A password transition is needed"),              // a retriable error
		nil, // success
		errors.New("450 Requested mail action not taken: mailbox unavailable"),    // a retriable error
		errors.New("451 Requested action aborted: local error in processing"),     // a retriable error
		errors.New("452 Requested action not taken: insufficient system storage"), // too many retriable errors
		nil, // success
	)
	smtpx.SetSender(sender)

	retries := smtpx.NewFixedRetries(time.Millisecond*100, time.Millisecond*100)

	err := smtpx.Send(c, msg, retries)
	assert.EqualError(t, err, "535 5.7.8 Username and Password not accepted")
	assert.Equal(t, 1, len(sender.Logs()))

	err = smtpx.Send(c, msg, retries)
	assert.EqualError(t, err, "oops can't send")
	assert.Equal(t, 2, len(sender.Logs()))

	err = smtpx.Send(c, msg, retries)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(sender.Logs()))

	err = smtpx.Send(c, msg, retries)
	assert.EqualError(t, err, "452 Requested action not taken: insufficient system storage")
	assert.Equal(t, 8, len(sender.Logs()))

	err = smtpx.Send(c, msg, retries)
	assert.NoError(t, err)
	assert.Equal(t, 9, len(sender.Logs()))
}
