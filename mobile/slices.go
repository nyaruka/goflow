package mobile

// Because gomobile currently only supports slices of bytes

// StringSlice wraps a slice of strings
type StringSlice struct {
	items []string
}

// NewStringSlice creates a new string list
func NewStringSlice(capacity int) *StringSlice {
	return &StringSlice{items: make([]string, 0, capacity)}
}

// Add adds an string to this slice
func (l *StringSlice) Add(item string) {
	l.items = append(l.items, item)
}

// Length gets the length of this slice
func (l *StringSlice) Length() int {
	return len(l.items)
}

// Get returns the string at the given index
func (l *StringSlice) Get(index int) string {
	return l.items[index]
}

// EventSlice wraps a slice of events
type EventSlice struct {
	items []*Event
}

// NewEventSlice creates a new slice of events
func NewEventSlice(capacity int) *EventSlice {
	return &EventSlice{items: make([]*Event, 0, capacity)}
}

// Add adds an event to this slice
func (l *EventSlice) Add(item *Event) {
	l.items = append(l.items, item)
}

// Length gets the length of this slice
func (l *EventSlice) Length() int {
	return len(l.items)
}

// Get returns the event at the given index
func (l *EventSlice) Get(index int) *Event {
	return l.items[index]
}

// ModifierSlice wraps a slice of modifiers
type ModifierSlice struct {
	items []*Modifier
}

// NewModifierSlice creates a new slice of modifiers
func NewModifierSlice(capacity int) *ModifierSlice {
	return &ModifierSlice{items: make([]*Modifier, 0, capacity)}
}

// Add adds an event to this slice
func (l *ModifierSlice) Add(item *Modifier) {
	l.items = append(l.items, item)
}

// Length gets the length of this slice
func (l *ModifierSlice) Length() int {
	return len(l.items)
}

// Get returns the event at the given index
func (l *ModifierSlice) Get(index int) *Modifier {
	return l.items[index]
}
