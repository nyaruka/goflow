# Context

The context is all the variables which are accessible in expressions:

<div class="context">
<a name="context:channel"></a>

## Channel

 * `__default__` the name of the channel ([text](#context:text))
 * `uuid` the UUID of the channel ([text](#context:text))
 * `name` the name of the channel ([text](#context:text))
 * `address` the address of the channel ([text](#context:text))

<a name="context:contact"></a>

## Contact

 * `__default__` the name or URN of the contact ([text](#context:text))
 * `uuid` the UUID of the contact ([text](#context:text))
 * `id` the numeric ID of the contact ([text](#context:text))
 * `first_name` the first name of the contact ([text](#context:text))
 * `name` the name of the contact ([text](#context:text))
 * `language` the language of the contact as 3 ([text](#context:text))
 * `created_on` the creation date of the contact ([datetime](#context:datetime))
 * `urns` the URNs belonging to the contact ([text](#context:text))
 * `urn` the preferred URN of the contact ([text](#context:text))
 * `groups` the groups the contact belongs to ([group](#context:group))
 * `fields` the custom field values of the contact ([fields](#context:fields))
 * `channel` the preferred channel of the contact ([channel](#context:channel))

<a name="context:flow"></a>

## Flow

 * `__default__` the name of the flow ([text](#context:text))
 * `uuid` the UUID of the flow ([text](#context:text))
 * `name` the name of the flow ([text](#context:text))
 * `revision` the revision number of the flow ([text](#context:text))

<a name="context:group"></a>

## Group

 * `uuid` the UUID of the group ([text](#context:text))
 * `name` the name of the group ([text](#context:text))

<a name="context:input"></a>

## Input

 * `__default__` the text and attachments of the input ([text](#context:text))
 * `uuid` the UUID of the input ([text](#context:text))
 * `created_on` the creation date of the input ([datetime](#context:datetime))
 * `channel` the channel that the input was received on ([channel](#context:channel))
 * `urn` the contact URN that the input was received on ([text](#context:text))
 * `text` the text part of the input ([text](#context:text))
 * `attachments` any attachments on the input ([text](#context:text))
 * `external_id` the external ID of the input ([text](#context:text))

<a name="context:related_run"></a>

## Related_run

 * `uuid` the UUID of the run ([text](#context:text))
 * `contact` the contact of the run ([contact](#context:contact))
 * `flow` the flow of the run ([flow](#context:flow))
 * `fields` the custom field values of the run ([fields](#context:fields))
 * `urns` the URN values of the run ([urns](#context:urns))
 * `results` the results saved by the run ([results](#context:results))
 * `status` the current status of the run ([text](#context:text))

<a name="context:result"></a>

## Result

 * `__default__` the value of the result ([text](#context:text))
 * `name` the name of the result ([text](#context:text))
 * `value` the value of the result ([text](#context:text))
 * `category` the category of the result ([text](#context:text))
 * `category_localized` the localized category of the result ([text](#context:text))
 * `input` the input of the result ([text](#context:text))
 * `extra` the extra data of the result such as a webhook response ([any](#context:any))
 * `node_uuid` the UUIF of the node in the flow that generated the result ([text](#context:text))
 * `created_on` the creation date of the result ([datetime](#context:datetime))

<a name="context:root"></a>

## Root

 * `contact` the current contact ([contact](#context:contact))
 * `fields` the custom field values of the current contact ([fields](#context:fields))
 * `urns` the URN values of the current contact ([urns](#context:urns))
 * `results` the current run results ([results](#context:results))
 * `input` the current input from the contact ([input](#context:input))
 * `run` the current run ([run](#context:run))
 * `child` the last child run ([related_run](#context:related_run))
 * `parent` the parent of the run ([related_run](#context:related_run))
 * `webhook` the parsed JSON response of the last webhook call ([any](#context:any))
 * `trigger` the trigger that started this session ([trigger](#context:trigger))

<a name="context:run"></a>

## Run

 * `uuid` the UUID of the run ([text](#context:text))
 * `contact` the contact of the run ([contact](#context:contact))
 * `flow` the flow of the run ([flow](#context:flow))
 * `status` the current status of the run ([text](#context:text))
 * `results` the results saved by the run ([results](#context:results))
 * `created_on` the creation date of the run ([datetime](#context:datetime))
 * `exited_on` the exit date of the run ([datetime](#context:datetime))

<a name="context:trigger"></a>

## Trigger

 * `type` the type of trigger that started this session ([text](#context:text))
 * `params` the parameters passed to the trigger ([any](#context:any))

<a name="context:trigger"></a>

## Trigger

 * `type` the type of trigger that started this session ([text](#context:text))
 * `params` the parameters passed to the trigger ([any](#context:any))

<a name="context:trigger"></a>

## Trigger

 * `type` the type of trigger that started this session ([text](#context:text))
 * `params` the parameters passed to the trigger ([any](#context:any))

<a name="context:trigger"></a>

## Trigger

 * `type` the type of trigger that started this session ([text](#context:text))
 * `params` the parameters passed to the trigger ([any](#context:any))

<a name="context:trigger"></a>

## Trigger

 * `type` the type of trigger that started this session ([text](#context:text))
 * `params` the parameters passed to the trigger ([any](#context:any))

<a name="context:trigger"></a>

## Trigger

 * `type` the type of trigger that started this session ([text](#context:text))
 * `params` the parameters passed to the trigger ([any](#context:any))


</div>

