# Container

Flow definitions are defined as a list of nodes, the first node being the entry into the flow. The simplest possible flow containing no nodes whatsoever (and therefore being a no-op) contains the following fields:

 * `uuid` the UUID
 * `name` the name
 * `language` the base authoring language used for localization
 * `type` the type, one of `messaging`, `messaging_offline`, `voice` whcih determines which actions are allowed in the flow
 * `nodes` the nodes (may be empty)

For example:

```json
{
    "uuid": "b7bb5e7c-ad49-4e65-9e24-bf7f1e4ff00a",
    "name": "Empty Flow",
    "language": "eng",
    "type": "messaging",
    "nodes": []
}
```

# Nodes

Flow definitions are composed of zero or more nodes, the first node is always the entry node.

A Node consists of:

 * `uuid` the UUID
 * `actions` a list of 0-n actions which will be executed upon first entering a node
 * `router` an optional router which determines which exit to take
 * `exit` a list of 0-n exits which can be used to link to other nodes

At its simplest, a node can be just a single action with no exits, wait or router, such as:

```json
{
    "uuid":"5a06445e-d790-4bd3-a10b-b47bdcc9abed",
    "actions":[{
        "uuid": "abc0a2bf-6b4a-4ee0-83e1-1eebae6948ac",
        "type": "send_msg",
        "text": "What is your name?"
    }]
}
```

If a node wishes to route to another node, it can do so by defining one or more exits, each with the UUID of the node that is next. Without a router defined, the first exit will always be taken. 

An exit consists of:

 * `uuid` the UUID
 * `destination_uuid` the uuid of the node that should be visited if this exit is chosen by the router (optional)

```json
{
    "uuid":"5a06445e-d790-4bd3-a10b-b47bdcc9abed",
    "actions":[{
        "uuid": "abc0a2bf-6b4a-4ee0-83e1-1eebae6948ac",
        "type": "send_msg",
        "text": "What is your name?"
    }],
    "exits": [{
        "uuid":"eb7defc9-3c66-4dfc-80bc-825567ccd9de",
        "destination_uuid":"ee0bee3f-34b3-4275-af78-f9ff52c82e6a"
    }]
}
```

# Actions

Actions on a node generate events which can then be ingested by the engine container. In some cases the actions cause an immediate action, such 
as calling a webhook, in others the engine container is responsible for taking the action based on the event that is output, such as sending 
messages or updating contact fields. In either case the internal state of the engine is always updated to represent the new state so that
flow execution is consistent. For example, while the engine itself does not have access to a contact store, it updates its internal 
representation of a contact's state based on action performed on a flow so that later references in the flow are correct.

<div class="actions">
<h2 class="item_title"><a name="action:add_contact_groups" href="#action:add_contact_groups">add_contact_groups</a></h2>

Can be used to add a contact to one or more groups. A [contact_groups_changed](sessions.html#event:contact_groups_changed) event will be created
for the groups which the contact has been added to.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "add_contact_groups",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "groups": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ]
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "contact_groups_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "8ed05195-68cc-47fa-8e78-3bde7b3370ae",
    "groups_added": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ]
}
```
</div>
<h2 class="item_title"><a name="action:add_contact_urn" href="#action:add_contact_urn">add_contact_urn</a></h2>

Can be used to add a URN to the current contact. A [contact_urns_changed](sessions.html#event:contact_urns_changed) event
will be created when this action is encountered.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "add_contact_urn",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "scheme": "tel",
    "path": "@results.phone_number.value"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "contact_urns_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "f3cbd795-9bb3-4331-ba82-c15b24dd577f",
    "urns": [
        "tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d",
        "twitterid:54784326227#nyaruka",
        "mailto:foo@bar.com",
        "tel:+12344563452"
    ]
}
```
</div>
<h2 class="item_title"><a name="action:add_input_labels" href="#action:add_input_labels">add_input_labels</a></h2>

