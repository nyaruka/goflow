Assets are objects which can be referenced in a flow definition or flow definitions themselves. For example 
the [set_contact_field](flows.html#action:set_contact_field) action references a field asset.

# Types

<div class="assets">
<a name="asset:channel"></a>

## Channel

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

<a name="asset:field"></a>

## Field

Is a custom contact property.


```objectivec
{
    "key": "gender",
    "name": "Gender",
    "type": "text"
}
```

<a name="asset:flow"></a>

## Flow

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

<a name="asset:group"></a>

## Group

Is a set of contacts which can be static or dynamic (i.e. based on a query).


```objectivec
{
    "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
    "name": "Youth",
    "query": "age <= 18"
}
```

<a name="asset:label"></a>

## Label

Is an organizational tag that can be applied to a message.


```objectivec
{
    "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
    "name": "Spam"
}
```

<a name="asset:location"></a>

## Location

Is a searchable hierachy of locations.


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

<a name="asset:resthook"></a>

## Resthook

Is a set of URLs which are subscribed to the named event.


```objectivec
{
    "slug": "new-registration",
    "subscribers": [
        "http://example.com/record.php?@contact.uuid"
    ]
}
```

<a name="asset:template"></a>

## Template

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


</div>