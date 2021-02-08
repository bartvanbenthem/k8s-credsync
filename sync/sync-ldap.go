package sync

import (
	"fmt"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/ldap"
)

func Contains(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}

func LDAP() error {
	// Get LDAP groups from configmap
	// Get Grafana ORG IDs from Grafana API
	// Check if ORGID exists in LDAP.toml
	// Check if DN =; if ORGID exists.
	// If ORGID !exists; or DN !=; add/replace to LDAP.toml

	nsgrafana := os.Getenv("K8S_GRAFANA_NAMESPACE")
	group := ldap.GetLDAPGroup(nsgrafana, "team-alpha-dev")
	fmt.Printf("%v\n", group)

	toml, err := ldap.GetLDAPToml(nsgrafana)

	id := ldap.GetOrgIDFromLDAPToml(nsgrafana, toml)
	fmt.Printf("%v\n", id)

	newtoml := ldap.UpdateLDAPTomlData("66",
		"cn=team alpha-dev,ou=Groups,ou=k8test.nl,ou=Hosting,dc=k8,dc=test,dc=nl",
		toml)
	for _, l := range newtoml {
		fmt.Printf("%v\n", l)
	}
	return err
}
