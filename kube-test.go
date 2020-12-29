package main

import "fmt"

func TestKubeFunctions() {
	var kube KubeCLient

	// get proxy secret data
	fmt.Printf("\n%v\n", string(kube.GetSecretData(kube.CreateClientSet(),
		"co-monitoring", "loki-multi-tenant-proxy-auth-config", "authn.yaml")))
	// get tenant secret data
	fmt.Printf("\n%v\n", string(kube.GetSecretData(kube.CreateClientSet(),
		"team-alpha-dev", "team-alpha-dev-log-recolector-config", "promtail.yaml")))

}
