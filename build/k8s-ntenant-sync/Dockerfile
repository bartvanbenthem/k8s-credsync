# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="bart.vanbenthem@kpn.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy files
COPY ./bin/k8s-ntenant ./
# create dirs for tls files when tls is required
RUN mkdir -p cert/server
RUN mkdir kubeconfig

# Expose port 3111 to the outside world
EXPOSE 3111

# Command to run the executable
CMD ["./k8s-ntenant"]
