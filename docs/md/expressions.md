# Overview

Excellent gets its name from borrowing some of the syntax and function names of formulas in Microsoft Excel™, 
though it has evolved over time and similarities are now much fewer. It is an expression based templating 
language which aims to make it easy to generate text from a context of values.

# Templates

Templates can contain single variables or more complex expressions. A single variable is embedded using the `@` 
character. For example the template `Hi @foo` contains a single variable which at runtime will be replaced with 
with the value of `foo` in the context.

More complex expressions can be embedded using the `@(...)` syntax. For example the template `Hi @("Dr " & upper(foo))` 
takes the value of `foo`, converts it to uppercase, and the prefixes it with another string. Note than within a 
complex expression you don't prefix variables with `@`.

The `@` symbol can be escaped in templates by repeating it, ie, `Hi @@twitter` would output `Hi @twitter`.

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
<a name="type:array"></a>

## Array

Is an array of items.


```objectivec
@(array(1, "x", true)) → [1, x, true]
@(array(1, "x", true)[1]) → x
@(count(array(1, "x", true))) → 3
@(json(array(1, "x", true))) → [1,"x",true]
```

<a name="type:boolean"></a>

## Boolean

Is a boolean `true` or `false`.


```objectivec
@(true) → true
@(1 = 1) → true
@(1 = 2) → false
@(json(true)) → true
```

<a name="type:date"></a>

## Date

Is a Gregorian calendar date value.


```objectivec
@(date_from_parts(2019, 4, 11)) → 2019-04-11
@(format_date(date_from_parts(2019, 4, 11))) → 11-04-2019
@(json(date_from_parts(2019, 4, 11))) → "2019-04-11"
```

<a name="type:datetime"></a>

## Datetime

Is a datetime value.


```objectivec
@(datetime("1979-07-18T10:30:45.123456Z")) → 1979-07-18T10:30:45.123456Z
@(format_datetime(datetime("1979-07-18T10:30:45.123456Z"))) → 18-07-1979 05:30
@(json(datetime("1979-07-18T10:30:45.123456Z"))) → "1979-07-18T10:30:45.123456Z"
```

<a name="type:function"></a>

## Function

Is a callable function.


```objectivec
@(upper) → function
@(array(upper)[0]("abc")) → ABC
@(json(upper)) → "function"
```

<a name="type:number"></a>

## Number

Is a whole or fractional number.


```objectivec
@(1234) → 1234
@(1234.5678) → 1234.5678
@(format_number(1234.5678)) → 1,234.57
@(json(1234.5678)) → 1234.5678
```

<a name="type:object"></a>

## Object

Is an object with named properties.


```objectivec
@(object("foo", 1, "bar", "x")) → {bar: x, foo: 1}
@(object("foo", 1, "bar", "x").bar) → x
@(object("foo", 1, "bar", "x")["bar"]) → x
@(count(object("foo", 1, "bar", "x"))) → 2
@(json(object("foo", 1, "bar", "x"))) → {"bar":"x","foo":1}
```

<a name="type:text"></a>

## Text

Is a string of characters.


```objectivec
@("abc") → abc
@(text_length("abc")) → 3
@(upper("abc")) → ABC
@(json("abc")) → "abc"
```

<a name="type:time"></a>

## Time

Is a time of day.


```objectivec
@(time_from_parts(16, 30, 45)) → 16:30:45.000000
@(format_time(time_from_parts(16, 30, 45))) → 16:30
@(json(time_from_parts(16, 30, 45))) → "16:30:45.000000"
```


</div>

# Operators

<div class="operators">
<a name="operator:add"></a>

## Add

Adds two numbers.


```objectivec
@(2 + 3) → 5
@(fields.age + 10) → 33
```

<a name="operator:concatenate"></a>

## Concatenate

Joins two text values together.


```objectivec
@("hello" & " " & "bar") → hello bar
@("hello" & null) → hello
```

<a name="operator:divide"></a>

## Divide

Divides a number by another.


```objectivec
@(4 / 2) → 2
@(3 / 2) → 1.5
@(46 / fields.age) → 2
@(3 / 0) → ERROR
```

<a name="operator:equal"></a>

## Equal

Returns true if two values are textually equal.


```objectivec
@("hello" = "hello") → true
@("hello" = "bar") → false
@(1 = 1) → true
```

<a name="operator:exponent"></a>

