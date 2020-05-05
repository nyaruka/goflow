# Triggers

Triggers start a new session with the flow engine. They describe why the session is being started and provide parameters which can
be accessed in expressions.

<div class="triggers">
<h2 class="item_title"><a name="trigger:campaign" href="#trigger:campaign">campaign</a></h2>

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
        "name": "Bob",
        "status": "active",
        "created_on": "2018-01-01T12:00:00Z"
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

<h2 class="item_title"><a name="trigger:channel" href="#trigger:channel">channel</a></h2>

Is used when a session was triggered by a channel event


```json
{
    "type": "channel",
    "flow": {
        "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
        "name": "Registration"
    },
    "contact": {
        "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
        "name": "Bob",
        "status": "active",
        "created_on": "2018-01-01T12:00:00Z"
    },
    "triggered_on": "2000-01-01T00:00:00Z",
    "event": {
        "type": "new_conversation",
        "channel": {
            "uuid": "58e9b092-fe42-4173-876c-ff45a14a24fe",
            "name": "Facebook"
        }
    }
}
```

<h2 class="item_title"><a name="trigger:flow_action" href="#trigger:flow_action">flow_action</a></h2>

Is used when another session triggered this run using a trigger_flow action.


```json
{
    "type": "flow_action",
    "flow": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "name": "Collect Age"
    },
    "triggered_on": "2000-01-01T00:00:00Z",
    "run_summary": {
        "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
        "flow": {
            "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
            "name": "Registration"
        },
        "contact": {
            "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
            "name": "Bob",
            "fields": {
                "gender": {
                    "text": "Male"
                }
            },
            "created_on": "2018-01-01T12:00:00.000000000-00:00"
        },
        "status": "active",
        "results": {
            "age": {
                "result_name": "Age",
                "value": "33",
                "node": "cd2be8c4-59bc-453c-8777-dec9a80043b8",
                "created_on": "2018-01-01T12:00:00.000000000-00:00"
            }
        }
    }
}
```

