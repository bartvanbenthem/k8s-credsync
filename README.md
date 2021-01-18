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

#### Set environment variables
```shell
export K8S_KUBECONFIG='/var/snap/microk8s/current/credentials/client.config'
export K8S_PROXY_SECRET_NAME='loki-multi-tenant-proxy-auth-config'
export K8S_PROXY_SECRET_NAMESPACE='co-monitoring'
export K8S_TENANT_SECRET_NAME='log-recolector-config'
export K8S_PROXY_POD_NAME='loki-multi-tenant-proxy-'
export K8S_PROXY_URL_PORT='http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100'

export K8S_GRAFANA_BA_USER='admin'
export K8S_GRAFANA_BA_PASSWORD='prom-operator'
export K8S_GRAFANA_API_URL='grafanatls/api'
export K8S_GRAFANA_CA_FILE='build/cert/rootCA.crt'

export K8S_SERVER_ADDRESS='localhost:8443'
export K8S_SERVER_CERT='build/cert/server/server.pem'
export K8S_SERVER_KEY='build/cert/server/server.key'
```

#### Deploy k8s-ntenant
```shell
$ git clone https://github.com/bartvanbenthem/k8s-ntenant.git
# Deploy Loki and proxy for a multi-tenant logging setup
$ cd k8s-ntenant/build/loki-ntenant-setup/
$ ./deploy.sh
```

#### Run the synchronization server
Run the k8s-ntenant synchronization server
```shell
# Start the k8s-ntenant server
$ ./k8s-ntenant/build/bin/k8s-ntenant
# Execute the synchronization services
$ curl -k https://$K8S_SERVER_ADDRESS/proxy/sync
$ curl -k https://$K8S_SERVER_ADDRESS/grafana/sync
```

# TODO
* Design and transform the current build to a kubernetes native build including monitoring.
* Design and create a function for snapshotting the proxy secret before change trough the sync functions.
* Design and create a update function only for the passwords.


