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

func GetAllLDAPGroups(namespace string) map[string]string {
	ldap := os.Getenv("K8S_GRAFANA_LDAP_GROUPS")
	var kube kube.KubeCLient
	cm := kube.GetConfigmap(kube.CreateClientSet(), namespace, ldap)
	return cm.Data
}

func GetLDAPGroup(namespace, tenantname string) string {
	ldap := os.Getenv("K8S_GRAFANA_LDAP_GROUPS")
	var kube kube.KubeCLient
	cm := kube.GetConfigmap(kube.CreateClientSet(), namespace, ldap)
	return cm.Data[tenantname]
}

func GetLDAPSecret(namespace string) *v1.Secret {
	toml := os.Getenv("K8S_GRAFANA_LDAP_SECRET")
	var kube kube.KubeCLient
	s := kube.GetSecret(kube.CreateClientSet(), namespace, toml)
	return s
}

func GetLDAPData(namespace string) ([]string, error) {
	toml := os.Getenv("K8S_GRAFANA_LDAP_SECRET")
	var kube kube.KubeCLient
	s := kube.GetSecretData(kube.CreateClientSet(),
		namespace, toml, "ldap.toml")
	cfg, err := utils.StringToLines(fmt.Sprintf("%v", string(s)))
	if err != nil {
		return nil, err
	}
	return cfg, err
}

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

func UpdateLDAPSecret(namespace string, toml *v1.Secret, tomldata []string) *v1.Secret {
	stom := strings.Join(tomldata, "\n")
	toml.Data["ldap.toml"] = []byte(stom)
	var kube kube.KubeCLient
	kube.UpdateSecret(kube.CreateClientSet(), namespace, toml)
	return toml
}
