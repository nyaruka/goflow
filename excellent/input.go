package excellent

import (
	"bufio"
)

const eof rune = rune(0)

// provides functionality for reading and unreading from a buffered source
type xinput struct {
	base        *bufio.Reader
	unreadRunes []rune
	unreadCount int
}

func newInput(base *bufio.Reader) *xinput {
	return &xinput{
		base:        base,
		unreadRunes: make([]rune, 4),
	}
}

// gets the next rune or EOF if we are at the end of the string
func (r *xinput) read() rune {
	// first see if we have any unread runes to return
	if r.unreadCount > 0 {
		ch := r.unreadRunes[r.unreadCount-1]
		r.unreadCount--
		return ch
	}

	// otherwise, read the next run
	ch, _, err := r.base.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// pops the passed in rune as the next rune to be returned
func (r *xinput) unread(ch rune) {
	r.unreadRunes[r.unreadCount] = ch
	r.unreadCount++
}
