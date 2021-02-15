package sync

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/proxy"
	"github.com/bartvanbenthem/k8s-ntenant/tenant"
)

func Proxy() error {
	// Collect all current tenant credentials
	tcreds, err := tenant.AllTenantCredentials()
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	// update proxy config with tenant credentials
	var pcreds proxy.ProxyCredentials
	var newcreds proxy.Users
	for _, tc := range tcreds {
		newcreds.Username = tc.Client.BasicAuth.Username
		newcreds.Password = tc.Client.BasicAuth.Password
		newcreds.Orgid = tc.Client.BasicAuth.Username
		pcreds.Users = append(pcreds.Users, newcreds)
	}
	// update proxy kubernetes secret
	err = proxy.UpdateProxySecret(os.Getenv("K8S_PROXY_SECRET_NAMESPACE"),
		"authn.yaml", pcreds)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	// restart proxy
	proxy.RestartProxy(os.Getenv("K8S_PROXY_SECRET_NAMESPACE"),
		os.Getenv("K8S_PROXY_POD_NAME"))
	log.Printf("Proxy \"%v\" has been restarted\n",
		os.Getenv("K8S_PROXY_POD_NAME"))
	// return err
	return err
}