Can be used to add labels to the last user input on a flow. An [input_labels_added](sessions.html#event:input_labels_added) event
will be created with the labels added when this action is encountered. If there is
no user input at that point then this action will be ignored.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "add_input_labels",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "labels": [
        {
            "uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
            "name": "Spam"
        }
    ]
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "input_labels_added",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "f54bbb16-3fba-491d-8768-7d0ba4d4f3e8",
    "input_uuid": "9bf91c2b-ce58-4cef-aacc-281e03f69ab5",
    "labels": [
        {
            "uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
            "name": "Spam"
        }
    ]
}
```
</div>
<h2 class="item_title"><a name="action:call_classifier" href="#action:call_classifier">call_classifier</a></h2>

Can be used to classify the intent and entities from a given input using an NLU classifier. It always
saves a result indicating whether the classification was successful, skipped or failed, and what the extracted intents
and entities were.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "call_classifier",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "classifier": {
        "uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000",
        "name": "Booking"
    },
    "input": "@input.text",
    "result_name": "Intent"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
[
    {
        "type": "classifier_called",
        "created_on": "2018-04-11T18:24:30.123456Z",
        "step_uuid": "8707af30-d50f-440b-9803-f4a851d20f2b",
        "url": "http://test.acme.ai?classifiy",
        "status": "success",
        "request": "GET /message?v=20170307&q=hello HTTP/1.1",
        "response": "HTTP/1.1 200 OK\r\n\r\n{\"intents\":[]}",
        "elapsed_ms": 1,
        "classifier": {
            "uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000",
            "name": "Booking"
        }
    },
    {
        "type": "run_result_changed",
        "created_on": "2018-04-11T18:24:30.123456Z",
        "step_uuid": "8707af30-d50f-440b-9803-f4a851d20f2b",
        "name": "Intent",
        "value": "book_flight",
        "category": "Success",
        "input": "Hi there",
        "extra": {
            "intents": [
                {
                    "name": "book_flight",
                    "confidence": 0.5
                },
                {
                    "name": "book_hotel",
                    "confidence": 0.25
                }
            ],
            "entities": {
                "location": [
                    {
                        "value": "Quito",
                        "confidence": 1
                    }
                ]
            }
        }
    }
]
```
</div>
<h2 class="item_title"><a name="action:call_resthook" href="#action:call_resthook">call_resthook</a></h2>

Can be used to call a resthook.

