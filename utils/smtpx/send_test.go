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
		smtpx.MockDialError("535 5.7.8 Username and Password not accepted"), // a non-retriable dial stage error
		errors.New("oops can't send"),                                       // a non-retriable send stage error
		smtpx.MockDialError("unable to connect to server"),                  // a retriable error
		smtpx.MockDialError("unable to connect to server"),                  // a retriable error
		nil, // success
		smtpx.MockDialError("unable to connect to server"), // a retriable error
		smtpx.MockDialError("unable to connect to server"), // a retriable error
		smtpx.MockDialError("unable to connect to server"), // too many retriable errors
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
	assert.EqualError(t, err, "unable to connect to server")
	assert.Equal(t, 8, len(sender.Logs()))
}
