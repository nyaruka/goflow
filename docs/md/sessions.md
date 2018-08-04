# Triggers

Triggers start a new session with the flow engine. They describe why the session is being started and provide parameters which can
be accessed in expressions.

<div class="triggers">
<a name="trigger:campaign"></a>

## campaign

Is used when a session was triggered by a campaign event


```json
{
    "type": "campaign",
    "flow": {
        "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
        "name": "Registration"
    },
    "contact": {
        "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
        "id": 0,
        "name": "Bob",
        "language": "",
        "timezone": "",
        "created_on": "0001-01-01T00:00:00Z",
        "urns": []
    },
    "triggered_on": "2000-01-01T00:00:00Z",
    "event": {
        "uuid": "34d16dbd-476d-4b77-bac3-9f3d597848cc",
        "campaign": {
            "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
            "name": "New Mothers"
        }
    }
}
```

<a name="trigger:flow_action"></a>

## flow_action

Is used when another session triggered this run using a trigger_flow action.


```json
{
    "type": "flow_action",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Collect Age"
    },
    "triggered_on": "2000-01-01T00:00:00Z",
    "run": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "flow": {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "Registration"
        },
        "contact": {
            "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
            "id": 0,
            "name": "Bob",
            "language": "",
            "timezone": "",
            "created_on": "0001-01-01T00:00:00Z",
            "urns": []
        },
        "status": "active",
        "results": {
            "age": {
                "name": "",
                "value": "33",
                "node_uuid": "",
                "created_on": "2000-01-01T00:00:00Z"
            }
        }
    }
}
```

<a name="trigger:manual"></a>

## manual

Is used when a session was triggered manually by a user


```json
{
    "type": "manual",
    "flow": {
        "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
        "name": "Registration"
    },
    "contact": {
        "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
        "id": 0,
        "name": "Bob",
        "language": "",
        "timezone": "",
        "created_on": "0001-01-01T00:00:00Z",
        "urns": []
    },
    "triggered_on": "2000-01-01T00:00:00Z"
}
```


</div>

# Events

Events are the output of a flow run and represent instructions to the engine container on what actions should be taken due to the flow execution.
All templates in events have been evaluated and can be used to create concrete messages, contact updates, emails etc by the container.

<div class="events">
<a name="event:broadcast_created"></a>

## broadcast_created

Events are created for outgoing broadcasts.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "broadcast_created",
    "created_on": "2006-01-02T15:04:05Z",
    "translations": {
        "eng": {
            "text": "hi, what's up",
            "quick_replies": [
                "All good",
                "Got 99 problems"
            ]
        },
        "spa": {
            "text": "Que pasa",
            "quick_replies": [
                "Todo bien",
                "Tengo 99 problemas"
            ]
        }
    },
    "base_language": "eng",
    "urns": [
        "tel:+12065551212"
    ],
    "contacts": [
        {
            "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
            "name": "Bob"
        }
    ]
}
```
</div>
<a name="event:contact_changed"></a>

## contact_changed

Events are created to set a contact on a session

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "contact": {
        "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
        "name": "Bob",
        "urns": [
            "tel:+11231234567"
        ]
    }
}
```
</div>
<a name="event:contact_channel_changed"></a>

## contact_channel_changed

Events are created when a contact's preferred channel is changed.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_channel_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "channel": {
        "uuid": "67a3ac69-e5e0-4ef0-8423-eddf71a71472",
        "name": "Twilio"
    }
}
```
</div>
<a name="event:contact_field_changed"></a>

## contact_field_changed

Events are created when a contact field is updated.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_field_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "field": {
        "key": "gender",
        "name": "Gender"
    },
    "value": "Male"
}
```
</div>
<a name="event:contact_groups_added"></a>

## contact_groups_added

Events will be created with the groups a contact was added to.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_groups_added",
    "created_on": "2006-01-02T15:04:05Z",
    "groups": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Reporters"
        }
    ]
}
```
</div>
<a name="event:contact_groups_removed"></a>

## contact_groups_removed

Events are created when a contact has been removed from one or more
groups.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_groups_removed",
    "created_on": "2006-01-02T15:04:05Z",
    "groups": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Reporters"
        }
    ]
}
```
</div>
<a name="event:contact_language_changed"></a>

