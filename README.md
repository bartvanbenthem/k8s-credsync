# k8s-ntenant-sync
Synchronise credentials for all tenants with the authentication proxy and grafana datasource.

### technical choices
* go-client sdk is used to interract with the kubernetes API
* a kubernetes service account is used to athenticate and authorize the credsync service

### tenant2proxy sync
* the tenant secret name to scan is set with an environment variable. 
* there is one tenant secret per namespace to scan.
* the tenant user-name should always be identical with the tenant namespace name.
* the auth-proxy secret name to scan is set with an environment variable. There is a single secret per cluster regarding the auth-proxy.
* if the tenant secret is not registered in the auth-proxy secret, the auth-proxy is updated
* proxy service needs to be restarted after config update (remove pod)

### grafana2proxy sync
* Provision Grafana Organizations and Datasources based on auth proxy credentials
* grafana organisation names must always match the tenants namespace name and auth-proxy orgid
* grafana datasource configurations need to be created or updated with the credentials from the auth-proxy secret

* get all proxy credentials.
* get all grafana organizations.
* compare proxy credential names with organization names.
* if proxy credential doesnt map to an organization, create a new organization.
* connect new organization to the loki multi tenant proxy with the proxy credentials.

## Set environment variables
```shell
export K8S_KUBECONFIG='/var/snap/microk8s/current/credentials/client.config'
export K8S_PROXY_SECRET_NAME='loki-multi-tenant-proxy-auth-config'
export K8S_PROXY_SECRET_NAMESPACE='co-monitoring'
export K8S_TENANT_SECRET_NAME='log-recolector-config'
export K8S_PROXY_POD_NAME='loki-multi-tenant-proxy-'
export K8S_PROXY_URL_PORT='http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100'

export K8S_GRAFANA_BA_USER='admin'
export K8S_GRAFANA_BA_PASSWORD='prom-operator'
export K8S_GRAFANA_API_URL='grafana/api'
```
## Run Sync services
```shell
$ git clone https://github.com/bartvanbenthem/k8s-ntenant-sync.git
$ ./k8s-ntenant-sync/bin/k8s-ntenant-sync
```

# TODO
* create a function that checks if a kubernetes resource object exists to replace temp wait func.
* expose both sync functions through 2 api endpoints (net/http).
* when a tenant password does not match the auth-proxy password, the auth proxy is updated.
* Update organization ID on Proxy with the generated Grafana organization ID.
