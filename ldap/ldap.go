package ldap

import (
	"fmt"
	"os"
	"strings"

	"github.com/bartvanbenthem/k8s-ntenant/kube"
	"github.com/bartvanbenthem/k8s-ntenant/utils"
	v1 "k8s.io/api/core/v1"
)

type GroupMapping struct {
	Header       string
	GroupDN      string
	OrgID        int
	OrgRole      string
	GrafanaAdmin bool
}

// get all the groups from the ldap-groups configmap
func GetAllLDAPGroups(namespace string) map[string]string {
	ldap := os.Getenv("K8S_GRAFANA_LDAP_GROUPS")
	var kube kube.KubeCLient
	cm := kube.GetConfigmap(kube.CreateClientSet(), namespace, ldap)
	return cm.Data
}

// get a specific group from the ldap-groups configmap
// the tenant/namespace name is input for search
func GetLDAPGroup(namespace, tenantname string) string {
	ldap := os.Getenv("K8S_GRAFANA_LDAP_GROUPS")
	var kube kube.KubeCLient
	cm := kube.GetConfigmap(kube.CreateClientSet(), namespace, ldap)
	return cm.Data[tenantname]
}

// get ldap-toml secret in given namespace
func GetLDAPSecret(namespace string) *v1.Secret {
	toml := os.Getenv("K8S_GRAFANA_LDAP_SECRET")
	var kube kube.KubeCLient
	s := kube.GetSecret(kube.CreateClientSet(), namespace, toml)
	return s
}

// get ldap-toml secret data field in given namespace
func GetLDAPData(namespace, datakey string) ([]string, error) {
	toml := os.Getenv("K8S_GRAFANA_LDAP_SECRET")
	var kube kube.KubeCLient
	s := kube.GetSecretData(kube.CreateClientSet(),
		namespace, toml, datakey)
	cfg, err := utils.StringToLines(fmt.Sprintf("%v", string(s)))
	if err != nil {
		return nil, err
	}
	return cfg, err
}

// get all organization ids from the ldap-toml secret
func GetOrgIDFromLDAPSecret(namespace string, toml []string) []string {
	var orgids []string
	for _, l := range toml {
		if strings.Contains(string(l), "org_id") {
			id := strings.Split(l, "=")
			orgids = append(orgids,
				strings.ReplaceAll(id[1], " ", ""))
		}
	}
	return orgids
}

// clean all the current group mappings from
// the ldap-toml secret
func CleanMappingsLDAPData(tomldata []string) []string {
	var n int
	var lines []int
	for i := 0; i < len(tomldata); i++ {
		if tomldata[i] == "[[servers.group_mappings]]" {
			n = i
			lines = append(lines, n)
		}
	}
	ctom := tomldata[:lines[0]]
	return ctom
}

// build the new group mappings object containing all the
// ldap-groups and group mappings for the ldap-toml config
func CreateGroupMappings(groupdn, role, header string, orgid int, root bool) []string {
	group := GroupMapping{Header: header,
		GroupDN:      groupdn,
		OrgID:        orgid,
		OrgRole:      role,
		GrafanaAdmin: root}

	var groups []string
	groups = append(groups, group.Header)
	groups = append(groups, fmt.Sprintf("group_dn = \"%v\"", group.GroupDN))
	groups = append(groups, fmt.Sprintf("org_id = %v", group.OrgID))
	groups = append(groups, fmt.Sprintf("org_role = \"%v\"", group.OrgRole))
	groups = append(groups, fmt.Sprintf("grafana_admin = %v", group.GrafanaAdmin))
	return groups
}

// Update the ldap secret with net ldap-toml data
func UpdateLDAPSecret(namespace, datakey string, toml *v1.Secret, tomldata []string) *v1.Secret {
	stom := strings.Join(tomldata, "\n")
	toml.Data[datakey] = []byte(stom)
	var kube kube.KubeCLient
	kube.UpdateSecret(kube.CreateClientSet(), namespace, toml)
	return toml
}
