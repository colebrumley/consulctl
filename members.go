package main

import (
	"github.com/codegangsta/cli"
	log "github.com/sirupsen/logrus"
)

func GetMemberInfo(c *cli.Context) {
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}
	members, err := cfg.client.Agent().Members(c.Bool("wan"))
	if err != nil {
		log.Printf("Could not list member info: %v", err)
		return
	}

	if c.GlobalBool("verbose") {
		dumpJson(members)
		return
	}
	prettyPrintMembers(members)
}
