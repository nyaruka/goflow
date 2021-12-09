package engine

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSprint(t *testing.T) {
	env := envs.NewBuilder().Build()

	defer dates.SetNowSource(dates.DefaultNowSource)

	source, err := static.NewSource([]byte(`{
		"flows": [
			{
				"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
				"name": "Empty Flow",
				"spec_version": "13.0",
				"language": "eng",
				"type": "messaging",
				"nodes": [
					{
						"uuid": "d6cdbd1b-d7db-4a38-a22b-9ec357fa228c",
						"exits": [
							{
								"uuid": "c0f31cdf-bc9a-404f-88c3-9d6c39d345c9",
								"destination_uuid": "1747f81b-3692-4ef0-81c9-921c1124cf61"
							}
						]
					},
					{
						"uuid": "1747f81b-3692-4ef0-81c9-921c1124cf61",
						"exits": [
							{
								"uuid": "fcf6d3b9-f611-4b37-96e2-655d2a46b049",
								"destination_uuid": "597fba02-b996-4f41-a842-d8962817fff9"
							}
						]
					},
					{
						"uuid": "597fba02-b996-4f41-a842-d8962817fff9",
						"exits": [
							{
								"uuid": "4ca632d5-67f1-41fa-9528-aa77a22a029b",
								"destination_uuid": null
							}
						]
					}
				]
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	assert.Equal(t, source, sa.Source())

	flow, err := sa.Flows().Get("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	require.NoError(t, err)
	node1 := flow.Nodes()[0]
	node1Exit1 := node1.Exits()[0]
	node2 := flow.Nodes()[1]
	node2Exit1 := node2.Exits()[0]
	node3 := flow.Nodes()[2]

	mod1 := modifiers.NewName("Bob")
	mod2 := modifiers.NewName("Joe")

	event1 := events.NewError(errors.New("error 1"))
	event2 := events.NewError(errors.New("error 1"))

	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2021, 12, 8, 10, 13, 30, 0, time.UTC)))

	sprint := newEmptySprint()
	sprint.logSegment(flow, node1, node1Exit1, "yes", node2)
	sprint.logModifier(mod1)
	sprint.logEvent(event1)
	sprint.logSegment(flow, node2, node2Exit1, "", node3)
	sprint.logModifier(mod2)
	sprint.logEvent(event2)

	var seg1 flows.Segment = &segment{flow: flow, node: node1, exit: node1Exit1, operand: "yes", destination: node2, time: time.Date(2021, 12, 8, 10, 13, 30, 0, time.UTC)}
	var seg2 flows.Segment = &segment{flow: flow, node: node2, exit: node2Exit1, destination: node3, time: time.Date(2021, 12, 8, 10, 13, 31, 0, time.UTC)}

	assert.Equal(t, []flows.Modifier{mod1, mod2}, sprint.Modifiers())
	assert.Equal(t, []flows.Event{event1, event2}, sprint.Events())
	assert.Equal(t, []flows.Segment{seg1, seg2}, sprint.Segments())

	seg := sprint.Segments()[0]
	assert.Equal(t, flow, seg.Flow())
	assert.Equal(t, node1, seg.Node())
	assert.Equal(t, node1Exit1, seg.Exit())
	assert.Equal(t, "yes", seg.Operand())
	assert.Equal(t, node2, seg.Destination())
	assert.Equal(t, time.Date(2021, 12, 8, 10, 13, 30, 0, time.UTC), seg.Time())
	assert.Equal(t,
		`{"flow_uuid":"76f0a02f-3b75-4b86-9064-e9195e1b3a02","node_uuid":"d6cdbd1b-d7db-4a38-a22b-9ec357fa228c","exit_uuid":"c0f31cdf-bc9a-404f-88c3-9d6c39d345c9","operand":"yes","destination_uuid":"1747f81b-3692-4ef0-81c9-921c1124cf61","time":"2021-12-08T10:13:30Z"}`,
		string(jsonx.MustMarshal(seg)),
	)

	assert.Equal(t, sprint, NewSprint(
		[]flows.Modifier{mod1, mod2},
		[]flows.Event{event1, event2},
		[]flows.Segment{seg1, seg2},
	))
}
