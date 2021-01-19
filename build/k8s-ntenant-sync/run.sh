#!/bin/bash

# build image
docker build -t k8s-ntenant .

# run container with environment variables
docker run -d --name k8s-ntenant \
-e K8S_KUBECONFIG='kubeconfig/client.config' \
-e K8S_PROXY_SECRET_NAME='loki-multi-tenant-proxy-auth-config' \
-e K8S_PROXY_SECRET_NAMESPACE='co-monitoring' \
-e K8S_TENANT_SECRET_NAME='log-recolector-config' \
-e K8S_PROXY_POD_NAME='loki-multi-tenant-proxy-' \
-e K8S_PROXY_URL_PORT='http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100' \
-e K8S_GRAFANA_BA_USER='admin' \
-e K8S_GRAFANA_BA_PASSWORD='prom-operator' \
-e K8S_GRAFANA_API_URL='grafanatls/api' \
-e K8S_GRAFANA_CA_FILE='grafana/rootCA.crt' \
-e K8S_SERVER_ADDRESS='0.0.0.0:3111' \
-e K8S_SERVER_CERT='' \
-e K8S_SERVER_KEY='' \
--add-host 'grafanatls:192.168.2.163' -p 8080:3111 k8s-ntenant 

# test from client
curl http://localhost:8080/
curl http://localhost:8080/proxy/sync
curl http://localhost:8080/grafana/sync

# interactive session
docker container exec -it k8s-ntenant /bin/bash