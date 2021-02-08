#!/bin/bash

# go get k8s.io/api kubernetes-1.19.0
# go get k8s.io/apimachinery kubernetes-1.19.0
# go get k8s.io/client-go kubernetes-1.19.0

export K8S_KUBECONFIG='/var/snap/microk8s/current/credentials/client.config'
export K8S_PROXY_SECRET_NAME='loki-multi-tenant-proxy-auth-config'
export K8S_PROXY_SECRET_NAMESPACE='co-monitoring'
export K8S_TENANT_SECRET_NAME='log-recolector-config'
export K8S_PROXY_POD_NAME='loki-multi-tenant-proxy-'
export K8S_PROXY_URL_PORT='http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100'
export K8S_GRAFANA_BA_USER='admin'
export K8S_GRAFANA_BA_PASSWORD='prom-operator'
export K8S_GRAFANA_API_URL='grafanatls/api'
export K8S_GRAFANA_CA_FILE='build/k8s-ntenant-sync/cert/grafana/rootCA.crt'
export K8S_GRAFANA_NAMESPACE='co-monitoring'
export K8S_GRAFANA_LDAP_SECRET='ldap-toml'
export K8S_GRAFANA_LDAP_GROUPS='grafana-ldap-groups'

# create build
cd cmd
go build .
cd ..
mv -f cmd/cmd build/k8s-ntenant-sync/bin/

# run ntenant-sync binary
./build/k8s-ntenant-sync/bin/cmd