package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant-sync/grafana"
	"github.com/bartvanbenthem/k8s-ntenant-sync/proxy"
	"github.com/bartvanbenthem/k8s-ntenant-sync/tenant"
)

func main() {
	//Start the tenant 2 proxy sync
	//Tenant2Proxy()

	// test by getting the credentials from the current proxy secret
	//TestGetProxyCredentials()

	//Start the Grafana 2 proxy sync
	Grafana2Proxy()

}
func Grafana2Proxy() {
	// Collect all current proxy credentials
	pcreds, err := proxy.AllProxyCredentials()
	if err != nil {
		log.Printf("\n%v\n")
	}

	for _, org := range pcreds.Users {
		o := grafana.GetOrganization(org.Username)
		if len(o.Name) != 0 {
			fmt.Printf("id: %v name: %v\n", o.ID, o.Name)
		} else {
			fmt.Printf("Organization: %v does not exist\n", org)
			organization := grafana.Organization{Name: org.Username}
			grafana.CreateOrganization(organization)
		}
	}
}

func Tenant2Proxy() {
	// Update and collect all current tenant credentials
	tcreds, err := tenant.AllTenantCredentials()
	if err != nil {
		log.Printf("\n%v\n")
	}
	// Update and collect all current proxy credentials
	pcreds, err := proxy.AllProxyCredentials()
	if err != nil {
		log.Printf("\n%v\n")
	}

	// create a slice with all the tenant usernames
	// slice is used to compare to the current proxy users
	var usernames []string
	for _, pc := range pcreds.Users {
		usernames = append(usernames, pc.Username)
	}

	// compare tenant credentials with proxy credentials
	// apply new credentials to the proxy credentials
	var newcreds proxy.Users
	for _, tc := range tcreds {
		b := Contains(usernames, tc.Client.BasicAuth.Username)
		if b != true {
			newcreds.Username = tc.Client.BasicAuth.Username
			newcreds.Password = tc.Client.BasicAuth.Password
			newcreds.Orgid = tc.Client.BasicAuth.Username
			pcreds.Users = append(pcreds.Users, newcreds)
		}
	}

	// update proxy kubernetes secret
	proxy.ReplaceProxySecret(os.Getenv("K8S_PROXY_SECRET_NAMESPACE"),
		"authn.yaml", pcreds)
	// restart proxy
	proxy.RestartProxy(os.Getenv("K8S_PROXY_SECRET_NAMESPACE"),
		os.Getenv("K8S_PROXY_POD_NAME"))
	fmt.Printf("\nproxy has been restarted\n")
}

func Contains(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}
