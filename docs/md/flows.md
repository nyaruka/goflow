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

Flow definitions are composed of zero or more nodes, the first node is always the entry node.

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
within actions by using the `@` symbol. For example, to greet a contact by their name in a [send_msg](#action:send_msg) action, the text of the action can be `Hi @contact.name!`.

The `@` symbol can be escaped in templates by repeating it, ie, `Hi @@twitter` would output `Hi @twitter`.

The context contains the following top-level variables:

 * `contact` the [contact](#context:contact) of the current flow run
 * `run` the current [run](#context:run)
 * `parent` the parent of the current [run](#context:run), i.e. the run that started the current run
 * `child` the child of the current [run](#context:run), i.e. the last subflow
 * `trigger` the [trigger](#context:trigger) that initiated this session

The following types appear in the context:

 * [Channel](#context:channel)
 * [Contact](#context:contact)
 * [Flow](#context:flow)
 * [Group](#context:group)
 * [Input](#context:input)
 * [Result](#context:result)
 * [Run](#context:run)
 * [Trigger](#context:trigger)
 * [URN](#context:urn)
 * [Webhook](#context:webhook)

<div class="context">
<a name="context:attachment"></a>

## Attachment

Is a media attachment on a message, and it has the following properties which can be accessed:

 * `content_type` the MIME type of the attachment
 * `url` the URL of the attachment

Examples:


```objectivec
@run.input.attachments.0.content_type → image/jpeg
@run.input.attachments.0.url → http://s3.amazon.com/bucket/test.jpg
@(json(run.input.attachments.0)) → {"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test.jpg"}
```

<a name="context:channel"></a>

## Channel

Represents a means for sending and receiving input during a flow run. It renders as its name in a template,
and has the following properties which can be accessed:

 * `uuid` the UUID of the channel
 * `name` the name of the channel
 * `address` the address of the channel

Examples:


```objectivec
@contact.channel → My Android Phone
@contact.channel.name → My Android Phone
@contact.channel.address → +12345671111
@run.input.channel.uuid → 57f1078f-88aa-46f4-a59a-948a5739c03d
@(json(contact.channel)) → {"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"}
```

<a name="context:contact"></a>

## Contact

Represents a person who is interacting with the flow. It renders as the person's name
(or perferred URN if name isn't set) in a template, and has the following properties which can be accessed:

 * `uuid` the UUID of the contact
 * `name` the full name of the contact
 * `first_name` the first name of the contact
 * `language` the [ISO-639-3](http://www-01.sil.org/iso639-3/) language code of the contact
 * `timezone` the timezone name of the contact
 * `created_on` the datetime when the contact was created
 * `urns` all [URNs](#context:urn) the contact has set
 * `urns.[scheme]` all the [URNs](#context:urn) the contact has set for the particular URN scheme
 * `urn` shorthand for `@(format_urn(c.urns.0))`, i.e. the contact's preferred [URN](#context:urn) in friendly formatting
 * `groups` all the [groups](#context:group) that the contact belongs to
 * `fields` all the custom contact fields the contact has set
 * `fields.[snaked_field_name]` the value of the specific field
 * `channel` shorthand for `contact.urns.0.channel`, i.e. the [channel](#context:channel) of the contact's preferred URN

Examples:


```objectivec
@contact → Ryan Lewis
@contact.name → Ryan Lewis
@contact.first_name → Ryan
@contact.language → eng
@contact.timezone → America/Guayaquil
@contact.created_on → 2018-06-20T11:40:30.123456Z
@contact.urns → ["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]
@contact.urns.0 → tel:+12065551212
@contact.urns.tel → ["tel:+12065551212"]
@contact.urns.mailto.0 → mailto:foo@bar.com
@contact.urn → (206) 555-1212
@contact.groups → ["Testers","Males"]
@contact.fields → {"activation_token":"AACC55","age":"23","gender":"Male","join_date":"2017-12-02T00:00:00.000000-02:00"}
@contact.fields.activation_token → AACC55
@contact.fields.gender → Male
```

<a name="context:flow"></a>

## Flow

Describes the ordered logic of actions and routers. It renders as its name in a template, and has the following
properties which can be accessed:

 * `uuid` the UUID of the flow
 * `name` the name of the flow
 * `revision` the revision number of the flow

Examples:


```objectivec
@run.flow → Registration
@child.flow → Collect Age
@run.flow.uuid → 50c3706e-fedb-42c0-8eab-dda3335714b7
@(json(run.flow)) → {"name":"Registration","revision":123,"uuid":"50c3706e-fedb-42c0-8eab-dda3335714b7"}
```

<a name="context:group"></a>

## Group

Represents a grouping of contacts. It can be static (contacts are added and removed manually through
[actions](#action:add_contact_groups)) or dynamic (contacts are added automatically by a query). It renders as its name in a
template, and has the following properties which can be accessed:

 * `uuid` the UUID of the group
 * `name` the name of the group

Examples:


```objectivec
@contact.groups → ["Testers","Males"]
@contact.groups.0.uuid → b7cf0d83-f1c9-411c-96fd-c511a4cfa86d
@contact.groups.1.name → Males
@(json(contact.groups.1)) → {"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}
```

<a name="context:input"></a>

## Input

Describes input from the contact and currently we only support one type of input: `msg`. Any input has the following
properties which can be accessed:

 * `uuid` the UUID of the input
 * `type` the type of the input, e.g. `msg`
 * `channel` the [channel](#context:channel) that the input was received on
 * `created_on` the time when the input was created

An input of type `msg` renders as its text and attachments in a template, and has the following additional properties:

 * `text` the text of the message
 * `attachments` any [attachments](#context:attachment) on the message
 * `urn` the [URN](#context:urn) that the input was received on

Examples:


```objectivec
@run.input → Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3
@run.input.type → msg
@run.input.text → Hi there
@run.input.attachments → ["http://s3.amazon.com/bucket/test.jpg","http://s3.amazon.com/bucket/test.mp3"]
@(json(run.input)) → {"attachments":[{"content_type":"image/jpeg","url":"http://s3.amazon.com/bucket/test.jpg"},{"content_type":"audio/mp3","url":"http://s3.amazon.com/bucket/test.mp3"}],"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2000-01-01T00:00:00.000000Z","text":"Hi there","type":"msg","urn":{"display":"","path":"+12065551212","scheme":"tel"},"uuid":"9bf91c2b-ce58-4cef-aacc-281e03f69ab5"}
```

<a name="context:result"></a>

## Result

Describes a value captured during a run's execution. It might have been implicitly created by a router, or explicitly
created by a [set_run_result](#action:set_run_result) action.It renders as its value in a template, and has the following
properties which can be accessed:

 * `value` the value of the result
 * `category` the category of the result
 * `category_localized` the localized category of the result
 * `input` the input associated with the result
 * `node_uuid` the UUID of the node where the result was created
 * `created_on` the time when the result was created

Examples:


```objectivec
@run.results.favorite_color → red
@run.results.favorite_color.value → red
@run.results.favorite_color.category → Red
```

<a name="context:run"></a>

## Run

Is a single contact's journey through a flow. It records the path they have taken, and the results that have been
collected. It has several properties which can be accessed in expressions:

 * `uuid` the UUID of the run
 * `flow` the [flow](#context:flow) of the run
 * `contact` the [contact](#context:contact) of the flow run
 * `input` the [input](#context:input) of the current run
 * `results` the results that have been saved for this run
 * `results.[snaked_result_name]` the value of the specific result, e.g. `run.results.age`
 * `webhook` the last [webhook](#context:webhook) call made in the current run

Examples:


```objectivec
@run.flow.name → Registration
```

<a name="context:trigger"></a>

## Trigger

Represents something which can initiate a session with the flow engine. It has several properties which can be
accessed in expressions:

 * `type` the type of the trigger, one of "manual" or "flow"
 * `params` the parameters passed to the trigger

Examples:


```objectivec
@trigger.type → flow_action
@trigger.params → {"source": "website","address": {"state": "WA"}}
@(json(trigger)) → {"params":{"source":"website","address":{"state":"WA"}},"type":"flow_action"}
```

<a name="context:urn"></a>

## Urn

Represents a destination for an outgoing message or a source of an incoming message. It is string composed of 3
components: scheme, path, and display (optional). For example:

 - _tel:+16303524567_
 - _twitterid:54784326227#nyaruka_
 - _telegram:34642632786#bobby_

It has several properties which can be accessed in expressions:

 * `scheme` the scheme of the URN, e.g. "tel", "twitter"
 * `path` the path of the URN, e.g. "+16303524567"
 * `display` the display portion of the URN, e.g. "+16303524567"
 * `channel` the preferred no link template for type context of the URN

To render a URN in a human friendly format, use the [format_urn](flows.html#function:format_urn) function.

Examples:


```objectivec
@contact.urns.0 → tel:+12065551212
@contact.urns.0.scheme → tel
@contact.urns.0.path → +12065551212
@contact.urns.1.display → nyaruka
@(format_urn(contact.urns.0)) → (206) 555-1212
@(json(contact.urns.0)) → {"display":"","path":"+12065551212","scheme":"tel"}
```

<a name="context:webhook"></a>

## Webhook

Describes a call made to an external service. It has several properties which can be accessed in expressions:

 * `status` the status of the webhook - one of "success", "connection_error" or "response_error"
 * `status_code` the status code of the response
 * `body` the body of the response
 * `json` the parsed JSON response (if response body was JSON)
 * `json.[key]` sub-elements of the parsed JSON response
 * `request` the raw request made, including headers
 * `response` the raw response received, including headers

Examples:


```objectivec
@run.webhook.status_code → 200
@run.webhook.json.results.0.state → WA
```


</div>

# Template Functions

In addition to simple substitutions, flows also have access to a set of functions which can be used in templates to further manipulate the context.
Functions are called using the `@(function_name(args..))` syntax. For example, to title case a contact's name in a message, you can use `@(title(contact.name))`. 
Context variables referred to within functions do not need a leading `@`. Functions can also use literal numbers or strings as arguments, for example
`@(length(split("1 2 3", " "))`.

<div class="functions">
<a name="function:abs"></a>

## abs(num)

Returns the absolute value of `num`


```objectivec
@(abs(-10)) → 10
@(abs(10.5)) → 10.5
@(abs("foo")) → ERROR
```

<a name="function:and"></a>

## and(tests...)

Returns whether all the passed in arguments are truthy


```objectivec
@(and(true)) → true
@(and(true, false, true)) → false
```

<a name="function:array"></a>

## array(values...)

Takes a list of `values` and returns them as an array


```objectivec
@(array("a", "b", 356)[1]) → b
@(join(array("a", "b", "c"), "|")) → a|b|c
@(length(array())) → 0
@(length(array("a", "b"))) → 2
```

<a name="function:boolean"></a>

## boolean(value)

Tries to convert `value` to a boolean. An error is returned if the value can't be converted.


```objectivec
@(boolean(array(1, 2))) → true
@(boolean("FALSE")) → false
@(boolean(1 / 0)) → ERROR
```

<a name="function:char"></a>

## char(num)

Returns the rune for the passed in codepoint, `num`, which may be unicode, this is the reverse of code


```objectivec
@(char(33)) → !
@(char(128512)) → 😀
@(char("foo")) → ERROR
```

<a name="function:clean"></a>

## clean(text)

Strips any non-printable characters from `text`


```objectivec
@(clean("😃 Hello \nwo\tr\rld")) → 😃 Hello world
@(clean(123)) → 123
```

<a name="function:code"></a>

## code(text)

Returns the numeric code for the first character in `text`, it is the inverse of char


```objectivec
@(code("a")) → 97
@(code("abc")) → 97
@(code("😀")) → 128512
@(code("15")) → 49
@(code(15)) → 49
@(code("")) → ERROR
```

<a name="function:datetime"></a>

## datetime(text)

Turns `text` into a date according to the environment's settings. It will return an error
if it is unable to convert the text to a date.


```objectivec
@(datetime("1979-07-18")) → 1979-07-18T00:00:00.000000-05:00
@(datetime("1979-07-18T10:30:45.123456Z")) → 1979-07-18T10:30:45.123456Z
@(datetime("2010 05 10")) → 2010-05-10T00:00:00.000000-05:00
@(datetime("NOT DATE")) → ERROR
```

<a name="function:datetime_add"></a>

## datetime_add(date, offset, unit)

Calculates the date value arrived at by adding `offset` number of `unit` to the `date`

Valid durations are "Y" for years, "M" for months, "W" for weeks, "D" for days, "h" for hour,
"m" for minutes, "s" for seconds


```objectivec
@(datetime_add("2017-01-15", 5, "D")) → 2017-01-20T00:00:00.000000-05:00
@(datetime_add("2017-01-15 10:45", 30, "m")) → 2017-01-15T11:15:00.000000-05:00
```

<a name="function:datetime_diff"></a>

## datetime_diff(date1, date2, unit)

Returns the integer duration between `date1` and `date2` in the `unit` specified.

Valid durations are "Y" for years, "M" for months, "W" for weeks, "D" for days, "h" for hour,
"m" for minutes, "s" for seconds


```objectivec
@(datetime_diff("2017-01-17", "2017-01-15", "D")) → 2
@(datetime_diff("2017-01-17 10:50", "2017-01-17 12:30", "h")) → -1
@(datetime_diff("2017-01-17", "2015-12-17", "Y")) → 2
```

<a name="function:datetime_from_parts"></a>

## datetime_from_parts(year, month, day)

Converts the passed in `year`, `month` and `day`


```objectivec
@(datetime_from_parts(2017, 1, 15)) → 2017-01-15T00:00:00.000000-05:00
@(datetime_from_parts(2017, 2, 31)) → 2017-03-03T00:00:00.000000-05:00
@(datetime_from_parts(2017, 13, 15)) → ERROR
```

<a name="function:default"></a>

## default(test, default)

Takes two arguments, returning `test` if not an error or nil or empty text, otherwise returning `default`


```objectivec
@(default(undeclared.var, "default_value")) → default_value
@(default("10", "20")) → 10
@(default("", "value")) → value
@(default(array(1, 2), "value")) → ["1","2"]
@(default(array(), "value")) → value
@(default(datetime("invalid-date"), "today")) → today
```

<a name="function:field"></a>

## field(text, offset, delimiter)

Splits `text` based on the passed in `delimiter` and returns the field at `offset`.  When splitting
with a space, the delimiter is considered to be all whitespace.  (first field is 0)


```objectivec
@(field("a,b,c", 1, ",")) → b
@(field("a,,b,c", 1, ",")) →
@(field("a   b c", 1, " ")) → b
@(field("a		b	c	d", 1, "	")) →
@(field("a\t\tb\tc\td", 1, " ")) →
@(field("a,b,c", "foo", ",")) → ERROR
```

<a name="function:format_date"></a>

## format_date(date, [,format])

Turns `date` into text according to the `format` specified.

The format string can consist of the following characters. The characters
' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.

* `YY`        - last two digits of year 0-99
* `YYYY`      - four digits of year 0000-9999
* `M`         - month 1-12
* `MM`        - month 01-12
* `D`         - day of month, 1-31
* `DD`        - day of month, zero padded 0-31


```objectivec
@(format_date("1979-07-18T15:00:00.000000Z")) → 1979-07-18
@(format_date("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) → 1979-07-18
@(format_date("2010-05-10T19:50:00.000000Z", "YYYY M DD")) → 2010 5 10
@(format_date("1979-07-18T15:00:00.000000Z", "YYYY")) → 1979
@(format_date("1979-07-18T15:00:00.000000Z", "M")) → 7
@(format_date("NOT DATE", "YYYY-MM-DD")) → ERROR
```

<a name="function:format_datetime"></a>

## format_datetime(date [,format [,timezone]])

Turns `date` into text according to the `format` specified and in
the optional `timezone`.

The format string can consist of the following characters. The characters
' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.

* `YY`        - last two digits of year 0-99
* `YYYY`      - four digits of year 0000-9999
* `M`         - month 1-12
* `MM`        - month 01-12
* `D`         - day of month, 1-31
* `DD`        - day of month, zero padded 0-31
* `h`         - hour of the day 1-12
* `hh`        - hour of the day 01-12
* `tt`        - twenty four hour of the day 01-23
* `m`         - minute 0-59
* `mm`        - minute 00-59
* `s`         - second 0-59
* `ss`        - second 00-59
* `fff`       - milliseconds
* `ffffff`    - microseconds
* `fffffffff` - nanoseconds
* `aa`        - am or pm
* `AA`        - AM or PM
* `Z`         - hour and minute offset from UTC, or Z for UTC
* `ZZZ`       - hour and minute offset from UTC

Timezone should be a location name as specified in the IANA Time Zone database, such
as "America/Guayaquil" or "America/Los_Angeles". If not specified the timezone of your
environment will be used. An error will be returned if the timezone is not recognized.


```objectivec
@(format_datetime("1979-07-18T15:00:00.000000Z")) → 1979-07-18 10:00
@(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) → 1979-07-18
@(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY M DD tt:mm")) → 2010 5 10 14:50
@(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY-MM-DD tt:mm AA", "America/Los_Angeles")) → 2010-05-10 12:50 PM
@(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY")) → 1979
@(format_datetime("1979-07-18T15:00:00.000000Z", "M")) → 7
@(format_datetime("NOT DATE", "YYYY-MM-DD")) → ERROR
```

<a name="function:format_location"></a>

## format_location(location)

Formats the given location as its name


```objectivec
@(format_location("Rwanda")) → Rwanda
@(format_location("Rwanda > Kigali")) → Kigali
```

<a name="function:format_number"></a>

## format_number(num, places, commas)

Returns `num` formatted with the passed in number of decimal `places` and optional `commas` dividing thousands separators


```objectivec
@(format_number(31337)) → 31,337.00
@(format_number(31337, 2)) → 31,337.00
@(format_number(31337, 2, true)) → 31,337.00
@(format_number(31337, 0, false)) → 31337
@(format_number("foo", 2, false)) → ERROR
```

<a name="function:format_urn"></a>

## format_urn(urn)

Turns `urn` into human friendly text


```objectivec
@(format_urn("tel:+250781234567")) → 0781 234 567
@(format_urn("twitter:134252511151#billy_bob")) → billy_bob
@(format_urn(contact.urns)) → (206) 555-1212
@(format_urn(contact.urns.2)) → foo@bar.com
@(format_urn(contact.urns.mailto)) → foo@bar.com
@(format_urn(contact.urns.mailto.0)) → foo@bar.com
@(format_urn(contact.urns.telegram)) →
@(format_urn("NOT URN")) → ERROR
```

<a name="function:from_epoch"></a>

## from_epoch(num)

Returns a new date created from `num` which represents number of nanoseconds since January 1st, 1970 GMT


```objectivec
@(from_epoch(1497286619000000000)) → 2017-06-12T11:56:59.000000-05:00
```

<a name="function:if"></a>

## if(test, true_value, false_value)

Evaluates the `test` argument, and if truthy returns `true_value`, if not returning `false_value`

If the first argument is an error that error is returned


```objectivec
@(if(1 = 1, "foo", "bar")) → foo
@(if("foo" > "bar", "foo", "bar")) → ERROR
```

<a name="function:join"></a>

## join(array, delimiter)

Joins the passed in `array` of strings with the passed in `delimiter`


```objectivec
@(join(array("a", "b", "c"), "|")) → a|b|c
@(join(split("a.b.c", "."), " ")) → a b c
```

<a name="function:json"></a>

## json(value)

Tries to return a JSON representation of `value`. An error is returned if there is
no JSON representation of that object.


```objectivec
@(json("string")) → "string"
@(json(10)) → 10
@(json(contact.uuid)) → "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
```

<a name="function:left"></a>

## left(text, count)

Returns the `count` most left characters of the passed in `text`


```objectivec
@(left("hello", 2)) → he
@(left("hello", 7)) → hello
@(left("😀😃😄😁", 2)) → 😀😃
@(left("hello", -1)) → ERROR
```

<a name="function:length"></a>

## length(value)

Returns the length of the passed in text or array.

length will return an error if it is passed an item which doesn't have length.


```objectivec
@(length("Hello")) → 5
@(length("😀😃😄😁")) → 4
@(length(array())) → 0
@(length(array("a", "b", "c"))) → 3
@(length(1234)) → ERROR
```

<a name="function:lower"></a>

## lower(text)

Lowercases the passed in `text`


```objectivec
@(lower("HellO")) → hello
@(lower("hello")) → hello
@(lower("123")) → 123
@(lower("😀")) → 😀
```

<a name="function:max"></a>

## max(values...)

Takes a list of `values` and returns the greatest of them


```objectivec
@(max(1, 2)) → 2
@(max(1, -1, 10)) → 10
@(max(1, 10, "foo")) → ERROR
```

<a name="function:mean"></a>

## mean(values)

Takes a list of `values` and returns the arithmetic mean of them


```objectivec
@(mean(1, 2)) → 1.5
@(mean(1, 2, 6)) → 3
@(mean(1, "foo")) → ERROR
```

<a name="function:min"></a>

## min(values)

Takes a list of `values` and returns the smallest of them


```objectivec
@(min(1, 2)) → 1
@(min(2, 2, -10)) → -10
@(min(1, 2, "foo")) → ERROR
```

<a name="function:mod"></a>

## mod(dividend, divisor)

Returns the remainder of the division of `divident` by `divisor`


```objectivec
@(mod(5, 2)) → 1
@(mod(4, 2)) → 0
@(mod(5, "foo")) → ERROR
```

<a name="function:now"></a>

## now()

Returns the current date and time in the environment timezone


```objectivec
@(now()) → 2018-04-11T13:24:30.123456-05:00
```

<a name="function:number"></a>

## number(value)

Tries to convert `value` to a number. An error is returned if the value can't be converted.


```objectivec
@(number(10)) → 10
@(number("123.45000")) → 123.45
@(number("what?")) → ERROR
```

<a name="function:or"></a>

## or(tests...)

Returns whether if any of the passed in arguments are truthy


```objectivec
@(or(true)) → true
@(or(true, false, true)) → true
```

<a name="function:parse_datetime"></a>

## parse_datetime(text, format [,timezone])

Turns `text` into a date according to the `format` and optional `timezone` specified

The format string can consist of the following characters. The characters
' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.

* `YY`        - last two digits of year 0-99
* `YYYY`      - four digits of year 0000-9999
* `M`         - month 1-12
* `MM`        - month 01-12
* `D`         - day of month, 1-31
* `DD`        - day of month, zero padded 0-31
* `h`         - hour of the day 1-12
* `hh`        - hour of the day 01-12
* `tt`        - twenty four hour of the day 01-23
* `m`         - minute 0-59
* `mm`        - minute 00-59
* `s`         - second 0-59
* `ss`        - second 00-59
* `fff`       - milliseconds
* `ffffff`    - microseconds
* `fffffffff` - nanoseconds
* `aa`        - am or pm
* `AA`        - AM or PM
* `Z`         - hour and minute offset from UTC, or Z for UTC
* `ZZZ`       - hour and minute offset from UTC

Timezone should be a location name as specified in the IANA Time Zone database, such
as "America/Guayaquil" or "America/Los_Angeles". If not specified the timezone of your
environment will be used. An error will be returned if the timezone is not recognized.

Note that fractional seconds will be parsed even without an explicit format identifier.
You should only specify fractional seconds when you want to assert the number of places
in the input format.

parse_datetime will return an error if it is unable to convert the text to a datetime.


```objectivec
@(parse_datetime("1979-07-18", "YYYY-MM-DD")) → 1979-07-18T00:00:00.000000-05:00
@(parse_datetime("2010 5 10", "YYYY M DD")) → 2010-05-10T00:00:00.000000-05:00
@(parse_datetime("2010 5 10 12:50", "YYYY M DD tt:mm", "America/Los_Angeles")) → 2010-05-10T12:50:00.000000-07:00
@(parse_datetime("NOT DATE", "YYYY-MM-DD")) → ERROR
```

<a name="function:parse_json"></a>

## parse_json(text)

Tries to parse `text` as JSON, returning a fragment you can index into

If the passed in value is not JSON, then an error is returned


```objectivec
@(parse_json("[1,2,3,4]").2) → 3
@(parse_json("invalid json")) → ERROR
```

<a name="function:percent"></a>

## percent(num)

Converts `num` to text represented as a percentage


```objectivec
@(percent(0.54234)) → 54%
@(percent(1.2)) → 120%
@(percent("foo")) → ERROR
```

<a name="function:rand"></a>

## rand()

Returns a single random number between [0.0-1.0).


```objectivec
@(rand()) → 0.3849275689214193274523267973563633859157562255859375
@(rand()) → 0.607552015674623913099594574305228888988494873046875
```

<a name="function:rand_between"></a>

## rand_between()

A single random integer in the given inclusive range.


```objectivec
@(rand_between(1, 10)) → 5
@(rand_between(1, 10)) → 10
```

<a name="function:read_chars"></a>

## read_chars(text)

Converts `text` into something that can be read by IVR systems

ReadChars will split the numbers such as they are easier to understand. This includes
splitting in 3s or 4s if appropriate.


```objectivec
@(read_chars("1234")) → 1 2 3 4
@(read_chars("abc")) → a b c
@(read_chars("abcdef")) → a b c , d e f
```

<a name="function:remove_first_word"></a>

## remove_first_word(text)

Removes the 1st word of `text`


```objectivec
@(remove_first_word("foo bar")) → bar
```

<a name="function:repeat"></a>

## repeat(text, count)

Return `text` repeated `count` number of times


```objectivec
@(repeat("*", 8)) → ********
@(repeat("*", "foo")) → ERROR
```

<a name="function:replace"></a>

## replace(text, needle, replacement)

Replaces all occurrences of `needle` with `replacement` in `text`


```objectivec
@(replace("foo bar", "foo", "zap")) → zap bar
@(replace("foo bar", "baz", "zap")) → foo bar
```

<a name="function:right"></a>

## right(text, count)

Returns the `count` most right characters of the passed in `text`


```objectivec
@(right("hello", 2)) → lo
@(right("hello", 7)) → hello
@(right("😀😃😄😁", 2)) → 😄😁
@(right("hello", -1)) → ERROR
```

<a name="function:round"></a>

## round(num [,places])

Rounds `num` to the nearest value. You can optionally pass in the number of decimal places to round to as `places`.

If places < 0, it will round the integer part to the nearest 10^(-places).


```objectivec
@(round(12)) → 12
@(round(12.141)) → 12
@(round(12.6)) → 13
@(round(12.141, 2)) → 12.14
@(round(12.146, 2)) → 12.15
@(round(12.146, -1)) → 10
@(round("notnum", 2)) → ERROR
```

<a name="function:round_down"></a>

## round_down(num [,places])

Rounds `num` down to the nearest integer value. You can optionally pass in the number of decimal places to round to as `places`.


```objectivec
@(round_down(12)) → 12
@(round_down(12.141)) → 12
@(round_down(12.6)) → 12
@(round_down(12.141, 2)) → 12.14
@(round_down(12.146, 2)) → 12.14
@(round_down("foo")) → ERROR
```

<a name="function:round_up"></a>

## round_up(num [,places])

Rounds `num` up to the nearest integer value. You can optionally pass in the number of decimal places to round to as `places`.


```objectivec
@(round_up(12)) → 12
@(round_up(12.141)) → 13
@(round_up(12.6)) → 13
@(round_up(12.141, 2)) → 12.15
@(round_up(12.146, 2)) → 12.15
@(round_up("foo")) → ERROR
```

<a name="function:split"></a>

## split(text, delimiters)

Splits `text` based on the characters in `delimiters`

Empty values are removed from the returned list


```objectivec
@(split("a b c", " ")) → ["a","b","c"]
@(split("a", " ")) → ["a"]
@(split("abc..d", ".")) → ["abc","d"]
@(split("a.b.c.", ".")) → ["a","b","c"]
@(split("a|b,c  d", " .|,")) → ["a","b","c","d"]
```

<a name="function:text"></a>

## text(value)

Tries to convert `value` to text. An error is returned if the value can't be converted.


```objectivec
@(text(3 = 3)) → true
@(json(text(123.45))) → "123.45"
@(text(1 / 0)) → ERROR
```

<a name="function:text_compare"></a>

## text_compare(text1, text2)

Returns the comparison between the strings `text1` and `text2`.
The return value will be -1 if str1 is smaller than str2, 0 if they
are equal and 1 if str1 is greater than str2


```objectivec
@(text_compare("abc", "abc")) → 0
@(text_compare("abc", "def")) → -1
@(text_compare("zzz", "aaa")) → 1
```

<a name="function:title"></a>

## title(text)

Titlecases the passed in `text`, capitalizing each word


```objectivec
@(title("foo")) → Foo
@(title("ryan lewis")) → Ryan Lewis
@(title(123)) → 123
```

<a name="function:to_epoch"></a>

## to_epoch(date)

Converts `date` to the number of nanoseconds since January 1st, 1970 GMT


```objectivec
@(to_epoch("2017-06-12T16:56:59.000000Z")) → 1497286619000000000
```

<a name="function:today"></a>

## today()

Returns the current date in the current timezone, time is set to midnight in the environment timezone


```objectivec
@(today()) → 2018-04-11T00:00:00.000000-05:00
```

<a name="function:tz"></a>

## tz(date)

Returns the timezone for `date``

If not timezone information is present in the date, then the environment's
timezone will be returned


```objectivec
@(tz("2017-01-15T02:15:18.123456Z")) → UTC
@(tz("2017-01-15 02:15:18PM")) → America/Guayaquil
@(tz("2017-01-15")) → America/Guayaquil
@(tz("foo")) → ERROR
```

<a name="function:tz_offset"></a>

## tz_offset(date)

Returns the offset for the timezone as text +/- HHMM for `date`

If no timezone information is present in the date, then the environment's
timezone offset will be returned


```objectivec
@(tz_offset("2017-01-15T02:15:18.123456Z")) → +0000
@(tz_offset("2017-01-15 02:15:18PM")) → -0500
@(tz_offset("2017-01-15")) → -0500
@(tz_offset("foo")) → ERROR
```

<a name="function:upper"></a>

## upper(text)

Uppercases all characters in the passed `text`


```objectivec
@(upper("Asdf")) → ASDF
@(upper(123)) → 123
```

<a name="function:url_encode"></a>

## url_encode(text)

URL encodes `text` for use in a URL parameter


```objectivec
@(url_encode("two words")) → two+words
@(url_encode(10)) → 10
```

<a name="function:weekday"></a>

## weekday(date)

Returns the day of the week for `date`, 0 is sunday, 1 is monday..


```objectivec
@(weekday("2017-01-15")) → 0
@(weekday("foo")) → ERROR
```

<a name="function:word"></a>

## word(text, index [,delimiters])

Returns the word at the passed in `index` for the passed in `text`. There is an optional final
parameter `delimiters` which is string of characters used to split the text into words.


```objectivec
@(word("bee cat dog", 0)) → bee
@(word("bee.cat,dog", 0)) → bee
@(word("bee.cat,dog", 1)) → cat
@(word("bee.cat,dog", 2)) → dog
@(word("bee.cat,dog", -1)) → dog
@(word("bee.cat,dog", -2)) → cat
@(word("bee.*cat,dog", 1, ".*=|")) → cat,dog
@(word("O'Grady O'Flaggerty", 1, " ")) → O'Flaggerty
```

<a name="function:word_count"></a>

## word_count(text [,delimiters])

Returns the number of words in `text`. There is an optional final parameter `delimiters`
which is string of characters used to split the text into words.


```objectivec
@(word_count("foo bar")) → 2
@(word_count(10)) → 1
@(word_count("")) → 0
@(word_count("😀😃😄😁")) → 4
@(word_count("bee.*cat,dog", ".*=|")) → 2
@(word_count("O'Grady O'Flaggerty", " ")) → 2
```

<a name="function:word_slice"></a>

## word_slice(text, start, end [,delimiters])

Extracts a substring from `text` spanning from `start` up to but not-including `end`. (first word is 0). A negative
end value means that all words after the start should be returned. There is an optional final parameter `delimiters`
which is string of characters used to split the text into words.


```objectivec
@(word_slice("bee cat dog", 0, 1)) → bee
@(word_slice("bee cat dog", 0, 2)) → bee cat
@(word_slice("bee cat dog", 1, -1)) → cat dog
@(word_slice("bee cat dog", 1)) → cat dog
@(word_slice("bee cat dog", 2, 3)) → dog
@(word_slice("bee cat dog", 3, 10)) →
@(word_slice("bee.*cat,dog", 1, -1, ".*=|,")) → cat dog
@(word_slice("O'Grady O'Flaggerty", 1, 2, " ")) → O'Flaggerty
```


</div>

# Router Tests

Router tests are a special class of functions which are used within the switch router. They are called in the same way as normal functions, but 
all return a test result object which by default evalutes to true or false, but can also be used to find the matching portion of the test by using
the `match` component of the result. The flow editor builds these expressions using UI widgets, but they can be used anywhere a normal template
function is used.

<div class="tests">
<a name="test:has_all_words"></a>

## has_all_words(text, words)

Tests whether all the `words` are contained in `text`

The words can be in any order and may appear more than once.


```objectivec
@(has_all_words("the quick brown FOX", "the fox")) → true
@(has_all_words("the quick brown FOX", "the fox").match) → the FOX
@(has_all_words("the quick brown fox", "red fox")) → false
```

<a name="test:has_any_word"></a>

## has_any_word(text, words)

Tests whether any of the `words` are contained in the `text`

Only one of the words needs to match and it may appear more than once.


```objectivec
@(has_any_word("The Quick Brown Fox", "fox quick")) → true
@(has_any_word("The Quick Brown Fox", "red fox")) → true
@(has_any_word("The Quick Brown Fox", "red fox").match) → Fox
```

<a name="test:has_beginning"></a>

## has_beginning(text, beginning)

Tests whether `text` starts with `beginning`

Both text values are trimmed of surrounding whitespace, but otherwise matching is strict
without any tokenization.


```objectivec
@(has_beginning("The Quick Brown", "the quick")) → true
@(has_beginning("The Quick Brown", "the quick").match) → The Quick
@(has_beginning("The Quick Brown", "the   quick")) → false
@(has_beginning("The Quick Brown", "quick brown")) → false
```

<a name="test:has_date"></a>

## has_date(text)

Tests whether `text` contains a date formatted according to our environment


```objectivec
@(has_date("the date is 2017-01-15")) → true
@(has_date("the date is 2017-01-15").match) → 2017-01-15T00:00:00.000000-05:00
@(has_date("there is no date here, just a year 2017")) → false
```

<a name="test:has_date_eq"></a>

## has_date_eq(text, date)

Tests whether `text` a date equal to `date`


```objectivec
@(has_date_eq("the date is 2017-01-15", "2017-01-15")) → true
@(has_date_eq("the date is 2017-01-15", "2017-01-15").match) → 2017-01-15T00:00:00.000000-05:00
@(has_date_eq("the date is 2017-01-15 15:00", "2017-01-15")) → false
@(has_date_eq("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_eq("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_date_gt"></a>

## has_date_gt(text, min)

Tests whether `text` a date after the date `min`


```objectivec
@(has_date_gt("the date is 2017-01-15", "2017-01-01")) → true
@(has_date_gt("the date is 2017-01-15", "2017-01-01").match) → 2017-01-15T00:00:00.000000-05:00
@(has_date_gt("the date is 2017-01-15", "2017-03-15")) → false
@(has_date_gt("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_gt("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_date_lt"></a>

## has_date_lt(text, max)

Tests whether `text` contains a date before the date `max`


```objectivec
@(has_date_lt("the date is 2017-01-15", "2017-06-01")) → true
@(has_date_lt("the date is 2017-01-15", "2017-06-01").match) → 2017-01-15T00:00:00.000000-05:00
@(has_date_lt("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_lt("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_district"></a>

## has_district(text, state)

Tests whether a district name is contained in the `text`. If `state` is also provided
then the returned district must be within that state.


```objectivec
@(has_district("Gasabo", "Kigali")) → true
@(has_district("I live in Gasabo", "Kigali")) → true
@(has_district("I live in Gasabo", "Kigali").match) → Rwanda > Kigali City > Gasabo
@(has_district("Gasabo", "Boston")) → false
@(has_district("Gasabo")) → true
```

<a name="test:has_email"></a>

## has_email(text)

Tests whether an email is contained in `text`


```objectivec
@(has_email("my email is foo1@bar.com, please respond")) → true
@(has_email("my email is foo1@bar.com, please respond").match) → foo1@bar.com
@(has_email("my email is <foo@bar2.com>")) → true
@(has_email("i'm not sharing my email")) → false
```

<a name="test:has_group"></a>

## has_group(contact, group_uuid)

Returns whether the `contact` is part of group with the passed in UUID


```objectivec
@(has_group(contact, "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d")) → true
@(has_group(contact, "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d").match) → Testers
@(has_group(contact, "97fe7029-3a15-4005-b0c7-277b884fc1d5")) → false
```

<a name="test:has_number"></a>

## has_number(text)

Tests whether `text` contains a number


```objectivec
@(has_number("the number is 42")) → true
@(has_number("the number is 42").match) → 42
@(has_number("the number is forty two")) → false
```

<a name="test:has_number_between"></a>

## has_number_between(text, min, max)

Tests whether `text` contains a number between `min` and `max` inclusive


```objectivec
@(has_number_between("the number is 42", 40, 44)) → true
@(has_number_between("the number is 42", 40, 44).match) → 42
@(has_number_between("the number is 42", 50, 60)) → false
@(has_number_between("the number is not there", 50, 60)) → false
@(has_number_between("the number is not there", "foo", 60)) → ERROR
```

<a name="test:has_number_eq"></a>

## has_number_eq(text, value)

Tests whether `text` contains a number equal to the `value`


```objectivec
@(has_number_eq("the number is 42", 42)) → true
@(has_number_eq("the number is 42", 42).match) → 42
@(has_number_eq("the number is 42", 40)) → false
@(has_number_eq("the number is not there", 40)) → false
@(has_number_eq("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_gt"></a>

## has_number_gt(text, min)

Tests whether `text` contains a number greater than `min`


```objectivec
@(has_number_gt("the number is 42", 40)) → true
@(has_number_gt("the number is 42", 40).match) → 42
@(has_number_gt("the number is 42", 42)) → false
@(has_number_gt("the number is not there", 40)) → false
@(has_number_gt("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_gte"></a>

## has_number_gte(text, min)

Tests whether `text` contains a number greater than or equal to `min`


```objectivec
@(has_number_gte("the number is 42", 42)) → true
@(has_number_gte("the number is 42", 42).match) → 42
@(has_number_gte("the number is 42", 45)) → false
@(has_number_gte("the number is not there", 40)) → false
@(has_number_gte("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_lt"></a>

## has_number_lt(text, max)

Tests whether `text` contains a number less than `max`


```objectivec
@(has_number_lt("the number is 42", 44)) → true
@(has_number_lt("the number is 42", 44).match) → 42
@(has_number_lt("the number is 42", 40)) → false
@(has_number_lt("the number is not there", 40)) → false
@(has_number_lt("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_lte"></a>

## has_number_lte(text, max)

Tests whether `text` contains a number less than or equal to `max`


```objectivec
@(has_number_lte("the number is 42", 42)) → true
@(has_number_lte("the number is 42", 44).match) → 42
@(has_number_lte("the number is 42", 40)) → false
@(has_number_lte("the number is not there", 40)) → false
@(has_number_lte("the number is not there", "foo")) → ERROR
```

<a name="test:has_only_phrase"></a>

## has_only_phrase(text, phrase)

Tests whether the `text` contains only `phrase`

The phrase must be the only text in the text to match


```objectivec
@(has_only_phrase("The Quick Brown Fox", "quick brown")) → false
@(has_only_phrase("Quick Brown", "quick brown")) → true
@(has_only_phrase("the Quick Brown fox", "")) → false
@(has_only_phrase("", "")) → true
@(has_only_phrase("Quick Brown", "quick brown").match) → Quick Brown
@(has_only_phrase("The Quick Brown Fox", "red fox")) → false
```

<a name="test:has_pattern"></a>

## has_pattern(text, pattern)

Tests whether `text` matches the regex `pattern`

Both text values are trimmed of surrounding whitespace and matching is case-insensitive.


```objectivec
@(has_pattern("Sell cheese please", "buy (\w+)")) → false
@(has_pattern("Buy cheese please", "buy (\w+)")) → true
@(has_pattern("Buy cheese please", "buy (\w+)").match) → Buy cheese
@(has_pattern("Buy cheese please", "buy (\w+)").match.groups[0]) → Buy cheese
@(has_pattern("Buy cheese please", "buy (\w+)").match.groups[1]) → cheese
```

<a name="test:has_phone"></a>

## has_phone(text, country_code)

Tests whether a phone number (in the passed in `country_code`) is contained in the `text`


```objectivec
@(has_phone("my number is 2067799294", "US")) → true
@(has_phone("my number is 206 779 9294", "US").match) → +12067799294
@(has_phone("my number is none of your business", "US")) → false
```

<a name="test:has_phrase"></a>

## has_phrase(text, phrase)

Tests whether `phrase` is contained in `text`

The words in the test phrase must appear in the same order with no other words
in between.


```objectivec
@(has_phrase("the quick brown fox", "brown fox")) → true
@(has_phrase("the Quick Brown fox", "quick fox")) → false
@(has_phrase("the Quick Brown fox", "")) → true
@(has_phrase("the.quick.brown.fox", "the quick").match) → the quick
```

<a name="test:has_state"></a>

## has_state(text)

Tests whether a state name is contained in the `text`


```objectivec
@(has_state("Kigali")) → true
@(has_state("Boston")) → false
@(has_state("¡Kigali!")) → true
@(has_state("¡Kigali!").match) → Rwanda > Kigali City
@(has_state("I live in Kigali")) → true
```

<a name="test:has_text"></a>

## has_text(text)

Tests whether there the text has any characters in it


```objectivec
@(has_text("quick brown")) → true
@(has_text("quick brown").match) → quick brown
@(has_text("")) → false
@(has_text(" \n")) → false
@(has_text(123)) → true
```

<a name="test:has_value"></a>

## has_value(value)

Returns whether `value` is non-nil and not an error

Note that `contact.fields` and `run.results` are considered dynamic, so it is not an error
to try to retrieve a value from fields or results which don't exist, rather these return an empty
value.


```objectivec
@(has_value(datetime("foo"))) → false
@(has_value(not.existing)) → false
@(has_value(contact.fields.unset)) → false
@(has_value("")) → false
@(has_value("hello")) → true
```

<a name="test:has_wait_timed_out"></a>

## has_wait_timed_out(run)

Returns whether the last wait timed out.


```objectivec
@(has_wait_timed_out(run)) → false
```

<a name="test:has_ward"></a>

## has_ward(text, district, state)

Tests whether a ward name is contained in the `text`


```objectivec
@(has_ward("Gisozi", "Gasabo", "Kigali")) → true
@(has_ward("I live in Gisozi", "Gasabo", "Kigali")) → true
@(has_ward("I live in Gisozi", "Gasabo", "Kigali").match) → Rwanda > Kigali City > Gasabo > Gisozi
@(has_ward("Gisozi", "Gasabo", "Brooklyn")) → false
@(has_ward("Gisozi", "Brooklyn", "Kigali")) → false
@(has_ward("Brooklyn", "Gasabo", "Kigali")) → false
@(has_ward("Gasabo")) → false
@(has_ward("Gisozi")) → true
```

<a name="test:has_webhook_status"></a>

## has_webhook_status(webhook, status)

Tests whether the passed in `webhook` call has the passed in `status`. If there is no
webhook set, then "success" will still match.


```objectivec
@(has_webhook_status(NULL, "success")) → true
@(has_webhook_status(run.webhook, "success")) → true
@(has_webhook_status(run.webhook, "connection_error")) → false
@(has_webhook_status(run.webhook, "success").match) → {"results":[{"state":"WA"},{"state":"IN"}]}
@(has_webhook_status("abc", "success")) → ERROR
```

<a name="test:is_error"></a>

## is_error(value)

Returns whether `value` is an error

Note that `contact.fields` and `run.results` are considered dynamic, so it is not an error
to try to retrieve a value from fields or results which don't exist, rather these return an empty
value.


```objectivec
@(is_error(datetime("foo"))) → true
@(is_error(run.not.existing)) → true
@(is_error(contact.fields.unset)) → true
@(is_error("hello")) → false
```

<a name="test:is_text_eq"></a>

## is_text_eq(text1, text2)

Returns whether two text values are equal (case sensitive). In the case that they
are, it will return the text as the match.


```objectivec
@(is_text_eq("foo", "foo")) → true
@(is_text_eq("foo", "FOO")) → false
@(is_text_eq("foo", "bar")) → false
@(is_text_eq("foo", " foo ")) → false
@(is_text_eq(run.status, "completed")) → true
@(is_text_eq(run.webhook.status, "success")) → true
@(is_text_eq(run.webhook.status, "connection_error")) → false
```


</div>

# Actions

Actions on a node generate events which can then be ingested by the engine container. In some cases the actions cause an immediate action, such 
as calling a webhook, in others the engine container is responsible for taking the action based on the event that is output, such as sending 
messages or updating contact fields. In either case the internal state of the engine is always updated to represent the new state so that
flow execution is consistent. For example, while the engine itself does not have access to a contact store, it updates its internal 
representation of a contact's state based on action performed on a flow so that later references in the flow are correct.

<div class="actions">
<a name="action:add_contact_groups"></a>

## add_contact_groups

Can be used to add a contact to one or more groups. An [contact_groups_added](sessions.html#event:contact_groups_added) event will be created
for the groups which the contact has been added to.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "add_contact_groups",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "groups": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ]
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_groups_added",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "4f15f627-b1e2-4851-8dbf-00ecf5d03034",
    "groups": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ]
}
```
</div>
<a name="action:add_contact_urn"></a>

## add_contact_urn

Can be used to add a URN to the current contact. An [contact_urn_added](sessions.html#event:contact_urn_added) event
will be created when this action is encountered. If there is no contact then this
action will be ignored.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "add_contact_urn",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "scheme": "tel",
    "path": "@run.results.phone_number"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_urn_added",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "b504fe9e-d8a8-47fd-af9c-ff2f1faac4db",
    "urn": "tel:+12344563452"
}
```
</div>
<a name="action:add_input_labels"></a>

## add_input_labels

Can be used to add labels to the last user input on a flow. An [input_labels_added](sessions.html#event:input_labels_added) event
will be created with the labels added when this action is encountered. If there is
no user input at that point then this action will be ignored.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "add_input_labels",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "labels": [
        {
            "uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
            "name": "Spam"
        }
    ]
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "input_labels_added",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "f3cbd795-9bb3-4331-ba82-c15b24dd577f",
    "input_uuid": "9bf91c2b-ce58-4cef-aacc-281e03f69ab5",
    "labels": [
        {
            "uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
            "name": "Spam"
        }
    ]
}
```
</div>
<a name="action:call_resthook"></a>

## call_resthook

Can be used to call a resthook.

A [resthook_called](sessions.html#event:resthook_called) event will be created based on the results of the HTTP call
to each subscriber of the resthook.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "call_resthook",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "resthook": "new-registration"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "resthook_called",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "229bd432-dac7-4a3f-ba91-c48ad8c50e6b",
    "resthook": "new-registration",
    "payload": "{\n\t\"contact\": {\"uuid\": \"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f\", \"name\": \"Ryan Lewis\", \"urn\": \"tel:+12065551212\"},\n\t\"flow\": {\"name\":\"Registration\",\"revision\":123,\"uuid\":\"50c3706e-fedb-42c0-8eab-dda3335714b7\"},\n\t\"path\": [{\"arrived_on\":\"2018-04-11T18:24:30.123456Z\",\"exit_uuid\":\"37d8813f-1402-4ad2-9cc2-e9054a96525b\",\"node_uuid\":\"72a1f5df-49f9-45df-94c9-d86f7ea064e5\",\"uuid\":\"347b55be-7be1-4e68-aaa3-04d3fbce5f9a\"},{\"arrived_on\":\"2018-04-11T18:24:30.123456Z\",\"exit_uuid\":\"d898f9a4-f0fc-4ac4-a639-c98c602bb511\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"uuid\":\"da339edd-083b-48cb-bef6-3979f99a96f9\"},{\"arrived_on\":\"2018-04-11T18:24:30.123456Z\",\"exit_uuid\":\"\",\"node_uuid\":\"c0781400-737f-4940-9a6c-1ec1c3df0325\",\"uuid\":\"229bd432-dac7-4a3f-ba91-c48ad8c50e6b\"}],\n\t\"results\": {\"favorite_color\":{\"category\":\"Red\",\"category_localized\":\"Red\",\"created_on\":\"2018-04-11T18:24:30.123456Z\",\"input\":null,\"name\":\"Favorite Color\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"value\":\"red\"},\"phone_number\":{\"category\":\"\",\"category_localized\":\"\",\"created_on\":\"2018-04-11T18:24:30.123456Z\",\"input\":null,\"name\":\"Phone Number\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"value\":\"+12344563452\"}},\n\t\"run\": {\"uuid\": \"4c9abf31-d821-4e97-ba7e-53c2263e32f8\", \"created_on\": \"2018-04-11T18:24:30.123456Z\"},\n\t\"input\": {\"attachments\":[{\"content_type\":\"image/jpeg\",\"url\":\"http://s3.amazon.com/bucket/test.jpg\"},{\"content_type\":\"audio/mp3\",\"url\":\"http://s3.amazon.com/bucket/test.mp3\"}],\"channel\":{\"address\":\"+12345671111\",\"name\":\"My Android Phone\",\"uuid\":\"57f1078f-88aa-46f4-a59a-948a5739c03d\"},\"created_on\":\"2000-01-01T00:00:00.000000Z\",\"text\":\"Hi there\",\"type\":\"msg\",\"urn\":{\"display\":\"\",\"path\":\"+12065551212\",\"scheme\":\"tel\"},\"uuid\":\"9bf91c2b-ce58-4cef-aacc-281e03f69ab5\"},\n\t\"channel\": {\"address\":\"+12345671111\",\"name\":\"My Android Phone\",\"uuid\":\"57f1078f-88aa-46f4-a59a-948a5739c03d\"}\n}",
    "calls": [
        {
            "url": "http://127.0.0.1:49998/?cmd=success",
            "status": "success",
            "status_code": 200,
            "response": "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }"
        }
    ]
}
```
</div>
<a name="action:call_webhook"></a>

## call_webhook

Can be used to call an external service. The body, header and url fields may be
templates and will be evaluated at runtime.

A [webhook_called](sessions.html#event:webhook_called) event will be created based on the results of the HTTP call.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "call_webhook",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "method": "GET",
    "url": "http://localhost:49998/?cmd=success",
    "headers": {
        "Authorization": "Token AAFFZZHH"
    }
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "webhook_called",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "e68a851e-6328-426b-a8fd-1537ca860f97",
    "url": "http://localhost:49998/?cmd=success",
    "status": "success",
    "status_code": 200,
    "request": "GET /?cmd=success HTTP/1.1\r\nHost: localhost:49998\r\nUser-Agent: goflow-testing\r\nAuthorization: Token AAFFZZHH\r\nAccept-Encoding: gzip\r\n\r\n",
    "response": "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }"
}
```
</div>
<a name="action:remove_contact_groups"></a>

## remove_contact_groups

Can be used to remove a contact from one or more groups. A [contact_groups_removed](sessions.html#event:contact_groups_removed) event will be created
for the groups which the contact is removed from. Groups can either be explicitly provided or `all_groups` can be set to true to remove
the contact from all non-dynamic groups.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "remove_contact_groups",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "groups": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Registered Users"
        }
    ],
    "all_groups": false
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_groups_removed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "5fa51f39-76ea-421c-a71b-fe4af29b871a",
    "groups": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Testers"
        }
    ]
}
```
</div>
<a name="action:send_broadcast"></a>

## send_broadcast

Can be used to send a message to one or more contacts. It accepts a list of URNs, a list of groups
and a list of contacts.

The URNs and text fields may be templates. A [broadcast_created](sessions.html#event:broadcast_created) event will be created for each unique urn, contact and group
with the evaluated text.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "send_broadcast",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "text": "Hi @contact.name, are you ready to complete today's survey?",
    "attachments": null,
    "urns": [
        "tel:+12065551212"
    ]
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "broadcast_created",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "8e64b588-d46e-4016-a5ef-59cf4d9d7a5b",
    "translations": {
        "": {
            "text": "Hi Ryan Lewis, are you ready to complete today's survey?"
        }
    },
    "base_language": "",
    "urns": [
        "tel:+12065551212"
    ]
}
```
</div>
<a name="action:send_email"></a>

## send_email

Can be used to send an email to one or more recipients. The subject, body and addresses
can all contain expressions.

An [email_created](sessions.html#event:email_created) event will be created for each email address.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "send_email",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "addresses": [
        "@contact.urns.mailto.0"
    ],
    "subject": "Here is your activation token",
    "body": "Your activation token is @contact.fields.activation_token"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "email_created",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "08eba586-0bb1-47ab-8c15-15a7c0c5228d",
    "addresses": [
        "foo@bar.com"
    ],
    "subject": "Here is your activation token",
    "body": "Your activation token is AACC55"
}
```
</div>
<a name="action:send_msg"></a>

## send_msg

Can be used to reply to the current contact in a flow. The text field may contain templates.

A [msg_created](sessions.html#event:msg_created) event will be created with the evaluated text.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "send_msg",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "text": "Hi @contact.name, are you ready to complete today's survey?",
    "attachments": []
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "msg_created",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "c1f115c7-bcf3-44ef-88b2-5d345629f07f",
    "msg": {
        "uuid": "10c62052-7db1-49d1-b8ba-60d66db82e39",
        "urn": "tel:+12065551212",
        "channel": {
            "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
            "name": "My Android Phone"
        },
        "text": "Hi Ryan Lewis, are you ready to complete today's survey?"
    }
}
```
</div>
<a name="action:set_contact_channel"></a>

## set_contact_channel

Can be used to update the preferred channel of the current contact.

A [contact_channel_changed](sessions.html#event:contact_channel_changed) event will be created with the set channel.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "set_contact_channel",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "channel": {
        "uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
        "name": "FAcebook Channel"
    }
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_channel_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "c174a241-6057-41a3-874b-f17fb8365c22",
    "channel": {
        "uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
        "name": "FAcebook Channel"
    }
}
```
</div>
<a name="action:set_contact_field"></a>

## set_contact_field

Can be used to update a field value on the contact. The value is a localizable
template and white space is trimmed from the final value. An empty string clears the value.
A [contact_field_changed](sessions.html#event:contact_field_changed) event will be created with the corresponding value.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "set_contact_field",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "field": {
        "key": "gender",
        "name": "Gender"
    },
    "value": "Male"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_field_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "a08b46fc-f057-4e9a-9bd7-277a6a165264",
    "field": {
        "key": "gender",
        "name": "Gender"
    },
    "value": "Male"
}
```
</div>
<a name="action:set_contact_language"></a>

## set_contact_language

Can be used to update the name of the contact. The language is a localizable
template and white space is trimmed from the final value. An empty string clears the language.
A [contact_language_changed](sessions.html#event:contact_language_changed) event will be created with the corresponding value.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "set_contact_language",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "language": "eng"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_language_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "7ca3fc1e-e652-4f5c-979e-17606f578787",
    "language": "eng"
}
```
</div>
<a name="action:set_contact_name"></a>

## set_contact_name

Can be used to update the name of the contact. The name is a localizable
template and white space is trimmed from the final value. An empty string clears the name.
A [contact_name_changed](sessions.html#event:contact_name_changed) event will be created with the corresponding value.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "set_contact_name",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "name": "Bob Smith"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_name_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "fbce9f1c-ddff-45f4-8d46-86b76f70a6a6",
    "name": "Bob Smith"
}
```
</div>
<a name="action:set_contact_timezone"></a>

## set_contact_timezone

Can be used to update the timezone of the contact. The timezone is a localizable
template and white space is trimmed from the final value. An empty string clears the timezone.
A [contact_timezone_changed](sessions.html#event:contact_timezone_changed) event will be created with the corresponding value.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "set_contact_timezone",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "timezone": "Africa/Kigali"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_timezone_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "e4be9d25-b3ab-4a47-8704-ab259cb52a5d",
    "timezone": "Africa/Kigali"
}
```
</div>
<a name="action:set_run_result"></a>

## set_run_result

Can be used to save a result for a flow. The result will be available in the context
for the run as @run.results.[name]. The optional category can be used as a way of categorizing results,
this can be useful for reporting or analytics.

Both the value and category fields may be templates. A [run_result_changed](sessions.html#event:run_result_changed) event will be created with the
final values.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "set_run_result",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "name": "Gender",
    "value": "m",
    "category": "Male"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "run_result_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "bb7de8fc-d0b0-41a6-bdf0-950b64bbbc6d",
    "name": "Gender",
    "value": "m",
    "category": "Male",
    "node_uuid": "c0781400-737f-4940-9a6c-1ec1c3df0325"
}
```
</div>
<a name="action:start_flow"></a>

## start_flow

Can be used to start a contact down another flow. The current flow will pause until the subflow exits or expires.

A [flow_triggered](sessions.html#event:flow_triggered) event will be created to record that the flow was started.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "start_flow",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Collect Language"
    }
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "flow_triggered",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "dda50da0-8fc0-4f22-9c96-61ebc05df996",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Collect Language"
    },
    "parent_run_uuid": "a8ff08ef-6f27-44bd-9029-066bfcb36cf8"
}
```
</div>
<a name="action:start_session"></a>

## start_session

Can be used to trigger sessions for other contacts and groups. A [session_triggered](sessions.html#event:session_triggered) event
will be created and it's the responsibility of the caller to act on that by initiating a new session with the flow engine.

<div class="input_action"><h3>Action</h3>```json
{
    "type": "start_session",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "groups": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ],
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Registration"
    }
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "session_triggered",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "636bcfe8-1dd9-4bbd-a2a5-6b6ffeeada26",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Registration"
    },
    "groups": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ],
    "run": {
        "uuid": "e6e30b78-f9c1-462b-9418-6d3e4ae5a100",
        "flow": {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "Registration"
        },
        "contact": {
            "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
            "id": 1234567,
            "name": "Ryan Lewis",
            "language": "eng",
            "timezone": "America/Guayaquil",
            "created_on": "2018-06-20T11:40:30.123456789Z",
            "urns": [
                "tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d",
                "twitterid:54784326227#nyaruka",
                "mailto:foo@bar.com"
            ],
            "groups": [
                {
                    "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
                    "name": "Testers"
                },
                {
                    "uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9",
                    "name": "Males"
                }
            ],
            "fields": {
                "activation_token": {
                    "text": "AACC55"
                },
                "age": {
                    "text": "23",
                    "number": 23
                },
                "gender": {
                    "text": "Male"
                },
                "join_date": {
                    "text": "2017-12-02",
                    "datetime": "2017-12-02T00:00:00-02:00"
                }
            }
        },
        "status": "active",
        "results": {
            "favorite_color": {
                "name": "Favorite Color",
                "value": "red",
                "category": "Red",
                "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                "created_on": "2018-04-11T18:24:30.123456Z"
            },
            "phone_number": {
                "name": "Phone Number",
                "value": "+12344563452",
                "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                "created_on": "2018-04-11T18:24:30.123456Z"
            }
        }
    }
}
```
</div>

</div>