A [webhook_called](sessions.html#event:webhook_called) event will be created for each subscriber of the resthook with the results
of the HTTP call. If the action has `result_name` set, a result will
be created with that name, and if the resthook returns valid JSON, that will be accessible
through `extra` on the result.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "call_resthook",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "resthook": "new-registration"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
[
    {
        "type": "resthook_called",
        "created_on": "2018-04-11T18:24:30.123456Z",
        "step_uuid": "8e64b588-d46e-4016-a5ef-59cf4d9d7a5b",
        "resthook": "new-registration",
        "payload": {
            "channel": {
                "address": "+12345671111",
                "name": "My Android Phone",
                "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
            },
            "contact": {
                "name": "Ryan Lewis",
                "urn": "tel:+12065551212",
                "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
            },
            "flow": {
                "name": "Registration",
                "revision": 123,
                "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7"
            },
            "input": {
                "attachments": [
                    {
                        "content_type": "image/jpeg",
                        "url": "http://s3.amazon.com/bucket/test.jpg"
                    },
                    {
                        "content_type": "audio/mp3",
                        "url": "http://s3.amazon.com/bucket/test.mp3"
                    }
                ],
                "channel": {
                    "address": "+12345671111",
                    "name": "My Android Phone",
                    "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
                },
                "created_on": "2017-12-31T11:35:10.035757-02:00",
                "text": "Hi there",
                "type": "msg",
                "urn": {
                    "display": "(206) 555-1212",
                    "path": "+12065551212",
                    "scheme": "tel"
                },
                "uuid": "9bf91c2b-ce58-4cef-aacc-281e03f69ab5"
            },
            "path": [
                {
                    "arrived_on": "2018-04-11T18:24:30.123456Z",
                    "exit_uuid": "d7a36118-0a38-4b35-a7e4-ae89042f0d3c",
                    "node_uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
                    "uuid": "5865a06e-6fcc-4db9-bfd7-d22404241e07"
                },
                {
                    "arrived_on": "2018-04-11T18:24:30.123456Z",
                    "exit_uuid": "100f2d68-2481-4137-a0a3-177620ba3c5f",
                    "node_uuid": "3dcccbb4-d29c-41dd-a01f-16d814c9ab82",
                    "uuid": "8f570cd6-588b-4b10-acdc-e181352d44cc"
                },
                {
                    "arrived_on": "2018-04-11T18:24:30.123456Z",
                    "exit_uuid": "d898f9a4-f0fc-4ac4-a639-c98c602bb511",
                    "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                    "uuid": "19ebde80-3969-47d3-a09e-6806aab9f510"
                },
                {
                    "arrived_on": "2018-04-11T18:24:30.123456Z",
                    "exit_uuid": "",
                    "node_uuid": "c0781400-737f-4940-9a6c-1ec1c3df0325",
                    "uuid": "8e64b588-d46e-4016-a5ef-59cf4d9d7a5b"
                }
            ],
            "results": {
                "2factor": {
                    "category": "",
                    "category_localized": "",
                    "created_on": "2018-04-11T18:24:30.123456Z",
                    "input": "",
                    "name": "2Factor",
                    "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                    "value": "34634624463525"
                },
                "favorite_color": {
                    "category": "Red",
                    "category_localized": "Red",
                    "created_on": "2018-04-11T18:24:30.123456Z",
                    "input": "",
                    "name": "Favorite Color",
                    "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                    "value": "red"
                },
                "intent": {
                    "category": "Success",
                    "category_localized": "Success",
                    "created_on": "2018-04-11T18:24:30.123456Z",
                    "input": "Hi there",
                    "name": "intent",
                    "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                    "value": "book_flight"
                },
                "phone_number": {
                    "category": "",
                    "category_localized": "",
                    "created_on": "2018-04-11T18:24:30.123456Z",
                    "input": "",
                    "name": "Phone Number",
                    "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                    "value": "+12344563452"
                },
                "webhook": {
                    "category": "Success",
                    "category_localized": "Success",
                    "created_on": "2018-04-11T18:24:30.123456Z",
                    "input": "GET http://localhost:49998/?content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D",
                    "name": "webhook",
                    "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                    "value": "200"
                }
            },
            "run": {
                "created_on": "2018-04-11T18:24:30.123456Z",
                "uuid": "5fa51f39-76ea-421c-a71b-fe4af29b871a"
            }
        }
    },
    {
        "type": "webhook_called",
        "created_on": "2018-04-11T18:24:30.123456Z",
        "step_uuid": "8e64b588-d46e-4016-a5ef-59cf4d9d7a5b",
        "url": "http://localhost:49998/?cmd=success",
        "status": "success",
        "request": "POST /?cmd=success HTTP/1.1\r\nHost: localhost:49998\r\nUser-Agent: goflow-testing\r\nContent-Length: 2785\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"channel\":{\"address\":\"+12345671111\",\"name\":\"My Android Phone\",\"uuid\":\"57f1078f-88aa-46f4-a59a-948a5739c03d\"},\"contact\":{\"name\":\"Ryan Lewis\",\"urn\":\"tel:+12065551212\",\"uuid\":\"5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f\"},\"flow\":{\"name\":\"Registration\",\"revision\":123,\"uuid\":\"50c3706e-fedb-42c0-8eab-dda3335714b7\"},\"input\":{\"attachments\":[{\"content_type\":\"image/jpeg\",\"url\":\"http://s3.amazon.com/bucket/test.jpg\"},{\"content_type\":\"audio/mp3\",\"url\":\"http://s3.amazon.com/bucket/test.mp3\"}],\"channel\":{\"address\":\"+12345671111\",\"name\":\"My Android Phone\",\"uuid\":\"57f1078f-88aa-46f4-a59a-948a5739c03d\"},\"created_on\":\"2017-12-31T11:35:10.035757-02:00\",\"text\":\"Hi there\",\"type\":\"msg\",\"urn\":{\"display\":\"(206) 555-1212\",\"path\":\"+12065551212\",\"scheme\":\"tel\"},\"uuid\":\"9bf91c2b-ce58-4cef-aacc-281e03f69ab5\"},\"path\":[{\"arrived_on\":\"2018-04-11T18:24:30.123456Z\",\"exit_uuid\":\"d7a36118-0a38-4b35-a7e4-ae89042f0d3c\",\"node_uuid\":\"72a1f5df-49f9-45df-94c9-d86f7ea064e5\",\"uuid\":\"5865a06e-6fcc-4db9-bfd7-d22404241e07\"},{\"arrived_on\":\"2018-04-11T18:24:30.123456Z\",\"exit_uuid\":\"100f2d68-2481-4137-a0a3-177620ba3c5f\",\"node_uuid\":\"3dcccbb4-d29c-41dd-a01f-16d814c9ab82\",\"uuid\":\"8f570cd6-588b-4b10-acdc-e181352d44cc\"},{\"arrived_on\":\"2018-04-11T18:24:30.123456Z\",\"exit_uuid\":\"d898f9a4-f0fc-4ac4-a639-c98c602bb511\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"uuid\":\"19ebde80-3969-47d3-a09e-6806aab9f510\"},{\"arrived_on\":\"2018-04-11T18:24:30.123456Z\",\"exit_uuid\":\"\",\"node_uuid\":\"c0781400-737f-4940-9a6c-1ec1c3df0325\",\"uuid\":\"8e64b588-d46e-4016-a5ef-59cf4d9d7a5b\"}],\"results\":{\"2factor\":{\"category\":\"\",\"category_localized\":\"\",\"created_on\":\"2018-04-11T18:24:30.123456Z\",\"input\":\"\",\"name\":\"2Factor\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"value\":\"34634624463525\"},\"favorite_color\":{\"category\":\"Red\",\"category_localized\":\"Red\",\"created_on\":\"2018-04-11T18:24:30.123456Z\",\"input\":\"\",\"name\":\"Favorite Color\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"value\":\"red\"},\"intent\":{\"category\":\"Success\",\"category_localized\":\"Success\",\"created_on\":\"2018-04-11T18:24:30.123456Z\",\"input\":\"Hi there\",\"name\":\"intent\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"value\":\"book_flight\"},\"phone_number\":{\"category\":\"\",\"category_localized\":\"\",\"created_on\":\"2018-04-11T18:24:30.123456Z\",\"input\":\"\",\"name\":\"Phone Number\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"value\":\"+12344563452\"},\"webhook\":{\"category\":\"Success\",\"category_localized\":\"Success\",\"created_on\":\"2018-04-11T18:24:30.123456Z\",\"input\":\"GET http://localhost:49998/?content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D\",\"name\":\"webhook\",\"node_uuid\":\"f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03\",\"value\":\"200\"}},\"run\":{\"created_on\":\"2018-04-11T18:24:30.123456Z\",\"uuid\":\"5fa51f39-76ea-421c-a71b-fe4af29b871a\"}}",
        "response": "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
        "elapsed_ms": 0,
        "resthook": "new-registration",
        "status_code": 200
    }
]
```
</div>
<h2 class="item_title"><a name="action:call_webhook" href="#action:call_webhook">call_webhook</a></h2>

Can be used to call an external service. The body, header and url fields may be
templates and will be evaluated at runtime. A [webhook_called](sessions.html#event:webhook_called) event will be created based on
the results of the HTTP call. If this action has a `result_name`, then addtionally it will create
a new result with that name. If the webhook returned valid JSON, that will be accessible
through `extra` on the result.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "call_webhook",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "method": "GET",
    "url": "http://localhost/?cmd=success",
    "headers": {
        "Authorization": "Token AAFFZZHH"
    },
    "result_name": "webhook"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
[
    {
        "type": "webhook_called",
        "created_on": "2018-04-11T18:24:30.123456Z",
        "step_uuid": "5861c68a-8201-4d77-9e14-020037f6ddea",
        "url": "http://localhost:49998/?cmd=success",
        "status": "success",
        "request": "GET /?cmd=success HTTP/1.1\r\nHost: localhost:49998\r\nUser-Agent: goflow-testing\r\nAuthorization: Token AAFFZZHH\r\nAccept-Encoding: gzip\r\n\r\n",
        "response": "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
        "elapsed_ms": 0,
        "status_code": 200
    },
    {
        "type": "run_result_changed",
        "created_on": "2018-04-11T18:24:30.123456Z",
        "step_uuid": "5861c68a-8201-4d77-9e14-020037f6ddea",
        "name": "webhook",
        "value": "200",
        "category": "Success",
        "input": "GET http://localhost:49998/?cmd=success",
        "extra": {
            "ok": "true"
        }
    }
]
```
</div>
<h2 class="item_title"><a name="action:enter_flow" href="#action:enter_flow">enter_flow</a></h2>

Can be used to start a contact down another flow. The current flow will pause until the subflow exits or expires.

A [flow_entered](sessions.html#event:flow_entered) event will be created to record that the flow was started.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "enter_flow",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Collect Language"
    }
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "flow_entered",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "dde64b44-09cf-4e6f-a52e-e58736ac73ba",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Collect Language"
    },
    "parent_run_uuid": "e3895066-303a-4b1f-be22-6e6983962829",
    "terminal": false
}
```
</div>
<h2 class="item_title"><a name="action:play_audio" href="#action:play_audio">play_audio</a></h2>

Can be used to play an audio recording in a voice flow. It will generate an
[ivr_created](sessions.html#event:ivr_created) event if there is a valid audio URL. This will contain a message which
the caller should handle as an IVR play command using the audio attachment.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "play_audio",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "audio_url": "http://uploads.temba.io/2353262.m4a"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "ivr_created",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "9972fa41-f437-4bbd-881a-ef06948e0f99",
    "msg": {
        "uuid": "aa863fa2-cb90-435f-802a-9fffea2a27fa",
        "urn": "tel:+12065551212",
        "channel": {
            "uuid": "fd47a886-451b-46fb-bcb6-242a4046c0c0",
            "name": "Nexmo"
        },
        "text": "",
        "attachments": [
            "audio:http://uploads.temba.io/2353262.m4a"
        ]
    }
}
```
</div>
<h2 class="item_title"><a name="action:remove_contact_groups" href="#action:remove_contact_groups">remove_contact_groups</a></h2>

