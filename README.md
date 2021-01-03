# k8s-ntenant-sync
Synchronise credentials for all tenants with the authentication proxy and grafana datasource.

### min requirements teant2proxy sync
* the tenant secret name to scan is set with an environment variable. 
* there is one tenant secret per namespace to scan.
* the tenant user-name should always be identical with the tenant namespace name.
* the auth-proxy secret name to scan is set with an environment variable. There is a single secret per cluster regarding the auth-proxy.

* if the tenant secret is not registered in the auth-proxy secret, the auth-proxy is updated

### technical choices
* go-client sdk is used to interract with the kubernetes API
* a kubernetes service account is used to athenticate and authorize the credsync service
* proxy service needs to be restarted after config update (remove pod)

# TODO
* create a function that checks if a kubernetes resource object exists to replace temp wait func
* when a tenant password does not match the auth-proxy password, the auth proxy is updated
* Provision Grafana Organizations and Datasources based on auth proxy credentials

### grafana datasource provisioning
* grafana organisation names must always match the tenants namespace name and auth-proxy orgid
* grafana datasource configurations need to be created or updated with the credentials from the auth-proxy secret

* get all proxy credentials.
* get all grafana organizations.
* compare proxy credential names with organization names.
* if proxy credential doesn`t map to an organization, create a new organization.
* connect new organization to the loki multi tenant proxy with the proxy credentials.
