#!/bin/bash

# create namespaces
kubectl create namespace 'co-monitoring'
kubectl create namespace 'team-alpha-dev'
kubectl create namespace 'team-beta-test'
kubectl create namespace 'team-charlie-test'

# apply the loki multi tenant setup
kubectl apply -f .

# install Grafana helmchart
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install grafana --namespace=co-monitoring grafana/grafana
# grafana password
kubectl get secret --namespace co-monitoring grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

# expose grafana to localhost run in seperate terminal
echo 'kubectl port-forward svc/grafana -n co-monitoring 3000:80 &'
# datasource url
echo 'http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100'



