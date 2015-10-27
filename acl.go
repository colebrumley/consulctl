package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

func CreateACL(c *cli.Context) {
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	t, _, err := cfg.client.ACL().Create(&api.ACLEntry{
		Name:  c.String("name"),
		Type:  c.String("type"),
		Rules: c.Args().First(),
		ID:    c.String("id"),
	}, cfg.writeOpts)

	if err != nil {
		log.Errorf("Could not create ACL: %v", err)
		return
	}

	fmt.Println(t)
}

func UpdateACL(c *cli.Context) {
	if len(c.String("id")) > 1 {
		log.Errorln("--id is required!")
		cli.ShowAppHelp(c)
		return
	}

	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	_, err = cfg.client.ACL().Update(&api.ACLEntry{
		Name:  c.String("name"),
		Type:  c.String("type"),
		Rules: c.Args().First(),
		ID:    c.String("id"),
	}, cfg.writeOpts)

	if err != nil {
		log.Errorf("Could not update ACL: %v", err)
		return
	}

	fmt.Println("Success")
}

func ListACL(c *cli.Context) {
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	acls, _, err := cfg.client.ACL().List(cfg.queryOpts)

	if err != nil {
		log.Errorf("Could not retrieve ACLs: %v", err)
		return
	}

	if c.GlobalBool("verbose") {
		dumpJson(acls)
		return
	}

	prettyPrintACLs(acls)
}

func DeregisterACL(c *cli.Context) {
	if len(c.Args().First()) < 1 {
		log.Errorln("ACL ID is required")
		cli.ShowAppHelp(c)
		return
	}
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	_, err = cfg.client.ACL().Destroy(c.Args().First(), cfg.writeOpts)
	if err != nil {
		log.Errorf("Could not destroy ACL: %v", err)
		return
	}
	log.Println("Success")
}

func CloneACL(c *cli.Context) {
	if len(c.Args().First()) < 1 {
		log.Errorln("ACL ID is required")
		cli.ShowAppHelp(c)
		return
	}

	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	newId, _, err := cfg.client.ACL().Clone(c.Args().First(), cfg.writeOpts)
	if err != nil {
		log.Errorf("Could not clone ACL: %v", err)
		return
	}

	fmt.Println(newId)
}
