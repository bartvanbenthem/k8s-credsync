package main

import "github.com/bartvanbenthem/k8s-ntenant-sync/sync"

func main() {
	//Start the tenant 2 proxy sync
	sync.Tenant2Proxy()
	//Start the Grafana 2 proxy sync
	sync.Grafana2Proxy()
}
