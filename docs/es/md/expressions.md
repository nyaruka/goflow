# Overview

Excellent gets its name from borrowing some of the syntax and function names of formulas in Microsoft Excel™, 
though it has evolved over time and similarities are now much fewer. It is an expression based templating 
language which aims to make it easy to generate text from a context of values.

# Templates

Templates can contain single variables or more complex expressions. A single variable is embedded using the `@` 
character. For example the template `Hi @foo` contains a single variable which at runtime will be replaced with 
the value of `foo` in the context.

More complex expressions can be embedded using the `@(...)` syntax. For example the template `Hi @("Dr " & upper(foo))` 
takes the value of `foo`, converts it to uppercase, and the prefixes it with another string. Note than within a 
complex expression you don't prefix variables with `@`.

The `@` symbol can be escaped in templates by repeating it, e.g, `Hi @@twitter` will output `Hi @twitter`.

# Types

Excellent has the following types:

 * [Array](#type:array)
 * [Boolean](#type:boolean)
 * [Date](#type:date)
 * [DateTime](#type:datetime)
 * [Function](#type:function)
 * [Number](#type:number)
 * [Object](#type:object)
 * [Text](#type:text)
 * [Time](#type:time)

<div class="types">
<h2 class="item_title"><a name="type:array" href="#type:array">array</a></h2>

Is an array of items.


```objectivec
@(array(1, "x", true)) → [1, x, true]
@(array(1, "x", true)[1]) → x
@(count(array(1, "x", true))) → 3
@(json(array(1, "x", true))) → [1,"x",true]
```

<h2 class="item_title"><a name="type:boolean" href="#type:boolean">boolean</a></h2>

Is a boolean `true` or `false`.


```objectivec
@(true) → true
@(1 = 1) → true
@(1 = 2) → false
@(json(true)) → true
```

<h2 class="item_title"><a name="type:date" href="#type:date">date</a></h2>

Is a Gregorian calendar date value.


```objectivec
@(date_from_parts(2019, 4, 11)) → 2019-04-11
@(format_date(date_from_parts(2019, 4, 11))) → 11-04-2019
@(json(date_from_parts(2019, 4, 11))) → "2019-04-11"
```

<h2 class="item_title"><a name="type:datetime" href="#type:datetime">datetime</a></h2>

Is a datetime value.


```objectivec
@(datetime("1979-07-18T10:30:45.123456Z")) → 1979-07-18T10:30:45.123456Z
@(format_datetime(datetime("1979-07-18T10:30:45.123456Z"))) → 18-07-1979 05:30
@(json(datetime("1979-07-18T10:30:45.123456Z"))) → "1979-07-18T10:30:45.123456Z"
```

<h2 class="item_title"><a name="type:function" href="#type:function">function</a></h2>

Is a callable function.


```objectivec
@(upper) → function
@(array(upper)[0]("abc")) → ABC
@(json(upper)) → null
```

<h2 class="item_title"><a name="type:number" href="#type:number">number</a></h2>

Is a whole or fractional number.


```objectivec
@(1234) → 1234
@(1234.5678) → 1234.5678
@(format_number(1234.5670)) → 1,234.567
@(json(1234.5678)) → 1234.5678
```

<h2 class="item_title"><a name="type:object" href="#type:object">object</a></h2>

Is an object with named properties.


```objectivec
@(object("foo", 1, "bar", "x")) → {bar: x, foo: 1}
@(object("foo", 1, "bar", "x").bar) → x
@(object("foo", 1, "bar", "x")["bar"]) → x
@(count(object("foo", 1, "bar", "x"))) → 2
@(json(object("foo", 1, "bar", "x"))) → {"bar":"x","foo":1}
```

<h2 class="item_title"><a name="type:text" href="#type:text">text</a></h2>

Is a string of characters.


```objectivec
@("abc") → abc
@(text_length("abc")) → 3
@(upper("abc")) → ABC
@(json("abc")) → "abc"
```

<h2 class="item_title"><a name="type:time" href="#type:time">time</a></h2>

Is a time of day.


```objectivec
@(time_from_parts(16, 30, 45)) → 16:30:45.000000
@(format_time(time_from_parts(16, 30, 45))) → 16:30
@(json(time_from_parts(16, 30, 45))) → "16:30:45.000000"
```


</div>

# Operators

<div class="operators">
<h2 class="item_title"><a name="operator:add" href="#operator:add">+</a></h2>

Adds two numbers.


```objectivec
@(2 + 3) → 5
@(fields.age + 10) → 33
```

<h2 class="item_title"><a name="operator:concatenate" href="#operator:concatenate">&</a></h2>

Joins two text values together.


```objectivec
@("hello" & " " & "bar") → hello bar
@("hello" & null) → hello
```

<h2 class="item_title"><a name="operator:divide" href="#operator:divide">/</a></h2>

Divides a number by another.


```objectivec
@(4 / 2) → 2
@(3 / 2) → 1.5
@(46 / fields.age) → 2
@(3 / 0) → ERROR
```

<h2 class="item_title"><a name="operator:equal" href="#operator:equal">=</a></h2>

Returns true if two values are textually equal.


```objectivec
@("hello" = "hello") → true
@("hello" = "bar") → false
@(1 = 1) → true
```

<h2 class="item_title"><a name="operator:exponent" href="#operator:exponent">^</a></h2>

Raises a number to the power of a another number.


```objectivec
@(2 ^ 8) → 256
```

<h2 class="item_title"><a name="operator:greaterthan" href="#operator:greaterthan">></a></h2>

Returns true if the first number is greater than the second.


```objectivec
@(2 > 3) → false
@(3 > 3) → false
@(4 > 3) → true
```

<h2 class="item_title"><a name="operator:greaterthanorequal" href="#operator:greaterthanorequal">>=</a></h2>

Returns true if the first number is greater than or equal to the second.


```objectivec
@(2 >= 3) → false
@(3 >= 3) → true
@(4 >= 3) → true
```

<h2 class="item_title"><a name="operator:lessthan" href="#operator:lessthan"><</a></h2>

Returns true if the first number is less than the second.


```objectivec
@(2 < 3) → true
@(3 < 3) → false
@(4 < 3) → false
```

<h2 class="item_title"><a name="operator:lessthanorequal" href="#operator:lessthanorequal"><=</a></h2>

Returns true if the first number is less than or equal to the second.


```objectivec
@(2 <= 3) → true
@(3 <= 3) → true
@(4 <= 3) → false
```

<h2 class="item_title"><a name="operator:multiply" href="#operator:multiply">*</a></h2>

Multiplies two numbers.


```objectivec
@(3 * 2) → 6
@(fields.age * 3) → 69
```

<h2 class="item_title"><a name="operator:negate" href="#operator:negate">- (unary)</a></h2>

Negates a number


```objectivec
@(-fields.age) → -23
```

<h2 class="item_title"><a name="operator:notequal" href="#operator:notequal">!=</a></h2>

Returns true if two values are textually not equal.


```objectivec
@("hello" != "hello") → false
@("hello" != "bar") → true
@(1 != 2) → true
```

<h2 class="item_title"><a name="operator:subtract" href="#operator:subtract">- (binary)</a></h2>

Subtracts two numbers.


```objectivec
@(3 - 2) → 1
@(2 - 3) → -1
```


</div>

# Functions

Expressions have access to a set of built-in functions which can be used to perform more complex tasks. Functions are called 
using the `@(function_name(args..))` syntax, and can take as arguments either literal values `@(length(split("1 2 3", " "))` 
or variables in the context `@(title(contact.name))`.

<div class="functions">
<h2 class="item_title"><a name="function:abs" href="#function:abs">abs(number)</a></h2>

Returns the absolute value of `number`.


```objectivec
@(abs(-10)) → 10
@(abs(10.5)) → 10.5
@(abs("foo")) → ERROR
```

<h2 class="item_title"><a name="function:and" href="#function:and">and(values...)</a></h2>

Returns whether all the given `values` are truthy.


```objectivec
@(and(true)) → true
@(and(true, false, true)) → false
```

<h2 class="item_title"><a name="function:array" href="#function:array">array(values...)</a></h2>

Takes multiple `values` and returns them as an array.


```objectivec
@(array("a", "b", 356)[1]) → b
@(join(array("a", "b", "c"), "|")) → a|b|c
@(count(array())) → 0
@(count(array("a", "b"))) → 2
```

<h2 class="item_title"><a name="function:attachment_parts" href="#function:attachment_parts">attachment_parts(attachment)</a></h2>

Parses an attachment into its different parts


```objectivec
@(attachment_parts("image/jpeg:https://example.com/test.jpg")) → {content_type: image/jpeg, url: https://example.com/test.jpg}
```

<h2 class="item_title"><a name="function:boolean" href="#function:boolean">boolean(value)</a></h2>

Tries to convert `value` to a boolean.

An error is returned if the value can't be converted.


```objectivec
@(boolean(array(1, 2))) → true
@(boolean("FALSE")) → false
@(boolean(1 / 0)) → ERROR
```

<h2 class="item_title"><a name="function:char" href="#function:char">char(code)</a></h2>

Returns the character for the given UNICODE `code`.

It is the inverse of [code](expressions.html#function:code).


```objectivec
@(char(33)) → !
@(char(128512)) → 😀
@(char("foo")) → ERROR
```

<h2 class="item_title"><a name="function:clean" href="#function:clean">clean(text)</a></h2>

Removes any non-printable characters from `text`.


```objectivec
@(clean("😃 Hello \nwo\tr\rld")) → 😃 Hello world
@(clean(123)) → 123
```

<h2 class="item_title"><a name="function:code" href="#function:code">code(text)</a></h2>

Returns the UNICODE code for the first character of `text`.

It is the inverse of [char](expressions.html#function:char).


```objectivec
@(code("a")) → 97
@(code("abc")) → 97
@(code("😀")) → 128512
@(code("15")) → 49
@(code(15)) → 49
@(code("")) → ERROR
```

<h2 class="item_title"><a name="function:count" href="#function:count">count(value)</a></h2>

Returns the number of items in the given array or properties on an object.

It will return an error if it is passed an item which isn't countable.


```objectivec
@(count(contact.fields)) → 5
@(count(array())) → 0
@(count(array("a", "b", "c"))) → 3
@(count(1234)) → ERROR
```

<h2 class="item_title"><a name="function:date" href="#function:date">date(value)</a></h2>

Tries to convert `value` to a date.

If it is text then it will be parsed into a date using the default date format.
An error is returned if the value can't be converted.


```objectivec
@(date("1979-07-18")) → 1979-07-18
@(date("1979-07-18T10:30:45.123456Z")) → 1979-07-18
@(date("10/05/2010")) → 2010-05-10
@(date("NOT DATE")) → ERROR
```

<h2 class="item_title"><a name="function:date_from_parts" href="#function:date_from_parts">date_from_parts(year, month, day)</a></h2>

Creates a date from `year`, `month` and `day`.


```objectivec
@(date_from_parts(2017, 1, 15)) → 2017-01-15
@(date_from_parts(2017, 2, 31)) → 2017-03-03
@(date_from_parts(2017, 13, 15)) → ERROR
```

<h2 class="item_title"><a name="function:datetime" href="#function:datetime">datetime(value)</a></h2>

Tries to convert `value` to a datetime.

If it is text then it will be parsed into a datetime using the default date
and time formats. An error is returned if the value can't be converted.


```objectivec
@(datetime("1979-07-18")) → 1979-07-18T00:00:00.000000-05:00
@(datetime("1979-07-18T10:30:45.123456Z")) → 1979-07-18T10:30:45.123456Z
@(datetime("10/05/2010")) → 2010-05-10T00:00:00.000000-05:00
@(datetime("NOT DATE")) → ERROR
```

<h2 class="item_title"><a name="function:datetime_add" href="#function:datetime_add">datetime_add(datetime, offset, unit)</a></h2>

Calculates the date value arrived at by adding `offset` number of `unit` to the `datetime`

Valid durations are "Y" for years, "M" for months, "W" for weeks, "D" for days, "h" for hour,
"m" for minutes, "s" for seconds


```objectivec
@(datetime_add("2017-01-15", 5, "D")) → 2017-01-20T00:00:00.000000-05:00
@(datetime_add("2017-01-15 10:45", 30, "m")) → 2017-01-15T11:15:00.000000-05:00
```

<h2 class="item_title"><a name="function:datetime_diff" href="#function:datetime_diff">datetime_diff(date1, date2, unit)</a></h2>

Returns the duration between `date1` and `date2` in the `unit` specified.

Valid durations are "Y" for years, "M" for months, "W" for weeks, "D" for days, "h" for hour,
"m" for minutes, "s" for seconds.


```objectivec
@(datetime_diff("2017-01-15", "2017-01-17", "D")) → 2
@(datetime_diff("2017-01-15", "2017-05-15", "W")) → 17
@(datetime_diff("2017-01-15", "2017-05-15", "M")) → 4
@(datetime_diff("2017-01-17 10:50", "2017-01-17 12:30", "h")) → 1
@(datetime_diff("2017-01-17", "2015-12-17", "Y")) → -2
```

<h2 class="item_title"><a name="function:datetime_from_epoch" href="#function:datetime_from_epoch">datetime_from_epoch(seconds)</a></h2>

Converts the UNIX epoch time `seconds` into a new date.


```objectivec
@(datetime_from_epoch(1497286619)) → 2017-06-12T11:56:59.000000-05:00
@(datetime_from_epoch(1497286619.123456)) → 2017-06-12T11:56:59.123456-05:00
```

<h2 class="item_title"><a name="function:default" href="#function:default">default(value, default)</a></h2>

Returns `value` if is not empty or an error, otherwise it returns `default`.


```objectivec
@(default(undeclared.var, "default_value")) → default_value
@(default("10", "20")) → 10
@(default("", "value")) → value
@(default("  ", "value")) → \x20\x20
@(default(datetime("invalid-date"), "today")) → today
@(default(format_urn("invalid-urn"), "ok")) → ok
```

<h2 class="item_title"><a name="function:epoch" href="#function:epoch">epoch(date)</a></h2>

Converts `date` to a UNIX epoch time.

The returned number can contain fractional seconds.


```objectivec
@(epoch("2017-06-12T16:56:59.000000Z")) → 1497286619
@(epoch("2017-06-12T18:56:59.000000+02:00")) → 1497286619
@(epoch("2017-06-12T16:56:59.123456Z")) → 1497286619.123456
@(round_down(epoch("2017-06-12T16:56:59.123456Z"))) → 1497286619
```

<h2 class="item_title"><a name="function:extract" href="#function:extract">extract(object, properties)</a></h2>

Takes an object and extracts the named property.


```objectivec
@(extract(contact, "name")) → Ryan Lewis
@(extract(contact.groups[0], "name")) → Testers
```

<h2 class="item_title"><a name="function:extract_object" href="#function:extract_object">extract_object(object, properties...)</a></h2>

Takes an object and returns a new object by extracting only the named properties.


```objectivec
@(extract_object(contact.groups[0], "name")) → {name: Testers}
```

<h2 class="item_title"><a name="function:field" href="#function:field">field(text, index, delimiter)</a></h2>

Splits `text` using the given `delimiter` and returns the field at `index`.

The index starts at zero. When splitting with a space, the delimiter is considered to be all whitespace.


```objectivec
@(field("a,b,c", 1, ",")) → b
@(field("a,,b,c", 1, ",")) →
@(field("a   b c", 1, " ")) → b
@(field("a		b	c	d", 1, "	")) →
@(field("a\t\tb\tc\td", 1, " ")) →
@(field("a,b,c", "foo", ",")) → ERROR
```

<h2 class="item_title"><a name="function:foreach" href="#function:foreach">foreach(values, func, [args...])</a></h2>

Creates a new array by applying `func` to each value in `values`.

If the given function takes more than one argument, you can pass additional arguments after the function.


```objectivec
@(foreach(array("a", "b", "c"), upper)) → [A, B, C]
@(foreach(array("the man", "fox", "jumped up"), word, 0)) → [the, fox, jumped]
```

<h2 class="item_title"><a name="function:foreach_value" href="#function:foreach_value">foreach_value(object, func, [args...])</a></h2>

Creates a new object by applying `func` to each property value of `object`.

If the given function takes more than one argument, you can pass additional arguments after the function.


```objectivec
@(foreach_value(object("a", "x", "b", "y"), upper)) → {a: X, b: Y}
@(foreach_value(object("a", "hi there", "b", "good bye"), word, 1)) → {a: there, b: bye}
```

<h2 class="item_title"><a name="function:format" href="#function:format">format(value)</a></h2>

Formats `value` according to its type.


```objectivec
@(format(1234.5670)) → 1,234.567
@(format(now())) → 11-04-2018 13:24
@(format(today())) → 11-04-2018
```

<h2 class="item_title"><a name="function:format_date" href="#function:format_date">format_date(date, [,format])</a></h2>

Formats `date` as text according to the given `format`.

If `format` is not specified then the environment's default format is used. The format
string can consist of the following characters. The characters ' ', ':', ',', 'T', '-'
and '_' are ignored. Any other character is an error.

* `YY`        - last two digits of year 0-99
* `YYYY`      - four digits of year 0000-9999
* `M`         - month 1-12
* `MM`        - month 01-12
* `D`         - day of month, 1-31
* `DD`        - day of month, zero padded 0-31


```objectivec
@(format_date("1979-07-18T15:00:00.000000Z")) → 18-07-1979
@(format_date("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) → 1979-07-18
@(format_date("2010-05-10T19:50:00.000000Z", "YYYY M DD")) → 2010 5 10
@(format_date("1979-07-18T15:00:00.000000Z", "YYYY")) → 1979
@(format_date("1979-07-18T15:00:00.000000Z", "M")) → 7
@(format_date("NOT DATE", "YYYY-MM-DD")) → ERROR
```

<h2 class="item_title"><a name="function:format_datetime" href="#function:format_datetime">format_datetime(datetime [,format [,timezone]])</a></h2>

Formats `datetime` as text according to the given `format`.

If `format` is not specified then the environment's default format is used. The format
string can consist of the following characters. The characters ' ', ':', ',', 'T', '-'
and '_' are ignored. Any other character is an error.

* `YY`        - last two digits of year 0-99
* `YYYY`      - four digits of year 0000-9999
* `M`         - month 1-12
* `MM`        - month 01-12
* `D`         - day of month, 1-31
* `DD`        - day of month, zero padded 0-31
* `h`         - hour of the day 1-12
* `hh`        - hour of the day 01-12
* `tt`        - twenty four hour of the day 00-23
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
as "America/Guayaquil" or "America/Los_Angeles". If not specified, the current timezone
will be used. An error will be returned if the timezone is not recognized.


```objectivec
@(format_datetime("1979-07-18T15:00:00.000000Z")) → 18-07-1979 10:00
@(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) → 1979-07-18
@(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY M DD tt:mm")) → 2010 5 10 14:50
@(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY-MM-DD hh:mm AA", "America/Los_Angeles")) → 2010-05-10 12:50 PM
@(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY")) → 1979
@(format_datetime("1979-07-18T15:00:00.000000Z", "M")) → 7
@(format_datetime("NOT DATE", "YYYY-MM-DD")) → ERROR
```

<h2 class="item_title"><a name="function:format_location" href="#function:format_location">format_location(location)</a></h2>

Formats the given `location` as its name.


```objectivec
@(format_location("Rwanda")) → Rwanda
@(format_location("Rwanda > Kigali")) → Kigali
```

<h2 class="item_title"><a name="function:format_number" href="#function:format_number">format_number(number, places [, humanize])</a></h2>

Formats `number` to the given number of decimal `places`.

An optional third argument `humanize` can be false to disable the use of thousand separators.


```objectivec
@(format_number(1234)) → 1,234
@(format_number(1234.5670)) → 1,234.567
@(format_number(1234.5670, 2, true)) → 1,234.57
@(format_number(1234.5678, 0, false)) → 1235
@(format_number("foo", 2, false)) → ERROR
```

<h2 class="item_title"><a name="function:format_time" href="#function:format_time">format_time(time [,format])</a></h2>

Formats `time` as text according to the given `format`.

If `format` is not specified then the environment's default format is used. The format
string can consist of the following characters. The characters ' ', ':', ',', 'T', '-'
and '_' are ignored. Any other character is an error.

* `h`         - hour of the day 1-12
* `hh`        - hour of the day 01-12
* `tt`        - twenty four hour of the day 00-23
* `m`         - minute 0-59
* `mm`        - minute 00-59
* `s`         - second 0-59
* `ss`        - second 00-59
* `fff`       - milliseconds
* `ffffff`    - microseconds
* `fffffffff` - nanoseconds
* `aa`        - am or pm
* `AA`        - AM or PM


```objectivec
@(format_time("14:50:30.000000")) → 14:50
@(format_time("14:50:30.000000", "h:mm aa")) → 2:50 pm
@(format_time("15:00:27.000000", "s")) → 27
@(format_time("NOT TIME", "hh:mm")) → ERROR
```

<h2 class="item_title"><a name="function:format_urn" href="#function:format_urn">format_urn(urn)</a></h2>

Formats `urn` into human friendly text.


```objectivec
@(format_urn("tel:+250781234567")) → 0781 234 567
@(format_urn("twitter:134252511151#billy_bob")) → billy_bob
@(format_urn(contact.urn)) → (202) 456-1111
@(format_urn(urns.tel)) → (202) 456-1111
@(format_urn(urns.mailto)) → foo@bar.com
@(format_urn("NOT URN")) → ERROR
```

<h2 class="item_title"><a name="function:html_decode" href="#function:html_decode">html_decode(text)</a></h2>

HTML decodes `text`


```objectivec
@(html_decode("Red &amp; Blue")) → Red & Blue
@(html_decode("5 + 10")) → 5 + 10
```

<h2 class="item_title"><a name="function:if" href="#function:if">if(test, value1, value2)</a></h2>

Returns `value1` if `test` is truthy or `value2` if not.

If the first argument is an error that error is returned.


```objectivec
@(if(1 = 1, "foo", "bar")) → foo
@(if("foo" > "bar", "foo", "bar")) → ERROR
```

<h2 class="item_title"><a name="function:is_error" href="#function:is_error">is_error(value)</a></h2>

Returns whether `value` is an error


```objectivec
@(is_error(datetime("foo"))) → true
@(is_error(run.not.existing)) → true
@(is_error("hello")) → false
```

<h2 class="item_title"><a name="function:join" href="#function:join">join(array, separator)</a></h2>

Joins the given `array` of strings with `separator` to make text.


```objectivec
@(join(array("a", "b", "c"), "|")) → a|b|c
@(join(split("a.b.c", "."), " ")) → a b c
```

<h2 class="item_title"><a name="function:json" href="#function:json">json(value)</a></h2>

Returns the JSON representation of `value`.


```objectivec
@(json("string")) → "string"
@(json(10)) → 10
@(json(null)) → null
@(json(contact.uuid)) → "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
```

<h2 class="item_title"><a name="function:lower" href="#function:lower">lower(text)</a></h2>

Converts `text` to lowercase.


```objectivec
@(lower("HellO")) → hello
@(lower("hello")) → hello
@(lower("123")) → 123
@(lower("😀")) → 😀
```

<h2 class="item_title"><a name="function:max" href="#function:max">max(numbers...)</a></h2>

Returns the maximum value in `numbers`.


```objectivec
@(max(1, 2)) → 2
@(max(1, -1, 10)) → 10
@(max(1, 10, "foo")) → ERROR
```

<h2 class="item_title"><a name="function:mean" href="#function:mean">mean(numbers...)</a></h2>

Returns the arithmetic mean of `numbers`.


```objectivec
@(mean(1, 2)) → 1.5
@(mean(1, 2, 6)) → 3
@(mean(1, "foo")) → ERROR
```

<h2 class="item_title"><a name="function:min" href="#function:min">min(numbers...)</a></h2>

Returns the minimum value in `numbers`.


```objectivec
@(min(1, 2)) → 1
@(min(2, 2, -10)) → -10
@(min(1, 2, "foo")) → ERROR
```

<h2 class="item_title"><a name="function:mod" href="#function:mod">mod(dividend, divisor)</a></h2>

Returns the remainder of the division of `dividend` by `divisor`.


```objectivec
@(mod(5, 2)) → 1
@(mod(4, 2)) → 0
@(mod(5, "foo")) → ERROR
```

<h2 class="item_title"><a name="function:now" href="#function:now">now()</a></h2>

Returns the current date and time in the current timezone.


```objectivec
@(now()) → 2018-04-11T13:24:30.123456-05:00
```

<h2 class="item_title"><a name="function:number" href="#function:number">number(value)</a></h2>

Tries to convert `value` to a number.

An error is returned if the value can't be converted.


```objectivec
@(number(10)) → 10
@(number("123.45000")) → 123.45
@(number("what?")) → ERROR
```

<h2 class="item_title"><a name="function:object" href="#function:object">object(pairs...)</a></h2>

Takes property name value pairs and returns them as a new object.


```objectivec
@(object()) → {}
@(object("a", 123, "b", "hello")) → {a: 123, b: hello}
@(object("a")) → ERROR
```

<h2 class="item_title"><a name="function:or" href="#function:or">or(values...)</a></h2>

Returns whether if any of the given `values` are truthy.


```objectivec
@(or(true)) → true
@(or(true, false, true)) → true
```

<h2 class="item_title"><a name="function:parse_datetime" href="#function:parse_datetime">parse_datetime(text, format [,timezone])</a></h2>

Parses `text` into a date using the given `format`.

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
as "America/Guayaquil" or "America/Los_Angeles". If not specified, the current timezone
will be used. An error will be returned if the timezone is not recognized.

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

<h2 class="item_title"><a name="function:parse_json" href="#function:parse_json">parse_json(text)</a></h2>

Tries to parse `text` as JSON.

If the given `text` is not valid JSON, then an error is returned


```objectivec
@(parse_json("{\"foo\": \"bar\"}").foo) → bar
@(parse_json("[1,2,3,4]")[2]) → 3
@(parse_json("invalid json")) → ERROR
```

<h2 class="item_title"><a name="function:parse_time" href="#function:parse_time">parse_time(text, format)</a></h2>

Parses `text` into a time using the given `format`.

The format string can consist of the following characters. The characters
' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.

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

Note that fractional seconds will be parsed even without an explicit format identifier.
You should only specify fractional seconds when you want to assert the number of places
in the input format.

parse_time will return an error if it is unable to convert the text to a time.


```objectivec
@(parse_time("15:28", "tt:mm")) → 15:28:00.000000
@(parse_time("2:40 pm", "h:mm aa")) → 14:40:00.000000
@(parse_time("NOT TIME", "tt:mm")) → ERROR
```

<h2 class="item_title"><a name="function:percent" href="#function:percent">percent(number)</a></h2>

Formats `number` as a percentage.


```objectivec
@(percent(0.54234)) → 54%
@(percent(1.2)) → 120%
@(percent("foo")) → ERROR
```

<h2 class="item_title"><a name="function:rand" href="#function:rand">rand()</a></h2>

Returns a single random number between [0.0-1.0).


```objectivec
@(rand()) → 0.607552015674623913099594574305228888988494873046875
@(rand()) → 0.484677570947340263796121462291921488940715789794921875
```

<h2 class="item_title"><a name="function:rand_between" href="#function:rand_between">rand_between()</a></h2>

A single random integer in the given inclusive range.


```objectivec
@(rand_between(1, 10)) → 10
@(rand_between(1, 10)) → 2
```

<h2 class="item_title"><a name="function:read_chars" href="#function:read_chars">read_chars(text)</a></h2>

Converts `text` into something that can be read by IVR systems.

ReadChars will split the numbers such as they are easier to understand. This includes
splitting in 3s or 4s if appropriate.


```objectivec
@(read_chars("1234")) → 1 2 3 4
@(read_chars("abc")) → a b c
@(read_chars("abcdef")) → a b c , d e f
```

<h2 class="item_title"><a name="function:regex_match" href="#function:regex_match">regex_match(text, pattern [,group])</a></h2>

Returns the first match of the regular expression `pattern` in `text`.

An optional third parameter `group` determines which matching group will be returned.


```objectivec
@(regex_match("sda34dfddg67", "\d+")) → 34
@(regex_match("Bob Smith", "(\w+) (\w+)", 1)) → Bob
@(regex_match("Bob Smith", "(\w+) (\w+)", 2)) → Smith
@(regex_match("Bob Smith", "(\w+) (\w+)", 5)) → ERROR
@(regex_match("abc", "[\.")) → ERROR
```

<h2 class="item_title"><a name="function:remove_first_word" href="#function:remove_first_word">remove_first_word(text)</a></h2>

Removes the first word of `text`.


```objectivec
@(remove_first_word("foo bar")) → bar
@(remove_first_word("Hi there. I'm a flow!")) → there. I'm a flow!
```

<h2 class="item_title"><a name="function:repeat" href="#function:repeat">repeat(text, count)</a></h2>

Returns `text` repeated `count` number of times.


```objectivec
@(repeat("*", 8)) → ********
@(repeat("*", "foo")) → ERROR
```

<h2 class="item_title"><a name="function:replace" href="#function:replace">replace(text, needle, replacement [, count])</a></h2>

Replaces up to `count` occurrences of `needle` with `replacement` in `text`.

If `count` is omitted or is less than 0 then all occurrences are replaced.


```objectivec
@(replace("foo bar foo", "foo", "zap")) → zap bar zap
@(replace("foo bar foo", "foo", "zap", 1)) → zap bar foo
@(replace("foo bar", "baz", "zap")) → foo bar
```

<h2 class="item_title"><a name="function:replace_time" href="#function:replace_time">replace_time(datetime)</a></h2>

Returns a new datetime with the time part replaced by the `time`.


```objectivec
@(replace_time(now(), "10:30")) → 2018-04-11T10:30:00.000000-05:00
@(replace_time("2017-01-15", "10:30")) → 2017-01-15T10:30:00.000000-05:00
@(replace_time("foo", "10:30")) → ERROR
```

<h2 class="item_title"><a name="function:round" href="#function:round">round(number [,places])</a></h2>

Rounds `number` to the nearest value.

You can optionally pass in the number of decimal places to round to as `places`. If `places` < 0,
it will round the integer part to the nearest 10^(-places).


```objectivec
@(round(12)) → 12
@(round(12.141)) → 12
@(round(12.6)) → 13
@(round(12.141, 2)) → 12.14
@(round(12.146, 2)) → 12.15
@(round(12.146, -1)) → 10
@(round("notnum", 2)) → ERROR
```

<h2 class="item_title"><a name="function:round_down" href="#function:round_down">round_down(number [,places])</a></h2>

Rounds `number` down to the nearest integer value.

You can optionally pass in the number of decimal places to round to as `places`.


```objectivec
@(round_down(12)) → 12
@(round_down(12.141)) → 12
@(round_down(12.6)) → 12
@(round_down(12.141, 2)) → 12.14
@(round_down(12.146, 2)) → 12.14
@(round_down("foo")) → ERROR
```

<h2 class="item_title"><a name="function:round_up" href="#function:round_up">round_up(number [,places])</a></h2>

Rounds `number` up to the nearest integer value.

You can optionally pass in the number of decimal places to round to as `places`.


```objectivec
@(round_up(12)) → 12
@(round_up(12.141)) → 13
@(round_up(12.6)) → 13
@(round_up(12.141, 2)) → 12.15
@(round_up(12.146, 2)) → 12.15
@(round_up("foo")) → ERROR
```

<h2 class="item_title"><a name="function:split" href="#function:split">split(text, [,delimiters])</a></h2>

Splits `text` into an array of separated words.

Empty values are removed from the returned list. There is an optional final parameter `delimiters` which
is string of characters used to split the text into words.


```objectivec
@(split("a b c")) → [a, b, c]
@(split("a", " ")) → [a]
@(split("abc..d", ".")) → [abc, d]
@(split("a.b.c.", ".")) → [a, b, c]
@(split("a|b,c  d", " .|,")) → [a, b, c, d]
```

<h2 class="item_title"><a name="function:text" href="#function:text">text(value)</a></h2>

Tries to convert `value` to text.

An error is returned if the value can't be converted.


```objectivec
@(text(3 = 3)) → true
@(json(text(123.45))) → "123.45"
@(text(1 / 0)) → ERROR
```

<h2 class="item_title"><a name="function:text_compare" href="#function:text_compare">text_compare(text1, text2)</a></h2>

Returns the dictionary order of `text1` and `text2`.

The return value will be -1 if `text1` comes before `text2`, 0 if they are equal
and 1 if `text1` comes after `text2`.


```objectivec
@(text_compare("abc", "abc")) → 0
@(text_compare("abc", "def")) → -1
@(text_compare("zzz", "aaa")) → 1
```

<h2 class="item_title"><a name="function:text_length" href="#function:text_length">text_length(value)</a></h2>

Returns the length (number of characters) of `value` when converted to text.


```objectivec
@(text_length("abc")) → 3
@(text_length(array(2, 3))) → 6
```

<h2 class="item_title"><a name="function:text_slice" href="#function:text_slice">text_slice(text, start [, end])</a></h2>

Returns the portion of `text` between `start` (inclusive) and `end` (exclusive).

If `end` is not specified then the entire rest of `text` will be included. Negative values
for `start` or `end` start at the end of `text`.


```objectivec
@(text_slice("hello", 2)) → llo
@(text_slice("hello", 1, 3)) → el
@(text_slice("hello😁", -3, -1)) → lo
@(text_slice("hello", 7)) →
```

<h2 class="item_title"><a name="function:time" href="#function:time">time(value)</a></h2>

Tries to convert `value` to a time.

If it is text then it will be parsed into a time using the default time format.
An error is returned if the value can't be converted.


```objectivec
@(time("10:30")) → 10:30:00.000000
@(time("10:30:45 PM")) → 22:30:45.000000
@(time(datetime("1979-07-18T10:30:45.123456Z"))) → 10:30:45.123456
@(time("what?")) → ERROR
```

<h2 class="item_title"><a name="function:time_from_parts" href="#function:time_from_parts">time_from_parts(hour, minute, second)</a></h2>

Creates a time from `hour`, `minute` and `second`


```objectivec
@(time_from_parts(14, 40, 15)) → 14:40:15.000000
@(time_from_parts(8, 10, 0)) → 08:10:00.000000
@(time_from_parts(25, 0, 0)) → ERROR
```

<h2 class="item_title"><a name="function:title" href="#function:title">title(text)</a></h2>

Capitalizes each word in `text`.


```objectivec
@(title("foo")) → Foo
@(title("ryan lewis")) → Ryan Lewis
@(title("RYAN LEWIS")) → Ryan Lewis
@(title(123)) → 123
```

<h2 class="item_title"><a name="function:today" href="#function:today">today()</a></h2>

Returns the current date in the environment timezone.


```objectivec
@(today()) → 2018-04-11
```

<h2 class="item_title"><a name="function:trim" href="#function:trim">trim(text, [,chars])</a></h2>

Removes whitespace from either end of `text`.

There is an optional final parameter `chars` which is string of characters to be removed instead of whitespace.


```objectivec
@(trim(" hello world    ")) → hello world
@(trim("+123157568", "+")) → 123157568
```

<h2 class="item_title"><a name="function:trim_left" href="#function:trim_left">trim_left(text, [,chars])</a></h2>

Removes whitespace from the start of `text`.

There is an optional final parameter `chars` which is string of characters to be removed instead of whitespace.


```objectivec
@("*" & trim_left(" hello world   ") & "*") → *hello world   *
@(trim_left("+12345+", "+")) → 12345+
```

<h2 class="item_title"><a name="function:trim_right" href="#function:trim_right">trim_right(text, [,chars])</a></h2>

Removes whitespace from the end of `text`.

There is an optional final parameter `chars` which is string of characters to be removed instead of whitespace.


```objectivec
@("*" & trim_right(" hello world   ") & "*") → * hello world*
@(trim_right("+12345+", "+")) → +12345
```

<h2 class="item_title"><a name="function:tz" href="#function:tz">tz(date)</a></h2>

Returns the name of the timezone of `date`.

If no timezone information is present in the date, then the current timezone will be returned.


```objectivec
@(tz("2017-01-15T02:15:18.123456Z")) → UTC
@(tz("2017-01-15 02:15:18PM")) → America/Guayaquil
@(tz("2017-01-15")) → America/Guayaquil
@(tz("foo")) → ERROR
```

<h2 class="item_title"><a name="function:tz_offset" href="#function:tz_offset">tz_offset(date)</a></h2>

Returns the offset of the timezone of `date`.

The offset is returned in the format `[+/-]HH:MM`. If no timezone information is present in the date,
then the current timezone offset will be returned.


```objectivec
@(tz_offset("2017-01-15T02:15:18.123456Z")) → +0000
@(tz_offset("2017-01-15 02:15:18PM")) → -0500
@(tz_offset("2017-01-15")) → -0500
@(tz_offset("foo")) → ERROR
```

<h2 class="item_title"><a name="function:upper" href="#function:upper">upper(text)</a></h2>

Converts `text` to uppercase.


```objectivec
@(upper("Asdf")) → ASDF
@(upper(123)) → 123
```

<h2 class="item_title"><a name="function:url_encode" href="#function:url_encode">url_encode(text)</a></h2>

Encodes `text` for use as a URL parameter.


```objectivec
@(url_encode("two & words")) → two%20%26%20words
@(url_encode(10)) → 10
```

<h2 class="item_title"><a name="function:urn_parts" href="#function:urn_parts">urn_parts(urn)</a></h2>

Parses a URN into its different parts


```objectivec
@(urn_parts("tel:+593979012345")) → {display: , path: +593979012345, scheme: tel}
@(urn_parts("twitterid:3263621177#bobby")) → {display: bobby, path: 3263621177, scheme: twitterid}
@(urn_parts("not a urn")) → ERROR
```

<h2 class="item_title"><a name="function:week_number" href="#function:week_number">week_number(date)</a></h2>

Returns the week number (1-54) of `date`.

The week is considered to start on Sunday and week containing Jan 1st is week number 1.


```objectivec
@(week_number("2019-01-01")) → 1
@(week_number("2019-07-23T16:56:59.000000Z")) → 30
@(week_number("xx")) → ERROR
```

<h2 class="item_title"><a name="function:weekday" href="#function:weekday">weekday(date)</a></h2>

Returns the day of the week for `date`.

The week is considered to start on Sunday so a Sunday returns 0, a Monday returns 1 etc.


```objectivec
@(weekday("2017-01-15")) → 0
@(weekday("foo")) → ERROR
```

<h2 class="item_title"><a name="function:word" href="#function:word">word(text, index [,delimiters])</a></h2>

Returns the word at `index` in `text`.

Indexes start at zero. There is an optional final parameter `delimiters` which
is string of characters used to split the text into words.


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

<h2 class="item_title"><a name="function:word_count" href="#function:word_count">word_count(text [,delimiters])</a></h2>

Returns the number of words in `text`.

There is an optional final parameter `delimiters` which is string of characters used
to split the text into words.


```objectivec
@(word_count("foo bar")) → 2
@(word_count(10)) → 1
@(word_count("")) → 0
@(word_count("😀😃😄😁")) → 4
@(word_count("bee.*cat,dog", ".*=|")) → 2
@(word_count("O'Grady O'Flaggerty", " ")) → 2
```

<h2 class="item_title"><a name="function:word_slice" href="#function:word_slice">word_slice(text, start, end [,delimiters])</a></h2>

Extracts a sub-sequence of words from `text`.

The returned words are those from `start` up to but not-including `end`. Indexes start at zero and a negative
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
