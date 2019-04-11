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
{{ .typeDocs }}
</div>
