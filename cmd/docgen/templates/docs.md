# Container

Flow definitions are defined as a list of nodes, the first node being the entry into the flow. The simplest possible flow containing no nodes whatsoever (and therefore being a no-op) can be defined as follows and includes only the UUID of the flow, its name and the authoring language for the flow:

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
        "type": "send_msg",
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
        "type": "send_msg",
        "text": "What is your name?"
    }],
    "exits": [{
        "uuid":"eb7defc9-3c66-4dfc-80bc-825567ccd9de",
        "destination_node_uuid":"ee0bee3f-34b3-4275-af78-f9ff52c82e6a"
    }]
}
```

# Routers

## Switch

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

# Context

Flows do not describe data flow but rather actions and logic branching. As such, variables collected in a flow and the state of the flow are accessed through
what is called the context. The context contains variables representing the current contact in a flow, the last input from that contact
as well as the results collected in a flow and any webhook requests made during the flow. Variables in the context may be referred to 
within actions by using the `@` symbol. For example, to greet a contact by their name in a [send_msg](#actions:send_msg) action, the text of the action can be `Hi @contact.name!`.

The `@` symbol can be escaped in templates by repeating it, ie, `Hi @@twitter` would output `Hi @twitter`.

The context contains the following top-level variables:

 * `contact` the [contact](#contacts) of the current flow run
 * `run` the current [run](#runs)
 * `parent` the parent of the current [run](#runs), i.e. the run that started the current run
 * `child` the child of the current [run](#runs), i.e. the last subflow
 * `trigger` the [trigger](#triggers) that initiated this session

The following types appear in the context:

 * [Channels](#channels)
 * [Contacts](#contacts)
 * [Flows](#flows)
 * [Groups](#groups)
 * [Inputs](#inputs)
 * [Results](#results)
 * [Runs](#runs)
 * [Triggers](#triggers)
 * [URNs](#urns)
 * [Webhooks](#webhooks)

<div class="context">

## Channels

A channel represents a means for sending and receiving input during a flow run.

A channel renders as its name in a template, and has the following properties which can be accessed:

 * `uuid` the UUID of the channel
 * `name` the name of the channel
 * `address` the address of the channel

### Examples

```
@contact.channel → My Android Phone
@contact.channel.name → My Android Phone
@contact.channel.address → +16303455678
@run.input.channel.uuid → c42528a5-8550-480e-ae4d-92995550e1d6
@(json(contact.channel)) → {"uuid": "c42528a5-8550-480e-ae4d-92995550e1d6", "name": "My Android Phone", "address": "+16303455678"}
```

## Contacts

A contact represents a person who is interacting with the flow.

A contact renders as the person's name (or perferred URN if name isn't set) in a template, and has the following properties which can be accessed:

 * `uuid` the UUID of the contact
 * `name` the full name of the contact
 * `first_name` the first name of the contact
 * `language` the [ISO-639-3](http://www-01.sil.org/iso639-3/) language code of the contact
 * `urns` all [URNs](#urns) the contact has set
 * `urns.[scheme]` all the [URNs](#urns) the contact has set for the particular URN scheme
 * `urn` shorthand for `@(format_urn(c.urns.0))`, i.e. the contact's preferred [URN](#urns) in friendly formatting
 * `groups` all the [groups](#groups) that the contact belongs to
 * `fields` all the custom contact fields the contact has set
 * `fields.[snaked_field_name]` the value of the specific field
 * `channel` shorthand for `contact.urns.0.channel`, i.e. the [channel](#channels) of the contact's preferred URN

### Examples

```
@contact → Bobby Smith
@contact.name → Bobby Smith
@contact.first_name → Bobby
@contact.language → eng
@contact.urns → tel:+12065551212, tel:+16302425788, mailto:foo@bar.com
@contact.urns.0 → tel:+12065551212
@contact.urns.tel → tel:+12065551212, tel:+16302425788
@contact.mailto.0 → mailto:foo@bar.com
@contact.urn → (206) 555 1212
@contact.groups → Males, Reporters
@contact.fields → age: 36\ngender: MALE
@contact.fields.age → 36
@contact.fields.gender → MALE
```

## Flows

A flow describes the ordered logic of actions and routers.

A flow renders as its name in a template, and has the following properties which can be accessed:

 * `uuid` the UUID of the flow
 * `name` the name of the flow

### Examples

```
@run.flow → Registration
@child.flow → Age Collection
@run.flow.uuid → 8eba5c7d-d7cb-4ebe-af7f-7d84bea870c5
@(json(run.flow)) → {"uuid": "8eba5c7d-d7cb-4ebe-af7f-7d84bea870c5", "name": "Registration"}
```

## Groups

A group represents a grouping of contacts. It can be static (contacts are added and removed manually through [actions](#actions:add_contact_group)) or dynamic (contacts are added automatically by a query).

A group renders as its name in a template, and has the following properties which can be accessed:

 * `uuid` the UUID of the group
 * `name` the name of the group

### Examples

```
@contact.groups → Males, Reporters
@contact.groups.0.uuid → 8ddfda9c-9ea7-451e-a812-1c3153f91a87
@contact.groups.1.name → Reporters
@(json(contact.groups.1)) → {"uuid": "8ddfda9c-9ea7-451e-a812-1c3153f91a87", "name": "Reporters"}
```

## Inputs

An input describes input from the contact and currently we only support one type of input: `msg`.

Any input has the following properties which can be accessed:

 * `uuid` the UUID of the input
 * `type` the type of the input, e.g. `msg`
 * `channel` the [channel](#channels) that the input was received on
 * `created_on` the time when the input was created

An input of type `msg` renders as its text and attachments in a template, and has the following additional properties:

 * `text` the text of the message
 * `attachments` any attachments on the message
 * `urn` the [URN](#urns) that the input was received on

### Examples

```
@run.input → Hello\nimage/jpeg:https://example.com/test.jpg
@run.input.type → msg
@run.input.text → Hello
@run.input.attachments → image/jpeg:https://example.com/test.jpg
@(json(run.input)) → {"uuid": "8ddfda9c-9ea7-451e-a812-1c3153f91a87", "text": "Hello", "attachments": ["image/jpeg:https://example.com/test.jpg"], "created_on": "2000-01-01T00:00:00.000000000-00:00"}
```

## Results

A result describes a value captured during a run's execution. It might have been implicitly created by a router, or explicitly created by a [set_run_result](#actions:set_run_result) action.

A result renders as its value in a template, and has the following properties which can be accessed:

 * `value` the value of the result
 * `category` the category of the result
 * `category_localized` the localized category of the result
 * `created_on` the time when the result was created

### Examples

```
@run.results.color → red
@run.results.color.value → red
@run.results.color.category → Red
@run.results.color.category_localized → Rojo
```

## Runs

A run is a single contact's journey through a flow. It records the path they have taken, and the results that have been collected.

A run has several properties which can be accessed in expressions:

 * `uuid` the UUID of the run
 * `flow` the [flow](#flows) of the run
 * `contact` the [contact](#contacts) of the flow run
 * `input` the [input](#inputs) of the current run
 * `results` the results that have been saved for this run
 * `results.[snaked_result_name]` the value of the specific result, e.g. `run.results.age`
 * `webhook` the last [webhook](#webhooks) call made in the current run

## Triggers

A trigger represents something which can initiate a session with the flow engine.

A trigger has several properties which can be accessed in expressions:

 * `type` the type of the trigger, one of `manual` or `flow`
 * `params` the parameters passed to the trigger

### Examples

```
@trigger.type → manual
@trigger.params → source: website\naddress:\n  state: WA
@(json(trigger.params)) → {"source": "website", "address": {"state": "WA"}}
```

## URNs

A URN represents a destination for an outgoing message or a source of an incoming message. It is string composed of 3 components: scheme, path, and display (optional). For example:

 * _tel:+16303524567_
 * _twitterid:54784326227#nyaruka_
 * _telegram:34642632786#bobby_

A URN has several properties which can be accessed in expressions:

 * `scheme` the scheme of the URN, e.g. "tel", "twitter"
 * `path` the path of the URN, e.g. "+16303524567"
 * `display` the display portion of the URN, e.g. "+16303524567"
 * `channel` the preferred [channel](#channels) of the URN

To render a URN in a human friendly format, use the [format_urn](#functions:format_urn) function.

```
@contact.urns.0 → tel:+12065551212
@contact.urns.0.scheme → tel
@contact.urns.0.path → +12065551212
@contact.urns.1.display → nyaruka
@(format_urn(contact.urns.0)) → (206) 555 1212
```

## Webhooks

A webhook describes a call made to an external service.

A webhook has several properties which can be accessed in expressions:

 * `status` the status of the webhook - one of "success", "connection_error" or "response_error"
 * `status_code` the status code of the response
 * `body` the body of the response
 * `json` the parsed JSON response (if response body was JSON)
 * `json.[key]` sub-elements of the parsed JSON response
 * `request` the raw request made, including headers
 * `response` the raw response received, including headers

### Examples

```
@run.webhook.status_code → 200
@run.webhook.json.results.0.state_name → Washington
```

</div>

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
{{ .EventDocs }}
</div>

</body>
</html>


