#!/bin/bash

#create rootCA cert
openssl genrsa -out rootCA.key 4096
openssl req -x509 -new -key rootCA.key -days 3650 -out rootCA.crt

openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr 
#In answer to question `Common Name (e.g. server FQDN or YOUR name) []:` you should set `secure.domain.com` (your real domain name)
openssl x509 -req -in server.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -days 365 -out server.crt