# Context

The context is all the variables which are accessible in expressions and contains the following top-level variables:

 * `run` the current [run](#context:run)
 * `parent` the parent of the current [run](#context:run), i.e. the run that started the current run
 * `child` the child of the current [run](#context:run), i.e. the last subflow
 * `contact` the current [contact](#context:contact), shortcut for `@run.contact`
 * `results` the current [results](#context:result), shortcut for `@run.results`
 * `trigger` the [trigger](#context:trigger) that initiated this session
 * `input` the last [input](#context:input) from the contact

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
@input.channel.uuid → 57f1078f-88aa-46f4-a59a-948a5739c03d
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
 * `channel` shorthand for `contact.urns[0].channel`, i.e. the [channel](#context:channel) of the contact's preferred URN

Examples:


```objectivec
@contact.name → Ryan Lewis
@contact.first_name → Ryan
@contact.language → eng
@contact.timezone → America/Guayaquil
@contact.created_on → 2018-06-20T11:40:30.123456Z
@contact.urns → [tel:+12065551212, twitterid:54784326227#nyaruka, mailto:foo@bar.com]
@(contact.urns[0]) → tel:+12065551212
@contact.urn → tel:+12065551212
@(foreach(contact.groups, extract, "name")) → [Testers, Males]
@contact.fields → {activation_token: AACC55, age: 23, gender: Male, join_date: 2017-12-02T00:00:00.000000-02:00, not_set: }
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
@child.run.flow.name → Collect Age
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
@(foreach(contact.groups, extract, "name")) → [Testers, Males]
@(contact.groups[0].uuid) → b7cf0d83-f1c9-411c-96fd-c511a4cfa86d
@(contact.groups[1].name) → Males
@(json(contact.groups[1])) → {"name":"Males","uuid":"4f1f98fc-27a7-4a69-bbdb-24744ba739a9"}
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
@input → Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3
@input.type → msg
@input.text → Hi there
@input.attachments → [image/jpeg:http://s3.amazon.com/bucket/test.jpg, audio/mp3:http://s3.amazon.com/bucket/test.mp3]
@(json(input)) → {"attachments":["image/jpeg:http://s3.amazon.com/bucket/test.jpg","audio/mp3:http://s3.amazon.com/bucket/test.mp3"],"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2017-12-31T11:35:10.035757-02:00","text":"Hi there","type":"msg","urn":"tel:+12065551212","uuid":"9bf91c2b-ce58-4cef-aacc-281e03f69ab5"}
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
@results → 2Factor: 34634624463525\nFavorite Color: red\nPhone Number: +12344563452\nwebhook: 200
@results.favorite_color → red
@results.favorite_color.value → red
@results.favorite_color.category → Red
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
 * `results.[snaked_result_name]` the value of the specific result, e.g. `results.age`

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
@trigger.params → {address: {state: WA}, source: website}
@(json(trigger)) → {"params":{"address":{"state":"WA"},"source":"website"},"type":"flow_action"}
```

<a name="context:urn"></a>

## Urn

Represents a destination for an outgoing message or a source of an incoming message. It is string composed of 3
components: scheme, path, and display (optional). For example:

 - _tel:+16303524567_
 - _twitterid:54784326227#nyaruka_
 - _telegram:34642632786#bobby_

To render a URN in a human friendly format, use the [format_urn](expressions.html#function:format_urn) function.

Examples:


```objectivec
@(urns.tel) → tel:+12065551212
@(urn_parts(urns.tel).scheme) → tel
@(format_urn(urns.tel)) → (206) 555-1212
@(json(contact.urns[0])) → "tel:+12065551212"
```


</div>