## Exponent

Raises a number to the power of a another number.


```objectivec
@(2 ^ 8) → 256
```

<a name="operator:greaterthan"></a>

## Greaterthan

Returns true if the first number is greater than the second.


```objectivec
@(2 > 3) → false
@(3 > 3) → false
@(4 > 3) → true
```

<a name="operator:greaterthanorequal"></a>

## Greaterthanorequal

Returns true if the first number is greater than or equal to the second.


```objectivec
@(2 >= 3) → false
@(3 >= 3) → true
@(4 >= 3) → true
```

<a name="operator:lessthan"></a>

## Lessthan

Returns true if the first number is less than the second.


```objectivec
@(2 < 3) → true
@(3 < 3) → false
@(4 < 3) → false
```

<a name="operator:lessthanorequal"></a>

## Lessthanorequal

Returns true if the first number is less than or equal to the second.


```objectivec
@(2 <= 3) → true
@(3 <= 3) → true
@(4 <= 3) → false
```

<a name="operator:multiply"></a>

## Multiply

Multiplies two numbers.


```objectivec
@(3 * 2) → 6
@(fields.age * 3) → 69
```

<a name="operator:negate"></a>

## Negate

Negates a number


```objectivec
@(-fields.age) → -23
```

<a name="operator:notequal"></a>

## Notequal

Returns true if two values are textually not equal.


```objectivec
@("hello" != "hello") → false
@("hello" != "bar") → true
@(1 != 2) → true
```

<a name="operator:subtract"></a>

## Subtract

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
<a name="function:abs"></a>

## abs(num)

Returns the absolute value of `num`.


```objectivec
@(abs(-10)) → 10
@(abs(10.5)) → 10.5
@(abs("foo")) → ERROR
```

<a name="function:and"></a>

## and(values...)

Returns whether all the given `values` are truthy.


```objectivec
@(and(true)) → true
@(and(true, false, true)) → false
```

<a name="function:array"></a>

## array(values...)

Takes multiple `values` and returns them as an array.


```objectivec
@(array("a", "b", 356)[1]) → b
@(join(array("a", "b", "c"), "|")) → a|b|c
@(count(array())) → 0
@(count(array("a", "b"))) → 2
```

<a name="function:attachment_parts"></a>

## attachment_parts(attachment)

Parses an attachment into its different parts


```objectivec
@(attachment_parts("image/jpeg:https://example.com/test.jpg")) → {content_type: image/jpeg, url: https://example.com/test.jpg}
```

<a name="function:boolean"></a>

## boolean(value)

Tries to convert `value` to a boolean.

An error is returned if the value can't be converted.


```objectivec
@(boolean(array(1, 2))) → true
@(boolean("FALSE")) → false
@(boolean(1 / 0)) → ERROR
```

<a name="function:char"></a>

## char(code)

Returns the character for the given UNICODE `code`.

