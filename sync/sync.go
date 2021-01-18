package sync

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/grafana"
	"github.com/bartvanbenthem/k8s-ntenant/proxy"
	"github.com/bartvanbenthem/k8s-ntenant/tenant"
)

func Contains(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}

func Proxy() error {
	// Update and collect all current tenant credentials
	tcreds, err := tenant.AllTenantCredentials()
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	// Update and collect all current proxy credentials
	pcreds, err := proxy.AllProxyCredentials()
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	// create a slice with all the tenant usernames
	// slice is used to compare to the current proxy users
	var usernames []string
	for _, pc := range pcreds.Users {
		usernames = append(usernames, pc.Username)
	}

	// compare tenant credentials with proxy credentials
	// apply new credentials to the proxy credentials
	var newcreds proxy.Users
	for _, tc := range tcreds {
		b := Contains(usernames, tc.Client.BasicAuth.Username)
		if b != true {
			newcreds.Username = tc.Client.BasicAuth.Username
			newcreds.Password = tc.Client.BasicAuth.Password
			newcreds.Orgid = tc.Client.BasicAuth.Username
			pcreds.Users = append(pcreds.Users, newcreds)
		}
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

func Grafana() error {
	// Collect all current proxy credentials
	pcreds, err := proxy.AllProxyCredentials()
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	// Scan and Create Organizations
	for _, p := range pcreds.Users {
		o, err := grafana.GetOrganization(p.Username)
		if err != nil {
			log.Printf("%v\n", err)
		}
		if len(o.Name) != 0 {
			log.Printf("Organization \"%v\" exists with ID \"%v\"\n", o.Name, o.ID)
		} else {
			log.Printf("Organization \"%v\" does not exist\n", p)
			organization := grafana.Organization{Name: p.Username}
			grafana.CreateOrganization(organization)
		}
	}

	// Scan an create datasources
	for _, p := range pcreds.Users {
		o, err := grafana.GetOrganization(p.Username)
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		if len(o.Name) == 0 {
			log.Printf("Organization \"%v\" Cannot be found\n", p.Username)
		}
		ds, err := grafana.GetDatasource(p.Username)
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		var datasource grafana.Datasource
		datasource.Name = p.Username
		datasource.Type = "loki"
		datasource.URL = os.Getenv("K8S_PROXY_URL_PORT")
		datasource.Access = "proxy"
		datasource.OrgID = o.ID
		datasource.BasicAuth = true
		datasource.BasicAuthUser = p.Username
		datasource.SecureJSONData.BasicAuthPassword = p.Password
		if len(ds.Name) != 0 {
			log.Printf("Datasource \"%v\" already exists\n", ds.Name)
		} else {
			// switch the user context to the correct organization
			err = grafana.SwitchUserContext(o)
			if err != nil {
				log.Printf("%v\n", err)
				return err
			}
			// create datasource in the current context
			err = grafana.CreateDatasource(datasource)
			if err != nil {
				log.Printf("%v\n", err)
				return err
			}
		}
	}
	// return err
	return err
}
