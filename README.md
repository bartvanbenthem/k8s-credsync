# k8s-ntenant
Multi-tenant logging solution with credential synchronization. Grafana is provisioned and configured automatically in a multi-tenant setup.

### Topology
![topology](/00-img/20210320-k8s-ntentant.png)

### Technical choices
* github.com/boz/kail is used for log streaming.
* go-client sdk is used to interract with the kubernetes API.
* A custom client with basic auth is used to interract with the Grafana API.
* Environment variables are set for dynamic configuration parameters.

### Requirements
* The Tenant username is a required field and should always be identical with the tenant namespace name.
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
# change password to prom-operator to run example 
```

## Install and run

#### Get the project source
```shell
$ git clone https://github.com/bartvanbenthem/k8s-ntenant.git
$ cd k8s-ntenant
```

### Deploy Loki on Kubernetes
```shell
# create required namespaces
$ kubectl create namespace 'co-monitoring'
$ kubectl create namespace 'team-alpha-dev'
$ kubectl create namespace 'team-beta-test'
$ kubectl create namespace 'team-charlie-test'

# deploy loki multi-tenant
$ kubectl apply -f build/loki-ntenant-setup/.

```

### Build k8s-ntenant-sync container and push to repo
```shell
# build the container
$ cd build/k8s-ntenant-sync
$ docker build -t bartvanbenthem/k8s-ntenant-sync .

# tag image
$ docker tag bartvanbenthem/k8s-ntenant-sync:latest bartvanbenthem/k8s-ntenant-sync:v5
$ docker image ls

# login and push image to dockerhub repo
$ docker login "docker.io"
$ docker push bartvanbenthem/k8s-ntenant-sync:v5

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
$ curl --resolve ntenant:127.0.0.1 http://ntenant/credential/sync
$ curl --resolve ntenant:127.0.0.1 http://ntenant/grafana/sync
$ curl --resolve ntenant:127.0.0.1 http://ntenant/ldap/sync

# view sync logs
$ kubectl logs k8s-ntenant-sync-
```

# TODO
* Design and create sync function for automatic removal of Grafana organization when a tenant is removed from the cluster.
* Design and create function for updating organization ids in the loki-ntenant-credentials after the organization has been created in grafana.