Can be used to remove a contact from one or more groups. A [contact_groups_changed](sessions.html#event:contact_groups_changed) event will be created
for the groups which the contact is removed from. Groups can either be explicitly provided or `all_groups` can be set to true to remove
the contact from all non-dynamic groups.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "remove_contact_groups",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "groups": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Registered Users"
        }
    ]
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "contact_groups_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "7dcc445a-83cf-432b-8188-76dd971a6205",
    "groups_removed": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Testers"
        }
    ]
}
```
</div>
<h2 class="item_title"><a name="action:say_msg" href="#action:say_msg">say_msg</a></h2>

Can be used to communicate with the contact in a voice flow by either reading
a message with TTS or playing a pre-recorded audio file. It will generate an [ivr_created](sessions.html#event:ivr_created)
event if there is a valid audio URL or backdown text. This will contain a message which
the caller should handle as an IVR play command if it has an audio attachment, or otherwise
an IVR say command using the message text.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "say_msg",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "text": "Hi @contact.name, are you ready to complete today's survey?",
    "audio_url": "http://uploads.temba.io/2353262.m4a"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "ivr_created",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "08b826b4-03e0-4528-b379-f167ef86d03e",
    "msg": {
        "uuid": "54e86e97-a008-4b6c-9ea1-dbfb3ec798f1",
        "urn": "tel:+12065551212",
        "channel": {
            "uuid": "fd47a886-451b-46fb-bcb6-242a4046c0c0",
            "name": "Nexmo"
        },
        "text": "Hi Ryan Lewis, are you ready to complete today's survey?",
        "attachments": [
            "audio:http://uploads.temba.io/2353262.m4a"
        ]
    }
}
```
</div>
<h2 class="item_title"><a name="action:send_broadcast" href="#action:send_broadcast">send_broadcast</a></h2>

