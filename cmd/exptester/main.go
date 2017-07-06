package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/utils"
)

func main() {
	vars := make(map[string]interface{})
	vars["int1"] = 1
	vars["string1"] = "string1"
	vars["int2"] = 2

	if len(os.Args) != 2 {
		log.Fatal("usage: exptester <expression>")
	}

	env := utils.NewDefaultEnvironment()

	val, err := excellent.EvaluateTemplateAsString(env, utils.NewMapResolver(vars), os.Args[1])

	fmt.Printf("Value: %s\n", val)
	if err != nil {
		fmt.Printf("Errors: %s\n", err.Error())
	}

}
