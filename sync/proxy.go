package sync

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/proxy"
	"github.com/bartvanbenthem/k8s-ntenant/tenant"
)

func Contains(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}

func Proxy() {
	// Update and collect all current tenant credentials
	tcreds, err := tenant.AllTenantCredentials()
	if err != nil {
		log.Printf("%v\n", err)
	}
	// Update and collect all current proxy credentials
	pcreds, err := proxy.AllProxyCredentials()
	if err != nil {
		log.Printf("%v\n", err)
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
	log.Printf("Proxy \"%v\" has been restarted\n",
		os.Getenv("K8S_PROXY_POD_NAME"))
	// check for errors
	if err == nil {
		log.Printf("Proxy synchronization finished without errors")
	}
}
