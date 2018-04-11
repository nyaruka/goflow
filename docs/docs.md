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
@(to_json(contact.channel)) → {"uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d","name":"My Android Phone","address":"+12345671111"}
```

<a name="context:contact"></a>

## Contact

Represents a person who is interacting with the flow. It renders as the person's name
(or perferred URN if name isn't set) in a template, and has the following properties which can be accessed:

 * `uuid` the UUID of the contact
 * `name` the full name of the contact
 * `first_name` the first name of the contact
 * `language` the [ISO-639-3](http://www-01.sil.org/iso639-3/) language code of the contact
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
@contact.urns → ["tel:+12065551212","twitterid:54784326227#nyaruka","mailto:foo@bar.com"]
@contact.urns.0 → tel:+12065551212
@contact.urns.tel → ["tel:+12065551212"]
@contact.urns.mailto.0 → mailto:foo@bar.com
@contact.urn → (206) 555-1212
@contact.groups → ["Testers","Males"]
@contact.fields → {"activation_token":"AACC55","gender":"Male"}
@contact.fields.activation_token → AACC55
@contact.fields.gender → Male
```

<a name="context:flow"></a>

## Flow

Describes the ordered logic of actions and routers. It renders as its name in a template, and has the following
properties which can be accessed:

 * `uuid` the UUID of the flow
 * `name` the name of the flow

Examples:


```objectivec
@run.flow → Registration
@child.flow → Collect Language
@run.flow.uuid → 50c3706e-fedb-42c0-8eab-dda3335714b7
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
@(to_json(contact.groups.1)) → {"uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9","name":"Males"}
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
 * `attachments` any attachments on the message
 * `urn` the [URN](#context:urn) that the input was received on

Examples:


```objectivec
@run.input → Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3
@run.input.type → msg
@run.input.text → Hi there
@run.input.attachments → ["http://s3.amazon.com/bucket/test.jpg","http://s3.amazon.com/bucket/test.mp3"]
```

<a name="context:result"></a>

## Result

Describes a value captured during a run's execution. It might have been implicitly created by a router, or explicitly
created by a [set_run_result](#action:set_run_result) action.It renders as its value in a template, and has the following
properties which can be accessed:

 * `value` the value of the result
 * `category` the category of the result
 * `category_localized` the localized category of the result
 * `created_on` the time when the result was created

Examples:


```objectivec
@run.results.color → red
@run.results.color.value → red
@run.results.color.category → Red
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
@trigger.type → manual
@trigger.params → {"source": "website","address": {"state": "WA"}}
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
 * `channel` the preferred [channel](#context:channel) of the URN

To render a URN in a human friendly format, use the [format_urn](#function:format_urn) function.

Examples:


```objectivec
@contact.urns.0 → tel:+12065551212
@contact.urns.0.scheme → tel
@contact.urns.0.path → +12065551212
@contact.urns.1.display → nyaruka
@(format_urn(contact.urns.0)) → (206) 555-1212
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

<a name="function:char"></a>

## char(num)

Returns the rune for the passed in codepoint, `num`, which may be unicode, this is the reverse of code


```objectivec
@(char(33)) → !
@(char(128512)) → 😀
@(char("foo")) → ERROR
```

<a name="function:clean"></a>

## clean(string)

Strips any leading or trailing whitespace from `string``


```objectivec
@(clean("\nfoo\t")) → foo
@(clean(" bar")) → bar
@(clean(123)) → 123
```

<a name="function:code"></a>

## code(string)

Returns the numeric code for the first character in `string`, it is the inverse of char


```objectivec
@(code("a")) → 97
@(code("abc")) → 97
@(code("😀")) → 128512
@(code("15")) → 49
@(code(15)) → 49
@(code("")) → ERROR
```

<a name="function:date"></a>

## date(string)

Turns `string` into a date according to the environment's settings

date will return an error if it is unable to convert the string to a date.


```objectivec
@(date("1979-07-18")) → 1979-07-18T00:00:00.000000Z
@(date("2010 05 10")) → 2010-05-10T00:00:00.000000Z
@(date("NOT DATE")) → ERROR
```

<a name="function:date_add"></a>

## date_add(date, offset, unit)

Calculates the date value arrived at by adding `offset` number of `unit` to the `date`

Valid durations are "y" for years, "M" for months, "w" for weeks, "d" for days, h" for hour,
"m" for minutes, "s" for seconds


```objectivec
@(date_add("2017-01-15", 5, "d")) → 2017-01-20T00:00:00.000000Z
@(date_add("2017-01-15 10:45", 30, "m")) → 2017-01-15T11:15:00.000000Z
```

<a name="function:date_diff"></a>

## date_diff(date1, date2, unit)

Returns the integer duration between `date1` and `date2` in the `unit` specified.

Valid durations are "y" for years, "M" for months, "w" for weeks, "d" for days, h" for hour,
"m" for minutes, "s" for seconds


```objectivec
@(date_diff("2017-01-17", "2017-01-15", "d")) → 2
@(date_diff("2017-01-17 10:50", "2017-01-17 12:30", "h")) → -1
@(date_diff("2017-01-17", "2015-12-17", "y")) → 2
```

<a name="function:date_from_parts"></a>

## date_from_parts(year, month, day)

Converts the passed in `year`, `month` and `day`


```objectivec
@(date_from_parts(2017, 1, 15)) → 2017-01-15T00:00:00.000000Z
@(date_from_parts(2017, 2, 31)) → 2017-03-03T00:00:00.000000Z
@(date_from_parts(2017, 13, 15)) → ERROR
```

<a name="function:default"></a>

## default(test, default)

Takes two arguments, returning `test` if not an error or nil, otherwise returning `default`


```objectivec
@(default(undeclared.var, "default_value")) → default_value
@(default("10", "20")) → 10
@(default(date("invalid-date"), "today")) → today
```

<a name="function:field"></a>

## field(string, offset, delimeter)

Splits `string` based on the passed in `delimiter` and returns the field at `offset`.  When splitting
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

## format_date(date, format [,timezone])

Turns `date` into a string according to the `format` specified and in
the optional `timezone`.

The format string can consist of the following characters. The characters
' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.

* `YY`        - last two digits of year 0-99
* `YYYY`      - four digits of your 0000-9999
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
@(format_date("1979-07-18T15:00:00.000000Z")) → 1979-07-18 15:00
@(format_date("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) → 1979-07-18
@(format_date("2010-05-10T19:50:00.000000Z", "YYYY M DD tt:mm")) → 2010 5 10 19:50
@(format_date("2010-05-10T19:50:00.000000Z", "YYYY-MM-DD tt:mm AA", "America/Los_Angeles")) → 2010-05-10 12:50 PM
@(format_date("1979-07-18T15:00:00.000000Z", "YYYY")) → 1979
@(format_date("1979-07-18T15:00:00.000000Z", "M")) → 7
@(format_date("NOT DATE", "YYYY-MM-DD")) → ERROR
```

<a name="function:format_num"></a>

## format_num(num, places, commas)

Returns `num` formatted with the passed in number of decimal `places` and optional `commas` dividing thousands separators


```objectivec
@(format_num(31337, 2, true)) → 31,337.00
@(format_num(31337, 0, false)) → 31337
@(format_num("foo", 2, false)) → ERROR
```

<a name="function:format_urn"></a>

## format_urn(urn)

Turns `urn` into a human friendly string


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
@(from_epoch(1497286619000000000)) → 2017-06-12T16:56:59.000000Z
```

<a name="function:from_json"></a>

## from_json(string)

Tries to parse `string` as JSON, returning a fragment you can index into

If the passed in value is not JSON, then an error is returned


```objectivec
@(from_json("[1,2,3,4]").2) → 3
@(from_json("invalid json")) → ERROR
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

## join(array, delimeter)

Joins the passed in `array` of strings with the passed in `delimeter`


```objectivec
@(join(array("a", "b", "c"), "|")) → a|b|c
@(join(split("a.b.c", "."), " ")) → a b c
```

<a name="function:left"></a>

## left(string, count)

Returns the `count` most left characters of the passed in `string`


```objectivec
@(left("hello", 2)) → he
@(left("hello", 7)) → hello
@(left("😀😃😄😁", 2)) → 😀😃
@(left("hello", -1)) → ERROR
```

<a name="function:length"></a>

## length(object)

Returns the length of the passed in string or array.

length will return an error if it is passed an item which doesn't have length.


```objectivec
@(length("Hello")) → 5
@(length("😀😃😄😁")) → 4
@(length(array())) → 0
@(length(array("a", "b", "c"))) → 3
@(length(1234)) → ERROR
```

<a name="function:lower"></a>

## lower(string)

Lowercases the passed in `string`


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
@(now()) → 2018-04-11T13:24:30.123456Z
```

<a name="function:or"></a>

## or(tests...)

Returns whether if any of the passed in arguments are truthy


```objectivec
@(or(true)) → true
@(or(true, false, true)) → true
```

<a name="function:parse_date"></a>

## parse_date(string, format [,timezone])

Turns `string` into a date according to the `format` and optional `timezone` specified

The format string can consist of the following characters. The characters
' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.

* `YY`        - last two digits of year 0-99
* `YYYY`      - four digits of your 0000-9999
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

parse_date will return an error if it is unable to convert the string to a date.


```objectivec
@(parse_date("1979-07-18", "YYYY-MM-DD")) → 1979-07-18T00:00:00.000000Z
@(parse_date("2010 5 10", "YYYY M DD")) → 2010-05-10T00:00:00.000000Z
@(parse_date("2010 5 10 12:50", "YYYY M DD tt:mm", "America/Los_Angeles")) → 2010-05-10T12:50:00.000000-07:00
@(parse_date("NOT DATE", "YYYY-MM-DD")) → ERROR
```

<a name="function:percent"></a>

## percent(num)

Converts `num` to a string represented as a percentage


```objectivec
@(percent(0.54234)) → 54%
@(percent(1.2)) → 120%
@(percent("foo")) → ERROR
```

<a name="function:rand"></a>

## rand(floor, ceiling)

Returns either a single random decimal between 0-1 or a random integer between `floor` and `ceiling` (inclusive)


```objectivec
@(rand() > 0) → true
@(rand(1, 5) <= 5) → true
```

<a name="function:read_code"></a>

## read_code(code)

Converts `code` into something that can be read by IVR systems

ReadCode will split the numbers such as they are easier to understand. This includes
splitting in 3s or 4s if appropriate.


```objectivec
@(read_code("1234")) → 1 2 3 4
@(read_code("abc")) → a b c
@(read_code("abcdef")) → a b c , d e f
```

<a name="function:remove_first_word"></a>

## remove_first_word(string)

Removes the 1st word of `string`


```objectivec
@(remove_first_word("foo bar")) → bar
```

<a name="function:repeat"></a>

## repeat(string, count)

Return `string` repeated `count` number of times


```objectivec
@(repeat("*", 8)) → ********
@(repeat("*", "foo")) → ERROR
```

<a name="function:replace"></a>

## replace(string, needle, replacement)

Replaces all occurrences of `needle` with `replacement` in `string`


```objectivec
@(replace("foo bar", "foo", "zap")) → zap bar
@(replace("foo bar", "baz", "zap")) → foo bar
```

<a name="function:right"></a>

## right(string, count)

Returns the `count` most right characters of the passed in `string`


```objectivec
@(right("hello", 2)) → lo
@(right("hello", 7)) → hello
@(right("😀😃😄😁", 2)) → 😄😁
@(right("hello", -1)) → ERROR
```

<a name="function:round"></a>

## round(num [,places])

Rounds `num` to the nearest value. You can optionally pass
in the number of decimal places to round to as `places`.

If places < 0, it will round the integer part to the nearest 10^(-places).


```objectivec
@(round(12.141)) → 12
@(round(12.6)) → 13
@(round(12.141, 2)) → 12.14
@(round(12.146, 2)) → 12.15
@(round(12.146, -1)) → 10
@(round("notnum", 2)) → ERROR
```

<a name="function:round_down"></a>

## round_down(num)

Rounds `num` down to the nearest integer value


```objectivec
@(round_down(12.141)) → 12
@(round_down(12.9)) → 12
@(round_down("foo")) → ERROR
```

<a name="function:round_up"></a>

## round_up(num)

Rounds `num` up to the nearest integer value, also good at fighting weeds


```objectivec
@(round_up(12.141)) → 13
@(round_up(12)) → 12
@(round_up("foo")) → ERROR
```

<a name="function:split"></a>

## split(string, delimiter)

Splits `string` based on the passed in `delimeter`

Empty values are removed from the returned list


```objectivec
@(split("a b c", " ")) → ["a","b","c"]
@(split("a", " ")) → ["a"]
@(split("abc..d", ".")) → ["abc","d"]
@(split("a.b.c.", ".")) → ["a","b","c"]
@(split("a && b && c", " && ")) → ["a","b","c"]
```

<a name="function:string_cmp"></a>

## string_cmp(str1, str2)

Returns the comparison between the strings `str1` and `str2`.
The return value will be -1 if str1 is smaller than str2, 0 if they
are equal and 1 if str1 is greater than str2


```objectivec
@(string_cmp("abc", "abc")) → 0
@(string_cmp("abc", "def")) → -1
@(string_cmp("zzz", "aaa")) → 1
```

<a name="function:title"></a>

## title(string)

Titlecases the passed in `string`, capitalizing each word


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

<a name="function:to_json"></a>

## to_json(value)

Tries to return a JSON representation of `value`. An error is returned if there is
no JSON representation of that object.


```objectivec
@(to_json("string")) → "string"
@(to_json(10)) → 10
@(to_json(contact.uuid)) → "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
```

<a name="function:today"></a>

## today()

Returns the current date in the current timezone, time is set to midnight in the environment timezone


```objectivec
@(today()) → 2018-04-11T00:00:00.000000Z
```

<a name="function:tz"></a>

## tz(date)

Returns the timezone for `date``

If not timezone information is present in the date, then the environment's
timezone will be returned


```objectivec
@(tz("2017-01-15 02:15:18PM UTC")) → UTC
@(tz("2017-01-15 02:15:18PM")) → UTC
@(tz("2017-01-15")) → UTC
@(tz("foo")) → ERROR
```

<a name="function:tz_offset"></a>

## tz_offset(date)

Returns the offset for the timezone as a string +/- HHMM for `date`

If no timezone information is present in the date, then the environment's
timezone offset will be returned


```objectivec
@(tz_offset("2017-01-15 02:15:18PM UTC")) → +0000
@(tz_offset("2017-01-15 02:15:18PM")) → +0000
@(tz_offset("2017-01-15")) → +0000
@(tz_offset("foo")) → ERROR
```

<a name="function:upper"></a>

## upper(string)

Uppercases all characters in the passed `string`


```objectivec
@(upper("Asdf")) → ASDF
@(upper(123)) → 123
```

<a name="function:url_encode"></a>

## url_encode(string)

URL encodes `string` for use in a URL parameter


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

## word(string, offset)

Returns the word at the passed in `offset` for the passed in `string`


```objectivec
@(word("foo bar", 0)) → foo
@(word("foo.bar", 0)) → foo
@(word("one two.three", 2)) → three
```

<a name="function:word_count"></a>

## word_count(string)

Returns the number of words in `string`


```objectivec
@(word_count("foo bar")) → 2
@(word_count(10)) → 1
@(word_count("")) → 0
@(word_count("😀😃😄😁")) → 4
```

<a name="function:word_slice"></a>

## word_slice(string, start, end)

Extracts a substring from `string` spanning from `start` up to but not-including `end`. (first word is 1)


```objectivec
@(word_slice("foo bar", 1, 1)) → foo
@(word_slice("foo bar", 1, 3)) → foo bar
@(word_slice("foo bar", 3, 4)) →
```


</div>

# Router Tests

Router tests are a special class of functions which are used within the switch router. They are called in the same way as normal functions, but 
all return a test result object which by default evalutes to true or false, but can also be used to find the matching portion of the test by using
the `match` component of the result. The flow editor builds these expressions using UI widgets, but they can be used anywhere a normal template
function is used.

<div class="tests">
<a name="test:has_all_words"></a>

## has_all_words(string, words)

Tests whether all the `words` are contained in `string`

The words can be in any order and may appear more than once.


```objectivec
@(has_all_words("the quick brown FOX", "the fox")) → true
@(has_all_words("the quick brown FOX", "the fox").match) → the FOX
@(has_all_words("the quick brown fox", "red fox")) → false
```

<a name="test:has_any_word"></a>

## has_any_word(string, words)

Tests whether any of the `words` are contained in the `string`

Only one of the words needs to match and it may appear more than once.


```objectivec
@(has_any_word("The Quick Brown Fox", "fox quick")) → true
@(has_any_word("The Quick Brown Fox", "red fox")) → true
@(has_any_word("The Quick Brown Fox", "red fox").match) → Fox
```

<a name="test:has_beginning"></a>

## has_beginning(string, beginning)

Tests whether `string` starts with `beginning`

Both strings are trimmed of surrounding whitespace, but otherwise matching is strict
without any tokenization.


```objectivec
@(has_beginning("The Quick Brown", "the quick")) → true
@(has_beginning("The Quick Brown", "the quick").match) → The Quick
@(has_beginning("The Quick Brown", "the   quick")) → false
@(has_beginning("The Quick Brown", "quick brown")) → false
```

<a name="test:has_date"></a>

## has_date(string)

Tests whether `string` contains a date formatted according to our environment


```objectivec
@(has_date("the date is 2017-01-15")) → true
@(has_date("the date is 2017-01-15").match) → 2017-01-15T00:00:00.000000Z
@(has_date("there is no date here, just a year 2017")) → false
```

<a name="test:has_date_eq"></a>

## has_date_eq(string, date)

Tests whether `string` a date equal to `date`


```objectivec
@(has_date_eq("the date is 2017-01-15", "2017-01-15")) → true
@(has_date_eq("the date is 2017-01-15", "2017-01-15").match) → 2017-01-15T00:00:00.000000Z
@(has_date_eq("the date is 2017-01-15 15:00", "2017-01-15")) → false
@(has_date_eq("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_eq("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_date_gt"></a>

## has_date_gt(string, min)

Tests whether `string` a date after the date `min`


```objectivec
@(has_date_gt("the date is 2017-01-15", "2017-01-01")) → true
@(has_date_gt("the date is 2017-01-15", "2017-01-01").match) → 2017-01-15T00:00:00.000000Z
@(has_date_gt("the date is 2017-01-15", "2017-03-15")) → false
@(has_date_gt("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_gt("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_date_lt"></a>

## has_date_lt(string, max)

Tests whether `value` contains a date before the date `max`


```objectivec
@(has_date_lt("the date is 2017-01-15", "2017-06-01")) → true
@(has_date_lt("the date is 2017-01-15", "2017-06-01").match) → 2017-01-15T00:00:00.000000Z
@(has_date_lt("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_lt("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_district"></a>

## has_district(string, state)

Tests whether a district name is contained in the `string`. If `state` is also provided
then the returned district must be within that state.


```objectivec
@(has_district("Gasabo", "Kigali")) → true
@(has_district("I live in Gasabo", "Kigali")) → true
@(has_district("Gasabo", "Boston")) → false
@(has_district("Gasabo")) → true
```

<a name="test:has_email"></a>

## has_email(string)

Tests whether an email is contained in `string`


```objectivec
@(has_email("my email is foo1@bar.com, please respond")) → true
@(has_email("my email is foo1@bar.com, please respond").match) → foo1@bar.com
@(has_email("my email is <foo@bar2.com>")) → true
@(has_email("i'm not sharing my email")) → false
```

<a name="test:has_group"></a>

## has_group(contact)

Returns whether the `contact` is part of group with the passed in UUID


```objectivec
@(has_group(contact, "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d")) → true
@(has_group(contact, "97fe7029-3a15-4005-b0c7-277b884fc1d5")) → false
```

<a name="test:has_number"></a>

## has_number(string)

Tests whether `string` contains a number


```objectivec
@(has_number("the number is 42")) → true
@(has_number("the number is 42").match) → 42
@(has_number("the number is forty two")) → false
```

<a name="test:has_number_between"></a>

## has_number_between(string, min, max)

Tests whether `string` contains a number between `min` and `max` inclusive


```objectivec
@(has_number_between("the number is 42", 40, 44)) → true
@(has_number_between("the number is 42", 40, 44).match) → 42
@(has_number_between("the number is 42", 50, 60)) → false
@(has_number_between("the number is not there", 50, 60)) → false
@(has_number_between("the number is not there", "foo", 60)) → ERROR
```

<a name="test:has_number_eq"></a>

## has_number_eq(string, value)

Tests whether `strung` contains a number equal to the `value`


```objectivec
@(has_number_eq("the number is 42", 42)) → true
@(has_number_eq("the number is 42", 42).match) → 42
@(has_number_eq("the number is 42", 40)) → false
@(has_number_eq("the number is not there", 40)) → false
@(has_number_eq("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_gt"></a>

## has_number_gt(string, min)

Tests whether `string` contains a number greater than `min`


```objectivec
@(has_number_gt("the number is 42", 40)) → true
@(has_number_gt("the number is 42", 40).match) → 42
@(has_number_gt("the number is 42", 42)) → false
@(has_number_gt("the number is not there", 40)) → false
@(has_number_gt("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_gte"></a>

## has_number_gte(string, min)

Tests whether `string` contains a number greater than or equal to `min`


```objectivec
@(has_number_gte("the number is 42", 42)) → true
@(has_number_gte("the number is 42", 42).match) → 42
@(has_number_gte("the number is 42", 45)) → false
@(has_number_gte("the number is not there", 40)) → false
@(has_number_gte("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_lt"></a>

## has_number_lt(string, max)

Tests whether `string` contains a number less than `max`


```objectivec
@(has_number_lt("the number is 42", 44)) → true
@(has_number_lt("the number is 42", 44).match) → 42
@(has_number_lt("the number is 42", 40)) → false
@(has_number_lt("the number is not there", 40)) → false
@(has_number_lt("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_lte"></a>

## has_number_lte(string, max)

Tests whether `value` contains a number less than or equal to `max`


```objectivec
@(has_number_lte("the number is 42", 42)) → true
@(has_number_lte("the number is 42", 44).match) → 42
@(has_number_lte("the number is 42", 40)) → false
@(has_number_lte("the number is not there", 40)) → false
@(has_number_lte("the number is not there", "foo")) → ERROR
```

<a name="test:has_only_phrase"></a>

## has_only_phrase(string, phrase)

Tests whether the `string` contains only `phrase`

The phrase must be the only text in the string to match


```objectivec
@(has_only_phrase("The Quick Brown Fox", "quick brown")) → false
@(has_only_phrase("Quick Brown", "quick brown")) → true
@(has_only_phrase("the Quick Brown fox", "")) → false
@(has_only_phrase("", "")) → true
@(has_only_phrase("Quick Brown", "quick brown").match) → Quick Brown
@(has_only_phrase("The Quick Brown Fox", "red fox")) → false
```

<a name="test:has_pattern"></a>

## has_pattern(string, pattern)

Tests whether `string` matches the regex `pattern`

Both strings are trimmed of surrounding whitespace and matching is case-insensitive.


```objectivec
@(has_pattern("Sell cheese please", "buy (\w+)")) → false
@(has_pattern("Buy cheese please", "buy (\w+)")) → true
@(has_pattern("Buy cheese please", "buy (\w+)").match) → Buy cheese
@(has_pattern("Buy cheese please", "buy (\w+)").match.groups[0]) → Buy cheese
@(has_pattern("Buy cheese please", "buy (\w+)").match.groups[1]) → cheese
```

<a name="test:has_phone"></a>

## has_phone(string, country_code)

Tests whether a phone number (in the passed in `country_code`) is contained in the `string`


```objectivec
@(has_phone("my number is 2067799294", "US")) → true
@(has_phone("my number is 206 779 9294", "US").match) → +12067799294
@(has_phone("my number is none of your business", "US")) → false
```

<a name="test:has_phrase"></a>

## has_phrase(string, phrase)

Tests whether `phrase` is contained in `string`

The words in the test phrase must appear in the same order with no other words
in between.


```objectivec
@(has_phrase("the quick brown fox", "brown fox")) → true
@(has_phrase("the Quick Brown fox", "quick fox")) → false
@(has_phrase("the Quick Brown fox", "")) → true
@(has_phrase("the.quick.brown.fox", "the quick").match) → the quick
```

<a name="test:has_state"></a>

## has_state(string)

Tests whether a state name is contained in the `string`


```objectivec
@(has_state("Kigali")) → true
@(has_state("Boston")) → false
@(has_state("¡Kigali!")) → true
@(has_state("I live in Kigali")) → true
```

<a name="test:has_text"></a>

## has_text(string)

Tests whether there the string has any characters in it


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
@(has_value(date("foo"))) → false
@(has_value(not.existing)) → false
@(has_value(contact.fields.unset)) → false
@(has_value("hello")) → true
```

<a name="test:has_wait_timed_out"></a>

## has_wait_timed_out(run)

Returns whether the last wait timed out.


```objectivec
@(has_wait_timed_out(run)) → false
```

<a name="test:has_ward"></a>

## has_ward(string, district, state)

Tests whether a ward name is contained in the `string`


```objectivec
@(has_ward("Gisozi", "Gasabo", "Kigali")) → true
@(has_ward("I live in Gisozi", "Gasabo", "Kigali")) → true
@(has_ward("Gisozi", "Gasabo", "Brooklyn")) → false
@(has_ward("Gisozi", "Brooklyn", "Kigali")) → false
@(has_ward("Brooklyn", "Gasabo", "Kigali")) → false
@(has_ward("Gasabo")) → false
@(has_ward("Gisozi")) → true
```

<a name="test:is_error"></a>

## is_error(value)

Returns whether `value` is an error

Note that `contact.fields` and `run.results` are considered dynamic, so it is not an error
to try to retrieve a value from fields or results which don't exist, rather these return an empty
value.


```objectivec
@(is_error(date("foo"))) → true
@(is_error(run.not.existing)) → true
@(is_error(contact.fields.unset)) → true
@(is_error("hello")) → false
```

<a name="test:is_string_eq"></a>

## is_string_eq(string, string)

Returns whether two strings are equal (case sensitive). In the case that they
are, it will return the string as the match.


```objectivec
@(is_string_eq("foo", "foo")) → true
@(is_string_eq("foo", "FOO")) → false
@(is_string_eq("foo", "bar")) → false
@(is_string_eq("foo", " foo ")) → false
@(is_string_eq(run.status, "completed")) → true
@(is_string_eq(run.webhook.status, "success")) → true
@(is_string_eq(run.webhook.status, "connection_error")) → false
```


</div>

# Action Definitions

Actions on a node generate events which can then be ingested by the engine container. In some cases the actions cause an immediate action, such 
as calling a webhook, in others the engine container is responsible for taking the action based on the event that is output, such as sending 
messages or updating contact fields. In either case the internal state of the engine is always updated to represent the new state so that
flow execution is consistent. For example, while the engine itself does not have access to a contact store, it updates its internal 
representation of a contact's state based on action performed on a flow so that later references in the flow are correct.

<div class="actions">
<a name="action:add_contact_groups"></a>

## add_contact_groups

Can be used to add a contact to one or more groups. An `contact_groups_added` event will be created
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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "970b8069-50f5-4f6f-8f41-6b2d9f33d623",
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

Can be used to add a URN to the current contact. An `contact_urn_added` event
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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "b88ce93d-4360-4455-a691-235cbe720980",
    "urn": "tel:+12344563452"
}
```
</div>
<a name="action:add_input_labels"></a>

## add_input_labels

Can be used to add labels to the last user input on a flow. An `input_labels_added` event
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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "688e64f9-2456-4b42-afcb-91a2073e5459",
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
<a name="action:call_webhook"></a>

## call_webhook

Can be used to call an external service and insert the results in @run.webhook
context variable. The body, header and url fields may be templates and will be evaluated at runtime.

A `webhook_called` event will be created based on the results of the HTTP call.

<div class="input_action"><h3>Action</h3>```json
{
  "type": "call_webhook",
  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
  "method": "GET",
  "url": "http://localhost:49999/?cmd=success",
  "headers": {
    "Authorization": "Token AAFFZZHH"
  }
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "webhook_called",
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "b504fe9e-d8a8-47fd-af9c-ff2f1faac4db",
    "url": "http://localhost:49999/?cmd=success",
    "status": "success",
    "status_code": 200,
    "request": "GET /?cmd=success HTTP/1.1\r\nHost: localhost:49999\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Token AAFFZZHH\r\nAccept-Encoding: gzip\r\n\r\n",
    "response": "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }"
}
```
</div>
<a name="action:remove_contact_groups"></a>

## remove_contact_groups

Can be used to remove a contact from one or more groups. A `contact_groups_removed` event will be created
for the groups which the contact is removed from. If no groups are specified, then the contact will be removed from
all groups.

<div class="input_action"><h3>Action</h3>```json
{
  "type": "remove_contact_groups",
  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
  "groups": [
    {
      "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
      "name": "Registered Users"
    }
  ]
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_groups_removed",
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "658fd57d-f132-4ae4-8ab7-4a517a86045c",
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

The URNs and text fields may be templates. A `send_broadcast` event will be created for each unique urn, contact and group
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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "347b55be-7be1-4e68-aaa3-04d3fbce5f9a",
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

An `email_created` event will be created for each email address.

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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "229bd432-dac7-4a3f-ba91-c48ad8c50e6b",
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

A `broadcast_created` event will be created with the evaluated text.

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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "951242a1-5333-4221-8f9d-465efd6fbb5e",
    "msg": {
        "uuid": "644592ee-11ad-4bc4-9566-6fb2598c32d6",
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

A `contact_channel_changed` event will be created with the set channel.

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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "dc47e96a-392b-429b-92ca-6e1d7f550554",
    "channel": {
        "uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
        "name": "FAcebook Channel"
    }
}
```
</div>
<a name="action:set_contact_field"></a>

## set_contact_field

Can be used to save a value to a contact. The value can be a template and will
be evaluated during the flow. A `contact_field_changed` event will be created with the corresponding value.

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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "5865a06e-6fcc-4db9-bfd7-d22404241e07",
    "field": {
        "key": "gender",
        "name": "Gender"
    },
    "value": "Male"
}
```
</div>
<a name="action:set_contact_property"></a>

## set_contact_property

Can be used to update one of the built in fields for a contact of "name" or
"language". An `contact_property_changed` event will be created with the corresponding values.

<div class="input_action"><h3>Action</h3>```json
{
  "type": "set_contact_property",
  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
  "property": "language",
  "value": "eng"
}
```
</div><div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_property_changed",
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "19ebde80-3969-47d3-a09e-6806aab9f510",
    "property": "language",
    "value": "eng"
}
```
</div>
<a name="action:set_run_result"></a>

## set_run_result

Can be used to save a result for a flow. The result will be available in the context
for the run as @run.results.[name]. The optional category can be used as a way of categorizing results,
this can be useful for reporting or analytics.

Both the value and category fields may be templates. A `run_result_changed` event will be created with the
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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "edbc66c0-53a8-4b2a-998e-ae5bd773804a",
    "name": "Gender",
    "value": "m",
    "category": "Male",
    "node_uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5"
}
```
</div>
<a name="action:start_flow"></a>

## start_flow

Can be used to start a contact down another flow. The current flow will pause until the subflow exits or expires.

A `flow_entered` event will be created when the flow is started, a `flow_exited` event will be created upon the subflows exit.

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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "40c152ee-c9ed-46ff-9c02-6222e1badc14",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Collect Language"
    },
    "parent_run_uuid": "08eba586-0bb1-47ab-8c15-15a7c0c5228d"
}
```
</div>
<a name="action:start_session"></a>

## start_session

Can be used to trigger sessions for other contacts and groups

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
    "created_on": "2018-04-11T13:24:30.123456Z",
    "step_uuid": "bcfb7b96-7c87-48ba-ad03-b49f80627da4",
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
        "uuid": "e3895066-303a-4b1f-be22-6e6983962829",
        "flow_uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
        "contact": {
            "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
            "name": "Ryan Lewis",
            "language": "eng",
            "timezone": "UTC",
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
                "gender": {
                    "text": "Male"
                }
            }
        },
        "status": "active",
        "results": {
            "color": {
                "name": "Color",
                "value": "red",
                "category": "Red",
                "node_uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
                "input": "",
                "created_on": "2018-04-11T13:24:30.123456Z"
            },
            "phone_number": {
                "name": "Phone Number",
                "value": "+12344563452",
                "node_uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
                "input": "",
                "created_on": "2018-04-11T13:24:30.123456Z"
            }
        }
    }
}
```
</div>

</div>

# Event Definitions

Events are the output of a flow run and represent instructions to the engine container on what actions should be taken due to the flow execution.
All templates in events have been evaluated and can be used to create concrete messages, contact updates, emails etc by the container.

<div class="events">
<a name="event:broadcast_created"></a>

## broadcast_created

Events are created for outgoing broadcasts.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "broadcast_created",
    "created_on": "2006-01-02T15:04:05Z",
    "translations": {
        "eng": {
            "text": "hi, what's up",
            "quick_replies": [
                "All good",
                "Got 99 problems"
            ]
        },
        "spa": {
            "text": "Que pasa",
            "quick_replies": [
                "Todo bien",
                "Tengo 99 problemas"
            ]
        }
    },
    "base_language": "eng",
    "urns": [
        "tel:+12065551212"
    ],
    "contacts": [
        {
            "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
            "name": "Bob"
        }
    ]
}
```
</div>
<a name="event:contact_changed"></a>

## contact_changed

Events are created to set a contact on a session

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "contact": {
        "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
        "name": "Bob",
        "urns": [
            "tel:+11231234567"
        ]
    }
}
```
</div>
<a name="event:contact_channel_changed"></a>

## contact_channel_changed

Events are created when a contact's preferred channel is changed.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_channel_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "channel": {
        "uuid": "67a3ac69-e5e0-4ef0-8423-eddf71a71472",
        "name": "Twilio"
    }
}
```
</div>
<a name="event:contact_field_changed"></a>

## contact_field_changed

Events are created when a contact field is updated.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_field_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "field": {
        "key": "gender",
        "name": "Gender"
    },
    "value": "Male"
}
```
</div>
<a name="event:contact_groups_added"></a>

## contact_groups_added

Events will be created with the groups a contact was added to.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_groups_added",
    "created_on": "2006-01-02T15:04:05Z",
    "groups": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Reporters"
        }
    ]
}
```
</div>
<a name="event:contact_groups_removed"></a>

## contact_groups_removed

Events are created when a contact has been removed from one or more
groups.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_groups_removed",
    "created_on": "2006-01-02T15:04:05Z",
    "groups": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Reporters"
        }
    ]
}
```
</div>
<a name="event:contact_property_changed"></a>

## contact_property_changed

Events are created when a property of a contact has been changed

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_property_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "property": "language",
    "value": "eng"
}
```
</div>
<a name="event:contact_urn_added"></a>

## contact_urn_added

Events will be created with the URN that should be added to the current contact.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_urn_added",
    "created_on": "2006-01-02T15:04:05Z",
    "urn": "tel:+12345678900"
}
```
</div>
<a name="event:email_created"></a>

## email_created

Events are created for each recipient which should receive an email.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "email_created",
    "created_on": "2006-01-02T15:04:05Z",
    "addresses": [
        "foo@bar.com"
    ],
    "subject": "Your activation token",
    "body": "Your activation token is AAFFKKEE"
}
```
</div>
<a name="event:environment_changed"></a>

## environment_changed

Events are created to set the environment on a session

<div class="output_event"><h3>Event</h3>```json
{
    "type": "environment_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "environment": {
        "date_format": "yyyy-MM-dd",
        "time_format": "hh:mm",
        "timezone": "Africa/Kigali",
        "languages": [
            "eng",
            "fra"
        ]
    }
}
```
</div>
<a name="event:error"></a>

## error

Events will be created whenever an error is encountered during flow execution. This
can vary from template evaluation errors to invalid actions.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "error",
    "created_on": "2006-01-02T15:04:05Z",
    "text": "invalid date format: '12th of October'",
    "fatal": false
}
```
</div>
<a name="event:flow_triggered"></a>

## flow_triggered

Events are created when an action wants to start a subflow

<div class="output_event"><h3>Event</h3>```json
{
    "type": "flow_triggered",
    "created_on": "2006-01-02T15:04:05Z",
    "flow": {
        "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
        "name": "Registration"
    },
    "parent_run_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd"
}
```
</div>
<a name="event:input_labels_added"></a>

## input_labels_added

Events will be created with the labels that were applied to the given input.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "input_labels_added",
    "created_on": "2006-01-02T15:04:05Z",
    "input_uuid": "4aef4050-1895-4c80-999a-70368317a4f5",
    "labels": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Spam"
        }
    ]
}
```
</div>
<a name="event:msg_created"></a>

