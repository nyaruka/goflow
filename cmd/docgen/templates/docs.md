
<html>
<center><h1>Flow Specification</h1></center>

# Container

Flow definitions are defined as a list of nodes, the first node being the entry into the flow. The simplest possible Flow containing no nodes whatsoever (and therefore being a no-op) can be defined as follows and includes only the UUID of the flow, its name and the authoring language for the flow:

```json
{
    "name": "Empty Flow",
    "uuid": "b7bb5e7c-ad49-4e65-9e24-bf7f1e4ff00a",
    "language": "eng",
    "nodes": []
}
```


# Nodes

Flow definitions are composed of zero or more Nodes, the first node is always the entry node.

A Node consists of:

 * `actions` a list of 0-n actions which will be executed upon first entering a Node
 * `wait` an optional pause in the flow waiting for some event to occur, such as a contact responding, a timeout for that response or a subflow completing
 * `exit` a list of 0-n exits which can be used to link to other Nodes
 * `router` an optional router which determines which exit to take

At its simplest, a node can be just a single action with no exits, wait or router, such as:

```json
{
    "uuid":"5a06445e-d790-4bd3-a10b-b47bdcc9abed",
    "actions":[{
        "uuid": "abc0a2bf-6b4a-4ee0-83e1-1eebae6948ac",
        "type": "reply",
        "text": "What is your name?"
    }]
}
```

If a node wishes to route to another node, it can do so by defining one or more exits, each with the UUID of the node that is next. Without a router defined, the first exit will always be taken. 

An exit consists of:

 * `uuid` the uuid of this exit 
 * `destination_node_uuid` the uuid of the node that should be visited if this exit is chosen by the router (optional)
 * `name` a name for this exit (optional)

```json
{
    "uuid":"5a06445e-d790-4bd3-a10b-b47bdcc9abed",
    "actions":[{
        "uuid": "abc0a2bf-6b4a-4ee0-83e1-1eebae6948ac",
        "type": "reply",
        "text": "What is your name?"
    }],
    "exits": [{
        "uuid":"eb7defc9-3c66-4dfc-80bc-825567ccd9de",
        "destination_node_uuid":"ee0bee3f-34b3-4275-af78-f9ff52c82e6a"
    }]
}
```

# Switch Router

If a node wishes to route differently based on some state, it can add a `switch` router which defines one or more `cases`. Each case defines a `type` which is the name 
of an expression function that is run by passing the evaluation of `operand` as the first argument. Cases may define additional arguments using the `arguments` array on a case.
If no case evaluates to true, then the `default_exit_uuid` will be used otherwise flow execution will stop.

A switch router may also define a `result_name` parameters which will save the result of the case which evaluated as true.

A switch router consists of:

 * `operand` the expression which will be evaluated against each of our cases
 * `default_exit_uuid` the uuid of the default exit to take if no case matches (optional)
 * `result_name` the name of the result which should be written when the switch is evaluated (optional)
 * `cases` a list of 1-n cases which are evaluated in order until one is true

Each case consists of:

 * `uuid` a unique uuid for this case
 * `type` the type of this test, this must be an excellent test (see below) and will be passed the value of the switch's operand as its first value
 * `arguments` an optional list of templates which can be passed as extra parameters to the test (after the initial operand)
 * `exit_uuid` the uuid of the exit that should be taken if this case evaluated to true

 An example switch router that tests for the input not being empty:

```json
{
    "uuid":"ee0bee3f-34b3-4275-af78-f9ff52c82e6a",
    "router": {
        "type":"switch",
        "operand": "@run.input",
        "default_exit_uuid": "9574fbfd-510f-4dfc-b989-97d2aecf50b9",
        "cases": [{
            "uuid": "6f78d564-029b-4715-b8d4-b28daeae4f24",
            "type": "has_text",
            "exit_uuid": "cab600f5-b54b-49b9-a7ea-5638f4cbf2b4"
        }]
    },
    "exits": [{
        "uuid":"cab600f5-b54b-49b9-a7ea-5638f4cbf2b4",
        "name":"Has Name",
        "destination_node_uuid":"deec1dd4-b727-4b21-800a-0b7bbd146a82"
    },{
        "uuid":"9574fbfd-510f-4dfc-b989-97d2aecf50b9",
        "name":"Other",
        "destination_node_uuid":"ee0bee3f-34b3-4275-af78-f9ff52c82e6a"
    }]
}
```

# Waits

A node can indicate that it needs more information to continue by containing a wait. The two types of waits currently are `msg` and `flow`. The former 
indicates that flow execution should pause until an incoming messages is received and also gives an optional timeout in seconds as to when the flow should
continue even if there is no reply:

```
{
    "type": "msg",
    "timeout": 600
}
```