## contact_language_changed

Events are created when a Language of a contact has been changed

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_language_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "language": "eng"
}
```
</div>
<a name="event:contact_name_changed"></a>

## contact_name_changed

Events are created when a name of a contact has been changed

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_name_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "name": "Bob Smith"
}
```
</div>
<a name="event:contact_timezone_changed"></a>

## contact_timezone_changed

Events are created when a timezone of a contact has been changed

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_timezone_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "timezone": "Africa/Kigali"
}
```
</div>
<a name="event:contact_urn_added"></a>

## contact_urn_added

Events will be created with the URN that should be added to the current contact.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "contact_urn_added",
    "created_on": "2006-01-02T15:04:05Z",
    "urn": "tel:+12345678900"
}
```
</div>
<a name="event:email_created"></a>

## email_created

Events are created for each recipient which should receive an email.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "email_created",
    "created_on": "2006-01-02T15:04:05Z",
    "addresses": [
        "foo@bar.com"
    ],
    "subject": "Your activation token",
    "body": "Your activation token is AAFFKKEE"
}
```
</div>
<a name="event:environment_changed"></a>

## environment_changed

Events are created to set the environment on a session

<div class="output_event"><h3>Event</h3>```json
{
    "type": "environment_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "environment": {
        "date_format": "YYYY-MM-DD",
        "time_format": "hh:mm",
        "timezone": "Africa/Kigali",
        "languages": [
            "eng",
            "fra"
        ]
    }
}
```
</div>
<a name="event:error"></a>

## error

Events will be created whenever an error is encountered during flow execution. This
can vary from template evaluation errors to invalid actions.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "error",
    "created_on": "2006-01-02T15:04:05Z",
    "text": "invalid date format: '12th of October'",
    "fatal": false
}
```
</div>
<a name="event:flow_triggered"></a>

## flow_triggered

Events are created when an action wants to start a subflow

<div class="output_event"><h3>Event</h3>```json
{
    "type": "flow_triggered",
    "created_on": "2006-01-02T15:04:05Z",
    "flow": {
        "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
        "name": "Registration"
    },
    "parent_run_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd"
}
```
</div>
<a name="event:input_labels_added"></a>

## input_labels_added

Events will be created with the labels that were applied to the given input.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "input_labels_added",
    "created_on": "2006-01-02T15:04:05Z",
    "input_uuid": "4aef4050-1895-4c80-999a-70368317a4f5",
    "labels": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Spam"
        }
    ]
}
```
</div>
<a name="event:msg_created"></a>

## msg_created

Events are used for replies to the session contact.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "msg_created",
    "created_on": "2006-01-02T15:04:05Z",
    "msg": {
        "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
        "urn": "tel:+12065551212",
        "channel": {
            "uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf",
            "name": "Twilio"
        },
        "text": "hi there",
        "attachments": [
            "https://s3.amazon.com/mybucket/attachment.jpg"
        ]
    }
}
```
</div>
<a name="event:msg_received"></a>

## msg_received

Events are used for starting flows or resuming flows which are waiting for a message.
They represent an incoming message for a contact.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "msg_received",
    "created_on": "2006-01-02T15:04:05Z",
    "msg": {
        "uuid": "2d611e17-fb22-457f-b802-b8f7ec5cda5b",
        "urn": "tel:+12065551212",
        "channel": {
            "uuid": "61602f3e-f603-4c70-8a8f-c477505bf4bf",
            "name": "Twilio"
        },
        "text": "hi there",
        "attachments": [
            "https://s3.amazon.com/mybucket/attachment.jpg"
        ]
    }
}
```
</div>
<a name="event:msg_wait"></a>

## msg_wait

Events are created when a flow pauses waiting for a response from
a contact. If a timeout is set, then the caller should resume the flow after
the number of seconds in the timeout to resume it.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "msg_wait",
    "created_on": "2006-01-02T15:04:05Z"
}
```
</div>
<a name="event:nothing_wait"></a>

