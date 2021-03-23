package sync

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/credential"
	"github.com/bartvanbenthem/k8s-ntenant/grafana"
)

func Grafana() error {
	// get environment variable
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	csec := os.Getenv("K8S_CRED_SECRET_NAME")
	cns := os.Getenv("K8S_CRED_SECRET_NAMESPACE")
	// Collect all current proxy credentials
	creds, err := credential.AllCredentials(cns, csec)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	// Scan and Create Organizations
	for _, p := range creds.Users {
		o, err := grafana.GetOrganization(grafanapi, p.Username)
		if err != nil {
			log.Printf("%v\n", err)
		}
		if len(o.Name) != 0 {
			log.Printf("Organization \"%v\" exists with ID \"%v\"\n", o.Name, o.ID)
		} else {
			log.Printf("Organization \"%v\" does not exist\n", p)
			organization := grafana.Organization{Name: p.Username}
			grafana.CreateOrganization(grafanapi, organization)
		}
	}

	// Scan an create datasources
	for _, p := range creds.Users {
		o, err := grafana.GetOrganization(grafanapi, p.Username)
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		if len(o.Name) == 0 {
			log.Printf("Organization \"%v\" Cannot be found\n", p.Username)
		}

		// switch the user context to the correct organization
		err = grafana.SwitchUserContext(grafanapi, o)
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}

		ds, err := grafana.GetDatasource(grafanapi, p.Username)
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}

		var datasource grafana.Datasource
		if len(ds.Name) != 0 {
			log.Printf("Datasource \"%v\" exists\n", ds.Name)
			log.Printf("Synchronize \"%v\" password\n", ds.Name)
			ds.SecureJSONData.BasicAuthPassword = p.Password
			ds.JSONData.HTTPHeaderName1 = "X-Scope-OrgID"
			ds.SecureJSONData.HTTPHeaderValue1 = p.TenantID
			// update existing datasource 'ds' in the current context
			err = grafana.UpdateDatasource(grafanapi, ds)
			if err != nil {
				log.Printf("%v\n", err)
				return err
			}
		} else {
			// create datasource object
			datasource.Name = p.Username
			datasource.Type = "loki"
			datasource.URL = os.Getenv("K8S_LOKI_URL_PORT")
			datasource.Access = "proxy"
			datasource.OrgID = o.ID
			datasource.BasicAuth = false
			datasource.ReadOnly = false
			if os.Getenv("K8S_DATASOURCE_BASIC_AUTH") == "true" {
				datasource.BasicAuthUser = p.Username
				datasource.SecureJSONData.BasicAuthPassword = p.Password
				datasource.BasicAuth = true
			}
			datasource.JSONData.HTTPHeaderName1 = "X-Scope-OrgID"
			datasource.SecureJSONData.HTTPHeaderValue1 = p.TenantID
			// create a new datasource in the current context
			err = grafana.CreateDatasource(grafanapi, datasource)
			if err != nil {
				log.Printf("%v\n", err)
				return err
			}
		}
	}
	// return err
	return err
}
