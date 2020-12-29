package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
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

	namespaces := kube.GetAllNamespaces(kube.CreateClientSet())
	for _, ns := range namespaces {
		s := kube.GetSecretData(kube.CreateClientSet(),
			ns, tenantsec, "promtail.yaml")
		if len(s) != 0 {
			updateTenantPassword(ns)
		}
	}
}

// scan the tenant credentials
// create and update password when empty
func updateTenantPassword(namespace string) {
	// import required environment variables
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")
	var kube KubeCLient
	// create and set password in the tenant credential when password is empty
	cred, err := GetTenantCredential(string(kube.GetSecretData(kube.CreateClientSet(),
		namespace, tenantsec, "promtail.yaml")))
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
	if len(cred.Client.BasicAuth.Username) != 0 {
		PasswordSetter(&cred)
	}

	credbyte, err := yaml.Marshal(&cred)
	if err != nil {
		fmt.Printf("\nerror encoding yaml: %v\n", err)
	}

	sec := kube.GetSecret(kube.CreateClientSet(), namespace, tenantsec)
	sec.Data["promtail.yaml"] = credbyte

	updatedsec := kube.UpdateSecret(kube.CreateClientSet(), namespace, sec)
	fmt.Printf("\n%v\n", string(updatedsec.Data["promtail.yaml"]))
}
