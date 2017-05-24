package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func handleEventDoc(prefix string, typeName string, docString string) {
	lines := strings.Split(docString, "\n")
	name := ""

	docs := make([]string, 0, len(lines))
	example := make([]string, 0, len(lines))
	inExample := false
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			name = l[len(prefix)+1:]
		} else if strings.HasPrefix(l, "```") {
			inExample = !inExample
		} else if inExample {
			example = append(example, l[2:])
		} else {
			docs = append(docs, l)
		}
	}

	// try to parse our example
	exampleJSON := []byte(strings.Join(example, "\n"))
	typed := &utils.TypedEnvelope{}
	err := json.Unmarshal(exampleJSON, typed)
	if err != nil {
		log.Fatalf("unable to parse example: %s\nHas err:", exampleJSON, err)
	}

	event, err := events.EventFromEnvelope(typed)
	if err != nil {
		log.Fatalf("unable to parse example: %s\nHas err:", exampleJSON, err)
	}

	// make sure types match
	if name != event.Type() {
		log.Fatalf("Mismatched types for example of %s", name)
	}

	// validate it
	err = utils.ValidateAll(event)
	if err != nil {
		log.Fatalf("unable to validate example: %s\nHad err: %s", exampleJSON, err)
	}

	typed, err = utils.EnvelopeFromTyped(event)
	if err != nil {
		log.Fatalf("unable to marshal example: %s\nHad err: %s", exampleJSON, err)
	}
	exampleJSON, err = json.MarshalIndent(typed, "", "  ")
	if err != nil {
		log.Fatalf("unable to marshal example: %s\nHad err: %s", exampleJSON, err)
	}

	if name != "" {
		if len(docs) > 0 && strings.HasPrefix(docs[0], typeName) {
			docs[0] = strings.Replace(docs[0], typeName, name, 1)
		}

		fmt.Printf("# %s\n\n", name)
		fmt.Printf("%s", strings.Join(docs, "\n"))
		if len(example) > 0 {
			fmt.Printf("```json\n")
			fmt.Printf("%s\n", exampleJSON)
			fmt.Printf("```\n")
		}
		fmt.Printf("\n")
	}
}
