package main

import (
	"flag"
	"github.com/knabben/distr53/pkg/conshash"
	"github.com/knabben/distr53/pkg/controller"
	"strings"
)

var (
	hostedzone string
	keyValue   string
	domain     string
)

const membersNumber = 5

func init() {
	flag.StringVar(&keyValue, "keyvalue", "", "send a key1=pair,key2=pair values.")
	flag.StringVar(&hostedzone, "hostedzone", "", "Hosted Zone ID used to access the domain via API.")
	flag.StringVar(&domain, "domain", "opssec.in", "domain used to store the TXT records.")
}

func main() {
	flag.Parse()

	// Initialize controller
	ctrl := controller.NewHashConfig(membersNumber)
	ctrl.StartMembersLL()

	c := conshash.New(ctrl.Members, *ctrl.Config)

	// Get record from consistent hash
	for _, value := range strings.Split(keyValue, ",") {
		record := c.LocateKey([]byte(value))
		ctrl.AddOnRecord(record, value)
	}

	// Persist records on Route53
	ctrl.DumpOnRoute53(hostedzone, domain)
}
