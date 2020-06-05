If a node has a single exit, the engine will pick that when leaving that node. If the node has more than one exit,
then we need a router to choose an exit. 

# Routers

Routers are primarily responsible for picking exits but can also generate events and save results. All routers have 
the following properties:

 * `type` router type
 * `wait` optional [wait](#waits)
 * `result_name` optional result name if router should save a result
 * `categories` possible categories of any result saved by this router

Different router types have different logic for how an exit will be chosen.

## Switch

If a node wishes to route differently based on some state in the session, it can add a `switch` router which defines one or more 
`cases`.  Each case defines a `type` which is the name of an expression function that is run by passing the evaluation of `operand` 
as the first argument. Cases may define additional arguments using the `arguments` array on a case. If no case evaluates 
to true, then the `default_category_uuid` will be used, otherwise flow execution will stop.

A `switch` router has these additional properties:

 * `operand` the template which will be evaluated against each of our cases
 * `cases` a list of 1-n cases which are evaluated in order until one is true
 * `default_category_uuid` the uuid of the default category to take if no case matches (optional)

Each case consists of:

 * `uuid` the UUID
 * `type` the type of this test which is the name of a [test function](#tests) - it will be called with the operand as the first argument
 * `arguments` an optional list of templates which can be passed as extra arguments to the test (after the initial operand)
 * `category_uuid` the uuid of the category that should be taken if this case evaluated to true

The following is an example switch router with 2 cases:

```json
{
    "uuid": "ee0bee3f-34b3-4275-af78-f9ff52c82e6a",
    "router": {
        "type": "switch",
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

## Random

A random router chooses one of its categories randomly and has no additional properties. For example:

```json
{
    "uuid": "ee0bee3f-34b3-4275-af78-f9ff52c82e6a",
    "router": {
        "type": "random",
        "categories": [
            {
                "uuid": "cab600f5-b54b-49b9-a7ea-5638f4cbf2b4",
                "name": "Bucket 1",
                "exit_uuid": "972fb580-54c2-4491-8438-09ace3500ba5"
            },
            {
                "uuid": "9574fbfd-510f-4dfc-b989-97d2aecf50b9",
                "name": "Bucket 2",
                "exit_uuid": "6981b1a9-af04-4e26-a248-1fc1f5e5c7eb"
            }
        ]
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

A wait tells the engine to hand back control to the caller and wait for the caller to resume execution by providing something.
The type of the wait indicates what is required to resume flow execution and currently we only support waits of type `msg`.

## Msg

This type indicates that flow execution should pause until an incoming message is received. It can have an optional timeout 
value which is the number of seconds after which execution can be resumed without a message, e.g.

```json
{
    "type": "msg",
    "timeout": 600
}
```

# Tests

Router tests are a special class of functions which are used within the switch router. They are called in the same way as normal functions, but 
all return a test result object which by default evalutes to true or false, but can also be used to find the matching portion of the test by using
the `match` component of the result. The flow editor builds these expressions using UI widgets, but they can be used anywhere a normal template
function is used.

<div class="tests">
<h2 class="item_title"><a name="test:has_all_words" href="#test:has_all_words">has_all_words(text, words)</a></h2>

Tests whether all the `words` are contained in `text`

The words can be in any order and may appear more than once.


```objectivec
@(has_all_words("the quick brown FOX", "the fox")) → true
@(has_all_words("the quick brown FOX", "the fox").match) → the FOX
@(has_all_words("the quick brown fox", "red fox")) → false
```

<h2 class="item_title"><a name="test:has_any_word" href="#test:has_any_word">has_any_word(text, words)</a></h2>

Tests whether any of the `words` are contained in the `text`

Only one of the words needs to match and it may appear more than once.


```objectivec
@(has_any_word("The Quick Brown Fox", "fox quick")) → true
@(has_any_word("The Quick Brown Fox", "fox quick").match) → Quick Fox
@(has_any_word("The Quick Brown Fox", "red fox").match) → Fox
```

<h2 class="item_title"><a name="test:has_beginning" href="#test:has_beginning">has_beginning(text, beginning)</a></h2>

Tests whether `text` starts with `beginning`

Both text values are trimmed of surrounding whitespace, but otherwise matching is strict
without any tokenization.


```objectivec
@(has_beginning("The Quick Brown", "the quick")) → true
@(has_beginning("The Quick Brown", "the quick").match) → The Quick
@(has_beginning("The Quick Brown", "the   quick")) → false
@(has_beginning("The Quick Brown", "quick brown")) → false
```

<h2 class="item_title"><a name="test:has_category" href="#test:has_category">has_category(result, categories...)</a></h2>

Tests whether the category of a result on of the passed in `categories`


```objectivec
@(has_category(results.webhook, "Success", "Failure")) → true
@(has_category(results.webhook, "Success", "Failure").match) → Success
@(has_category(results.webhook, "Failure")) → false
```

<h2 class="item_title"><a name="test:has_date" href="#test:has_date">has_date(text)</a></h2>

Tests whether `text` contains a date formatted according to our environment


```objectivec
@(has_date("the date is 15/01/2017")) → true
@(has_date("the date is 15/01/2017").match) → 2017-01-15T13:24:30.123456-05:00
@(has_date("there is no date here, just a year 2017")) → false
```

<h2 class="item_title"><a name="test:has_date_eq" href="#test:has_date_eq">has_date_eq(text, date)</a></h2>

Tests whether `text` a date equal to `date`


```objectivec
@(has_date_eq("the date is 15/01/2017", "2017-01-15")) → true
@(has_date_eq("the date is 15/01/2017", "2017-01-15").match) → 2017-01-15T13:24:30.123456-05:00
@(has_date_eq("the date is 15/01/2017 15:00", "2017-01-15").match) → 2017-01-15T15:00:00.000000-05:00
@(has_date_eq("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_eq("there is no date here, just a year 2017", "not date")) → ERROR
```

<h2 class="item_title"><a name="test:has_date_gt" href="#test:has_date_gt">has_date_gt(text, min)</a></h2>

Tests whether `text` a date after the date `min`


```objectivec
@(has_date_gt("the date is 15/01/2017", "2017-01-01")) → true
@(has_date_gt("the date is 15/01/2017", "2017-01-01").match) → 2017-01-15T13:24:30.123456-05:00
@(has_date_gt("the date is 15/01/2017", "2017-03-15")) → false
@(has_date_gt("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_gt("there is no date here, just a year 2017", "not date")) → ERROR
```

<h2 class="item_title"><a name="test:has_date_lt" href="#test:has_date_lt">has_date_lt(text, max)</a></h2>

Tests whether `text` contains a date before the date `max`


```objectivec
@(has_date_lt("the date is 15/01/2017", "2017-06-01")) → true
@(has_date_lt("the date is 15/01/2017", "2017-06-01").match) → 2017-01-15T13:24:30.123456-05:00
@(has_date_lt("there is no date here, just a year 2017", "2017-06-01")) → false
@(has_date_lt("there is no date here, just a year 2017", "not date")) → ERROR
```

<h2 class="item_title"><a name="test:has_district" href="#test:has_district">has_district(text, state)</a></h2>

Tests whether a district name is contained in the `text`. If `state` is also provided
then the returned district must be within that state.


```objectivec
@(has_district("Gasabo", "Kigali").match) → Rwanda > Kigali City > Gasabo
@(has_district("I live in Gasabo", "Kigali").match) → Rwanda > Kigali City > Gasabo
@(has_district("Gasabo", "Boston")) → false
@(has_district("Gasabo").match) → Rwanda > Kigali City > Gasabo
```

<h2 class="item_title"><a name="test:has_email" href="#test:has_email">has_email(text)</a></h2>

Tests whether an email is contained in `text`


```objectivec
@(has_email("my email is foo1@bar.com, please respond")) → true
@(has_email("my email is foo1@bar.com, please respond").match) → foo1@bar.com
@(has_email("my email is <foo@bar2.com>").match) → foo@bar2.com
@(has_email("i'm not sharing my email")) → false
```

<h2 class="item_title"><a name="test:has_error" href="#test:has_error">has_error(value)</a></h2>

Returns whether `value` is an error


```objectivec
@(has_error(datetime("foo"))) → true
@(has_error(datetime("foo")).match) → error calling DATETIME: unable to convert "foo" to a datetime
@(has_error(run.not.existing).match) → object has no property 'not'
@(has_error(contact.fields.unset).match) → object has no property 'unset'
@(has_error("hello")) → false
```

<h2 class="item_title"><a name="test:has_group" href="#test:has_group">has_group(contact, group_uuid)</a></h2>

Returns whether the `contact` is part of group with the passed in UUID


```objectivec
@(has_group(contact.groups, "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d").match) → {name: Testers, uuid: b7cf0d83-f1c9-411c-96fd-c511a4cfa86d}
@(has_group(array(), "97fe7029-3a15-4005-b0c7-277b884fc1d5")) → false
```

<h2 class="item_title"><a name="test:has_intent" href="#test:has_intent">has_intent(result, name, confidence)</a></h2>

Tests whether any intent in a classification result has `name` and minimum `confidence`


```objectivec
@(has_intent(results.intent, "book_flight", 0.5)) → true
@(has_intent(results.intent, "book_hotel", 0.2)) → true
```

<h2 class="item_title"><a name="test:has_number" href="#test:has_number">has_number(text)</a></h2>

Tests whether `text` contains a number


```objectivec
@(has_number("the number is 42")) → true
@(has_number("the number is 42").match) → 42
@(has_number("the number is forty two")) → false
```

<h2 class="item_title"><a name="test:has_number_between" href="#test:has_number_between">has_number_between(text, min, max)</a></h2>

Tests whether `text` contains a number between `min` and `max` inclusive


```objectivec
@(has_number_between("the number is 42", 40, 44)) → true
@(has_number_between("the number is 42", 40, 44).match) → 42
@(has_number_between("the number is 42", 50, 60)) → false
@(has_number_between("the number is not there", 50, 60)) → false
@(has_number_between("the number is not there", "foo", 60)) → ERROR
```

<h2 class="item_title"><a name="test:has_number_eq" href="#test:has_number_eq">has_number_eq(text, value)</a></h2>

Tests whether `text` contains a number equal to the `value`


```objectivec
@(has_number_eq("the number is 42", 42)) → true
@(has_number_eq("the number is 42", 42).match) → 42
@(has_number_eq("the number is 42", 40)) → false
@(has_number_eq("the number is not there", 40)) → false
@(has_number_eq("the number is not there", "foo")) → ERROR
```

<h2 class="item_title"><a name="test:has_number_gt" href="#test:has_number_gt">has_number_gt(text, min)</a></h2>

Tests whether `text` contains a number greater than `min`


```objectivec
@(has_number_gt("the number is 42", 40)) → true
@(has_number_gt("the number is 42", 40).match) → 42
@(has_number_gt("the number is 42", 42)) → false
@(has_number_gt("the number is not there", 40)) → false
@(has_number_gt("the number is not there", "foo")) → ERROR
```

<h2 class="item_title"><a name="test:has_number_gte" href="#test:has_number_gte">has_number_gte(text, min)</a></h2>

Tests whether `text` contains a number greater than or equal to `min`


```objectivec
@(has_number_gte("the number is 42", 42)) → true
@(has_number_gte("the number is 42", 42).match) → 42
@(has_number_gte("the number is 42", 45)) → false
@(has_number_gte("the number is not there", 40)) → false
@(has_number_gte("the number is not there", "foo")) → ERROR
```

<h2 class="item_title"><a name="test:has_number_lt" href="#test:has_number_lt">has_number_lt(text, max)</a></h2>

Tests whether `text` contains a number less than `max`


```objectivec
@(has_number_lt("the number is 42", 44)) → true
@(has_number_lt("the number is 42", 44).match) → 42
@(has_number_lt("the number is 42", 40)) → false
@(has_number_lt("the number is not there", 40)) → false
@(has_number_lt("the number is not there", "foo")) → ERROR
```

<h2 class="item_title"><a name="test:has_number_lte" href="#test:has_number_lte">has_number_lte(text, max)</a></h2>

Tests whether `text` contains a number less than or equal to `max`


```objectivec
@(has_number_lte("the number is 42", 42)) → true
@(has_number_lte("the number is 42", 42).match) → 42
@(has_number_lte("the number is 42", 40)) → false
@(has_number_lte("the number is not there", 40)) → false
@(has_number_lte("the number is not there", "foo")) → ERROR
```

<h2 class="item_title"><a name="test:has_only_phrase" href="#test:has_only_phrase">has_only_phrase(text, phrase)</a></h2>

Tests whether the `text` contains only `phrase`

The phrase must be the only text in the text to match


```objectivec
@(has_only_phrase("Quick Brown", "quick brown")) → true
@(has_only_phrase("Quick Brown", "quick brown").match) → Quick Brown
@(has_only_phrase("The Quick Brown Fox", "quick brown")) → false
@(has_only_phrase("the Quick Brown fox", "")) → false
@(has_only_phrase("", "").match) →
@(has_only_phrase("The Quick Brown Fox", "red fox")) → false
```

<h2 class="item_title"><a name="test:has_only_text" href="#test:has_only_text">has_only_text(text1, text2)</a></h2>

Returns whether two text values are equal (case sensitive). In the case that they
are, it will return the text as the match.


```objectivec
@(has_only_text("foo", "foo")) → true
@(has_only_text("foo", "foo").match) → foo
@(has_only_text("foo", "FOO")) → false
@(has_only_text("foo", "bar")) → false
@(has_only_text("foo", " foo ")) → false
@(has_only_text(run.status, "completed").match) → completed
@(has_only_text(results.webhook.category, "Success").match) → Success
@(has_only_text(results.webhook.category, "Failure")) → false
```

<h2 class="item_title"><a name="test:has_pattern" href="#test:has_pattern">has_pattern(text, pattern)</a></h2>

Tests whether `text` matches the regex `pattern`

Both text values are trimmed of surrounding whitespace and matching is case-insensitive.


```objectivec
@(has_pattern("Buy cheese please", "buy (\w+)")) → true
@(has_pattern("Buy cheese please", "buy (\w+)").match) → Buy cheese
@(has_pattern("Buy cheese please", "buy (\w+)").extra) → {0: Buy cheese, 1: cheese}
@(has_pattern("Sell cheese please", "buy (\w+)")) → false
```

<h2 class="item_title"><a name="test:has_phone" href="#test:has_phone">has_phone(text, country_code)</a></h2>

Tests whether `text` contains a phone number. The optional `country_code` argument specifies
the country to use for parsing.


```objectivec
@(has_phone("my number is +12067799294 thanks")) → true
@(has_phone("my number is +12067799294").match) → +12067799294
@(has_phone("my number is 2067799294", "US").match) → +12067799294
@(has_phone("my number is 206 779 9294", "US").match) → +12067799294
@(has_phone("my number is none of your business", "US")) → false
```

<h2 class="item_title"><a name="test:has_phrase" href="#test:has_phrase">has_phrase(text, phrase)</a></h2>

Tests whether `phrase` is contained in `text`

The words in the test phrase must appear in the same order with no other words
in between.


```objectivec
@(has_phrase("the quick brown fox", "brown fox")) → true
@(has_phrase("the quick brown fox", "brown fox").match) → brown fox
@(has_phrase("the Quick Brown fox", "quick fox")) → false
@(has_phrase("the Quick Brown fox", "").match) →
```

<h2 class="item_title"><a name="test:has_state" href="#test:has_state">has_state(text)</a></h2>

Tests whether a state name is contained in the `text`


```objectivec
@(has_state("Kigali").match) → Rwanda > Kigali City
@(has_state("¡Kigali!").match) → Rwanda > Kigali City
@(has_state("I live in Kigali").match) → Rwanda > Kigali City
@(has_state("Boston")) → false
```

<h2 class="item_title"><a name="test:has_text" href="#test:has_text">has_text(text)</a></h2>

Tests whether there the text has any characters in it


```objectivec
@(has_text("quick brown")) → true
@(has_text("quick brown").match) → quick brown
@(has_text("")) → false
@(has_text(" \n")) → false
@(has_text(123).match) → 123
@(has_text(contact.fields.not_set)) → false
```

<h2 class="item_title"><a name="test:has_time" href="#test:has_time">has_time(text)</a></h2>

Tests whether `text` contains a time.


```objectivec
@(has_time("the time is 10:30")) → true
@(has_time("the time is 10:30").match) → 10:30:00.000000
@(has_time("the time is 10 PM").match) → 22:00:00.000000
@(has_time("the time is 10:30:45").match) → 10:30:45.000000
@(has_time("there is no time here, just the number 25")) → false
```

<h2 class="item_title"><a name="test:has_top_intent" href="#test:has_top_intent">has_top_intent(result, name, confidence)</a></h2>

Tests whether the top intent in a classification result has `name` and minimum `confidence`


```objectivec
@(has_top_intent(results.intent, "book_flight", 0.5)) → true
@(has_top_intent(results.intent, "book_hotel", 0.5)) → false
```

<h2 class="item_title"><a name="test:has_ward" href="#test:has_ward">has_ward(text, district, state)</a></h2>

Tests whether a ward name is contained in the `text`


```objectivec
@(has_ward("Gisozi", "Gasabo", "Kigali").match) → Rwanda > Kigali City > Gasabo > Gisozi
@(has_ward("I live in Gisozi", "Gasabo", "Kigali").match) → Rwanda > Kigali City > Gasabo > Gisozi
@(has_ward("Gisozi", "Gasabo", "Brooklyn")) → false
@(has_ward("Gisozi", "Brooklyn", "Kigali")) → false
@(has_ward("Brooklyn", "Gasabo", "Kigali")) → false
@(has_ward("Gasabo")) → false
@(has_ward("Gisozi").match) → Rwanda > Kigali City > Gasabo > Gisozi
```


</div>