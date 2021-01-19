# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="bart.vanbenthem@kpn.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy files
COPY ./bin/k8s-ntenant ./
RUN mkdir server
COPY cert/server/. ./server
RUN mkdir grafana
COPY cert/grafana/. ./grafana
RUN mkdir kubeconfig
COPY context/client.config ./kubeconfig/

# Expose port 3111 to the outside world
EXPOSE 3111

# set environment variables
ENV K8S_KUBECONFIG='kubeconfig/client.config'
ENV K8S_PROXY_SECRET_NAME='loki-multi-tenant-proxy-auth-config'
ENV K8S_PROXY_SECRET_NAMESPACE='co-monitoring'
ENV K8S_TENANT_SECRET_NAME='log-recolector-config'
ENV K8S_PROXY_POD_NAME='loki-multi-tenant-proxy-'
ENV K8S_PROXY_URL_PORT='http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100'
ENV K8S_GRAFANA_BA_USER='admin'
ENV K8S_GRAFANA_BA_PASSWORD='prom-operator'
ENV K8S_GRAFANA_API_URL='grafanatls/api'
ENV K8S_GRAFANA_CA_FILE='grafana/rootCA.crt'
ENV K8S_SERVER_ADDRESS='0.0.0.0:3111'
ENV K8S_SERVER_CERT='server/server.pem'
ENV K8S_SERVER_KEY='server/server.key'

# Command to run the executable
CMD ["./k8s-ntenant"]
