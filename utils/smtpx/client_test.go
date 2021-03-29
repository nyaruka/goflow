package smtpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientFromURL(t *testing.T) {
	_, err := NewClientFromURL(":")
	assert.EqualError(t, err, "malformed connection URL")

	_, err = NewClientFromURL("http://")
	assert.EqualError(t, err, "connection URL must use SMTP scheme")

	_, err = NewClientFromURL("smtp://temba.io:1234567890")
	assert.EqualError(t, err, "1234567890 is not a valid port number")

	_, err = NewClientFromURL("smtp://temba.io:25")
	assert.EqualError(t, err, "missing credentials in connection URL")

	// from and port are optional
	s, err := NewClientFromURL("smtp://leah:pass123%21@temba.io")
	assert.NoError(t, err)
	assert.Equal(t, "temba.io", s.host)
	assert.Equal(t, 25, s.port)
	assert.Equal(t, "leah", s.username)
	assert.Equal(t, "pass123!", s.password)
	assert.Equal(t, "leah@temba.io", s.from)

	s, err = NewClientFromURL("smtp://leah%40nyaruka.com:pass123%21@temba.io:255?from=Leah+%3Cupdates%40temba.io%3E")
	assert.NoError(t, err)
	assert.Equal(t, "temba.io", s.host)
	assert.Equal(t, 255, s.port)
	assert.Equal(t, "leah@nyaruka.com", s.username)
	assert.Equal(t, "pass123!", s.password)
	assert.Equal(t, "Leah <updates@temba.io>", s.from)
}