Can be used to send a message to one or more contacts. It accepts a list of URNs, a list of groups
and a list of contacts.

The URNs and text fields may be templates. A [broadcast_created](sessions.html#event:broadcast_created) event will be created for each unique urn, contact and group
with the evaluated text.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "send_broadcast",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "urns": [
        "tel:+12065551212"
    ],
    "text": "Hi @contact.name, are you ready to complete today's survey?"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "broadcast_created",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "e70444e4-a4a4-407e-a8fa-2f769abd7608",
    "translations": {
        "eng": {
            "text": "Hi Ryan Lewis, are you ready to complete today's survey?"
        }
    },
    "base_language": "eng",
    "urns": [
        "tel:+12065551212"
    ]
}
```
</div>
<h2 class="item_title"><a name="action:send_email" href="#action:send_email">send_email</a></h2>

Can be used to send an email to one or more recipients. The subject, body and addresses
can all contain expressions.

An [email_created](sessions.html#event:email_created) event will be created for each email address.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "send_email",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "addresses": [
        "@urns.mailto"
    ],
    "subject": "Here is your activation token",
    "body": "Your activation token is @contact.fields.activation_token"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "email_created",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "a8ff08ef-6f27-44bd-9029-066bfcb36cf8",
    "addresses": [
        "foo@bar.com"
    ],
    "subject": "Here is your activation token",
    "body": "Your activation token is AACC55"
}
```
</div>
<h2 class="item_title"><a name="action:send_msg" href="#action:send_msg">send_msg</a></h2>

