# Goflow [![Build Status](https://github.com/nyaruka/goflow/workflows/CI/badge.svg)](https://github.com/nyaruka/goflow/actions?query=workflow%3ACI) [![codecov](https://codecov.io/gh/nyaruka/goflow/branch/master/graph/badge.svg)](https://codecov.io/gh/nyaruka/goflow) [![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/goflow)](https://goreportcard.com/report/github.com/nyaruka/goflow)

## Specification

See [here](https://nyaruka.github.io/goflow/en_US/) for the complete specification docs.

## Basic Usage

```go
import (
    "github.com/nyaruka/goflow/assets/static"
    "github.com/nyaruka/goflow/flows"
    "github.com/nyaruka/goflow/flows/engine"
    "github.com/nyaruka/goflow/utils"
)

env := envs.NewBuilder().Build()
source, _ := static.LoadSource("myassets.json")
assets, _ := engine.NewSessionAssets(env, source, nil)
contact := flows.NewContact(assets, ...)
trigger := triggers.NewManual(env, contact, flow.Reference(), nil, false, nil, time.Now())
eng := engine.NewBuilder().Build()
session, sprint, err := eng.NewSession(assets, trigger)
```

## Sessions

Sessions can be persisted between waits by calling `json.Marshal` on the `Session` instance to marshal it as JSON. You can inspect this JSON at https://sessions.temba.io/.

## Utilities

### Flow Runner 

Provides a command line interface for stepping through a given flow.

```
% go install github.com/nyaruka/goflow/cmd/flowrunner
% $GOPATH/bin/flowrunner test/testdata/runner/two_questions.json 615b8a0f-588c-4d20-a05f-363b0b4ce6f4
Starting flow 'U-Report Registration Flow'....
---------------------------------------
💬 "Hi Ben Haggerty! What is your favorite color? (red/blue) Your number is (206) 555-1212"
⏳ waiting for message....
```

By default it will use a manual trigger to create a session, but the `-msg` flag can be used
to start the session with a message trigger:

```
% $GOPATH/bin/flowrunner -msg "hi there" cmd/flowrunner/testdata/two_questions.json 615b8a0f-588c-4d20-a05f-363b0b4ce6f4
```

If the `-repro` flag is set, it will dump the triggers and resumes it used which can be used to reproduce the session in a test:

```
% $GOPATH/bin/flowrunner -repro cmd/flowrunner/testdata/two_questions.json 615b8a0f-588c-4d20-a05f-363b0b4ce6f4
```

### Flow Migrator

Takes a legacy flow definition as piped input and outputs the migrated definition:

```
% go install github.com/nyaruka/goflow/cmd/flowmigrate
% cat legacy_flow.json | $GOPATH/bin/flowmigrate
% cat legacy_export.json | jq '.flows[0]' | $GOPATH/bin/flowmigrate
```

### Expression Tester

Provides a quick way to test evaluation of expressions which can be used in flows:

```
% go install github.com/nyaruka/goflow/cmd/exptester
% $GOPATH/bin/exptester '@(10 / 5 >= 2)'
% $GOPATH/bin/exptester '@(TITLE("foo"))'
```

## Development

You can run all the tests with:

```
% go test github.com/nyaruka/goflow/...
```

If you've made changes to the flow engine output, regenerate the test files with:

```
% go test github.com/nyaruka/goflow/test -update
```
