#!/bin/bash

#create rootCA cert
openssl genrsa -out rootCA.key 4096
openssl req -x509 -new -key rootCA.key -days 3650 -out rootCA.crt

openssl genrsa -out grafanatls.key 2048
openssl req -new -key grafanatls.key -out grafanatls.csr 
#In answer to question `Common Name (e.g. server FQDN or YOUR name) []:` you should set `secure.domain.com` (your real domain name)
openssl x509 -req -in grafanatls.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -days 365 -out grafanatls.crt

# create secret
kubectl create secret tls grafanatls-secret --cert=grafanatls.crt --key=grafanatls.key --namespace co-monitoring
kubectl describe secret grafanatls-secret  --namespace co-monitoring