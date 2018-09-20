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
 * `input` the current [input](#context:input), shortcut for `@run.input`
 * `results` the current [results](#context:result), shortcut for `@run.results`
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

<div class="context">
{{ .contextDocs }}
</div>

# Functions

Templates also have access to a set of functions which can be used to further manipulate the context. Functions are called 
using the `@(function_name(args..))` syntax. For example, to title-case a contact's name in a message, you can use `@(title(contact.name))`. 
Context variables referred to within functions do not need a leading `@`. Functions can also use literal numbers or strings as arguments, for example
`@(length(split("1 2 3", " "))`.

<div class="functions">
{{ .functionDocs }}
</div>

# Router Tests

Router tests are a special class of functions which are used within the switch router. They are called in the same way as normal functions, but 
all return a test result object which by default evalutes to true or false, but can also be used to find the matching portion of the test by using
the `match` component of the result. The flow editor builds these expressions using UI widgets, but they can be used anywhere a normal template
function is used.

<div class="tests">
{{ .testDocs }}
</div>