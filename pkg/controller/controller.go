// Package controller holds the logic across components to
// real persistence of the records
package controller

import (
	"fmt"
	"github.com/knabben/distr53/pkg/conshash"
	"github.com/knabben/distr53/pkg/linkedlist"
	r "github.com/knabben/distr53/pkg/route53"
	"hash/fnv"
	"log"
)

type Controller struct {
	Config  *conshash.Config
	Members []conshash.Member
	Mapper  map[string]linkedlist.Storage
}

type Record string

func (m Record) String() string {
	return string(m)
}

type HasherFunc struct{}

func (hs HasherFunc) Sum64(data []byte) uint64 {
	h := fnv.New64()
	h.Write(data)
	return h.Sum64()
}

var names = []string{"a", "b", "c", "d", "e"}

func NewHashConfig(membersCount int) *Controller {
	members := make([]conshash.Member, membersCount)
	for i, name := range names {
		members[i] = Record(fmt.Sprintf("mb%s", name))
	}
	return &Controller{
		Config: &conshash.Config{
			PartitionCount:    23,
			ReplicationFactor: 20,
			Load:              1.25,
			Hasher:            HasherFunc{},
		},
		Members: members,
	}
}

func (c *Controller) StartMembersLL() {
	c.Mapper = make(map[string]linkedlist.Storage, len(c.Members))
	for _, member := range c.Members {
		c.Mapper[member.String()] = linkedlist.NewDLL()
	}
}

func (c *Controller) AddOnRecord(record conshash.Member, value string) {
	node := linkedlist.NewDLLNode(value)
	c.Mapper[record.String()].Append(node)
}

func (c *Controller) DumpOnRoute53(hostedzone, domain string) {
	svc := r.GetRoute53Service(r.GetConfig())
	for record, v := range c.Mapper {
		if len(v.GetAllValues()) > 0 {
			record = fmt.Sprintf("%s.%s.", record, domain)
			status := r.SaveRecord(svc, hostedzone, record, v.GetAllValues())
			log.Printf("Added the following values: %s on %s finished with %s\n\n", v.GetAllValues(), record, status)
		}
	}
}
