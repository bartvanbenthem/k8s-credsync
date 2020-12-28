package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type Credential struct {
	Namespace string `yaml:"namespace"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

func main() {

	t, err := getTenantCredentials(secretTenant)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("user:%v pass:%v orgid:%v \n",
		t.Username, t.Password, t.Namespace)
}

func getProxyCredentials(secretTenant string) ([]Credential, error) {
	var err error
	var a []Credential

	err = nil
	return a, err
}

func getTenantCredentials(secretTenant string) (Credential, error) {
	var err error
	var a Credential
	m := make(map[interface{}]interface{})

	// unmarshall entire tenant JSON into a map
	err = yaml.Unmarshal([]byte(secretTenant), &m)
	if err != nil {
		return a, err
	}

	// unmarshall the namespace value into Credential type
	ns, err := yaml.Marshal(m["metadata"])
	err = yaml.Unmarshal([]byte(ns), &a)
	if err != nil {
		return a, err
	}

	// Marshall the stringData from the secret into byte slice
	sd, err := yaml.Marshal(m["stringData"])
	if err != nil {
		return a, err
	}

	// Scan all the lines in sd byte slice
	// append every line to the lines slice of string
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(sd)))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err != nil {
		return a, err
	}

	// check the lines slice of strings for the username and password values
	// append values to Credential type
	for _, line := range lines {
		if strings.Contains(line, "username") {
			val := strings.Split(line, ":")
			a.Username = strings.ReplaceAll(val[1], " ", "")
		} else if strings.Contains(line, "password") {
			val := strings.Split(line, ":")
			a.Password = strings.ReplaceAll(val[1], " ", "")
		}
	}

	return a, err
}

// KUBERNETES TEST SECRETS

var secretProxy = `
apiVersion: v1
kind: Secret
metadata:
  name: loki-multi-tenant-proxy-auth-config
  namespace: co-monitoring
  labels:
    app: loki-multi-tenant-proxy
stringData:
  authn.yaml: |-
    users:
      - username: alpha
        password: alpha
        orgid: team-alpha-dev
      - username: beta
        password: beta
        orgid: team-beta-test
`

var secretTenant = `
apiVersion: v1
kind: Secret
metadata:
  name: team-alpha-dev-log-recolector-config
  namespace: team-alpha-dev
stringData:
  promtail.yaml:  |
    server:
      http_listen_port: 9080
      grpc_listen_port: 0
    client:
      url: http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100/api/prom/push
      basic_auth:
        username: alpha
        password: alpha
    scrape_configs:
      - job_name: containers
        static_configs:
          - targets:
              - localhost
            labels:
              job: containers
              __path__: /loki/logs/containers
        pipeline_stages:
        - regex:
            expression: '^(?P<namespace>.*)\/(?P<pod>.*)\[(?P<container>.*)\]: (?P<content>.*)'
        - labels:
            namespace:
            pod:
            container:
        - output:
            source: content
      - job_name: kail
        static_configs:
          - targets:
              - localhost
            labels:
              job: kail
              __path__: /loki/logs/kail
        pipeline_stages:
        - regex:
            expression: '^time="(?P<time>.*)" level=(?P<level>.*) msg="(?P<content>.*)" cmp=(?P<component>.*)'
        - labels:
            time:
            level:
            component:
        - timestamp:
            source: time
            format: RFC3339
        - output:
            source: content
`
