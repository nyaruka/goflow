package engine_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSprint(t *testing.T) {
	mod1 := modifiers.NewName("Bob")
	mod2 := modifiers.NewName("Joe")

	event1 := events.NewError(errors.New("error 1"))
	event2 := events.NewError(errors.New("error 1"))

	sprint := engine.NewSprint([]flows.Modifier{mod1}, []flows.Event{event1})
	sprint.LogModifier(mod2)
	sprint.LogEvent(event2)

	assert.Equal(t, []flows.Modifier{mod1, mod2}, sprint.Modifiers())
	assert.Equal(t, []flows.Event{event1, event2}, sprint.Events())
}
