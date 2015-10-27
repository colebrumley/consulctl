package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

func ListEvents(c *cli.Context) {
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	events, _, err := cfg.client.Event().List(c.Args().First(), cfg.queryOpts)
	if err != nil {
		log.Errorf("Failed to list events: %v", err)
		return
	}

	if c.GlobalBool("verbose") {
		dumpJson(events)
		return
	}
	prettyPrintEvents(events)
}

func FireEvent(c *cli.Context) {
	if len(c.Args().First()) < 1 || len(c.Args().Tail()) < 1 {
		log.Errorln("name and payload are required")
		cli.ShowAppHelp(c)
		return
	}
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	event := &api.UserEvent{
		Name:          c.Args().First(),
		Payload:       []byte(c.Args().Tail()[0]),
		NodeFilter:    c.String("node"),
		ServiceFilter: c.String("service"),
		TagFilter:     c.String("tag"),
	}
	eid, _, err := cfg.client.Event().Fire(event, cfg.writeOpts)

	if err != nil {
		log.Errorf("Could not fire event: %v", err)
		return
	}
	fmt.Println(eid)
}
