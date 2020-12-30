#!/bin/bash

#require (
#	k8s.io/api kubernetes-1.19.0
#	k8s.io/apimachinery kubernetes-1.19.0
#	k8s.io/client-go kubernetes-1.19.0
#)

export K8S_KUBECONFIG='/var/snap/microk8s/current/credentials/client.config'
export K8S_PROXY_SECRET_NAME='loki-multi-tenant-proxy-auth-config'
export K8S_PROXY_SECRET_NAMESPACE='co-monitoring'
export K8S_TENANT_SECRET_NAME='log-recolector-config'
export K8S_TENANT_POD_NAME='loki-multi-tenant-proxy-'
export K8S_TENANT_TOKEN_NAME='log-recolector-token-'