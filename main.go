package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bartvanbenthem/k8s-ntenant-sync/kube"
	"github.com/bartvanbenthem/k8s-ntenant-sync/proxy"
	"github.com/bartvanbenthem/k8s-ntenant-sync/tenant"
)

func main() {
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
	proxy.ReplaceProxySecret(os.Getenv("K8S_PROXY_SECRET_NAMESPACE"), "authn.yaml", pcreds)
	// restart proxy
	RestartProxy(os.Getenv("K8S_PROXY_SECRET_NAMESPACE"), os.Getenv("K8S_PROXY_POD_NAME"))
	fmt.Printf("\nproxy has been restarted\n")

	// test by getting the credentials from the current proxy and tenant secrets
	TestGetProxyCredentials()
}

func Contains(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}

func RestartProxy(namespace, podname string) {
	// initiate kube client
	var kube kube.KubeCLient
	// restart proxy pod by deleting pod
	// the replicaset will create a new pod with updated config
	pods := kube.GetAllPodNames(kube.CreateClientSet(), namespace)
	for _, p := range pods {
		if strings.Contains(p, podname) {
			kube.DeletePod(kube.CreateClientSet(), namespace, p)
		}
	}
}
