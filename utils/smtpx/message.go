package smtpx

// Message is email message
type Message struct {
	recipients []string
	subject    string
	text       string
	html       string
}

// NewMessage creates a new message
func NewMessage(recipients []string, subject, text, html string) *Message {
	return &Message{
		recipients: recipients,
		subject:    subject,
		text:       text,
		html:       html,
	}
}
