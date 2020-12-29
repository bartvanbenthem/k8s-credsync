package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	//TestCredentialFunctions()
	//TestKubeFunctions()

	// import required environment variables
	proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	proxyns := os.Getenv("K8S_PROXY_SECRET_NAMESPACE")
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")

	// initiate kubeclient
	var kube KubeCLient

	// get the proxy credentials
	proxycred, err := GetProxyCredentials(string(kube.GetSecretData(kube.CreateClientSet(),
		proxyns, proxysec, "authn.yaml")))
	if err != nil {
		log.Printf("\nError: %v\n", err)
	}
	for _, p := range proxycred.Users {
		fmt.Printf("\nUser:%v Password:%v org:%v\n",
			p.Username, p.Password, p.Orgid)
	}

	// get the credentials from all tenants
	var tenantcred []TenantCredential
	namespaces := kube.GetAllNamespaces(kube.CreateClientSet())
	for _, ns := range namespaces {
		cred, err := GetTenantCredential(string(kube.GetSecretData(kube.CreateClientSet(),
			ns, tenantsec, "promtail.yaml")))
		if err != nil {
			fmt.Printf("\n%v\n", err)
		}
		if len(cred.Client.BasicAuth.Username) != 0 {
			tenantcred = append(tenantcred, cred)
		}
	}

	for _, t := range tenantcred {
		fmt.Printf("\nUser:%v Password:%v\n",
			t.Client.BasicAuth.Username, t.Client.BasicAuth.Password)
	}
}
