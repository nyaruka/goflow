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
 * `wait` an optional pause in the flow waiting for some event to occur, such as a contact responding, a timeout for that response or a subflow completing
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

# Routers

The primary responsibility of a router is to choose an exit on the node, but they can also create results. Different router types
have different logic for how an exit will be chosen.

## Switch

If a node wishes to route differently based on some state in the session, it can add a `switch` router which defines one or more 
`cases`.  Each case defines a `type` which is the name of an expression function that is run by passing the evaluation of `operand` 
as the first argument. Cases may define additional arguments using the `arguments` array on a case. If no case evaluates 
to true, then the `default_category_uuid` will be used otherwise flow execution will stop.

A switch router may also define a `result_name` parameters which will save the result of the case which evaluated as true.

A switch router consists of:

 * `result_name` the name of the result which should be written when the switch is evaluated (optional)
 * `operand` the expression which will be evaluated against each of our cases
 * `cases` a list of 1-n cases which are evaluated in order until one is true
 * `default_category_uuid` the uuid of the default category to take if no case matches (optional)

Each case consists of:

 * `uuid` the UUID
 * `type` the type of this test, this must be an excellent test (see below) and will be passed the value of the switch's operand as its first value
 * `arguments` an optional list of templates which can be passed as extra parameters to the test (after the initial operand)
 * `category_uuid` the uuid of the category that should be taken if this case evaluated to true

 An example switch router that tests for the input not being empty:

```json
{
    "uuid": "ee0bee3f-34b3-4275-af78-f9ff52c82e6a",
    "router": {
        "type":"switch",
        "categories": [
            {
                "uuid": "cab600f5-b54b-49b9-a7ea-5638f4cbf2b4",
                "name": "Has Name",
                "exit_uuid": "972fb580-54c2-4491-8438-09ace3500ba5"
            },
            {
                "uuid": "9574fbfd-510f-4dfc-b989-97d2aecf50b9",
                "name": "Other",
                "exit_uuid": "6981b1a9-af04-4e26-a248-1fc1f5e5c7eb"
            }
        ],
        "operand": "@input",
        "cases": [
            {
                "uuid": "6f78d564-029b-4715-b8d4-b28daeae4f24",
                "type": "has_text",
                "category_uuid": "cab600f5-b54b-49b9-a7ea-5638f4cbf2b4"
            }
        ],
        "default_category_uuid": "9574fbfd-510f-4dfc-b989-97d2aecf50b9"
    },
    "exits": [
        {
            "uuid": "972fb580-54c2-4491-8438-09ace3500ba5",
            "destination_uuid": "deec1dd4-b727-4b21-800a-0b7bbd146a82"
        },
        {
            "uuid": "6981b1a9-af04-4e26-a248-1fc1f5e5c7eb",
            "destination_uuid": "ee0bee3f-34b3-4275-af78-f9ff52c82e6a"
        }
    ]
}
```

# Waits

A node can indicate that it needs more information to continue by containing a wait.

## Msg

This wait type indicates that flow execution should pause until an incoming message is received and also gives an optional timeout in seconds as to when the flow 
should continue even if there is no reply:

```json
{
    "type": "msg",
    "timeout": 600
}
```

## Nothing

This wait type indicates that the caller can resume the session immediately with no incoming message or any other input. This type of
wait enables the caller to commit changes in the session up to that point in the flow.

```json
{
    "type": "nothing"
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
