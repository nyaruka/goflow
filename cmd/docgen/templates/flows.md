# Container

Flow definitions are defined as a list of nodes, the first node being the entry into the flow. The simplest possible flow containing no nodes whatsoever (and therefore being a no-op) contains the following fields:

 * `uuid` the UUID
 * `name` the name
 * `language` the base authoring language used for localization
 * `type` the type, one of `messaging`, `messaging_offline`, `voice` whcih determines which actions are allowed in the flow
 * `nodes` the nodes (may be empty)

For example:

```json
{
    "uuid": "b7bb5e7c-ad49-4e65-9e24-bf7f1e4ff00a",
    "name": "Empty Flow",
    "language": "eng",
    "type": "messaging",
    "nodes": []
}
```

# Nodes

Flow definitions are composed of zero or more nodes, the first node is always the entry node.

A Node consists of:

 * `uuid` the UUID
 * `actions` a list of 0-n actions which will be executed upon first entering a node
 * `router` an optional router which determines which exit to take
 * `exit` a list of 0-n exits which can be used to link to other nodes

At its simplest, a node can be just a single action with no exits, wait or router, such as:

```json
{
    "uuid":"5a06445e-d790-4bd3-a10b-b47bdcc9abed",
    "actions":[{
        "uuid": "abc0a2bf-6b4a-4ee0-83e1-1eebae6948ac",
        "type": "send_msg",
        "text": "What is your name?"
    }]
}
```

If a node wishes to route to another node, it can do so by defining one or more exits, each with the UUID of the node that is next. Without a router defined, the first exit will always be taken. 

An exit consists of:

 * `uuid` the UUID
 * `destination_uuid` the uuid of the node that should be visited if this exit is chosen by the router (optional)

```json
{
    "uuid":"5a06445e-d790-4bd3-a10b-b47bdcc9abed",
    "actions":[{
        "uuid": "abc0a2bf-6b4a-4ee0-83e1-1eebae6948ac",
        "type": "send_msg",
        "text": "What is your name?"
    }],
    "exits": [{
        "uuid":"eb7defc9-3c66-4dfc-80bc-825567ccd9de",
        "destination_uuid":"ee0bee3f-34b3-4275-af78-f9ff52c82e6a"
    }]
}
```

# Actions

Actions on a node generate events which can then be ingested by the engine container. In some cases the actions cause an immediate action, such 
as calling a webhook, in others the engine container is responsible for taking the action based on the event that is output, such as sending 
messages or updating contact fields. In either case the internal state of the engine is always updated to represent the new state so that
flow execution is consistent. For example, while the engine itself does not have access to a contact store, it updates its internal 
representation of a contact's state based on action performed on a flow so that later references in the flow are correct.

<div class="actions">
{{ .actionDocs }}
</div>
