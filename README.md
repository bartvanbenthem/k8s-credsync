# k8s-credsync
generate and synchronise basic auth credentials for all tenants with the authentication proxy 

### min requirements credsync
* the tenant secret name to scan is set with an environment variable. 
* there is one tenant secret per namespace to scan.
* the tenant user-name should always be identical with the tenant namespace name.
* the auth-proxy secret name to scan is set with an environment variable. There is a single secret per cluster regarding the auth-proxy.

* the credsync service watches the cluster for new namespaces
* if the tenant password is empty a random password is first generated and added to the tenant secret on the cluster.
* if the tenant secret is not registered in the auth-proxy secret, the auth-proxy is updated
* when a tenant password does not match the auth-proxy password, the auth -roxy is updated

### grafana datasource requiremenst
* grafana organisation names must always match the tenants namespace name and auth-proxy orgid
* grafana datasource configurations need to be created or updated with the credentials from the auth-proxy secret

### technical choices
* go-client sdk is used to interract with the kubernetes API
* a kubernetes service account is used to athenticate and authorize the credsync service
* proxy service needs to be restarted after config update (remove proxy pod)

# TODO
* create a function that checks if a kubernetes resource object exists
* log-recollector service needs to be restarted after config update (remove proxy pod)

