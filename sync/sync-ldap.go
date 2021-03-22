package sync

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/grafana"
	"github.com/bartvanbenthem/k8s-ntenant/ldap"
)

func LDAP() error {
	// Get environment variables
	nsgrafana := os.Getenv("K8S_GRAFANA_NAMESPACE")
	tomlsecret := os.Getenv("K8S_GRAFANA_LDAP_SECRET")
	// Get Grafana ORG IDs from Grafana API
	// Get LDAP groups from configmap
	gs, err := GetAllMappings(nsgrafana)
	if err != nil {
		return err
	}
	// get ldap.toml data from ldap-toml secret
	datakey := os.Getenv("K8S_GRAFANA_LDAP_SECRET_DATA")
	if len(datakey) == 0 {
		datakey = "ldap-toml"
	}

	tomldata, err := ldap.GetLDAPData(nsgrafana, tomlsecret, datakey)
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
	toml := ldap.GetLDAPSecret(nsgrafana, tomlsecret)
	_ = ldap.UpdateLDAPSecret(nsgrafana, datakey, toml, update)
	// return err
	return err
}

func GetAllMappings(nsgrafana string) ([]ldap.GroupMapping, error) {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	ldapgroups := os.Getenv("K8S_GRAFANA_LDAP_GROUPS")
	var mapping ldap.GroupMapping
	var mappings []ldap.GroupMapping
	// get all the current organizations from the grafana api
	orgs, err := grafana.GetAllOrganizations(grafanapi)
	if err != nil {
		return nil, err
	}
	// for every organization get the ldap group from the
	// ldap-groups config map
	for _, o := range orgs {
		mapping.GroupDN = ldap.GetLDAPGroup(nsgrafana, o.Name, ldapgroups)
		mapping.OrgID = o.ID
		mapping.Header = "[[servers.group_mappings]]"
		mapping.OrgRole = "Admin"
		mappings = append(mappings, mapping)
	}
	//return all group mappings and err
	return mappings, err
}

// build the grafana admin group mapping object
func GrafanaAdmin(nsgrafana string) ldap.GroupMapping {
	ldapgroups := os.Getenv("K8S_GRAFANA_LDAP_GROUPS")
	var mapping ldap.GroupMapping
	mapping.GroupDN = ldap.GetLDAPGroup(nsgrafana, "grafana-admin", ldapgroups)
	mapping.OrgID = 1
	mapping.Header = "[[servers.group_mappings]]"
	mapping.OrgRole = "Admin"
	mapping.GrafanaAdmin = true
	// return the admin group mapping
	return mapping
}
