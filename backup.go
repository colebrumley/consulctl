package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type ConsulBackup struct {
	NodeName string                `json:"node_name"`
	Members  []*api.AgentMember    `json:"cluster_members"`
	KV       api.KVPairs           `json:"kv,omitempty"`
	Services []*api.CatalogService `json:"services,omitempty"`
}

func RestoreConsul(c *cli.Context) {
	fileBytes, err := ioutil.ReadFile(c.Args().First())
	if err != nil {
		log.Errorf("Could not load file: %v", err)
		return
	}

	restore := &ConsulBackup{}
	if err = json.Unmarshal(fileBytes, &restore); err != nil {
		log.Errorf("Could not parse JSON: %v", err)
		return
	}

	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	kv := cfg.client.KV()

	for _, k := range restore.KV {
		log.Printf("[KV] Restoring %s", k.Key)
		if _, err := kv.Put(k, cfg.writeOpts); err != nil {
			log.Errorf("[KV] Could not restore key %s: %v\n", k.Key, err)
			return
		}
	}

	cat := cfg.client.Catalog()
	for _, s := range restore.Services {
		log.Printf("[SVC] Restoring %s", s.ServiceID)
		if _, err := cat.Register(&api.CatalogRegistration{
			Node:       s.Node,
			Address:    s.Address,
			Datacenter: c.GlobalString("datacenter"),
			Service: &api.AgentService{
				ID:      s.ServiceID,
				Service: s.ServiceName,
				Tags:    s.ServiceTags,
				Port:    s.ServicePort,
				Address: s.ServiceAddress,
			},
		}, cfg.writeOpts); err != nil {
			log.Errorf("[SVC] Could not restore service %s: %v\n", s.ServiceID, err)
		}
	}

	log.Println("Restore complete!")
}

func BackupConsul(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	// Dump KV list
	kv := cfg.client.KV()
	keys, _, err := kv.List("/", cfg.queryOpts)
	if err != nil {
		log.Errorf("Failed to fetch KV pairs: %v", err)
		return
	}

	// Extract api.CatalogService from service catalog
	services, err := parseCatalogServices(cfg.client.Catalog(), cfg.queryOpts)
	if err != nil {
		log.Error(err)
	}

	// Get member list
	members, err := cfg.client.Agent().Members(false)
	if err != nil {
		log.Error(err)
	}

	// Get node name
	name, err := cfg.client.Agent().NodeName()
	if err != nil {
		log.Error(err)
	}

	// Create the JSON struct
	backup := &ConsulBackup{
		NodeName: name,
		Members:  members,
		KV:       keys,
		Services: services,
	}

	var (
		v []byte
	)

	// Indent json based on flag
	if c.Bool("indent") {
		v, err = json.MarshalIndent(backup, "", "  ")
	} else {
		v, err = json.Marshal(backup)
	}
	if err != nil {
		log.Error(err)
		return
	}

	if len(c.String("outfile")) > 0 {
		log.Infof("Writing backup to %s", c.String("outfile"))
		if err := ioutil.WriteFile(c.String("outfile"), v, 0664); err != nil {
			log.Error(err)
		}
		return
	}

	fmt.Println(string(v))
}