## msg_created

Events are used for replies to the session contact.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "msg_created",
    "created_on": "2006-01-02T15:04:05Z",
    "msg": {
        "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
        "urn": "tel:+12065551212",
        "channel": {
            "uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf",
            "name": "Twilio"
        },
        "text": "hi there",
        "attachments": [
            "https://s3.amazon.com/mybucket/attachment.jpg"
        ]
    }
}
```
</div>
<a name="event:msg_received"></a>

## msg_received

Events are used for starting flows or resuming flows which are waiting for a message.
They represent an incoming message for a contact.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "msg_received",
    "created_on": "2006-01-02T15:04:05Z",
    "msg": {
        "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
        "urn": "tel:+12065551212",
        "channel": {
            "uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf",
            "name": "Twilio"
        },
        "text": "hi there",
        "attachments": [
            "https://s3.amazon.com/mybucket/attachment.jpg"
        ]
    }
}
```
</div>
<a name="event:msg_wait"></a>

## msg_wait

Events are created when a flow pauses waiting for a response from
a contact. If a timeout is set, then the caller should resume the flow after
the number of seconds in the timeout to resume it.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "msg_wait",
    "created_on": "2006-01-02T15:04:05Z"
}
```
</div>
<a name="event:nothing_wait"></a>

## nothing_wait

Events are created when a flow requests to hand back control to the caller but isn't
waiting for anything from the caller.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "nothing_wait",
    "created_on": "2006-01-02T15:04:05.234532Z"
}
```
</div>
<a name="event:run_expired"></a>

