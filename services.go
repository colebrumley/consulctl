package main

import (
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func ServiceMaint(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}
	action := c.Args().First()
	service := c.Args().Tail()[0]
	reason := ""

	if len(c.Args().Tail()) > 1 {
		reason = c.Args().Tail()[1]
	}

	switch action {
	case "enable":
		if err = cfg.client.Agent().EnableServiceMaintenance(service, reason); err != nil {
			log.Errorf("Error setting maintenance mode: %v", err)
			return
		}
	case "disable":
		if err = cfg.client.Agent().DisableServiceMaintenance(service); err != nil {
			log.Errorf("Error disabling maintenance mode: %v", err)
			return
		}
	default:
		cli.ShowAppHelp(c)
		return
	}

	log.Println("Success")
}

func DeleteService(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	if err = cfg.client.Agent().ServiceDeregister(c.Args().First()); err != nil {
		log.Errorf("Error removing service: %v", err)
		return
	}
	log.Println("Success")
}

func SetService(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	reg := new(api.AgentServiceRegistration)

	for _, arg := range append(c.Args().Tail(), c.Args().First()) {
		split := strings.SplitN(arg, "=", 2)
		switch strings.ToLower(split[0]) {
		case "address":
			reg.Address = split[1]
		case "name":
			reg.Name = split[1]
		case "id":
			reg.ID = split[1]
		case "port":
			port, _ := strconv.ParseInt(split[1], 10, 32)
			reg.Port = int(port)
		case "tags":
			tags := strings.Split(split[1], ",")
			reg.Tags = tags
		}
	}

	err = cfg.client.Agent().ServiceRegister(reg)
	if err != nil {
		log.Errorf("%v\n", err)
		return
	}
	log.Println("Success")
}

// func GetService(c *cli.Context) {
// 	if !c.Args().Present() {
// 		cli.ShowAppHelp(c)
// 		return
// 	}

// 	// Get client
// 	cfg, err := NewAppConfig(c)
// 	if err != nil {
// 		log.Errorf("Failed to get client: %v", err)
// 		return
// 	}

// 	service, _, err := cfg.client.Catalog().
// 		Service(c.Args().First(), "", cfg.queryOpts)
// 	if err != nil {
// 		log.Errorf("Failed to load services: %v", err)
// 		return
// 	}

// 	if c.GlobalBool("verbose") {
// 		dumpJson(service)
// 		return
// 	}

// 	prettyPrintServices(service)
// }

func ListServices(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	services, err := parseCatalogServices(cfg.client.Catalog(), cfg.queryOpts)
	if err != nil {
		log.Errorf("Failed to list services: %v", err)
		return
	}

	if c.GlobalBool("verbose") {
		dumpJson(services)
		return
	}

	prettyPrintServices(services)
}

func parseCatalogServices(catalog *api.Catalog, qOpts *api.QueryOptions) ([]*api.CatalogService, error) {
	results := []*api.CatalogService{}
	svcs, _, err := catalog.Services(&api.QueryOptions{})
	if err != nil {
		return results, err
	}
	for k, v := range svcs {
		if len(v) > 0 {
			for _, i := range v {
				sInfo, _, err := catalog.Service(k, i, qOpts)
				if err != nil {
					return results, err
				}
				for _, s := range sInfo {
					results = append(results, s)
				}
			}
			continue
		}
		sInfo, _, err := catalog.Service(k, "", qOpts)
		if err != nil {
			return results, err
		}
		for _, s := range sInfo {
			results = append(results, s)
		}
	}
	return results, err
}
