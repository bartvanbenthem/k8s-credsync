#!/bin/bash

# go get k8s.io/api kubernetes-1.19.0
# go get k8s.io/apimachinery kubernetes-1.19.0
# go get k8s.io/client-go kubernetes-1.19.0

# Run TLS sync server and request a TLS encrypted Grafana API
export K8S_KUBECONFIG='/home/bartb/.kube/config'
export K8S_CRED_SECRET_NAME='loki-ntenant-credentials'
export K8S_CRED_SECRET_NAMESPACE='co-monitoring'
export K8S_TENANT_SECRET_NAME='log-recolector-config'
export K8S_LOKI_URL_PORT='http://loki.co-monitoring.svc.cluster.local:3100'
export K8S_DATASOURCE_BASIC_AUTH='false'
export K8S_GRAFANA_BA_USER='admin'
export K8S_GRAFANA_BA_PASSWORD='prom-operator'
export K8S_GRAFANA_API_URL='http://grafana/api'
export K8S_GRAFANA_NAMESPACE='co-monitoring'
export K8S_GRAFANA_LDAP_SECRET='ldap-toml'
export K8S_GRAFANA_LDAP_SECRET_DATA='ldap-toml'
export K8S_GRAFANA_LDAP_GROUPS='grafana-ldap-groups'
export K8S_SERVER_ADDRESS='0.0.0.0:3111'
export K8S_SERVER_CERT='build/k8s-ntenant-sync/certificate/server.pem'
export K8S_SERVER_KEY='build/k8s-ntenant-sync/certificate/server.key'

# create build
go build .
mv -f k8s-ntenant build/k8s-ntenant-sync/bin

# run ntenant-sync binary
./build/k8s-ntenant-sync/bin/k8s-ntenant