package sync

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/grafana"
	"github.com/bartvanbenthem/k8s-ntenant/ldap"
)

func LDAP() error {
	// Get Grafana namesace from environment variable
	nsgrafana := os.Getenv("K8S_GRAFANA_NAMESPACE")
	// Get Grafana ORG IDs from Grafana API
	// Get LDAP groups from configmap
	gs, err := GetAllMappings(nsgrafana)
	if err != nil {
		return err
	}
	// get ldap.toml data from ldap-toml secret
	tomldata, err := ldap.GetLDAPData(nsgrafana)
	if err != nil {
		return err
	}

	// generate a list of all group mappings
	var update []string
	for _, g := range gs {
		if g.OrgID == 1 {
			admin := GrafanaAdmin(nsgrafana)
			newdata := ldap.CreateGroupMappings(admin.GroupDN,
				"Admin", "[[servers.group_mappings]]", admin.OrgID,
				admin.GrafanaAdmin)
			update = append(update, newdata...)
			log.Printf("Append group \"%v\" with orgid \"%v\" to list\n",
				admin.GroupDN, admin.OrgID)
		} else {
			newdata := ldap.CreateGroupMappings(g.GroupDN,
				"Admin", "[[servers.group_mappings]]", g.OrgID, false)
			update = append(update, newdata...)
			log.Printf("Append group \"%v\" with orgid \"%v\" to list\n",
				g.GroupDN, g.OrgID)
		}
	}

	// update ldap.toml with all new group mappings
	update = append(ldap.CleanMappingsLDAPData(tomldata), update...)
	toml := ldap.GetLDAPSecret(nsgrafana)
	_ = ldap.UpdateLDAPSecret(nsgrafana, toml, update)

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

func GrafanaAdmin(nsgrafana string) ldap.GroupMapping {
	var mapping ldap.GroupMapping
	mapping.GroupDN = ldap.GetLDAPGroup(nsgrafana, "grafana-admin")
	mapping.OrgID = 1
	mapping.Header = "[[servers.group_mappings]]"
	mapping.OrgRole = "Admin"
	mapping.GrafanaAdmin = true
	return mapping
}
