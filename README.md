# Goflow [![Build Status](https://travis-ci.org/nyaruka/goflow.svg?branch=master)](https://travis-ci.org/nyaruka/goflow) [![Coverage Status](https://coveralls.io/repos/github/nyaruka/goflow/badge.svg?branch=master)](https://coveralls.io/github/nyaruka/goflow?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/goflow)](https://goreportcard.com/report/github.com/nyaruka/goflow)

## Runner 

This program provides a command line interface for stepping through a given flow.

```
% go install github.com/nyaruka/goflow/cmd/flowrunner
% $GOPATH/bin/flowrunner cmd/flowrunner/testdata/flows/two_questions.json
```

## Server

The server provides an HTTP endpoint for stepping through a given flow:

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

## Running Tests

You can run all the tests (excluding tests in vendor packages) with:

```
% go test $(go list ./... | grep -v /vendor/)
```

If you've made changes to the flow server response format, you should regenerate the test files with:

```
% go test github.com/nyaruka/goflow/cmd/flowrunner -write
```

If you've made changes to the flow server static files, you should regenerate the statik module with:

```
% go generate github.com/nyaruka/goflow/cmd/flowserver
```