# k8s-credsync
generate and synchronise basic auth credentials for all tenants with the authentication proxy 

### min requirements credsync
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


