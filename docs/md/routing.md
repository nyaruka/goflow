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
@(has_error(run.not.existing)) → {match: object has no property 'not'}
@(has_error(contact.fields.unset)) → {match: object has no property 'unset'}
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