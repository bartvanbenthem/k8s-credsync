package sync

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bartvanbenthem/k8s-ntenant/grafana"
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

	gs, _ := GetAllMappings(nsgrafana)
	// update ldap.toml with all new group mappings
	for _, g := range gs {
		if !utils.Contains(ids, strconv.Itoa(g.OrgID)) {
			tomldata, _ = ldap.GetLDAPTomlData(nsgrafana)
			newdata := ldap.UpdateLDAPTomlData(g.GroupDN,
				"Admin", "[[servers.group_mappings]]", tomldata, g.OrgID)
			toml := ldap.GetLDAPToml(nsgrafana)
			_ = ldap.UpdateLDAPTomlSecret(nsgrafana, toml, newdata)
			log.Printf("Added group \"%v\" to ldap.toml with orgid \"%v\"\n",
				g.GroupDN, g.OrgID)
		} else {
			log.Printf("No Updates for \"ldap.toml\"\n")
		}

	}

	// print updated toml file
	tomldata, _ = ldap.GetLDAPTomlData(nsgrafana)
	ids = ldap.GetOrgIDFromLDAPToml(nsgrafana, tomldata)
	fmt.Printf("%v\n", ids)

	return err

}

func GetAllMappings(nsgrafana string) ([]ldap.GroupMapping, error) {
	var mapping ldap.GroupMapping
	var mappings []ldap.GroupMapping
	orgs, err := grafana.GetAllOrganizations()
	if err != nil {
		return nil, err
	}
	for _, o := range orgs {
		mapping.GroupDN = ldap.GetLDAPGroup(nsgrafana, o.Name)
		mapping.OrgID = o.ID
		mapping.Header = "[[servers.group_mappings]]"
		mapping.OrgRole = "Admin"
		mappings = append(mappings, mapping)
	}
	return mappings, err
}
