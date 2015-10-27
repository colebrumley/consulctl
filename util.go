package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"text/tabwriter"
)

func getTabwriter() *tabwriter.Writer {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 10, 0, 3, ' ', 0)
	return w
}

func prettyPrintMembers(members []*api.AgentMember) {
	w := getTabwriter()
	fmt.Fprintf(w, "Name\tAddress\tPort\tStatus\tProto\tProtoMax\tProtoMin\tTags\n")
	for _, m := range members {
		if len(m.Tags) < 1 {
			fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\t%v\t%v\t%v\n",
				m.Name, m.Addr, m.Port, m.Status, m.ProtocolCur, m.ProtocolMax, m.ProtocolMin)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\t%v\t%v\t%v\n",
			m.Name, m.Addr, m.Port, m.Status, m.ProtocolCur, m.ProtocolMax, m.ProtocolMin, m.Tags)
	}
	w.Flush()
}

func dumpJson(v interface{}) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Errorf("Could not marshal JSON: %v", err)
		return
	}
	fmt.Printf("%s\n", string(out))
}

func prettyPrintACLs(acls []*api.ACLEntry) {
	w := getTabwriter()
	fmt.Fprintf(w, "Name\tType\tID\tCreateIndex\tModifyIndex\tRules\n")
	for _, a := range acls {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%v\t%s\n", a.Name, a.Type, a.ID, a.CreateIndex, a.ModifyIndex, a.Rules)
	}
	w.Flush()
}

func prettyPrintEvents(events []*api.UserEvent) {
	w := getTabwriter()
	fmt.Fprintf(w, "ID\tName\tPayload\tNodeFilter\tServiceFilter\tTagFilter\n")
	for _, e := range events {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", e.ID, e.Name, string(e.Payload), e.NodeFilter, e.ServiceFilter, e.TagFilter)
	}
	w.Flush()
}

func prettyPrintChecks(checks map[string]*api.AgentCheck) {
	w := getTabwriter()
	fmt.Fprintf(w, "ID\tName\tNode\tService\tStatus\tNotes\n")
	for _, c := range checks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", c.CheckID, c.Name, c.Node, c.ServiceID, c.Status, c.Notes)
	}
	w.Flush()
}

func prettyPrintServices(svcs []*api.CatalogService) {
	w := getTabwriter()
	fmt.Fprintf(w, "ID\tName\tNode\tAddress\tPort\tTags\n")
	ignoreList := []string{}
	for _, s := range svcs {
		skip := false
		for _, i := range ignoreList {
			if s.ServiceID == i {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		addr := s.Address
		if len(s.ServiceAddress) > 0 {
			addr = s.ServiceAddress
		}
		if len(s.ServiceTags) > 1 {
			ignoreList = append(ignoreList, s.ServiceID)
		}

		if len(s.ServiceTags) > 0 {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%v\n",
				s.ServiceID, s.ServiceName, s.Node, addr, s.ServicePort, s.ServiceTags)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n",
			s.ServiceID, s.ServiceName, s.Node, addr, s.ServicePort)
	}

	w.Flush()
}

func prettyPrintKeyList(pairs []*api.KVPair, prefix string, recurse bool, root bool) {
	resultList := []string{}

	for _, p := range pairs {
		v := strings.TrimPrefix(p.Key, prefix)

		if len(v) > 0 {
			subKeys := strings.Split(v, "/")
			if len(subKeys) > 2 && !recurse {
				str := strings.Join(subKeys[:2], "/") + "..."
				if root {
					resultList = appendUnique("/"+str, resultList)
					// fmt.Println("/" + str)
					continue
				}
				resultList = appendUnique(str, resultList)
				// fmt.Println(str)
				continue
			}
			if root {
				resultList = appendUnique("/"+v, resultList)
				// fmt.Println("/" + v)
				continue
			}
			resultList = appendUnique(v, resultList)
			// fmt.Println(v)
			continue
		}
		resultList = appendUnique(".", resultList)
		// fmt.Println(".")
	}

	for _, str := range resultList {
		fmt.Println(str)
	}
}

func appendUnique(s string, l []string) []string {
	if stringArrayContains(s, l) {
		return l
	}
	return append(l, s)
}

func stringArrayContains(s string, l []string) bool {
	for _, i := range l {
		if s == i {
			return true
		}
	}
	return false
}

func marshalPrettyKey(p *api.KVPair) ([]byte, error) {
	xmog := &struct {
		Key         string
		CreateIndex uint64
		ModifyIndex uint64
		LockIndex   uint64
		Flags       uint64
		Value       string
		Session     string
	}{
		Key:         p.Key,
		CreateIndex: p.CreateIndex,
		ModifyIndex: p.ModifyIndex,
		LockIndex:   p.LockIndex,
		Flags:       p.Flags,
		Value:       string(p.Value),
		Session:     p.Session,
	}
	return json.MarshalIndent(xmog, "", "  ")
}
