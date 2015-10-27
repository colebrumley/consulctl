package main

import (
	"github.com/codegangsta/cli"
	log "github.com/sirupsen/logrus"
)

func GetSelfInfo(c *cli.Context) {
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	info, err := cfg.client.Agent().Self()
	if err != nil {
		log.Errorf("Could not retrieve agent info: %v", err)
		return
	}
	dumpJson(info)
}

func ToggleMaintenanceMode(c *cli.Context) {
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	action := c.Args().First()

	switch action {
	case "enable":
		if err = cfg.client.Agent().EnableNodeMaintenance(c.String("reason")); err != nil {
			log.Errorf("Could not set maintenance mode: %v", err)
			return
		}
	case "disable":
		if err = cfg.client.Agent().DisableNodeMaintenance(); err != nil {
			log.Errorf("Could not unset maintenance mode: %v", err)
			return
		}
	default:
		log.Warningf("Must choose either enable or disable")
		cli.ShowAppHelp(c)
		return
	}
	log.Println("Success")
}
