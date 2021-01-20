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
COPY cert/server/. ./server/
RUN mkdir grafana
COPY cert/grafana/. ./grafana/
RUN mkdir kubeconfig
COPY kubeconfig/client.config ./kubeconfig/

# Expose port 3111 to the outside world
EXPOSE 3111

# Command to run the executable
CMD ["./k8s-ntenant"]