Another wait type is `flow` which indicates that flow execution should pause until the flow with the passed in UUID has completed. The flow should have
have been started by a `start_flow` action previously in the flow:

```
{
    "type": "flow",
    "flow_uuid": "235e70c4-c808-44e5-8f85-0a2021e9dbc9"
}
```

# Context

Flows do not describe data flow but rather actions and logic branching. As such, variables collected in a flow and the state of the flow is accessed through
what is called the Context. The context contains variables representing the current contact in a flow, the last input from that contact
as well as the results collected in a flow and any webhook requests made during the flow. Variables in the context may be referred to 
within actions by using the `@` symbol. For example, to greet a contact by their name in a reply action, the text of the reply can be `Hi @contact.name!`.

The `@` symbol can be escaped in templates by repeating it, ie, `Hi @@twitter` would output `Hi @twitter`.

The flow context contains the following variables:

 * `contact` variables on the contact, shorthand for the contact name
 * `contact.uuid` the uuid of the contact
 * `contact.name` the name of the contact
 * `contact.language` the language of the contact
 * `contact.urns` all URNs the contact has set, shorthand for the preferred URN
 * `contact.urns.[scheme]` all the URNs the contact has set for the particular URN scheme, shorthand for preferred URN of that scheme
 * `contact.urns.tel` the preferred phone number fo the contact (tel is the scheme)
 * `contact.fields` all custom contact fields the contact has set
 * `contact.fields.[snaked_field_name]` the value of the specific field, ex: contact.fields.age or contact.fields.first_name

 * `run.uuid` the uuid of the current flow run
 * `run.results` the results that have been saved for this run
 * `run.results.[snaked_result_name]` the value of the specific result, ex: run.results.age
 * `run.results.[snaked_result_name].category` the category (if any) of the specific result, ex: run.results.age.category

 * `run.child` values associated with the last subflow started
 * `run.child.results` the result values saved by the last subflow started
 * `run.child.results.[snaked_result_name]` the value of a specific value of the last subflow started
 * `run.child.results.[snaked_result_name].category` the category (if any) of a specific value of the last subflow started

 * `run.parent` values associated with the run that started this run
 * `run.parent.results` the result values saved by run that started this run
 * `run.parent.results.[snaked_result_name]` the value of a specific value of the run that started this run
 * `run.parent.results.[snaked_result_name].category` the category (if any) of a specific value of the run that started this run

 * `run.input` values of the last input from the contact
 * `run.input.text` the text of the last input from the contact
 * `run.input.urn` the urn that was used to send the last input
 * `run.input.channel` the channel that was used to send the last input

 * `webhook` values associated with the last webhook called, shorthand for `webhook.json`
 * `webhook.status` the status of the last webhook called one of "S" for success, "E" for error or "F" for failure
 * `webhook.status_code` the status code returned by the last webhook called
 * `webhook.json` the parsed JSON response to the last webhook (if response was JSON)
 * `webhook.json.[keys]` sub-elements of the parsed JSON response to the last webhook, ex: webhook.json.results.0.state_name
 * `webhook.request` the raw request made to the webhook, including headers
 * `webhook.response` the raw response from the webhook, including headers

# Template Functions

In addition to simple substitutions, flows also have access to a set of functions which can be used in templates to further manipulate the context.
Functions are called using the `@(function_name(args..))` syntax. For example, to title case a contact's name in a message, you can use `@(title(contact.name))`. 
Context variables referred to within functions do not need a leading `@`. Functions can also use literal numbers or strings as arguments, for example
`@(array_length(split("1 2 3", " "))`.

<div class="excellent_functions">
{{ .ExcellentFunctionDocs }}
</div>

# Template Tests

Template tests are a special class of functions which are used within the switch router. They are called in the same way as normal functions, but 
all return a test result object which by default evalutes to true or false, but can also be used to find the matching portion of the test by using
the `match` component of the result. The flow editor builds these expressions using UI widgets, but they can be used anywhere a normal template
function is used.

<div class="excellent_tests">
{{ .ExcellentTestDocs }}
</div>

# Action Definitions

Actions on a node generate events which can then be ingested by the engine container. In some cases the actions cause an immediate action, such 
as calling a webhook, in others the engine container is responsible for taking the action based on the event that is output, such as sending 
messages or updating contact fields. In either case the internal state of the engine is always updated to represent the new state so that
flow execution is consistent. For example, while the engine itself does not have access to a contact store, it updates its internal 
representation of a contact's state based on action performed on a flow so that later references in the flow are correct.

<div class="actions">
{{ .ActionDocs }}
</div>

# Event Definitions

Events are the output of a flow run and represent instructions to the engine container on what actions should be taken due to the flow execution.
All templates in events have been evaluated and can be used to create concrete messages, contact updates, emails etc by the container.

<div class="events">
{{ .EventDocs}}
</div>

</body>
</html>