<h2 class="item_title"><a name="trigger:manual" href="#trigger:manual">manual</a></h2>

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
        "name": "Bob",
        "status": "active",
        "created_on": "2018-01-01T12:00:00Z"
    },
    "triggered_on": "2000-01-01T00:00:00Z"
}
```

<h2 class="item_title"><a name="trigger:msg" href="#trigger:msg">msg</a></h2>

Is used when a session was triggered by a message being received by the caller


```json
{
    "type": "msg",
    "flow": {
        "uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
        "name": "Registration"
    },
    "contact": {
        "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
        "name": "Bob",
        "status": "active",
        "created_on": "2018-01-01T12:00:00Z"
    },
    "triggered_on": "2000-01-01T00:00:00Z",
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
    },
    "keyword_match": {
        "type": "first_word",
        "keyword": "start"
    }
}
```


</div>

# Resumes

Resumes resume an existing session with the flow engine and describe why the session is being resumed.

<div class="resumes">
<h2 class="item_title"><a name="resume:msg" href="#resume:msg">msg</a></h2>

Is used when a session is resumed with a new message from the contact


```json
{
    "type": "msg",
    "contact": {
        "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
        "name": "Bob",
        "language": "fra",
        "status": "active",
        "created_on": "2018-01-01T12:00:00Z",
        "fields": {
            "gender": {
                "text": "Male"
            }
        }
    },
    "resumed_on": "2000-01-01T00:00:00Z",
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

<h2 class="item_title"><a name="resume:run_expiration" href="#resume:run_expiration">run_expiration</a></h2>

Is used when a session is resumed because the waiting run has expired


```json
{
    "type": "run_expiration",
    "contact": {
        "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
        "name": "Bob",
        "language": "fra",
        "status": "active",
        "created_on": "2018-01-01T12:00:00Z",
        "fields": {
            "gender": {
                "text": "Male"
            }
        }
    },
    "resumed_on": "2000-01-01T00:00:00Z"
}
```

<h2 class="item_title"><a name="resume:wait_timeout" href="#resume:wait_timeout">wait_timeout</a></h2>

Is used when a session is resumed because a wait has timed out


```json
{
    "type": "wait_timeout",
    "contact": {
        "uuid": "9f7ede93-4b16-4692-80ad-b7dc54a1cd81",
        "name": "Bob",
        "language": "fra",
        "status": "active",
        "created_on": "2018-01-01T12:00:00Z",
        "fields": {
            "gender": {
                "text": "Male"
            }
        }
    },
    "resumed_on": "2000-01-01T00:00:00Z"
}
```


</div>

# Events

Events are the output of a flow run and represent instructions to the engine container on what actions should be taken due to the flow execution.
All templates in events have been evaluated and can be used to create concrete messages, contact updates, emails etc by the container.

<div class="events">
<h2 class="item_title"><a name="event:airtime_transferred" href="#event:airtime_transferred">airtime_transferred</a></h2>

Events are created when airtime has been transferred to the contact.

<div class="output_event">

```json
{
    "type": "airtime_transferred",
    "created_on": "2006-01-02T15:04:05Z",
    "sender": "tel:4748",
    "recipient": "tel:+1242563637",
    "currency": "RWF",
    "desired_amount": 120,
    "actual_amount": 100,
    "http_logs": [
        {
            "url": "https://airtime-api.dtone.com/cgi-bin/shop/topup",
            "status": "success",
            "request": "POST /topup HTTP/1.1\r\n\r\naction=ping",
            "response": "HTTP/1.1 200 OK\r\n\r\ninfo_txt=pong\r\n",
            "created_on": "2006-01-02T15:04:05Z",
            "elapsed_ms": 123
        }
    ]
}
```
</div>
<h2 class="item_title"><a name="event:broadcast_created" href="#event:broadcast_created">broadcast_created</a></h2>

Events are created when an action wants to send a message to other contacts.

<div class="output_event">

```json
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
    "contacts": [
        {
            "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
            "name": "Bob"
        }
    ],
    "urns": [
        "tel:+12065551212"
    ]
}
```
</div>
<h2 class="item_title"><a name="event:contact_field_changed" href="#event:contact_field_changed">contact_field_changed</a></h2>

Events are created when a custom field value of the contact has been changed.
A null values indicates that the field value has been cleared.

<div class="output_event">

```json
{
    "type": "contact_field_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "field": {
        "key": "gender",
        "name": "Gender"
    },
    "value": {
        "text": "Male"
    }
}
```
</div>
<h2 class="item_title"><a name="event:contact_groups_changed" href="#event:contact_groups_changed">contact_groups_changed</a></h2>

Events are created when a contact is added or removed to/from one or more groups.

<div class="output_event">

```json
{
    "type": "contact_groups_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "groups_added": [
        {
            "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
            "name": "Reporters"
        }
    ],
    "groups_removed": [
        {
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        }
    ]
}
```
</div>
<h2 class="item_title"><a name="event:contact_language_changed" href="#event:contact_language_changed">contact_language_changed</a></h2>

Events are created when the language of the contact has been changed.

<div class="output_event">

```json
{
    "type": "contact_language_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "language": "eng"
}
```
</div>
<h2 class="item_title"><a name="event:contact_name_changed" href="#event:contact_name_changed">contact_name_changed</a></h2>

Events are created when the name of the contact has been changed.

<div class="output_event">

```json
{
    "type": "contact_name_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "name": "Bob Smith"
}
```
</div>
<h2 class="item_title"><a name="event:contact_refreshed" href="#event:contact_refreshed">contact_refreshed</a></h2>

Events are generated when the resume has a contact with differences to the current session contact.

<div class="output_event">

```json
{
    "type": "contact_refreshed",
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
<h2 class="item_title"><a name="event:contact_status_changed" href="#event:contact_status_changed">contact_status_changed</a></h2>

Events are created when the status of the contact has been changed.

<div class="output_event">

```json
{
    "type": "contact_timezone_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "timezone": ""
}
```
</div>
<h2 class="item_title"><a name="event:contact_timezone_changed" href="#event:contact_timezone_changed">contact_timezone_changed</a></h2>

Events are created when the timezone of the contact has been changed.

<div class="output_event">

```json
{
    "type": "contact_timezone_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "timezone": "Africa/Kigali"
}
```
</div>
<h2 class="item_title"><a name="event:contact_urns_changed" href="#event:contact_urns_changed">contact_urns_changed</a></h2>

Events are created when a contact's URNs have changed.

<div class="output_event">

```json
{
    "type": "contact_urns_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "urns": [
        "tel:+12345678900",
        "twitter:bob"
    ]
}
```
</div>
<h2 class="item_title"><a name="event:email_sent" href="#event:email_sent">email_sent</a></h2>

Events are created when an action has sent an email.

<div class="output_event">

```json
{
    "type": "email_sent",
    "created_on": "2006-01-02T15:04:05Z",
    "to": [
        "foo@bar.com"
    ],
    "subject": "Your activation token",
    "body": "Your activation token is AAFFKKEE"
}
```
</div>
<h2 class="item_title"><a name="event:environment_refreshed" href="#event:environment_refreshed">environment_refreshed</a></h2>

Events are sent by the caller to tell the engine to update the session environment.

<div class="output_event">

```json
{
    "type": "environment_refreshed",
    "created_on": "2006-01-02T15:04:05Z",
    "environment": {
        "date_format": "YYYY-MM-DD",
        "time_format": "hh:mm",
        "timezone": "Africa/Kigali",
        "default_language": "eng",
        "allowed_languages": [
            "eng",
            "fra"
        ]
    }
}
```
</div>
<h2 class="item_title"><a name="event:error" href="#event:error">error</a></h2>

Events are created when an error occurs during flow execution.

<div class="output_event">

```json
{
    "type": "error",
    "created_on": "2006-01-02T15:04:05Z",
    "text": "invalid date format: '12th of October'"
}
```
</div>
<h2 class="item_title"><a name="event:failure" href="#event:failure">failure</a></h2>

Events are created when an error occurs during flow execution which prevents continuation of the session.

<div class="output_event">

```json
{
    "type": "failure",
    "created_on": "2006-01-02T15:04:05Z",
    "text": "unable to read flow"
}
```
</div>
<h2 class="item_title"><a name="event:flow_entered" href="#event:flow_entered">flow_entered</a></h2>

Events are created when an action has entered a sub-flow.

<div class="output_event">

```json
{
    "type": "flow_entered",
    "created_on": "2006-01-02T15:04:05Z",
    "flow": {
        "uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a",
        "name": "Registration"
    },
    "parent_run_uuid": "95eb96df-461b-4668-b168-727f8ceb13dd",
    "terminal": false
}
```
</div>
<h2 class="item_title"><a name="event:input_labels_added" href="#event:input_labels_added">input_labels_added</a></h2>

Events are created when an action wants to add labels to the current input.

<div class="output_event">

```json
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
<h2 class="item_title"><a name="event:ivr_created" href="#event:ivr_created">ivr_created</a></h2>

