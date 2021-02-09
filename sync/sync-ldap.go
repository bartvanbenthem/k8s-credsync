package sync

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bartvanbenthem/k8s-ntenant/ldap"
	"github.com/bartvanbenthem/k8s-ntenant/utils"
)

func LDAP() error {
	// Get LDAP groups from configmap
	// Get Grafana ORG IDs from Grafana API
	// Check if ORGID exists in LDAP.toml
	// Check if DN =; if ORGID exists.
	// If ORGID !exists; or DN !=; add/replace to LDAP.toml
	nsgrafana := os.Getenv("K8S_GRAFANA_NAMESPACE")

	tomldata, err := ldap.GetLDAPTomlData(nsgrafana)
	ids := ldap.GetOrgIDFromLDAPToml(nsgrafana, tomldata)

	neworg := strconv.Itoa(44)
	dn := ldap.GetLDAPGroup(nsgrafana, "team-beta-test")

	if !utils.Contains(ids, neworg) {
		tomldata, _ = ldap.GetLDAPTomlData(nsgrafana)
		newdata := ldap.UpdateLDAPTomlData(neworg, dn, tomldata)
		toml := ldap.GetLDAPToml(nsgrafana)
		_ = ldap.UpdateLDAPTomlSecret(nsgrafana, toml, newdata)
		log.Printf("Added group \"%v\" to ldap.toml with orgid \"%v\"\n", dn, neworg)
	} else {
		log.Printf("No Updates regarding the ldap.toml\n")
	}

	// print updated toml file
	tomldata, _ = ldap.GetLDAPTomlData(nsgrafana)
	ids = ldap.GetOrgIDFromLDAPToml(nsgrafana, tomldata)
	fmt.Printf("%v\n", ids)

	return err

}