Can be used to reply to the current contact in a flow. The text field may contain templates. The action
will attempt to find pairs of URNs and channels which can be used for sending. If it can't find such a pair, it will
create a message without a channel or URN.

A [msg_created](sessions.html#event:msg_created) event will be created with the evaluated text.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "send_msg",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "text": "Hi @contact.name, are you ready to complete today's survey?",
    "templating": {
        "template": {
            "uuid": "3ce100b7-a734-4b4e-891b-350b1279ade2",
            "name": "revive_issue"
        },
        "variables": [
            "@contact.name"
        ]
    }
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "msg_created",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "e6e30b78-f9c1-462b-9418-6d3e4ae5a100",
    "msg": {
        "uuid": "0df5d5bc-99aa-466a-b715-6b60849cfb2b",
        "urn": "tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d",
        "channel": {
            "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
            "name": "My Android Phone"
        },
        "text": "Hi Ryan Lewis, are you ready to complete today's survey?"
    }
}
```
</div>
<h2 class="item_title"><a name="action:set_contact_channel" href="#action:set_contact_channel">set_contact_channel</a></h2>

Can be used to change or clear the preferred channel of the current contact.

Because channel affinity is a property of a contact's URNs, a [contact_urns_changed](sessions.html#event:contact_urns_changed) event will be created if any
changes are made to the contact's URNs.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "set_contact_channel",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "channel": {
        "uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
        "name": "FAcebook Channel"
    }
}
```
</div><div class="output_event"><h3>Event</h3>

