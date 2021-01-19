#!/bin/bash
docker build -t k8s-ntenant .
docker run -d --name k8s-ntenant --add-host 'grafanatls:192.168.2.163' -p 8080:3111 k8s-ntenant 

# test from client
curl http://localhost:8080/
curl http://localhost:8080/proxy/sync
curl http://localhost:8080/grafana/sync

# interactive session
docker container exec -it container-name /bin/bash