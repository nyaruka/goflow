# goflow 

## runner

```
% go install github.com/nyaruka/goflow/cmd/flowrunner
% $GOPATH/bin/flowrunner sample.json
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
