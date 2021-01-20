# k8s-ntenant
Multi tenant monitoring and logging solution with credential synchronization on the authentication proxies and grafana datasource configurations. Grafana Organizations and Datasources will be provisioned and configured automatically in a multi tenant setup.

### Technical choices
* github.com/k8spin/loki-multi-tenant-proxy is used as Loki auth proxy.
* github.com/boz/kail is used for log streaming.
* go-client sdk is used to interract with the kubernetes API.
* A custom tls client with basic auth is used to interract with the Grafana API.
* Environment variables are set for dynamic configuration parameters.

### Technical requirements
* The Tenant username should always be identical with the tenant namespace name.
* The multi tenant auth proxy needs to be restarted after secret data update.

## Prerequisites
Install kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl/

## Install and run

#### Deploy Loki in multi-tenant setup
```shell
$ git clone https://github.com/bartvanbenthem/k8s-ntenant.git
# Deploy Loki and proxy for a multi-tenant logging setup
$ cd k8s-ntenant/build/loki-ntenant-setup/
$ ./deploy.sh
```

#### build and run synchronization container
```shell
# change dir
cd build/k8s-ntenant-sync
# build the container
docker build -t k8s-ntenant .
# back to project root
cd ../..
# run container with env variables
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
--add-host 'grafanatls:192.168.2.163' -p 8080:3111 k8s-ntenant 
#-e K8S_SERVER_CERT='<provide .crt file for starting a TLS server>' \
#-e K8S_SERVER_KEY='<provide .pem file for starting a TLS server>' \

```

#### test sync execution from client
Run the k8s-ntenant synchronization server
```shell
# test from client
curl http://localhost:8080/
curl http://localhost:8080/proxy/sync
curl http://localhost:8080/grafana/sync
# view sync logs
docker container logs k8s-ntenant
# interactive session
docker container exec -it k8s-ntenant /bin/bash
```

# TODO
* Design and transform the current build to a kubernetes native build including monitoring.
* Design and create a function for snapshotting the proxy secret before change trough the sync functions.
* Design and create a update function only for the passwords.


