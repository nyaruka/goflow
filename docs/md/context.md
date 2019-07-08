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
 * `groups` the groups the contact belongs to ([group](#context:group))
 * `fields` the custom field values of the contact ([fields](#context:fields))
 * `channel` the preferred channel of the contact ([channel](#context:channel))

<a name="context:group"></a>

## Group

 * `uuid` the UUID of the group ([text](#context:text))
 * `name` the name of the group ([text](#context:text))

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
 * `fields` the current contact custom fields values ([fields](#context:fields))
 * `urns` the current contact URN values ([urns](#context:urns))
 * `results` the current run results ([results](#context:results))


</div>

