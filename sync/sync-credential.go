package sync

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/credential"
	"github.com/bartvanbenthem/k8s-ntenant/tenant"
)

func Credential() error {
	// Collect all current tenant credentials
	tcreds, err := tenant.AllTenantCredentials()
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	// update Credential secret with tenant credentials
	var pcreds credential.Credentials
	var newcreds credential.Users
	for _, tc := range tcreds {
		newcreds.Username = tc.Client.BasicAuth.Username
		newcreds.Password = tc.Client.BasicAuth.Password
		newcreds.Orgid = tc.Client.BasicAuth.Username
		newcreds.TenantID = tc.Client.TenantID
		pcreds.Users = append(pcreds.Users, newcreds)
	}
	// update Credential kubernetes secret
	err = credential.UpdateCredentialSecret(os.Getenv("K8S_CRED_SECRET_NAMESPACE"),
		"authn.yaml", pcreds)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	return err
}
