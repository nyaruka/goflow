package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/koding/multiconfig"
	_ "github.com/lib/pq"
	"github.com/nyaruka/goflow/mailroom"
	"github.com/nyaruka/goflow/mailroom/config"
)

func main() {
	m := multiconfig.NewWithPath("mailroom.toml")
	config := &config.Mailroom{}

	err := m.Load(config)
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	mailroom := mailroom.New(config)
	err = mailroom.Start()
	if err != nil {
		log.Fatalf("Error starting mailroom: %s", err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	mailroom.Stop()
}
