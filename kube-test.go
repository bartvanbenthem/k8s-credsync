package main

import "fmt"

func TestKubeFunctions() {
	var kube KubeCLient

	fmt.Printf("\n%v\n", string(kube.getTenantSecret(kube.CreateClientSet(),
		"team-alpha-dev", "team-alpha-dev-log-recolector-config")))

	fmt.Printf("\n%v\n", string(kube.getTenantSecret(kube.CreateClientSet(),
		"team-alpha-dev", "team-alpha-dev-log-recolector-config")))

}
