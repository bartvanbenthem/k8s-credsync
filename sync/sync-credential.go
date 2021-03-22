package sync

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/credential"
	"github.com/bartvanbenthem/k8s-ntenant/grafana"
	"github.com/bartvanbenthem/k8s-ntenant/tenant"
)

func Credential() error {
	// import environment variable
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	csec := os.Getenv("K8S_CRED_SECRET_NAME")
	cns := os.Getenv("K8S_CRED_SECRET_NAME")
	// Collect all current tenant credentials
	tcreds, err := tenant.AllTenantCredentials(tenantsec)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	// update Credential secret with tenant credentials
	var pcreds credential.Credentials
	var newcreds credential.Users
	var org grafana.Organization

	for _, tc := range tcreds {
		newcreds.Username = tc.Client.BasicAuth.Username
		newcreds.Password = tc.Client.BasicAuth.Password
		newcreds.TenantID = tc.Client.TenantID
		org, err = grafana.GetOrganization(grafanapi, tc.Client.BasicAuth.Username)
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		newcreds.Orgid = org.ID
		pcreds.Users = append(pcreds.Users, newcreds)
	}
	// update Credential kubernetes secret
	err = credential.UpdateCredentialSecret(cns,
		csec, "authn.yaml", pcreds)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	return err
}
