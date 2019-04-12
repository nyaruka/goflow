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
{{ .contextDocs }}
</div>