```json
[]
```
</div>
<h2 class="item_title"><a name="action:set_contact_field" href="#action:set_contact_field">set_contact_field</a></h2>

Can be used to update a field value on the contact. The value is a localizable
template and white space is trimmed from the final value. An empty string clears the value.
A [contact_field_changed](sessions.html#event:contact_field_changed) event will be created with the corresponding value.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "set_contact_field",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "field": {
        "key": "gender",
        "name": "Gender"
    },
    "value": "Female"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "contact_field_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "77edcde6-5d7c-4ae9-b660-52c755c4d15c",
    "field": {
        "key": "gender",
        "name": "Gender"
    },
    "value": {
        "text": "Female"
    }
}
```
</div>
<h2 class="item_title"><a name="action:set_contact_language" href="#action:set_contact_language">set_contact_language</a></h2>

Can be used to update the name of the contact. The language is a localizable
template and white space is trimmed from the final value. An empty string clears the language.
A [contact_language_changed](sessions.html#event:contact_language_changed) event will be created with the corresponding value.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "set_contact_language",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "language": "eng"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
[]
```
</div>
<h2 class="item_title"><a name="action:set_contact_name" href="#action:set_contact_name">set_contact_name</a></h2>

Can be used to update the name of the contact. The name is a localizable
template and white space is trimmed from the final value. An empty string clears the name.
A [contact_name_changed](sessions.html#event:contact_name_changed) event will be created with the corresponding value.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "set_contact_name",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "name": "Bob Smith"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "contact_name_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "962f49fb-f56d-40fd-98dc-b94fc84d107e",
    "name": "Bob Smith"
}
```
</div>
<h2 class="item_title"><a name="action:set_contact_timezone" href="#action:set_contact_timezone">set_contact_timezone</a></h2>

Can be used to update the timezone of the contact. The timezone is a localizable
template and white space is trimmed from the final value. An empty string clears the timezone.
A [contact_timezone_changed](sessions.html#event:contact_timezone_changed) event will be created with the corresponding value.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "set_contact_timezone",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "timezone": "Africa/Kigali"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "contact_timezone_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "72f5d1cd-fa2b-4313-b79f-52e1499b6db5",
    "timezone": "Africa/Kigali"
}
```
</div>
<h2 class="item_title"><a name="action:set_run_result" href="#action:set_run_result">set_run_result</a></h2>

Can be used to save a result for a flow. The result will be available in the context
for the run as @results.[name]. The optional category can be used as a way of categorizing results,
this can be useful for reporting or analytics.

Both the value and category fields may be templates. A [run_result_changed](sessions.html#event:run_result_changed) event will be created with the
final values.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "set_run_result",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "name": "Gender",
    "value": "m",
    "category": "Male"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "run_result_changed",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "7dbc8996-c0f0-460b-a23f-a7c1e36a9103",
    "name": "Gender",
    "value": "m",
    "category": "Male"
}
```
</div>
<h2 class="item_title"><a name="action:start_session" href="#action:start_session">start_session</a></h2>

Can be used to trigger sessions for other contacts and groups. A [session_triggered](sessions.html#event:session_triggered) event
will be created and it's the responsibility of the caller to act on that by initiating a new session with the flow engine.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "start_session",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "groups": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ],
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Registration"
    }
}
```
</div><div class="output_event"><h3>Event</h3>

