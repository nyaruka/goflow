# Root

These are the top-level variables that can be accessed in the context:

 * `contact` the contact ([contact](context.html#context:contact))
 * `fields` the custom field values of the contact (fields)
 * `urns` the URN values of the contact (urns)
 * `results` the current run results (results)
 * `input` the current input from the contact ([input](context.html#context:input))
 * `run` the current run ([run](context.html#context:run))
 * `child` the last child run ([related_run](context.html#context:related_run))
 * `parent` the parent of the run ([related_run](context.html#context:related_run))
 * `webhook` the parsed JSON response of the last webhook call (any)
 * `globals` the global values (globals)
 * `trigger` the trigger that started this session ([trigger](context.html#context:trigger))



# Types

The following types are found in the context:

<div class="context">
<h2 class="item_title"><a name="context:channel" href="#context:channel">channel</a></h2>

Defaults to the name ([text](expressions.html#type:text))

 * `uuid` the UUID of the channel ([text](expressions.html#type:text))
 * `name` the name of the channel ([text](expressions.html#type:text))
 * `address` the address of the channel ([text](expressions.html#type:text))

<h2 class="item_title"><a name="context:contact" href="#context:contact">contact</a></h2>

Defaults to the name or URN ([text](expressions.html#type:text))

 * `uuid` the UUID of the contact ([text](expressions.html#type:text))
 * `id` the numeric ID of the contact ([text](expressions.html#type:text))
 * `first_name` the first name of the contact ([text](expressions.html#type:text))
 * `name` the name of the contact ([text](expressions.html#type:text))
 * `language` the language of the contact as 3-letter ISO code ([text](expressions.html#type:text))
 * `created_on` the creation date of the contact ([datetime](expressions.html#type:datetime))
 * `urns` the URNs belonging to the contact ([text](expressions.html#type:text))
 * `urn` the preferred URN of the contact ([text](expressions.html#type:text))
 * `groups` the groups the contact belongs to ([group](context.html#context:group))
 * `fields` the custom field values of the contact (fields)
 * `channel` the preferred channel of the contact ([channel](context.html#context:channel))

<h2 class="item_title"><a name="context:flow" href="#context:flow">flow</a></h2>

Defaults to the name ([text](expressions.html#type:text))

 * `uuid` the UUID of the flow ([text](expressions.html#type:text))
 * `name` the name of the flow ([text](expressions.html#type:text))
 * `revision` the revision number of the flow ([text](expressions.html#type:text))

<h2 class="item_title"><a name="context:group" href="#context:group">group</a></h2>

 * `uuid` the UUID of the group ([text](expressions.html#type:text))
 * `name` the name of the group ([text](expressions.html#type:text))

<h2 class="item_title"><a name="context:input" href="#context:input">input</a></h2>

Defaults to the text and attachments ([text](expressions.html#type:text))

 * `uuid` the UUID of the input ([text](expressions.html#type:text))
 * `created_on` the creation date of the input ([datetime](expressions.html#type:datetime))
 * `channel` the channel that the input was received on ([channel](context.html#context:channel))
 * `urn` the contact URN that the input was received on ([text](expressions.html#type:text))
 * `text` the text part of the input ([text](expressions.html#type:text))
 * `attachments` any attachments on the input ([text](expressions.html#type:text))
 * `external_id` the external ID of the input ([text](expressions.html#type:text))

<h2 class="item_title"><a name="context:related_run" href="#context:related_run">related_run</a></h2>

Defaults to the contact name and flow UUID ([text](expressions.html#type:text))

 * `uuid` the UUID of the run ([text](expressions.html#type:text))
 * `contact` the contact of the run ([contact](context.html#context:contact))
 * `flow` the flow of the run ([flow](context.html#context:flow))
 * `fields` the custom field values of the run (fields)
 * `urns` the URN values of the run (urns)
 * `results` the results saved by the run (any)
 * `status` the current status of the run ([text](expressions.html#type:text))

<h2 class="item_title"><a name="context:result" href="#context:result">result</a></h2>

Defaults to the value ([text](expressions.html#type:text))

 * `name` the name of the result ([text](expressions.html#type:text))
 * `value` the value of the result ([text](expressions.html#type:text))
 * `category` the category of the result ([text](expressions.html#type:text))
 * `category_localized` the localized category of the result ([text](expressions.html#type:text))
 * `input` the input of the result ([text](expressions.html#type:text))
 * `extra` the extra data of the result such as a webhook response (any)
 * `node_uuid` the UUID of the node in the flow that generated the result ([text](expressions.html#type:text))
 * `created_on` the creation date of the result ([datetime](expressions.html#type:datetime))

<h2 class="item_title"><a name="context:run" href="#context:run">run</a></h2>

Defaults to the contact name and flow UUID ([text](expressions.html#type:text))

 * `uuid` the UUID of the run ([text](expressions.html#type:text))
 * `contact` the contact of the run ([contact](context.html#context:contact))
 * `flow` the flow of the run ([flow](context.html#context:flow))
 * `status` the current status of the run ([text](expressions.html#type:text))
 * `results` the results saved by the run (results)
 * `created_on` the creation date of the run ([datetime](expressions.html#type:datetime))
 * `exited_on` the exit date of the run ([datetime](expressions.html#type:datetime))

<h2 class="item_title"><a name="context:trigger" href="#context:trigger">trigger</a></h2>

 * `type` the type of trigger that started this session ([text](expressions.html#type:text))
 * `params` the parameters passed to the trigger (any)
 * `keyword` the keyword match if this is a keyword trigger (any)


</div>

