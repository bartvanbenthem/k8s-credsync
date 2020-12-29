package main

import (
	"fmt"
	"os"
)

func TestKubeFunctions() {
	// import environment variables
	proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	proxyns := os.Getenv("K8S_PROXY_SECRET_NAMESPACE")
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")

	var kube KubeCLient
	// get proxy secret data
	fmt.Printf("\n%v\n", string(kube.GetSecretData(kube.CreateClientSet(),
		proxyns, proxysec, "authn.yaml")))
	// get tenant secret data
	fmt.Printf("\n%v\n", string(kube.GetSecretData(kube.CreateClientSet(),
		"team-alpha-dev", tenantsec, "promtail.yaml")))

}