It is the inverse of [code](expressions.html#function:code).


```objectivec
@(char(33)) → !
@(char(128512)) → 😀
@(char("foo")) → ERROR
```

<a name="function:clean"></a>

## clean(text)

Strips any non-printable characters from `text`.


```objectivec
@(clean("😃 Hello \nwo\tr\rld")) → 😃 Hello world
@(clean(123)) → 123
```

<a name="function:code"></a>

## code(text)

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

<a name="function:count"></a>

## count(value)

Returns the number of items in the given array or properties on an object.

It will return an error if it is passed an item which isn't countable.


```objectivec
@(count(contact.fields)) → 5
@(count(array())) → 0
@(count(array("a", "b", "c"))) → 3
@(count(1234)) → ERROR
```

<a name="function:date"></a>

## date(value)

Tries to convert `value` to a date.

If it is text then it will be parsed into a date using the default date format.
An error is returned if the value can't be converted.


```objectivec
@(date("1979-07-18")) → 1979-07-18
@(date("1979-07-18T10:30:45.123456Z")) → 1979-07-18
@(date("10/05/2010")) → 2010-05-10
@(date("NOT DATE")) → ERROR
```

<a name="function:date_from_parts"></a>

## date_from_parts(year, month, day)

Creates a date from `year`, `month` and `day`.


```objectivec
@(date_from_parts(2017, 1, 15)) → 2017-01-15
@(date_from_parts(2017, 2, 31)) → 2017-03-03
@(date_from_parts(2017, 13, 15)) → ERROR
```

<a name="function:datetime"></a>

## datetime(value)

Tries to convert `value` to a datetime.

If it is text then it will be parsed into a datetime using the default date
and time formats. An error is returned if the value can't be converted.


```objectivec
@(datetime("1979-07-18")) → 1979-07-18T00:00:00.000000-05:00
@(datetime("1979-07-18T10:30:45.123456Z")) → 1979-07-18T10:30:45.123456Z
@(datetime("10/05/2010")) → 2010-05-10T00:00:00.000000-05:00
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

<a name="function:datetime_from_epoch"></a>

## datetime_from_epoch(seconds)

Converts the UNIX epoch time `seconds` into a new date.


```objectivec
@(datetime_from_epoch(1497286619)) → 2017-06-12T11:56:59.000000-05:00
@(datetime_from_epoch(1497286619.123456)) → 2017-06-12T11:56:59.123456-05:00
```

<a name="function:default"></a>

## default(value, default)

Returns `value` if is not empty or an error, otherwise it returns `default`.


```objectivec
@(default(undeclared.var, "default_value")) → default_value
@(default("10", "20")) → 10
@(default("", "value")) → value
@(default(array(1, 2), "value")) → [1, 2]
@(default(array(), "value")) → value
@(default(datetime("invalid-date"), "today")) → today
@(default(format_urn("invalid-urn"), "ok")) → ok
```

<a name="function:epoch"></a>

## epoch(date)

Converts `date` to a UNIX epoch time.

The returned number can contain fractional seconds.


```objectivec
@(epoch("2017-06-12T16:56:59.000000Z")) → 1497286619
@(epoch("2017-06-12T18:56:59.000000+02:00")) → 1497286619
@(epoch("2017-06-12T16:56:59.123456Z")) → 1497286619.123456
@(round_down(epoch("2017-06-12T16:56:59.123456Z"))) → 1497286619
```

<a name="function:extract"></a>

## extract(object, properties...)

Takes an object and extracts the named property.


```objectivec
@(extract(contact, "name")) → Ryan Lewis
@(extract(contact.groups[0], "name")) → Testers
```

<a name="function:extract_object"></a>

## extract_object(object, properties...)

Takes an object and returns a new object by extracting only the named properties.


```objectivec
@(extract_object(contact.groups[0], "name")) → {name: Testers}
```

<a name="function:field"></a>

## field(text, index, delimiter)

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

<a name="function:foreach"></a>

## foreach(array, func, [args...])

Takes an array of objects and returns a new array by applying the given function to each item.

If the given function takes more than one argument, you can pass additional arguments after the function.


```objectivec
@(foreach(array("a", "b", "c"), upper)) → [A, B, C]
@(foreach(array("the man", "fox", "jumped up"), word, 0)) → [the, fox, jumped]
```

<a name="function:format_date"></a>

## format_date(date, [,format])

Formats `date` as text according to the given `format`. If `format` is not
specified then the environment's default format is used.

The format string can consist of the following characters. The characters
' ', ':', ',', 'T', '-' and '_' are ignored. Any other character is an error.

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

<a name="function:format_datetime"></a>

## format_datetime(date [,format [,timezone]])

Formats `date` as text according to the given `format`. If `format` is not
specified then the environment's default format is used.

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


```objectivec
@(format_datetime("1979-07-18T15:00:00.000000Z")) → 18-07-1979 10:00
@(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY-MM-DD")) → 1979-07-18
@(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY M DD tt:mm")) → 2010 5 10 14:50
@(format_datetime("2010-05-10T19:50:00.000000Z", "YYYY-MM-DD tt:mm AA", "America/Los_Angeles")) → 2010-05-10 12:50 PM
@(format_datetime("1979-07-18T15:00:00.000000Z", "YYYY")) → 1979
@(format_datetime("1979-07-18T15:00:00.000000Z", "M")) → 7
@(format_datetime("NOT DATE", "YYYY-MM-DD")) → ERROR
```

<a name="function:format_input"></a>

## format_input(urn)

Formats `input` to be the text followed by the URLs of any attachment, separated by newlines.


```objectivec
@(format_input(input)) → Hi there\nhttp://s3.amazon.com/bucket/test.jpg\nhttp://s3.amazon.com/bucket/test.mp3
@(format_input("NOT INPUT")) → ERROR
```

<a name="function:format_location"></a>

## format_location(location)

Formats the given `location` as its name.


```objectivec
@(format_location("Rwanda")) → Rwanda
@(format_location("Rwanda > Kigali")) → Kigali
```

<a name="function:format_number"></a>

## format_number(number, places [, humanize])

Formats `number` to the given number of decimal `places`.

An optional third argument `humanize` can be false to disable the use of thousand separators.


```objectivec
@(format_number(31337)) → 31,337.00
@(format_number(31337, 2)) → 31,337.00
@(format_number(31337, 2, true)) → 31,337.00
@(format_number(31337, 0, false)) → 31337
@(format_number("foo", 2, false)) → ERROR
```

<a name="function:format_time"></a>

## format_time(time [,format])

Formats `time` as text according to the given `format`. If `format` is not
specified then the environment's default format is used.

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


```objectivec
@(format_time("14:50:30.000000")) → 14:50
@(format_time("14:50:30.000000", "h:mm aa")) → 2:50 pm
@(format_time("15:00:27.000000", "s")) → 27
@(format_time("NOT TIME", "hh:mm")) → ERROR
```

<a name="function:format_urn"></a>

## format_urn(urn)

Formats `urn` into human friendly text.


```objectivec
@(format_urn("tel:+250781234567")) → 0781 234 567
@(format_urn("twitter:134252511151#billy_bob")) → billy_bob
@(format_urn(contact.urn)) → (206) 555-1212
@(format_urn(urns.tel)) → (206) 555-1212
@(format_urn(urns.mailto)) → foo@bar.com
@(format_urn("NOT URN")) → ERROR
```

<a name="function:if"></a>

## if(test, value1, value2)

Returns `value1` if `test` is truthy or `value2` if not.

If the first argument is an error that error is returned.


```objectivec
@(if(1 = 1, "foo", "bar")) → foo
@(if("foo" > "bar", "foo", "bar")) → ERROR
```

<a name="function:is_error"></a>

## is_error(value)

Returns whether `value` is an error


```objectivec
@(is_error(datetime("foo"))) → true
@(is_error(run.not.existing)) → true
@(is_error("hello")) → false
```

<a name="function:join"></a>

## join(array, separator)

Joins the given `array` of strings with `separator` to make text.


```objectivec
@(join(array("a", "b", "c"), "|")) → a|b|c
@(join(split("a.b.c", "."), " ")) → a b c
```

<a name="function:json"></a>

## json(value)

Returns the JSON representation of `value`.


```objectivec
@(json("string")) → "string"
@(json(10)) → 10
@(json(null)) → null
@(json(contact.uuid)) → "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
```

<a name="function:left"></a>

## left(text, count)

Returns the `count` left-most characters in `text`


```objectivec
@(left("hello", 2)) → he
@(left("hello", 7)) → hello
@(left("😀😃😄😁", 2)) → 😀😃
@(left("hello", -1)) → ERROR
```

<a name="function:lower"></a>

## lower(text)

Converts `text` to lowercase.


```objectivec
@(lower("HellO")) → hello
@(lower("hello")) → hello
@(lower("123")) → 123
@(lower("😀")) → 😀
```

<a name="function:max"></a>

## max(values...)

Returns the maximum value in `values`.


```objectivec
@(max(1, 2)) → 2
@(max(1, -1, 10)) → 10
@(max(1, 10, "foo")) → ERROR
```

<a name="function:mean"></a>

## mean(values)

Returns the arithmetic mean of the numbers in `values`.


```objectivec
@(mean(1, 2)) → 1.5
@(mean(1, 2, 6)) → 3
@(mean(1, "foo")) → ERROR
```

<a name="function:min"></a>

## min(values)

Returns the minimum value in `values`.


```objectivec
@(min(1, 2)) → 1
@(min(2, 2, -10)) → -10
@(min(1, 2, "foo")) → ERROR
```

<a name="function:mod"></a>

## mod(dividend, divisor)

Returns the remainder of the division of `dividend` by `divisor`.


```objectivec
@(mod(5, 2)) → 1
@(mod(4, 2)) → 0
@(mod(5, "foo")) → ERROR
```

<a name="function:now"></a>

## now()

Returns the current date and time in the current timezone.


```objectivec
@(now()) → 2018-04-11T13:24:30.123456-05:00
```

<a name="function:number"></a>

## number(value)

Tries to convert `value` to a number.

An error is returned if the value can't be converted.


```objectivec
@(number(10)) → 10
@(number("123.45000")) → 123.45
@(number("what?")) → ERROR
```

<a name="function:object"></a>

## object(pairs...)

Takes property name value pairs and returns them as a new object.


```objectivec
@(object()) → {}
@(object("a", 123, "b", "hello")) → {a: 123, b: hello}
@(object("a")) → ERROR
```

<a name="function:or"></a>

## or(values...)

Returns whether if any of the given `values` are truthy.


```objectivec
@(or(true)) → true
@(or(true, false, true)) → true
```

<a name="function:parse_datetime"></a>

## parse_datetime(text, format [,timezone])

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

<a name="function:parse_json"></a>

## parse_json(text)

Tries to parse `text` as JSON.

If the given `text` is not valid JSON, then an error is returned


```objectivec
@(parse_json("{\"foo\": \"bar\"}").foo) → bar
@(parse_json("[1,2,3,4]")[2]) → 3
@(parse_json("invalid json")) → ERROR
```

<a name="function:parse_time"></a>

## parse_time(text, format)

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

<a name="function:percent"></a>

## percent(num)

Formats `num` as a percentage.


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

Converts `text` into something that can be read by IVR systems.

ReadChars will split the numbers such as they are easier to understand. This includes
splitting in 3s or 4s if appropriate.


```objectivec
@(read_chars("1234")) → 1 2 3 4
@(read_chars("abc")) → a b c
@(read_chars("abcdef")) → a b c , d e f
```

<a name="function:regex_match"></a>

## regex_match(text, pattern [,group])

Returns the first match of the regular expression `pattern` in `text`.

An optional third parameter `group` determines which matching group will be returned.


```objectivec
@(regex_match("sda34dfddg67", "\d+")) → 34
@(regex_match("Bob Smith", "(\w+) (\w+)", 1)) → Bob
@(regex_match("Bob Smith", "(\w+) (\w+)", 2)) → Smith
@(regex_match("Bob Smith", "(\w+) (\w+)", 5)) → ERROR
@(regex_match("abc", "[\.")) → ERROR
```

<a name="function:remove_first_word"></a>

## remove_first_word(text)

Removes the first word of `text`.


```objectivec
@(remove_first_word("foo bar")) → bar
@(remove_first_word("Hi there. I'm a flow!")) → there. I'm a flow!
```

<a name="function:repeat"></a>

## repeat(text, count)

Returns `text` repeated `count` number of times.


```objectivec
@(repeat("*", 8)) → ********
@(repeat("*", "foo")) → ERROR
```

<a name="function:replace"></a>

## replace(text, needle, replacement)

Replaces all occurrences of `needle` with `replacement` in `text`.


```objectivec
@(replace("foo bar", "foo", "zap")) → zap bar
@(replace("foo bar", "baz", "zap")) → foo bar
```

<a name="function:replace_time"></a>

## replace_time(date)

Returns the a new date time with the time part replaced by the `time`.


```objectivec
@(replace_time(now(), "10:30")) → 2018-04-11T10:30:00.000000-05:00
@(replace_time("2017-01-15", "10:30")) → 2017-01-15T10:30:00.000000-05:00
@(replace_time("foo", "10:30")) → ERROR
```

<a name="function:right"></a>

## right(text, count)

Returns the `count` right-most characters in `text`


```objectivec
@(right("hello", 2)) → lo
@(right("hello", 7)) → hello
@(right("😀😃😄😁", 2)) → 😄😁
@(right("hello", -1)) → ERROR
```

<a name="function:round"></a>

## round(num [,places])

Rounds `num` to the nearest value.

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

<a name="function:round_down"></a>

## round_down(num [,places])

Rounds `num` down to the nearest integer value.

You can optionally pass in the number of decimal places to round to as `places`.


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

Rounds `num` up to the nearest integer value.

You can optionally pass in the number of decimal places to round to as `places`.


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

Splits `text` based on the given characters in `delimiters`.

Empty values are removed from the returned list.


```objectivec
@(split("a b c", " ")) → [a, b, c]
@(split("a", " ")) → [a]
@(split("abc..d", ".")) → [abc, d]
@(split("a.b.c.", ".")) → [a, b, c]
@(split("a|b,c  d", " .|,")) → [a, b, c, d]
```

<a name="function:text"></a>

## text(value)

Tries to convert `value` to text.

An error is returned if the value can't be converted.


```objectivec
@(text(3 = 3)) → true
@(json(text(123.45))) → "123.45"
@(text(1 / 0)) → ERROR
```

<a name="function:text_compare"></a>

## text_compare(text1, text2)

Returns the dictionary order of `text1` and `text2`.

The return value will be -1 if `text1` comes before `text2`, 0 if they are equal
and 1 if `text1` comes after `text2`.


```objectivec
@(text_compare("abc", "abc")) → 0
@(text_compare("abc", "def")) → -1
@(text_compare("zzz", "aaa")) → 1
```

<a name="function:text_length"></a>

## text_length(value)

Returns the length (number of characters) of `value` when converted to text.


```objectivec
@(text_length("abc")) → 3
@(text_length(array(2, 3))) → 6
```

<a name="function:time"></a>

## time(value)

Tries to convert `value` to a time.

If it is text then it will be parsed into a time using the default time format.
An error is returned if the value can't be converted.


```objectivec
@(time("10:30")) → 10:30:00.000000
@(time("10:30:45 PM")) → 22:30:45.000000
@(time(datetime("1979-07-18T10:30:45.123456Z"))) → 10:30:45.123456
@(time("what?")) → ERROR
```

<a name="function:time_from_parts"></a>

## time_from_parts(hour, minute, second)

Creates a time from `hour`, `minute` and `second`


```objectivec
@(time_from_parts(14, 40, 15)) → 14:40:15.000000
@(time_from_parts(8, 10, 0)) → 08:10:00.000000
@(time_from_parts(25, 0, 0)) → ERROR
```

<a name="function:title"></a>

## title(text)

Capitalizes each word in `text`.


```objectivec
@(title("foo")) → Foo
@(title("ryan lewis")) → Ryan Lewis
@(title("RYAN LEWIS")) → Ryan Lewis
@(title(123)) → 123
```

<a name="function:today"></a>

## today()

Returns the current date in the environment timezone.


```objectivec
@(today()) → 2018-04-11
```

<a name="function:tz"></a>

## tz(date)

Returns the name of the timezone of `date`.

If no timezone information is present in the date, then the current timezone will be returned.


```objectivec
@(tz("2017-01-15T02:15:18.123456Z")) → UTC
@(tz("2017-01-15 02:15:18PM")) → America/Guayaquil
@(tz("2017-01-15")) → America/Guayaquil
@(tz("foo")) → ERROR
```

<a name="function:tz_offset"></a>

## tz_offset(date)

Returns the offset of the timezone of `date`.

The offset is returned in the format `[+/-]HH:MM`. If no timezone information is present in the date,
then the current timezone offset will be returned.


```objectivec
@(tz_offset("2017-01-15T02:15:18.123456Z")) → +0000
@(tz_offset("2017-01-15 02:15:18PM")) → -0500
@(tz_offset("2017-01-15")) → -0500
@(tz_offset("foo")) → ERROR
```

<a name="function:upper"></a>

## upper(text)

Converts `text` to lowercase.


```objectivec
@(upper("Asdf")) → ASDF
@(upper(123)) → 123
```

<a name="function:url_encode"></a>

## url_encode(text)

Encodes `text` for use as a URL parameter.


```objectivec
@(url_encode("two & words")) → two%20%26%20words
@(url_encode(10)) → 10
```

<a name="function:urn_parts"></a>

## urn_parts(urn)

Parses a URN into its different parts


```objectivec
@(urn_parts("tel:+593979012345")) → {display: , path: +593979012345, scheme: tel}
@(urn_parts("twitterid:3263621177#bobby")) → {display: bobby, path: 3263621177, scheme: twitterid}
```

<a name="function:weekday"></a>

## weekday(date)

Returns the day of the week for `date`.

The week is considered to start on Sunday so a Sunday returns 0, a Monday returns 1 etc.


```objectivec
@(weekday("2017-01-15")) → 0
@(weekday("foo")) → ERROR
```

<a name="function:word"></a>

## word(text, index [,delimiters])

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

<a name="function:word_count"></a>

## word_count(text [,delimiters])

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

<a name="function:word_slice"></a>

## word_slice(text, start, end [,delimiters])

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