## run_expired

Events are sent by the caller to notify the engine that a run has expired

<div class="output_event"><h3>Event</h3>```json
{
    "type": "run_expired",
    "created_on": "2006-01-02T15:04:05Z",
    "run_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a"
}
```
</div>
<a name="event:run_result_changed"></a>

## run_result_changed

Events are created when a result is saved. They contain not only
the name, value and category of the result, but also the UUID of the node where
the result was generated.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "run_result_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "name": "Gender",
    "value": "m",
    "category": "Male",
    "category_localized": "Homme",
    "node_uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
}
```
</div>
<a name="event:session_triggered"></a>

## session_triggered

Events are created when an action wants to start a subflow

<div class="output_event"><h3>Event</h3>```json
{
    "type": "session_triggered",
    "created_on": "2006-01-02T15:04:05Z",
    "flow": {
        "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
        "name": "Registration"
    },
    "groups": [
        {
            "uuid": "8f8e2cae-3c8d-4dce-9c4b-19514437e427",
            "name": "New contacts"
        }
    ],
    "run": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "flow_uuid": "93c554a1-b90d-4892-b029-a2a87dec9b87",
        "contact": {
            "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
            "name": "Bob",
            "fields": {
                "state": {
                    "value": "Azuay",
                    "created_on": "2000-01-01T00:00:00.000000000-00:00"
                }
            }
        },
        "results": {
            "age": {
                "result_name": "Age",
                "value": "33",
                "node": "cd2be8c4-59bc-453c-8777-dec9a80043b8",
                "created_on": "2000-01-01T00:00:00.000000000-00:00"
            }
        }
    }
}
```
</div>
<a name="event:webhook_called"></a>

## webhook_called

Events are created when a webhook is called. The event contains
the status and status code of the response, as well as a full dump of the
request and response.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "webhook_called",
    "created_on": "2006-01-02T15:04:05Z",
    "url": "https://api.ipify.org?format=json",
    "status": "success",
    "status_code": 200,
    "request": "GET https://api.ipify.org?format=json",
    "response": "HTTP/1.1 200 OK {\"ip\":\"190.154.48.130\"}"
}
```
</div>

</div>

</body>
</html>


