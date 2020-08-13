Assets are objects which can be referenced in a flow definition or flow definitions themselves. For example 
the [set_contact_field](flows.html#action:set_contact_field) action references a field asset.

# Types

<div class="assets">
<h2 class="item_title"><a name="asset:channel" href="#asset:channel">channel</a></h2>

Is something that can send/receive messages.


```objectivec
{
    "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
    "name": "My Android",
    "address": "+593979011111",
    "schemes": [
        "tel"
    ],
    "roles": [
        "send",
        "receive"
    ],
    "country": "EC"
}
```

<h2 class="item_title"><a name="asset:classifier" href="#asset:classifier">classifier</a></h2>

Is an NLU classifier.


```objectivec
{
    "uuid": "37657cf7-5eab-4286-9cb0-bbf270587bad",
    "name": "Booking",
    "type": "wit",
    "intents": [
        "book_flight",
        "book_hotel"
    ]
}
```

<h2 class="item_title"><a name="asset:field" href="#asset:field">field</a></h2>

Is a custom contact property.


```objectivec
{
    "uuid": "d66a7823-eada-40e5-9a3a-57239d4690bf",
    "key": "gender",
    "name": "Gender",
    "type": "text"
}
```

<h2 class="item_title"><a name="asset:flow" href="#asset:flow">flow</a></h2>

Is graph of nodes with actions and routers.


```objectivec
{
    "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
    "name": "Registration",
    "definition": {
        "nodes": []
    }
}
```

<h2 class="item_title"><a name="asset:global" href="#asset:global">global</a></h2>

Is a named constant.


```objectivec
{
    "key": "organization_name",
    "name": "Organization Name",
    "value": "U-Report"
}
```

<h2 class="item_title"><a name="asset:group" href="#asset:group">group</a></h2>

Is a set of contacts which can be static or dynamic (i.e. based on a query).


```objectivec
{
    "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
    "name": "Youth",
    "query": "age <= 18"
}
```

<h2 class="item_title"><a name="asset:label" href="#asset:label">label</a></h2>

Is an organizational tag that can be applied to a message.


```objectivec
{
    "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
    "name": "Spam"
}
```

<h2 class="item_title"><a name="asset:location" href="#asset:location">location</a></h2>

Is a searchable hierarchy of locations.


```objectivec
{
    "name": "Rwanda",
    "aliases": [
        "Ruanda"
    ],
    "children": [
        {
            "name": "Kigali City",
            "aliases": [
                "Kigali",
                "Kigari"
            ],
            "children": [
                {
                    "name": "Gasabo",
                    "children": [
                        {
                            "id": "575743222",
                            "name": "Gisozi"
                        },
                        {
                            "id": "457378732",
                            "name": "Ndera"
                        }
                    ]
                },
                {
                    "name": "Nyarugenge",
                    "children": []
                }
            ]
        },
        {
            "name": "Eastern Province"
        }
    ]
}
```

<h2 class="item_title"><a name="asset:resthook" href="#asset:resthook">resthook</a></h2>

Is a set of URLs which are subscribed to the named event.


```objectivec
{
    "slug": "new-registration",
    "subscribers": [
        "http://example.com/record.php?@contact.uuid"
    ]
}
```

<h2 class="item_title"><a name="asset:template" href="#asset:template">template</a></h2>

Is a message template, currently only used by WhatsApp channels


```objectivec
{
    "name": "revive-issue",
    "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
    "translations": [
        {
            "language": "eng",
            "content": "Hi {{1}}, are you still experiencing your issue?",
            "channel": {
                "uuid": "cf26be4c-875f-4094-9e08-162c3c9dcb5b",
                "name": "Twilio Channel"
            }
        },
        {
            "language": "fra",
            "content": "Bonjour {{1}}",
            "channel": {
                "uuid": "cf26be4c-875f-4094-9e08-162c3c9dcb5b",
                "name": "Twilio Channel"
            }
        }
    ]
}
```

<h2 class="item_title"><a name="asset:ticketer" href="#asset:ticketer">ticketer</a></h2>

Is a system which can open or close tickets


```objectivec
{
    "uuid": "37657cf7-5eab-4286-9cb0-bbf270587bad",
    "name": "Support Tickets",
    "type": "mailgun"
}
```


</div>