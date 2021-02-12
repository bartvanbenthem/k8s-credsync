package main

import (
	"log"

	"github.com/bartvanbenthem/k8s-ntenant/sync"
)

func main() {
	var err error
	//Start the tenant 2 proxy sync
	err = sync.Proxy()
	if err != nil {
		log.Printf("%v\n", err)
	}
	//Start the Grafana 2 proxy sync
	err = sync.Grafana()
	if err != nil {
		log.Printf("%v\n", err)
	}
	// start the grafana 2 ldap sync
	err = sync.LDAP()
	if err != nil {
		log.Printf("%v\n", err)
	}
}
