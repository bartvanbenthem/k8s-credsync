#!/bin/bash

# basic auth for organizations api
AUTH=$(echo -ne "admin:prom-operator" | base64 --wrap 0)

# get organization
curl --header "Content-Type: application/json" --header "Authorization: Basic $AUTH" --request GET --url http://grafana/api/orgs/name/team-alpha-dev
# create organization
curl --header "Content-Type: application/json" --header "Authorization: Basic $AUTH" --request POST --data @grafana-organization.json  --url http://grafana/api/orgs

# get datasource
curl --header "Content-Type: application/json" --header "Authorization: Basic $AUTH" --request GET  --url http://grafana/api/datasources/name/team-alpha-dev
# create datasource
curl --header "Content-Type: application/json" --header "Authorization: Basic $AUTH" --request POST --data @grafana-datasource.json  --url http://grafana/api/datasources