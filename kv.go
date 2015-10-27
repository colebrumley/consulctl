package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"strings"
)

func ListKeys(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	kv := cfg.client.KV()

	prefix := "/"
	root := true
	arg := strings.TrimSuffix(strings.TrimPrefix(c.Args().First(), "/"), "/")
	if len(arg) > 0 {
		prefix = arg
		root = false
	}

	pairs, _, err := kv.List(prefix, cfg.queryOpts)
	if err != nil {
		log.Debugf("Could not list keys: %v", err)
		return
	}

	if c.GlobalBool("verbose") {
		printKeyJson(pairs, prefix, c.Bool("recurse"))
		return
	}

	prettyPrintKeyList(pairs, prefix, c.Bool("recurse"), root)
}

func GetKvKey(c *cli.Context) {
	if !c.Args().Present() {
		cli.ShowAppHelp(c)
		return
	}

	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	kv := cfg.client.KV()
	var results []*api.KVPair
	for _, a := range append([]string{c.Args().First()}, c.Args().Tail()...) {
		if c.Bool("recurse") {
			pairs, _, err := kv.List(a, cfg.queryOpts)
			if err != nil {
				log.Debugf("Could not list keys: %v", err)
				continue
			}

			for _, p := range pairs {
				results = append(results, p)
			}
			continue
		}

		pair, _, err := kv.Get(a, cfg.queryOpts)
		if err != nil {
			log.Debugf("Could not retrieve key: %v", err)
			continue
		}
		if pair != nil {
			results = append(results, pair)
		}
	}

	if len(results) > 0 {
		if c.GlobalBool("verbose") {
			for _, r := range results {
				bytes, err := marshalPrettyKey(r)
				if err != nil {
					log.Debugf("Could not marshal JSON: %v\n", err)
				}
				fmt.Println(string(bytes))
			}
			return
		}

		for _, r := range results {
			fmt.Printf("%s\n", r.Value)
		}

		return
	}

	log.Errorln("Key not found")
}

func SetKvKey(c *cli.Context) {
	// Get client
	cfg, err := NewAppConfig(c)
	if err != nil {
		log.Errorf("Failed to get client: %v", err)
		return
	}

	if c.Args().Present() && len(c.Args().Tail()) > 0 {
		kv := cfg.client.KV()
		setKey := strings.TrimPrefix(c.Args().First(), "/")
		keyVal := c.Args().Tail()[0]
		if len(setKey) > 0 && len(keyVal) > 0 {
			_, err = kv.Put(&api.KVPair{
				Key:   setKey,
				Value: []byte(keyVal),
			}, cfg.writeOpts)

			if err != nil {
				log.Errorf("Failed to set key: %v", err)
				return
			}
			if !c.Bool("quiet") {
				log.Println("Success")
			}
			return
		}
	}

	log.Errorln("Key or value is empty!")
}

func printKeyJson(pairs []*api.KVPair, prefix string, recurse bool) {
	for _, k := range pairs {
		v := strings.TrimPrefix(k.Key, prefix)
		subKeys := strings.Split(v, "/")
		if len(subKeys) > 2 && !recurse {
			continue
		}
		o, _ := marshalPrettyKey(k)
		fmt.Println(string(o))
	}
}
