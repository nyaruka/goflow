# goflow [![Build Status](https://travis-ci.org/nyaruka/goflow.svg?branch=master)](https://travis-ci.org/nyaruka/goflow) [![Coverage Status](https://coveralls.io/repos/github/nyaruka/goflow/badge.svg?branch=master)](https://coveralls.io/github/nyaruka/goflow?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/goflow)](https://goreportcard.com/report/github.com/nyaruka/goflow)

## runner

```
% go install github.com/nyaruka/goflow/cmd/flowrunner
% $GOPATH/bin/flowrunner cmd/flowrunner/testdata/flows/two_questions.json
```

## server

```
% go install github.com/nyaruka/goflow/cmd/flowserver
% $GOPATH/bin/flowserver
```

## expression tester

```
% go install github.com/nyaruka/goflow/cmd/exptester
% $GOPATH/bin/exptester '@(10 / 5 >= 2)'
% $GOPATH/bin/exptester '@(TITLE("foo"))'
```

## running tests

```
% go test github.com/nyaruka/goflow/...
```
