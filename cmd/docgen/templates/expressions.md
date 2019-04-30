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
{{ .typeDocs }}
</div>

# Operators

<div class="operators">
{{ .operatorDocs }}
</div>

# Functions

Expressions have access to a set of built-in functions which can be used to perform more complex tasks. Functions are called 
using the `@(function_name(args..))` syntax, and can take as arguments either literal values `@(length(split("1 2 3", " "))` 
or variables in the context `@(title(contact.name))`.

<div class="functions">
{{ .functionDocs }}
</div>
