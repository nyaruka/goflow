# Templates

Some properties of entities in the flow specification are _templates_ - that is their values are dynamic and are evaluated at runtime. 
Templates can contain single variables or more complex expressions. A single variable is embedded using the `@` character. For example 
the template `Hi @contact.name` contains a single variable which at runtime will be replaced with the name of the current contact.

More complex expressions can be embedded using the `@(...)` syntax. For example the template `Hi @("Dr " & upper(contact.name))` takes 
the contact name, converts it to uppercase, and the prefixes it with another string.

The `@` symbol can be escaped in templates by repeating it, ie, `Hi @@twitter` would output `Hi @twitter`.

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
@contact.channel → {address: +12345671111, name: My Android Phone, uuid: 57f1078f-88aa-46f4-a59a-948a5739c03d}
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
@run.flow → {name: Registration, revision: 123, uuid: 50c3706e-fedb-42c0-8eab-dda3335714b7}
@child.flow.name → Collect Age
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
@input → {attachments: [image/jpeg:http://s3.amazon.com/bucket/test.jpg, audio/mp3:http://s3.amazon.com/bucket/test.mp3], channel: {address: +12345671111, name: My Android Phone, uuid: 57f1078f-88aa-46f4-a59a-948a5739c03d}, created_on: 2017-12-31T11:35:10.035757-02:00, text: Hi there, type: msg, urn: tel:+12065551212, uuid: 9bf91c2b-ce58-4cef-aacc-281e03f69ab5}
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
@results.favorite_color → {category: Red, category_localized: Red, created_on: 2018-04-11T18:24:30.123456Z, input: , name: Favorite Color, node_uuid: f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03, value: red}
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
@(has_all_words("the quick brown FOX", "the fox")) → {match: the FOX}
@(has_all_words("the quick brown fox", "red fox")) →
```

<a name="test:has_any_word"></a>

## has_any_word(text, words)

Tests whether any of the `words` are contained in the `text`

Only one of the words needs to match and it may appear more than once.


```objectivec
@(has_any_word("The Quick Brown Fox", "fox quick")) → {match: Quick Fox}
@(has_any_word("The Quick Brown Fox", "red fox")) → {match: Fox}
```

<a name="test:has_beginning"></a>

## has_beginning(text, beginning)

Tests whether `text` starts with `beginning`

Both text values are trimmed of surrounding whitespace, but otherwise matching is strict
without any tokenization.


```objectivec
@(has_beginning("The Quick Brown", "the quick")) → {match: The Quick}
@(has_beginning("The Quick Brown", "the   quick")) →
@(has_beginning("The Quick Brown", "quick brown")) →
```

<a name="test:has_date"></a>

## has_date(text)

Tests whether `text` contains a date formatted according to our environment


```objectivec
@(has_date("the date is 15/01/2017")) → {match: 2017-01-15T13:24:30.123456-05:00}
@(has_date("there is no date here, just a year 2017")) →
```

<a name="test:has_date_eq"></a>

## has_date_eq(text, date)

Tests whether `text` a date equal to `date`


```objectivec
@(has_date_eq("the date is 15/01/2017", "2017-01-15")) → {match: 2017-01-15T13:24:30.123456-05:00}
@(has_date_eq("the date is 15/01/2017 15:00", "2017-01-15")) → {match: 2017-01-15T15:00:00.000000-05:00}
@(has_date_eq("there is no date here, just a year 2017", "2017-06-01")) →
@(has_date_eq("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_date_gt"></a>

## has_date_gt(text, min)

Tests whether `text` a date after the date `min`


```objectivec
@(has_date_gt("the date is 15/01/2017", "2017-01-01")) → {match: 2017-01-15T13:24:30.123456-05:00}
@(has_date_gt("the date is 15/01/2017", "2017-03-15")) →
@(has_date_gt("there is no date here, just a year 2017", "2017-06-01")) →
@(has_date_gt("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_date_lt"></a>

## has_date_lt(text, max)

Tests whether `text` contains a date before the date `max`


```objectivec
@(has_date_lt("the date is 15/01/2017", "2017-06-01")) → {match: 2017-01-15T13:24:30.123456-05:00}
@(has_date_lt("there is no date here, just a year 2017", "2017-06-01")) →
@(has_date_lt("there is no date here, just a year 2017", "not date")) → ERROR
```

<a name="test:has_district"></a>

## has_district(text, state)

Tests whether a district name is contained in the `text`. If `state` is also provided
then the returned district must be within that state.


```objectivec
@(has_district("Gasabo", "Kigali")) → {match: Rwanda > Kigali City > Gasabo}
@(has_district("I live in Gasabo", "Kigali")) → {match: Rwanda > Kigali City > Gasabo}
@(has_district("Gasabo", "Boston")) →
@(has_district("Gasabo")) → {match: Rwanda > Kigali City > Gasabo}
```

<a name="test:has_email"></a>

## has_email(text)

Tests whether an email is contained in `text`


```objectivec
@(has_email("my email is foo1@bar.com, please respond")) → {match: foo1@bar.com}
@(has_email("my email is <foo@bar2.com>")) → {match: foo@bar2.com}
@(has_email("i'm not sharing my email")) →
```

<a name="test:has_error"></a>

## has_error(value)

Returns whether `value` is an error


```objectivec
@(has_error(datetime("foo"))) → {match: error calling DATETIME: unable to convert "foo" to a datetime}
@(has_error(run.not.existing)) → {match: dict has no property 'not'}
@(has_error(contact.fields.unset)) → {match: dict has no property 'unset'}
@(has_error("hello")) →
```

<a name="test:has_group"></a>

## has_group(contact, group_uuid)

Returns whether the `contact` is part of group with the passed in UUID


```objectivec
@(has_group(contact.groups, "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d")) → {match: {name: Testers, uuid: b7cf0d83-f1c9-411c-96fd-c511a4cfa86d}}
@(has_group(array(), "97fe7029-3a15-4005-b0c7-277b884fc1d5")) →
```

<a name="test:has_number"></a>

## has_number(text)

Tests whether `text` contains a number


```objectivec
@(has_number("the number is 42")) → {match: 42}
@(has_number("the number is forty two")) →
```

<a name="test:has_number_between"></a>

## has_number_between(text, min, max)

Tests whether `text` contains a number between `min` and `max` inclusive


```objectivec
@(has_number_between("the number is 42", 40, 44)) → {match: 42}
@(has_number_between("the number is 42", 50, 60)) →
@(has_number_between("the number is not there", 50, 60)) →
@(has_number_between("the number is not there", "foo", 60)) → ERROR
```

<a name="test:has_number_eq"></a>

## has_number_eq(text, value)

Tests whether `text` contains a number equal to the `value`


```objectivec
@(has_number_eq("the number is 42", 42)) → {match: 42}
@(has_number_eq("the number is 42", 40)) →
@(has_number_eq("the number is not there", 40)) →
@(has_number_eq("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_gt"></a>

## has_number_gt(text, min)

Tests whether `text` contains a number greater than `min`


```objectivec
@(has_number_gt("the number is 42", 40)) → {match: 42}
@(has_number_gt("the number is 42", 42)) →
@(has_number_gt("the number is not there", 40)) →
@(has_number_gt("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_gte"></a>

## has_number_gte(text, min)

Tests whether `text` contains a number greater than or equal to `min`


```objectivec
@(has_number_gte("the number is 42", 42)) → {match: 42}
@(has_number_gte("the number is 42", 45)) →
@(has_number_gte("the number is not there", 40)) →
@(has_number_gte("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_lt"></a>

## has_number_lt(text, max)

Tests whether `text` contains a number less than `max`


```objectivec
@(has_number_lt("the number is 42", 44)) → {match: 42}
@(has_number_lt("the number is 42", 40)) →
@(has_number_lt("the number is not there", 40)) →
@(has_number_lt("the number is not there", "foo")) → ERROR
```

<a name="test:has_number_lte"></a>

## has_number_lte(text, max)

Tests whether `text` contains a number less than or equal to `max`


```objectivec
@(has_number_lte("the number is 42", 42)) → {match: 42}
@(has_number_lte("the number is 42", 40)) →
@(has_number_lte("the number is not there", 40)) →
@(has_number_lte("the number is not there", "foo")) → ERROR
```

<a name="test:has_only_phrase"></a>

## has_only_phrase(text, phrase)

Tests whether the `text` contains only `phrase`

The phrase must be the only text in the text to match


```objectivec
@(has_only_phrase("The Quick Brown Fox", "quick brown")) →
@(has_only_phrase("Quick Brown", "quick brown")) → {match: Quick Brown}
@(has_only_phrase("the Quick Brown fox", "")) →
@(has_only_phrase("", "")) → {match: }
@(has_only_phrase("The Quick Brown Fox", "red fox")) →
```

<a name="test:has_pattern"></a>

## has_pattern(text, pattern)

Tests whether `text` matches the regex `pattern`

Both text values are trimmed of surrounding whitespace and matching is case-insensitive.


```objectivec
@(has_pattern("Sell cheese please", "buy (\w+)")) →
@(has_pattern("Buy cheese please", "buy (\w+)")) → {extra: {0: Buy cheese, 1: cheese}, match: Buy cheese}
```

<a name="test:has_phone"></a>

## has_phone(text, country_code)

Tests whether `text` contains a phone number. The optional `country_code` argument specifies
the country to use for parsing.


```objectivec
@(has_phone("my number is +12067799294")) → {match: +12067799294}
@(has_phone("my number is 2067799294", "US")) → {match: +12067799294}
@(has_phone("my number is 206 779 9294", "US")) → {match: +12067799294}
@(has_phone("my number is none of your business", "US")) →
```

<a name="test:has_phrase"></a>

## has_phrase(text, phrase)

Tests whether `phrase` is contained in `text`

The words in the test phrase must appear in the same order with no other words
in between.


```objectivec
@(has_phrase("the quick brown fox", "brown fox")) → {match: brown fox}
@(has_phrase("the Quick Brown fox", "quick fox")) →
@(has_phrase("the Quick Brown fox", "")) → {match: }
```

<a name="test:has_state"></a>

## has_state(text)

Tests whether a state name is contained in the `text`


```objectivec
@(has_state("Kigali")) → {match: Rwanda > Kigali City}
@(has_state("¡Kigali!")) → {match: Rwanda > Kigali City}
@(has_state("I live in Kigali")) → {match: Rwanda > Kigali City}
@(has_state("Boston")) →
```

<a name="test:has_text"></a>

## has_text(text)

Tests whether there the text has any characters in it


```objectivec
@(has_text("quick brown")) → {match: quick brown}
@(has_text("")) →
@(has_text(" \n")) →
@(has_text(123)) → {match: 123}
@(has_text(contact.fields.not_set)) →
```

<a name="test:has_time"></a>

## has_time(text)

Tests whether `text` contains a time.


```objectivec
@(has_time("the time is 10:30")) → {match: 10:30:00.000000}
@(has_time("the time is 10 PM")) → {match: 22:00:00.000000}
@(has_time("the time is 10:30:45")) → {match: 10:30:45.000000}
@(has_time("there is no time here, just the number 25")) →
```

<a name="test:has_value"></a>

## has_value(value)

Returns whether `value` is non-nil and not an error

Note that `contact.fields` and `run.results` are considered dynamic, so it is not an error
to try to retrieve a value from fields or results which don't exist, rather these return an empty
value.


```objectivec
@(has_value(datetime("foo"))) →
@(has_value(not.existing)) →
@(has_value(contact.fields.unset)) →
@(has_value("")) →
@(has_value("hello")) → {match: hello}
```

<a name="test:has_ward"></a>

## has_ward(text, district, state)

Tests whether a ward name is contained in the `text`


```objectivec
@(has_ward("Gisozi", "Gasabo", "Kigali")) → {match: Rwanda > Kigali City > Gasabo > Gisozi}
@(has_ward("I live in Gisozi", "Gasabo", "Kigali")) → {match: Rwanda > Kigali City > Gasabo > Gisozi}
@(has_ward("Gisozi", "Gasabo", "Brooklyn")) →
@(has_ward("Gisozi", "Brooklyn", "Kigali")) →
@(has_ward("Brooklyn", "Gasabo", "Kigali")) →
@(has_ward("Gasabo")) →
@(has_ward("Gisozi")) → {match: Rwanda > Kigali City > Gasabo > Gisozi}
```

<a name="test:is_text_eq"></a>

## is_text_eq(text1, text2)

Returns whether two text values are equal (case sensitive). In the case that they
are, it will return the text as the match.


```objectivec
@(is_text_eq("foo", "foo")) → {match: foo}
@(is_text_eq("foo", "FOO")) →
@(is_text_eq("foo", "bar")) →
@(is_text_eq("foo", " foo ")) →
@(is_text_eq(run.status, "completed")) → {match: completed}
@(is_text_eq(results.webhook.category, "Success")) → {match: Success}
@(is_text_eq(results.webhook.category, "Failure")) →
```


</div>