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

	updateTenantPassword("team-beta-test")

}

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