```json
{
    "type": "session_triggered",
    "created_on": "2018-04-11T18:24:30.123456Z",
    "step_uuid": "1348fd9e-c478-42de-b8bf-413ebe9265fa",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Registration"
    },
    "groups": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ],
    "run_summary": {
        "uuid": "66595216-3739-4a5e-a225-4e488c77a340",
        "flow": {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "Registration"
        },
        "contact": {
            "uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
            "id": 1234567,
            "name": "Ryan Lewis",
            "language": "eng",
            "timezone": "America/Guayaquil",
            "created_on": "2018-06-20T11:40:30.123456789Z",
            "urns": [
                "tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d",
                "twitterid:54784326227#nyaruka",
                "mailto:foo@bar.com"
            ],
            "groups": [
                {
                    "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
                    "name": "Testers"
                },
                {
                    "uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9",
                    "name": "Males"
                }
            ],
            "fields": {
                "activation_token": {
                    "text": "AACC55"
                },
                "age": {
                    "text": "23",
                    "number": 23
                },
                "gender": {
                    "text": "Male"
                },
                "join_date": {
                    "text": "2017-12-02",
                    "datetime": "2017-12-02T00:00:00.000000-02:00"
                }
            }
        },
        "status": "active",
        "results": {
            "2factor": {
                "name": "2Factor",
                "value": "34634624463525",
                "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                "created_on": "2018-04-11T18:24:30.123456Z"
            },
            "favorite_color": {
                "name": "Favorite Color",
                "value": "red",
                "category": "Red",
                "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                "created_on": "2018-04-11T18:24:30.123456Z"
            },
            "intent": {
                "name": "intent",
                "value": "book_flight",
                "category": "Success",
                "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                "input": "Hi there",
                "extra": {
                    "intents": [
                        {
                            "name": "book_flight",
                            "confidence": 0.5
                        },
                        {
                            "name": "book_hotel",
                            "confidence": 0.25
                        }
                    ],
                    "entities": {
                        "location": [
                            {
                                "value": "Quito",
                                "confidence": 1
                            }
                        ]
                    }
                },
                "created_on": "2018-04-11T18:24:30.123456Z"
            },
            "phone_number": {
                "name": "Phone Number",
                "value": "+12344563452",
                "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                "created_on": "2018-04-11T18:24:30.123456Z"
            },
            "webhook": {
                "name": "webhook",
                "value": "200",
                "category": "Success",
                "node_uuid": "f5bb9b7a-7b5e-45c3-8f0e-61b4e95edf03",
                "input": "GET http://localhost:49998/?content=%7B%22results%22%3A%5B%7B%22state%22%3A%22WA%22%7D%2C%7B%22state%22%3A%22IN%22%7D%5D%7D",
                "extra": {
                    "results": [
                        {
                            "state": "WA"
                        },
                        {
                            "state": "IN"
                        }
                    ]
                },
                "created_on": "2018-04-11T18:24:30.123456Z"
            }
        }
    }
}
```
</div>
<h2 class="item_title"><a name="action:transfer_airtime" href="#action:transfer_airtime">transfer_airtime</a></h2>

Attempts to make an airtime transfer to the contact.

An [email_created](sessions.html#event:email_created) event will be created for each email address.

<div class="input_action"><h3>Action</h3>

```json
{
    "type": "transfer_airtime",
    "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
    "amounts": {
        "RWF": 500,
        "USD": 0.5
    },
    "result_name": "reward_transfer"
}
```
</div><div class="output_event"><h3>Event</h3>

```json
[
    {
        "type": "airtime_transferred",
        "created_on": "2018-04-11T18:24:30.123456Z",
        "step_uuid": "8fe599f1-bf4f-43be-83dd-3099689741ac",
        "sender": "tel:+12345671111",
        "recipient": "tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d",
        "currency": "RWF",
        "amount": 500,
        "status": "success"
    },
    {
        "type": "run_result_changed",
        "created_on": "2018-04-11T18:24:30.123456Z",
        "step_uuid": "8fe599f1-bf4f-43be-83dd-3099689741ac",
        "name": "reward_transfer",
        "value": "500",
        "category": "Success"
    }
]
```
</div>

</div>
