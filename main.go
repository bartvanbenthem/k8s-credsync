package main

import (
	"fmt"
	"os"
)

func main() {
	//TestCredentialFunctions()
	//TestKubeFunctions()
	//TestMainFunctions()

	// initiate kube client
	var kube KubeCLient

	// import environment variables
	//proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	//proxyns := os.Getenv("K8S_PROXY_SECRET_NAMESPACE")
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")

	var tcreds []TenantCredential
	namespaces := kube.GetAllNamespaces(kube.CreateClientSet())
	for _, ns := range namespaces {
		s := kube.GetSecretData(kube.CreateClientSet(),
			ns, tenantsec, "promtail.yaml")
		if len(s) != 0 {
			updateTenantPassword(ns)
			// get updated tenant credential
			// append updated credentials to slice of tenant credential
			upd, err := GetTenantCredential(string(kube.GetSecretData(
				kube.CreateClientSet(), ns, tenantsec, "promtail.yaml")))
			if err != nil {
				fmt.Printf("\n%v\n", err)
			}
			tcreds = append(tcreds, upd)
		}
	}

	// test tenant credentials
	fmt.Printf("\n%v\n", len(tcreds))
}
