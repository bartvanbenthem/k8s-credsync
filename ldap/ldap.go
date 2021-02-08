package ldap

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bartvanbenthem/k8s-ntenant/kube"
)

type GroupMapping struct {
	Header  string
	GroupDN string
	OrgID   string
	OrgRole string
}

func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return lines, err
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

func GetLDAPToml(namespace string) ([]string, error) {
	toml := os.Getenv("K8S_GRAFANA_LDAP_SECRET")
	var kube kube.KubeCLient

	s := kube.GetSecretData(kube.CreateClientSet(),
		namespace, toml, "ldap.toml")
	cfg, err := StringToLines(fmt.Sprintf("%v", string(s)))
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

func UpdateLDAPTomlData(orgid, groupdn string, toml []string) []string {
	group := GroupMapping{Header: "[[servers.group_mappings]]",
		GroupDN: groupdn,
		OrgID:   orgid,
		OrgRole: "Admin"}

	newtoml := toml
	newtoml = append(newtoml, group.Header)
	newtoml = append(newtoml, fmt.Sprintf("group_dn = \"%v\"", group.GroupDN))
	newtoml = append(newtoml, fmt.Sprintf("org_id = %v", group.OrgID))
	newtoml = append(newtoml, fmt.Sprintf("org_role = \"%v\"", group.OrgRole))

	return newtoml
}

func UpdateLDAPTomlSecret(namespace string, toml []byte) {}
