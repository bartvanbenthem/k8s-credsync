package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	//TestCredentialFunctions()
	//TestKubeFunctions()
	//TestMainFunctions()

	// Update and collect all tenant credentials
	tcreds, err := AllTenantCredentials()
	if err != nil {
		log.Printf("\n%v\n")
	}
	fmt.Printf("\n%v\n", len(tcreds))

	pcreds, err := AllProxyCredentials()
	if err != nil {
		log.Printf("\n%v\n")
	}

	fmt.Printf("\nProxy\n-----\n")
	for _, p := range pcreds.Users {
		fmt.Printf("User:%v Password:%v org:%v\n",
			p.Username, p.Password, p.Orgid)
	}
}

// collects all proxy credentials
func AllProxyCredentials() (ProxyCredentials, error) {
	var err error
	// import environment variables
	proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	proxyns := os.Getenv("K8S_PROXY_SECRET_NAMESPACE")
	// initiate kube client
	var kube KubeCLient

	// get the proxy credentials
	proxycred, err := GetProxyCredentials(string(kube.GetSecretData(kube.CreateClientSet(),
		proxyns, proxysec, "authn.yaml")))
	if err != nil {
		return proxycred, err
	}
	return proxycred, err
}

// collects all tenant credentials
// updates credentials when password is an empty string
func AllTenantCredentials() ([]TenantCredential, error) {
	var err error
	// import environment variable
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")
	// initiate kube client
	var kube KubeCLient
	//set slice of tenant credential
	var tcreds []TenantCredential

	namespaces := kube.GetAllNamespaces(kube.CreateClientSet())
	for _, ns := range namespaces {
		s := kube.GetSecretData(kube.CreateClientSet(),
			ns, tenantsec, "promtail.yaml")
		if len(s) != 0 {
			updateTenantSecret(ns, "promtail.yaml")
			// get updated tenant credential
			// append updated credentials to slice of tenant credential
			upd, err := GetTenantCredential(string(kube.GetSecretData(
				kube.CreateClientSet(), ns, tenantsec, "promtail.yaml")))
			if err != nil {
				return nil, err
			}
			tcreds = append(tcreds, upd)
		}
	}
	return tcreds, err
}
