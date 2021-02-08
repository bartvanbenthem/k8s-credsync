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
	Header  string
	GroupDN string
	OrgID   string
	OrgRole string
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

func GetLDAPToml(namespace string) *v1.Secret {
	toml := os.Getenv("K8S_GRAFANA_LDAP_SECRET")
	var kube kube.KubeCLient
	s := kube.GetSecret(kube.CreateClientSet(), namespace, toml)
	return s
}

func GetLDAPTomlData(namespace string) ([]string, error) {
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

func GetOrgIDFromLDAPToml(namespace string, toml []string) []string {
	var orgids []string
	for _, l := range toml {
		if strings.Contains(string(l), "org_id") {
			id := strings.Split(l, "=")
			orgids = append(orgids, id[1])
		}
	}
	return orgids
}

func UpdateLDAPTomlData(orgid, groupdn string, tomldata []string) []string {
	group := GroupMapping{Header: "[[servers.group_mappings]]",
		GroupDN: groupdn,
		OrgID:   orgid,
		OrgRole: "Admin"}

	newtoml := tomldata
	newtoml = append(newtoml, group.Header)
	newtoml = append(newtoml, fmt.Sprintf("group_dn = \"%v\"", group.GroupDN))
	newtoml = append(newtoml, fmt.Sprintf("org_id = %v", group.OrgID))
	newtoml = append(newtoml, fmt.Sprintf("org_role = \"%v\"", group.OrgRole))

	return newtoml
}

func UpdateLDAPTomlSecret(namespace string, toml *v1.Secret, tomldata []string) *v1.Secret {
	stom := strings.Join(tomldata, "\n")
	toml.Data["ldap.toml"] = []byte(stom)
	var kube kube.KubeCLient
	kube.UpdateSecret(kube.CreateClientSet(), namespace, toml)
	return toml
}