Events are created when an action wants to send an IVR response to the current contact.

<div class="output_event">

```json
{
    "type": "ivr_created",
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
            "audio:https://s3.amazon.com/mybucket/attachment.m4a"
        ]
    }
}
```
</div>
<h2 class="item_title"><a name="event:msg_created" href="#event:msg_created">msg_created</a></h2>

Events are created when an action wants to send a reply to the current contact.

<div class="output_event">

```json
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
            "image/jpeg:https://s3.amazon.com/mybucket/attachment.jpg"
        ]
    }
}
```
</div>
<h2 class="item_title"><a name="event:msg_received" href="#event:msg_received">msg_received</a></h2>

Events are sent by the caller to tell the engine that a message was received from
the contact and that it should try to resume the session.

<div class="output_event">

```json
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
<h2 class="item_title"><a name="event:msg_wait" href="#event:msg_wait">msg_wait</a></h2>

Events are created when a flow pauses waiting for a response from
a contact. If a timeout is set, then the caller should resume the flow after
the number of seconds in the timeout to resume it.

<div class="output_event">

```json
{
    "type": "msg_wait",
    "created_on": "2019-01-02T15:04:05Z",
    "timeout_seconds": 300,
    "hint": {
        "type": "image"
    }
}
```
</div>
<h2 class="item_title"><a name="event:resthook_called" href="#event:resthook_called">resthook_called</a></h2>

Events are created when a resthook is called. The event contains
the payload that will be sent to any subscribers of that resthook. Note that this event is
created regardless of whether there any subscriberes for that resthook.

<div class="output_event">

```json
{
    "type": "resthook_called",
    "created_on": "2006-01-02T15:04:05Z",
    "resthook": "success",
    "payload": {
        "contact:": {
            "name": "Bob"
        }
    }
}
```
</div>
<h2 class="item_title"><a name="event:run_expired" href="#event:run_expired">run_expired</a></h2>

Events are sent by the caller to tell the engine that a run has expired.

<div class="output_event">

```json
{
    "type": "run_expired",
    "created_on": "2006-01-02T15:04:05Z",
    "run_uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a"
}
```
</div>
<h2 class="item_title"><a name="event:run_result_changed" href="#event:run_result_changed">run_result_changed</a></h2>

Events are created when a run result is saved. They contain not only
the name, value and category of the result, but also the UUID of the node where
the result was generated.

<div class="output_event">

```json
{
    "type": "run_result_changed",
    "created_on": "2006-01-02T15:04:05Z",
    "name": "Gender",
    "value": "m",
    "category": "Male",
    "category_localized": "Homme",
    "input": "M"
}
```
</div>
<h2 class="item_title"><a name="event:service_called" href="#event:service_called">service_called</a></h2>

Events are created when an engine service is called.

<div class="output_event">

```json
{
    "type": "service_called",
    "created_on": "2006-01-02T15:04:05Z",
    "service": "classifier",
    "classifier": {
        "uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000",
        "name": "Booking"
    },
    "http_logs": [
        {
            "url": "https://api.wit.ai/message?v=20170307&q=hello",
            "status": "success",
            "request": "GET /message?v=20170307&q=hello HTTP/1.1",
            "response": "HTTP/1.1 200 OK\r\n\r\n{\"intents\":[]}",
            "created_on": "2006-01-02T15:04:05Z",
            "elapsed_ms": 123
        }
    ]
}
```
</div>
<h2 class="item_title"><a name="event:session_triggered" href="#event:session_triggered">session_triggered</a></h2>

Events are created when an action wants to start other people in a flow.

<div class="output_event">

```json
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
    "run_summary": {
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
                "name": "Age",
                "value": "33",
                "node_uuid": "cd2be8c4-59bc-453c-8777-dec9a80043b8",
                "created_on": "2000-01-01T00:00:00.000000000-00:00"
            }
        }
    }
}
```
</div>
<h2 class="item_title"><a name="event:ticket_opened" href="#event:ticket_opened">ticket_opened</a></h2>

Events are created when a new ticket is opened.

<div class="output_event">

```json
{
    "type": "ticket_opened",
    "created_on": "2006-01-02T15:04:05Z",
    "ticket": {
        "uuid": "2e677ae6-9b57-423c-b022-7950503eef35",
        "ticketer": {
            "uuid": "d605bb96-258d-4097-ad0a-080937db2212",
            "name": "Support Tickets"
        },
        "subject": "Need help",
        "body": "Where are my cookies?",
        "external_id": "32526523"
    }
}
```
</div>
<h2 class="item_title"><a name="event:wait_timed_out" href="#event:wait_timed_out">wait_timed_out</a></h2>

Events are sent by the caller when a wait has timed out - i.e. they are sent instead of
the item that the wait was waiting for.

<div class="output_event">

```json
{
    "type": "wait_timed_out",
    "created_on": "2006-01-02T15:04:05Z"
}
```
</div>
<h2 class="item_title"><a name="event:webhook_called" href="#event:webhook_called">webhook_called</a></h2>

Events are created when a webhook is called. The event contains
the URL and the status of the response, as well as a full dump of the
request and response.

<div class="output_event">

```json
{
    "type": "webhook_called",
    "created_on": "2006-01-02T15:04:05Z",
    "url": "http://localhost:49998/?cmd=success",
    "status": "success",
    "request": "GET /?format=json HTTP/1.1",
    "response": "HTTP/1.1 200 OK\r\n\r\n{\"ip\":\"190.154.48.130\"}",
    "elapsed_ms": 123,
    "status_code": 200
}
```
</div>

</div>
