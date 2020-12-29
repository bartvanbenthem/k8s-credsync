package main

import (
	"fmt"
	"log"
	"os"
)

func TestMainFunctions() {
	// import required environment variables
	proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	proxyns := os.Getenv("K8S_PROXY_SECRET_NAMESPACE")
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")

	// initiate kubeclient
	var kube KubeCLient

	// get proxy secret data
	fmt.Printf("\n%v\n", string(kube.GetSecretData(kube.CreateClientSet(),
		proxyns, proxysec, "authn.yaml")))
	// get tenant secret data
	fmt.Printf("\n%v\n", string(kube.GetSecretData(kube.CreateClientSet(),
		"team-alpha-dev", tenantsec, "promtail.yaml")))

	// get the proxy credentials
	proxycred, err := GetProxyCredentials(string(kube.GetSecretData(kube.CreateClientSet(),
		proxyns, proxysec, "authn.yaml")))
	if err != nil {
		log.Printf("\nError: %v\n", err)
	}
	fmt.Printf("\nProxy\n-----\n")
	for _, p := range proxycred.Users {
		fmt.Printf("User:%v Password:%v org:%v\n",
			p.Username, p.Password, p.Orgid)
	}

	// get the credentials from all tenants
	// if the password field is empty generate a random password
	var tenantcred []TenantCredential
	namespaces := kube.GetAllNamespaces(kube.CreateClientSet())
	for _, ns := range namespaces {
		cred, err := GetTenantCredential(string(kube.GetSecretData(kube.CreateClientSet(),
			ns, tenantsec, "promtail.yaml")))
		if err != nil {
			fmt.Printf("\n%v\n", err)
		}
		if len(cred.Client.BasicAuth.Username) != 0 {
			PasswordSetter(&cred)
			tenantcred = append(tenantcred, cred)
		}
	}
	fmt.Printf("\nTenant\n-----\n")
	for _, t := range tenantcred {
		fmt.Printf("User:%v Password:%v\n",
			t.Client.BasicAuth.Username, t.Client.BasicAuth.Password)
	}
}
