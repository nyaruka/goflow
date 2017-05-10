package main

import (
	"fmt"
	"os"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/utils"
)

type Vars struct {
}

func (vars *Vars) Resolve(key string) interface{} {
	switch key {
	case "int1":
		return 1
	case "string1":
		return "string1"
	case "int2":
		return 2
	}
	return nil
}

func (vars *Vars) Default() interface{} {
	return nil
}

func main() {
	vars := Vars{}
	env := utils.NewDefaultEnvironment()

	val, err := excellent.EvaluateTemplateAsString(env, &vars, os.Args[1])

	fmt.Printf("Value: %s\n", val)
	if err != nil {
		fmt.Printf("Errors: %s\n", err.Error())
	}

}
