package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/classification/luis"
	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/pkg/errors"
)

const usage = `usage: classify [flags] <input>`

func main() {
	var witToken, luisEndpoint, luisAppID, luisKey, luisSlot string
	var printLogs bool
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.StringVar(&witToken, "wit.token", "", "wit.ai: access token")
	flags.StringVar(&luisEndpoint, "luis.endpoint", "https://luismm2.cognitiveservices.azure.com/", "luis.ai: endpoint URL")
	flags.StringVar(&luisAppID, "luis.appid", "", "luis.ai: application ID")
	flags.StringVar(&luisKey, "luis.key", "production", "luis.ai: subscription key")
	flags.StringVar(&luisSlot, "luis.slot", "production", "luis.ai: slot")
	flags.BoolVar(&printLogs, "logs", false, "whether to print HTTP logs")
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) != 1 {
		fmt.Println(usage)
		flags.PrintDefaults()
		os.Exit(1)
	}

	svcs := make(map[string]flows.ClassificationService)

	if witToken != "" {
		c := flows.NewClassifier(static.NewClassifier("72a82155-deee-471a-97c0-02f36cf6a7e5", "Test", "wit", nil))
		svcs["wit"] = wit.NewService(http.DefaultClient, nil, c, witToken)
	}

	if luisAppID != "" && luisKey != "" {
		c := flows.NewClassifier(static.NewClassifier("ea166a58-a71d-404e-91c9-d28aeb396bc5", "Test", "luis", nil))
		svcs["luis"] = luis.NewService(http.DefaultClient, nil, nil, c, luisEndpoint, luisAppID, luisKey, luisSlot)
	}

	classifications, logs, err := classify(svcs, args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cj, _ := jsonx.MarshalPretty(classifications)
	fmt.Println(string(cj))

	if printLogs {
		fmt.Printf("================ logs ================\n")
		for _, l := range logs {
			lj, _ := jsonx.MarshalPretty(l)
			fmt.Println(string(lj))
		}
	}
}

func classify(svcs map[string]flows.ClassificationService, input string) (map[string]*flows.Classification, []*flows.HTTPLog, error) {
	res := make(map[string]*flows.Classification, len(svcs))
	log := &flows.HTTPLogger{}

	for t, s := range svcs {
		c, err := s.Classify(nil, input, log.Log)
		if err != nil {
			return nil, log.Logs, errors.Wrapf(err, "error classifying with %s", t)
		}
		res[t] = c
	}

	return res, log.Logs, nil
}
