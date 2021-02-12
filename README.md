# k8s-ntenant
Multi tenant monitoring and logging solution with credential synchronization on the authentication proxies and grafana datasource configurations. Grafana Organizations and Datasources will be provisioned and configured automatically in a multi tenant setup.

### Technical choices
* github.com/k8spin/loki-multi-tenant-proxy is used as Loki auth proxy.
* github.com/boz/kail is used for log streaming.
* go-client sdk is used to interract with the kubernetes API.
* A custom tls client with basic auth is used to interract with the Grafana API.
* Environment variables are set for dynamic configuration parameters.

### Requirements
* The Tenant username should always be identical with the tenant namespace name.
* Grafana Organization with ID 1 needs to be prestaged for Grafana administrators.

## Prerequisites
Install kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl/

Run and expose Grafana instance on the Kubernetes cluster:
```shell
# create namespace
$ kubectl create namespace 'co-monitoring'

# install Grafana helmchart
$ helm repo add grafana https://grafana.github.io/helm-charts
$ helm repo update
$ helm install grafana --namespace=co-monitoring grafana/grafana
# grafana password
$ kubectl get secret --namespace co-monitoring grafana \
  -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

# create secret with the grafanatls certificates
$ kubectl create secret tls grafanatls-secret \
   --cert=build/k8s-ntenant-sync/cert/grafana/grafanatls.crt \
   --key=build/k8s-ntenant-sync/cert/grafana/grafanatls.key \
   --namespace co-monitoring
$ kubectl describe secret grafanatls-secret --namespace co-monitoring

# expose grafanatls trough the ingress controller
$ kubectl apply -f build/k8s-ntenant-sync/cert/grafana/grafana-tls.yaml
# make sure the grafanatls hostname is resolvable to the Loadbalancer/node
```

## Install and run

#### Get the project source
```shell
$ git clone https://github.com/bartvanbenthem/k8s-ntenant.git
$ cd k8s-ntenant
```

### Deploy Loki with authentication proxy on Kubernetes
```shell
# create required namespaces
$ kubectl create namespace 'co-monitoring'
$ kubectl create namespace 'team-alpha-dev'
$ kubectl create namespace 'team-beta-test'
$ kubectl create namespace 'team-charlie-test'

# apply the loki multi tenant setup and print the datasource url
$ kubectl apply -f build/loki-ntenant-setup/.
$ echo 'http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100'
```

### Build k8s-ntenant-sync container and push to repo
```shell
# build the container
$ cd build/k8s-ntenant-sync
$ docker build -t bartvanbenthem/k8s-ntenant-sync .

# tag image
$ docker tag bartvanbenthem/k8s-ntenant-sync bartvanbenthem/k8s-ntenant-sync:v2
$ docker image ls

# login and push image to dockerhub repo
$ docker login "docker.io"
$ docker push bartvanbenthem/k8s-ntenant-sync:v2

# back to project root
$ cd ../..
```

### Run k8s-ntenant-sync server on Kubernetes
```shell
# Deploy k8s-ntenant sync server on kubernetes
$ kubectl apply -f build/k8s-ntenant-sync/kubernetes/.
```

#### Execute synchronization from webclient
```shell
# test from client
$ curl --resolve ntenant:127.0.0.1 http://ntenant
$ curl --resolve ntenant:127.0.0.1 http://ntenant/proxy/sync
$ curl --resolve ntenant:127.0.0.1 http://ntenant/grafana/sync

# view sync logs
$ kubectl logs k8s-ntenant-sync-
```

# TODO
* Design and create a function for snapshotting the proxy secret before change trough the sync functions.
* Design and create ldap sync function for automating Grafana LDAP configuration.


