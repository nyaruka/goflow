# Goflow [![Build Status](https://travis-ci.org/nyaruka/goflow.svg?branch=master)](https://travis-ci.org/nyaruka/goflow) [![codecov](https://codecov.io/gh/nyaruka/goflow/branch/master/graph/badge.svg)](https://codecov.io/gh/nyaruka/goflow) [![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/goflow)](https://goreportcard.com/report/github.com/nyaruka/goflow)

## Specification

See https://nyaruka.github.io/goflow/ for the complete specification docs.

## Basic Usage

```go
import (
    "github.com/nyaruka/goflow/assets/static"
    "github.com/nyaruka/goflow/flows"
    "github.com/nyaruka/goflow/flows/engine"
    "github.com/nyaruka/goflow/utils"
)

source, _ := static.LoadStaticSource("myassets.json")
assets, _ := engine.NewSessionAssets(source)
session := engine.NewSession(assets, engine.NewDefaultConfig(), utils.NewHTTPClient("goflow-flowrunner"))
contact := flows.NewContact(...)
trigger := triggers.NewManualTrigger(utils.NewDefaultEnvironment(), contact, flow.Reference(), nil, time.Now())
session.Start(trigger)
```

## Runner 

This program provides a command line interface for stepping through a given flow.

```
% go install github.com/nyaruka/goflow/cmd/flowrunner
% $GOPATH/bin/flowrunner cmd/flowrunner/testdata/two_questions.json 615b8a0f-588c-4d20-a05f-363b0b4ce6f4
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

## Server

This server provides an HTTP endpoint for stepping through a given flow:

```
% go install github.com/nyaruka/goflow/cmd/flowserver
% $GOPATH/bin/flowserver
```

## Expression Tester

This utility provides a quick way to test evaluation of expressions which can be used in flows:

```
% go install github.com/nyaruka/goflow/cmd/exptester
% $GOPATH/bin/exptester '@(10 / 5 >= 2)'
% $GOPATH/bin/exptester '@(TITLE("foo"))'
```

## Development

You can run the flow server with detailed output of actions being executed and events being applied with:

```
% $GOPATH/bin/flowserver --log-level=debug
```

You can run all the tests with:

```
% go test github.com/nyaruka/goflow/...
```

If you've made changes to the flow server response format, regenerate the test files with:

```
% go test github.com/nyaruka/goflow/test -write
```

If you've made changes to the flow server static files, you should regenerate the statik module with:

```
% go generate github.com/nyaruka/goflow/cmd/flowserver
```

To make a new release:

```
% git tag v?.?.?
% git push origin v?.?.?
% goreleaser
```
