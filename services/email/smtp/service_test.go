package smtp

import (
	"testing"

	"github.com/nyaruka/goflow/utils/smtpx"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	defer smtpx.SetSender(smtpx.DefaultSender)

	sender := smtpx.NewMockSender("")
	smtpx.SetSender(sender)

	svc := NewService("mail.temba.io", 255, "leah", "pass123", "updates@temba.io")
	err := svc.Send(nil, []string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Updates", "Have a great week")

	assert.NoError(t, err)
	assert.Equal(t, []string{"HELO localhost\nMAIL FROM:<updates@temba.io>\nRCPT TO:<bob@nyaruka.com>\nRCPT TO:<jim@nyaruka.com>\nDATA\nHave a great week\n.\nQUIT\n"}, sender.Logs())
}

func TestNewServiceFromURL(t *testing.T) {
	defer smtpx.SetSender(smtpx.DefaultSender)

	sender := smtpx.NewMockSender("")
	smtpx.SetSender(sender)

	_, err := NewServiceFromURL(":")
	assert.EqualError(t, err, "malformed connection URL")

	_, err = NewServiceFromURL("http://")
	assert.EqualError(t, err, "connection URL must use SMTP scheme")

	_, err = NewServiceFromURL("smtp://temba.io:1234567890")
	assert.EqualError(t, err, "1234567890 is not a valid port number")

	_, err = NewServiceFromURL("smtp://temba.io:25")
	assert.EqualError(t, err, "missing credentials in connection URL")

	// from and port are optional
	svc, err := NewServiceFromURL("smtp://leah:pass123@temba.io")
	assert.NoError(t, err)
	assert.Equal(t, "temba.io", svc.(*service).host)
	assert.Equal(t, 25, svc.(*service).port)
	assert.Equal(t, "leah", svc.(*service).username)
	assert.Equal(t, "pass123", svc.(*service).password)
	assert.Equal(t, "leah@temba.io", svc.(*service).from)

	svc, err = NewServiceFromURL("smtp://leah:pass123@temba.io:255?from=updates@temba.io")
	assert.NoError(t, err)
	assert.Equal(t, "temba.io", svc.(*service).host)
	assert.Equal(t, 255, svc.(*service).port)
	assert.Equal(t, "leah", svc.(*service).username)
	assert.Equal(t, "pass123", svc.(*service).password)
	assert.Equal(t, "updates@temba.io", svc.(*service).from)
}
