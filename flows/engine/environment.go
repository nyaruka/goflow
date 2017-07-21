package engine

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// NewSessionEnvironment creates and returns a new NewSessionEnvironment given the passed in environment and flow map
func NewSessionEnvironment(env utils.Environment, flowList []flows.Flow, channelList []flows.Channel, contactList []*flows.Contact) flows.SessionEnvironment {
	flowMap := make(map[flows.FlowUUID]flows.Flow, len(flowList))
	for _, f := range flowList {
		flowMap[f.UUID()] = f
	}

	channelMap := make(map[flows.ChannelUUID]flows.Channel, len(channelList))
	for _, c := range channelList {
		channelMap[c.UUID()] = c
	}

	contactMap := make(map[flows.ContactUUID]*flows.Contact, len(contactList))
	for _, c := range contactList {
		contactMap[c.UUID()] = c
	}

	return &sessionEnvironment{env, flowMap, channelMap, contactMap}
}

type sessionEnvironment struct {
	utils.Environment
	flows    map[flows.FlowUUID]flows.Flow
	channels map[flows.ChannelUUID]flows.Channel
	contacts map[flows.ContactUUID]*flows.Contact
}

func (e *sessionEnvironment) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	flow, exists := e.flows[uuid]
	if exists {
		return flow, nil
	}
	return nil, fmt.Errorf("unable to find flow with UUID: %s", uuid)
}

func (e *sessionEnvironment) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	channel, exists := e.channels[uuid]
	if exists {
		return channel, nil
	}
	return nil, fmt.Errorf("unable to find channel with UUID: %s %d", uuid)
}

func (e *sessionEnvironment) GetContact(uuid flows.ContactUUID) (*flows.Contact, error) {
	contact, exists := e.contacts[uuid]
	if exists {
		return contact, nil
	}
	return nil, fmt.Errorf("unable to find contact with UUID: %s", uuid)
}
