package main

import (
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

func UpdateCheck(c *cli.Context) {
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}
	if len(c.Args().Tail()) < 1 || len(c.Args().First()) < 1 {
		cli.ShowAppHelp(c)
		return
	}
	if err := cfg.client.Agent().UpdateTTL(c.Args().First(), c.String("note"), c.Args().Tail()[0]); err != nil {
		log.Errorf("Could not update check: %v", err)
		return
	}
	log.Println("Success")
}

func RegisterCheck(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	asc := api.AgentServiceCheck{
		Script:   c.String("script"),
		HTTP:     c.String("http"),
		TCP:      c.String("tcp"),
		Interval: c.String("interval"),
		TTL:      c.String("ttl"),
		Timeout:  c.String("timeout"),
	}

	acr := &api.AgentCheckRegistration{AgentServiceCheck: asc}
	acr.ServiceID = c.String("service")
	acr.Notes = c.String("notes")
	acr.Name = c.Args().First()

	if err = cfg.client.Agent().CheckRegister(acr); err != nil {
		log.Printf("Error registering check: %v", err)
		return
	}
	log.Println("Success")
}

func DeregisterCheck(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	if err = cfg.client.Agent().CheckDeregister(c.Args().First()); err != nil {
		log.Errorf("Failed to deregister check: %v", err)
		return
	}
	log.Println("Success")
}

func ListChecks(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	checks, err := cfg.client.Agent().Checks()
	if err != nil {
		log.Errorf("Error listing checks: %v", err)
		return
	}

	if c.GlobalBool("verbose") {
		dumpJson(checks)
		return
	}

	prettyPrintChecks(checks)
}