## nothing_wait

Events are created when a flow requests to hand back control to the caller but isn't
waiting for anything from the caller.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "nothing_wait",
    "created_on": "2006-01-02T15:04:05.234532Z"
}
```
</div>
<a name="event:resthook_called"></a>

## resthook_called

Events are created when a resthook is called. The event contains the status and status code
of each call to the resthook's subscribers, as well as the payload sent to each subscriber. Applying this event
updates @run.webhook in the context to the results of the last subscriber call. However if one of the subscriber
calls fails, then it is used to update @run.webhook instead.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "resthook_called",
    "created_on": "2006-01-02T15:04:05Z",
    "resthook": "new-registration",
    "payload": "{...}",
    "calls": [
        {
            "url": "http://localhost:49998/?cmd=success",
            "status": "success",
            "status_code": 200,
            "response": "{\"errors\":[]}"
        },
        {
            "url": "https://api.ipify.org?format=json",
            "status": "success",
            "status_code": 410,
            "response": "{\"errors\":[\"Unsubscribe\"]}"
        }
    ]
}
```
</div>
<a name="event:run_expired"></a>

## run_expired

Events are sent by the caller to notify the engine that a run has expired

<div class="output_event"><h3>Event</h3>```json
{
    "type": "run_expired",
    "created_on": "2006-01-02T15:04:05Z",
    "run_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a"
}
```
</div>
<a name="event:run_result_changed"></a>

## run_result_changed

Events are created when a result is saved. They contain not only
the name, value and category of the result, but also the UUID of the node where
the result was generated.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "run_result_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "name": "Gender",
    "value": "m",
    "category": "Male",
    "category_localized": "Homme",
    "node_uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
    "input": "M"
}
```
</div>
<a name="event:session_triggered"></a>

## session_triggered

Events are created when an action wants to start a subflow

<div class="output_event"><h3>Event</h3>```json
{
    "type": "session_triggered",
    "created_on": "2006-01-02T15:04:05Z",
    "flow": {
        "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
        "name": "Registration"
    },
    "groups": [
        {
            "uuid": "8f8e2cae-3c8d-4dce-9c4b-19514437e427",
            "name": "New contacts"
        }
    ],
    "run": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "flow": {
            "uuid": "93c554a1-b90d-4892-b029-a2a87dec9b87",
            "name": "Other Flow"
        },
        "contact": {
            "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
            "name": "Bob",
            "fields": {
                "state": {
                    "value": "Azuay",
                    "created_on": "2000-01-01T00:00:00.000000000-00:00"
                }
            }
        },
        "results": {
            "age": {
                "result_name": "Age",
                "value": "33",
                "node": "cd2be8c4-59bc-453c-8777-dec9a80043b8",
                "created_on": "2000-01-01T00:00:00.000000000-00:00"
            }
        }
    }
}
```
</div>
<a name="event:wait_timed_out"></a>

## wait_timed_out

Events are sent by the caller when a wait has timed out - i.e. they are sent instead of
the item that the wait was waiting for

<div class="output_event"><h3>Event</h3>```json
{
    "type": "wait_timed_out",
    "created_on": "2006-01-02T15:04:05Z"
}
```
</div>
<a name="event:webhook_called"></a>

## webhook_called

Events are created when a webhook is called. The event contains
the status and status code of the response, as well as a full dump of the
request and response. Applying this event updates @run.webhook in the context.

<div class="output_event"><h3>Event</h3>```json
{
    "type": "webhook_called",
    "created_on": "2006-01-02T15:04:05Z",
    "url": "https://api.ipify.org?format=json",
    "status": "success",
    "status_code": 200,
    "request": "GET https://api.ipify.org?format=json",
    "response": "HTTP/1.1 200 OK {\"ip\":\"190.154.48.130\"}"
}
```
</div>

</div>
