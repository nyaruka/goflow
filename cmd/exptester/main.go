package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

func main() {
	vars := types.NewXMap(map[string]types.XValue{
		"int1":    types.NewXNumberFromInt(1),
		"int2":    types.NewXNumberFromInt(2),
		"string1": types.NewXString("string1"),
	})

	if len(os.Args) != 2 {
		log.Fatal("usage: exptester <expression>")
	}

	env := utils.NewDefaultEnvironment()

	val, err := excellent.EvaluateTemplateAsString(env, vars, os.Args[1], false, nil)

	fmt.Printf("Value: %s\n", val)
	if err != nil {
		fmt.Printf("Errors: %s\n", err.Error())
	}
}
