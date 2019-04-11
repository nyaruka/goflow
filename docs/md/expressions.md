# Overview

Excellent gets its name from borrowing some of the syntax and function names of formulas in Microsoft Excel™, 
though it has evolved over time and similarities are now much fewer. It is an expression based templating 
language which aims to make it easy to generate text from a context of values.

# Templates

Templates can contain single variables or more complex expressions. A single variable is embedded using the `@` 
character. For example the template `Hi @foo` contains a single variable which at runtime will be replaced with 
with the value of `foo` in the context.

More complex expressions can be embedded using the `@(...)` syntax. For example the template `Hi @("Dr " & upper(foo))` 
takes the value of `foo`, converts it to uppercase, and the prefixes it with another string.

The `@` symbol can be escaped in templates by repeating it, ie, `Hi @@twitter` would output `Hi @twitter`.

# Types

Excellent has the following types:

 * [Array](#type:array)
 * [Boolean](#type:boolean)
 * [Date](#type:date)
 * [DateTime](#type:datetime)
 * [Dict](#type:dict)
 * [Function](#type:function)
 * [Number](#type:number)
 * [Text](#type:text)
 * [Time](#type:time)

<div class="types">
<a name="type:array"></a>

## Array

Is an array of items.


```objectivec
@(array(1, "x", true)) → [1, x, true]
@(length(array(1, "x", true))) → 3
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

<a name="type:dict"></a>

## Dict

Is a dictionary of keys and values.


```objectivec
@(dict("foo", 1, "bar", "x")) → {bar: x, foo: 1}
@(length(dict("foo", 1, "bar", "x"))) → 2
@(json(dict("foo", 1, "bar", "x"))) → {"bar":"x","foo":1}
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

<a name="type:text"></a>

## Text

Is a string of characters.


```objectivec
@("abc") → abc
@(length("abc")) → 3
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
