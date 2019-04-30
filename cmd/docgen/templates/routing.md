If a node has a single exit, the engine will pick that when leaving that node. If the node has more than one exit,
then we need a router to choose an exit. 

# Routers

Routers are primarily responsible for picking exits but can also generate events and save results. All routers have 
the following properties:

 * `type` router type
 * `wait` optional [wait](#waits)
 * `result_name` optional result name if router should save a result
 * `categories` possible categories of any result saved by this router

Different router types have different logic for how an exit will be chosen.

## Switch

If a node wishes to route differently based on some state in the session, it can add a `switch` router which defines one or more 
`cases`.  Each case defines a `type` which is the name of an expression function that is run by passing the evaluation of `operand` 
as the first argument. Cases may define additional arguments using the `arguments` array on a case. If no case evaluates 
to true, then the `default_category_uuid` will be used, otherwise flow execution will stop.

A `switch` router has these additional properties:

 * `operand` the template which will be evaluated against each of our cases
 * `cases` a list of 1-n cases which are evaluated in order until one is true
 * `default_category_uuid` the uuid of the default category to take if no case matches (optional)

Each case consists of:

 * `uuid` the UUID
 * `type` the type of this test which is the name of a [test function](#tests) - it will be called with the operand as the first argument
 * `arguments` an optional list of templates which can be passed as extra arguments to the test (after the initial operand)
 * `category_uuid` the uuid of the category that should be taken if this case evaluated to true

The following is an example switch router with 2 cases:

```json
{
    "uuid": "ee0bee3f-34b3-4275-af78-f9ff52c82e6a",
    "router": {
        "type": "switch",
        "categories": [
            {
                "uuid": "cab600f5-b54b-49b9-a7ea-5638f4cbf2b4",
                "name": "Has Name",
                "exit_uuid": "972fb580-54c2-4491-8438-09ace3500ba5"
            },
            {
                "uuid": "9574fbfd-510f-4dfc-b989-97d2aecf50b9",
                "name": "Other",
                "exit_uuid": "6981b1a9-af04-4e26-a248-1fc1f5e5c7eb"
            }
        ],
        "operand": "@input",
        "cases": [
            {
                "uuid": "6f78d564-029b-4715-b8d4-b28daeae4f24",
                "type": "has_text",
                "category_uuid": "cab600f5-b54b-49b9-a7ea-5638f4cbf2b4"
            }
        ],
        "default_category_uuid": "9574fbfd-510f-4dfc-b989-97d2aecf50b9"
    },
    "exits": [
        {
            "uuid": "972fb580-54c2-4491-8438-09ace3500ba5",
            "destination_uuid": "deec1dd4-b727-4b21-800a-0b7bbd146a82"
        },
        {
            "uuid": "6981b1a9-af04-4e26-a248-1fc1f5e5c7eb",
            "destination_uuid": "ee0bee3f-34b3-4275-af78-f9ff52c82e6a"
        }
    ]
}
```

## Random

A random router chooses one of its categories randomly and has no additional properties. For example:

```json
{
    "uuid": "ee0bee3f-34b3-4275-af78-f9ff52c82e6a",
    "router": {
        "type": "random",
        "categories": [
            {
                "uuid": "cab600f5-b54b-49b9-a7ea-5638f4cbf2b4",
                "name": "Bucket 1",
                "exit_uuid": "972fb580-54c2-4491-8438-09ace3500ba5"
            },
            {
                "uuid": "9574fbfd-510f-4dfc-b989-97d2aecf50b9",
                "name": "Bucket 2",
                "exit_uuid": "6981b1a9-af04-4e26-a248-1fc1f5e5c7eb"
            }
        ]
    },
    "exits": [
        {
            "uuid": "972fb580-54c2-4491-8438-09ace3500ba5",
            "destination_uuid": "deec1dd4-b727-4b21-800a-0b7bbd146a82"
        },
        {
            "uuid": "6981b1a9-af04-4e26-a248-1fc1f5e5c7eb",
            "destination_uuid": "ee0bee3f-34b3-4275-af78-f9ff52c82e6a"
        }
    ]
}
```

# Waits

A wait tells the engine to hand back control to the caller and wait for the caller to resume execution by providing something.
The type of the wait indicates what is required to resume flow execution and currently we only support waits of type `msg`.

## Msg

This type indicates that flow execution should pause until an incoming message is received. It can have an optional timeout 
value which is the number of seconds after which execution can be resumed without a message, e.g.

```json
{
    "type": "msg",
    "timeout": 600
}
```

# Tests

Router tests are a special class of functions which are used within the switch router. They are called in the same way as normal functions, but 
all return a test result object which by default evalutes to true or false, but can also be used to find the matching portion of the test by using
the `match` component of the result. The flow editor builds these expressions using UI widgets, but they can be used anywhere a normal template
function is used.

<div class="tests">
{{ .testDocs }}
</div>