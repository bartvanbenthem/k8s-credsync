package main

import (
	"fmt"
	"log"
)

// Prints the current tenant and proxy credentials
func TestMainFunctions() {
	// Prints the current tenant credentials
	tcurrent, err := AllTenantCredentials()
	fmt.Printf("\nTenant\n------\n")
	for _, tc := range tcurrent {
		fmt.Printf("User:%v Password:%v\n",
			tc.Client.BasicAuth.Username,
			tc.Client.BasicAuth.Password)
	}

	// Prints the current proxycredentials
	pcurrent, err := AllProxyCredentials()
	if err != nil {
		log.Printf("\n%v\n")
	}
	fmt.Printf("\nProxy\n-----\n")
	for _, pc := range pcurrent.Users {
		fmt.Printf("User:%v Password:%v org:%v\n",
			pc.Username, pc.Password, pc.Orgid)
	}
}

func TestCredentialFunctions() {
	// get base64 encoded proxy secret
	proxy, err := GetEncodedSecret(secretProxy, "\"authn.yaml\":")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// get base64 encoded tenant secret
	tenant, err := GetEncodedSecret(secretTenantEmptyPassword, "\"promtail.yaml\":")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// test by printing decoded values
	//fmt.Printf("\nproxy\n-----\n%v\n", decodeSecret(proxy))
	//fmt.Printf("\ntenant\n------\n%v\n", decodeSecret(tenant))

	proxycred, err := GetProxyCredentials(DecodeSecret(proxy))
	if err != nil {
		log.Printf("error: %v", err)
	}

	tenantcred, err := GetTenantCredential(DecodeSecret(tenant))
	if err != nil {
		log.Printf("error: %v", err)
	}

	// test by printing proxy struct values
	fmt.Printf("proxy\n-----\n")
	for _, c := range proxycred.Users {
		fmt.Printf("username:%v password:%v org:%v\n",
			c.Username, c.Password, c.Orgid)
	}
	// test by printing tenant struct values
	// when password is an empty string generate password
	fmt.Printf("\ntenant\n------\n")
	PasswordSetter(&tenantcred)
	fmt.Printf("username:%v password:%v\n",
		tenantcred.Client.BasicAuth.Username,
		tenantcred.Client.BasicAuth.Password)
}

// KUBERNETES TEST SECRETS IN RAW JSON RESPONSE
var secretProxy = `
{
    "apiVersion": "v1",
    "data": {
        "authn.yaml": "dXNlcnM6CiAgLSB1c2VybmFtZTogYWxwaGEKICAgIHBhc3N3b3JkOiBhbHBoYQogICAgb3JnaWQ6IHRlYW0tYWxwaGEtZGV2CiAgLSB1c2VybmFtZTogYmV0YQogICAgcGFzc3dvcmQ6IGJldGEKICAgIG9yZ2lkOiB0ZWFtLWJldGEtdGVzdA=="
    },
    "kind": "Secret",
    "metadata": {
        "annotations": {
            "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Secret\",\"metadata\":{\"annotations\":{},\"labels\":{\"app\":\"loki-multi-tenant-proxy\"},\"name\":\"loki-multi-tenant-proxy-auth-config\",\"namespace\":\"co-monitoring\"},\"stringData\":{\"authn.yaml\":\"users:\\n  - username: alpha\\n    password: alpha\\n    orgid: team-alpha-dev\\n  - username: beta\\n    password: beta\\n    orgid: team-beta-test\"}}\n"
        },
        "creationTimestamp": "2020-12-28T13:19:16Z",
        "labels": {
            "app": "loki-multi-tenant-proxy"
        },
        "managedFields": [
            {
                "apiVersion": "v1",
                "fieldsType": "FieldsV1",
                "fieldsV1": {
                    "f:data": {
                        ".": {},
                        "f:authn.yaml": {}
                    },
                    "f:metadata": {
                        "f:annotations": {
                            ".": {},
                            "f:kubectl.kubernetes.io/last-applied-configuration": {}
                        },
                        "f:labels": {
                            ".": {},
                            "f:app": {}
                        }
                    },
                    "f:type": {}
                },
                "manager": "kubectl-client-side-apply",
                "operation": "Update",
                "time": "2020-12-28T13:19:16Z"
            }
        ],
        "name": "loki-multi-tenant-proxy-auth-config",
        "namespace": "co-monitoring",
        "resourceVersion": "63827",
        "selfLink": "/api/v1/namespaces/co-monitoring/secrets/loki-multi-tenant-proxy-auth-config",
        "uid": "debfd1e0-6812-40ae-bc76-8716438bad09"
    },
    "type": "Opaque"
}
`

var secretTenant = `
{
    "apiVersion": "v1",
    "data": {
        "promtail.yaml": "c2VydmVyOgogIGh0dHBfbGlzdGVuX3BvcnQ6IDkwODAKICBncnBjX2xpc3Rlbl9wb3J0OiAwCmNsaWVudDoKICB1cmw6IGh0dHA6Ly9sb2tpLW11bHRpLXRlbmFudC1wcm94eS5jby1tb25pdG9yaW5nLnN2Yy5jbHVzdGVyLmxvY2FsOjMxMDAvYXBpL3Byb20vcHVzaAogIGJhc2ljX2F1dGg6CiAgICB1c2VybmFtZTogYWxwaGEKICAgIHBhc3N3b3JkOiBhbHBoYQpzY3JhcGVfY29uZmlnczoKICAtIGpvYl9uYW1lOiBjb250YWluZXJzCiAgICBzdGF0aWNfY29uZmlnczoKICAgICAgLSB0YXJnZXRzOgogICAgICAgICAgLSBsb2NhbGhvc3QKICAgICAgICBsYWJlbHM6CiAgICAgICAgICBqb2I6IGNvbnRhaW5lcnMKICAgICAgICAgIF9fcGF0aF9fOiAvbG9raS9sb2dzL2NvbnRhaW5lcnMKICAgIHBpcGVsaW5lX3N0YWdlczoKICAgIC0gcmVnZXg6CiAgICAgICAgZXhwcmVzc2lvbjogJ14oP1A8bmFtZXNwYWNlPi4qKVwvKD9QPHBvZD4uKilcWyg/UDxjb250YWluZXI+LiopXF06ICg/UDxjb250ZW50Pi4qKScKICAgIC0gbGFiZWxzOgogICAgICAgIG5hbWVzcGFjZToKICAgICAgICBwb2Q6CiAgICAgICAgY29udGFpbmVyOgogICAgLSBvdXRwdXQ6CiAgICAgICAgc291cmNlOiBjb250ZW50CiAgLSBqb2JfbmFtZToga2FpbAogICAgc3RhdGljX2NvbmZpZ3M6CiAgICAgIC0gdGFyZ2V0czoKICAgICAgICAgIC0gbG9jYWxob3N0CiAgICAgICAgbGFiZWxzOgogICAgICAgICAgam9iOiBrYWlsCiAgICAgICAgICBfX3BhdGhfXzogL2xva2kvbG9ncy9rYWlsCiAgICBwaXBlbGluZV9zdGFnZXM6CiAgICAtIHJlZ2V4OgogICAgICAgIGV4cHJlc3Npb246ICdedGltZT0iKD9QPHRpbWU+LiopIiBsZXZlbD0oP1A8bGV2ZWw+LiopIG1zZz0iKD9QPGNvbnRlbnQ+LiopIiBjbXA9KD9QPGNvbXBvbmVudD4uKiknCiAgICAtIGxhYmVsczoKICAgICAgICB0aW1lOgogICAgICAgIGxldmVsOgogICAgICAgIGNvbXBvbmVudDoKICAgIC0gdGltZXN0YW1wOgogICAgICAgIHNvdXJjZTogdGltZQogICAgICAgIGZvcm1hdDogUkZDMzMzOQogICAgLSBvdXRwdXQ6CiAgICAgICAgc291cmNlOiBjb250ZW50Cg=="
    },
    "kind": "Secret",
    "metadata": {
        "annotations": {
            "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Secret\",\"metadata\":{\"annotations\":{},\"name\":\"team-alpha-dev-log-recolector-config\",\"namespace\":\"team-alpha-dev\"},\"stringData\":{\"promtail.yaml\":\"server:\\n  http_listen_port: 9080\\n  grpc_listen_port: 0\\nclient:\\n  url: http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100/api/prom/push\\n  basic_auth:\\n    username: alpha\\n    password: alpha\\nscrape_configs:\\n  - job_name: containers\\n    static_configs:\\n      - targets:\\n          - localhost\\n        labels:\\n          job: containers\\n          __path__: /loki/logs/containers\\n    pipeline_stages:\\n    - regex:\\n        expression: '^(?P\\u003cnamespace\\u003e.*)\\\\/(?P\\u003cpod\\u003e.*)\\\\[(?P\\u003ccontainer\\u003e.*)\\\\]: (?P\\u003ccontent\\u003e.*)'\\n    - labels:\\n        namespace:\\n        pod:\\n        container:\\n    - output:\\n        source: content\\n  - job_name: kail\\n    static_configs:\\n      - targets:\\n          - localhost\\n        labels:\\n          job: kail\\n          __path__: /loki/logs/kail\\n    pipeline_stages:\\n    - regex:\\n        expression: '^time=\\\"(?P\\u003ctime\\u003e.*)\\\" level=(?P\\u003clevel\\u003e.*) msg=\\\"(?P\\u003ccontent\\u003e.*)\\\" cmp=(?P\\u003ccomponent\\u003e.*)'\\n    - labels:\\n        time:\\n        level:\\n        component:\\n    - timestamp:\\n        source: time\\n        format: RFC3339\\n    - output:\\n        source: content\\n\"}}\n"
        },
        "creationTimestamp": "2020-12-28T13:19:27Z",
        "managedFields": [
            {
                "apiVersion": "v1",
                "fieldsType": "FieldsV1",
                "fieldsV1": {
                    "f:data": {
                        ".": {},
                        "f:promtail.yaml": {}
                    },
                    "f:metadata": {
                        "f:annotations": {
                            ".": {},
                            "f:kubectl.kubernetes.io/last-applied-configuration": {}
                        }
                    },
                    "f:type": {}
                },
                "manager": "kubectl-client-side-apply",
                "operation": "Update",
                "time": "2020-12-28T13:19:27Z"
            }
        ],
        "name": "team-alpha-dev-log-recolector-config",
        "namespace": "team-alpha-dev",
        "resourceVersion": "63853",
        "selfLink": "/api/v1/namespaces/team-alpha-dev/secrets/team-alpha-dev-log-recolector-config",
        "uid": "d6ca05ba-c265-4cb8-a23b-7a3b32453f72"
    },
    "type": "Opaque"
}
`

var secretTenantEmptyPassword = `
{
    "apiVersion": "v1",
    "data": {
        "promtail.yaml": "c2VydmVyOgogIGh0dHBfbGlzdGVuX3BvcnQ6IDkwODAKICBncnBjX2xpc3Rlbl9wb3J0OiAwCmNsaWVudDoKICB1cmw6IGh0dHA6Ly9sb2tpLW11bHRpLXRlbmFudC1wcm94eS5jby1tb25pdG9yaW5nLnN2Yy5jbHVzdGVyLmxvY2FsOjMxMDAvYXBpL3Byb20vcHVzaAogIGJhc2ljX2F1dGg6CiAgICB1c2VybmFtZTogYWxwaGEKICAgIHBhc3N3b3JkOiAKc2NyYXBlX2NvbmZpZ3M6CiAgLSBqb2JfbmFtZTogY29udGFpbmVycwogICAgc3RhdGljX2NvbmZpZ3M6CiAgICAgIC0gdGFyZ2V0czoKICAgICAgICAgIC0gbG9jYWxob3N0CiAgICAgICAgbGFiZWxzOgogICAgICAgICAgam9iOiBjb250YWluZXJzCiAgICAgICAgICBfX3BhdGhfXzogL2xva2kvbG9ncy9jb250YWluZXJzCiAgICBwaXBlbGluZV9zdGFnZXM6CiAgICAtIHJlZ2V4OgogICAgICAgIGV4cHJlc3Npb246ICdeKD9QPG5hbWVzcGFjZT4uKilcLyg/UDxwb2Q+LiopXFsoP1A8Y29udGFpbmVyPi4qKVxdOiAoP1A8Y29udGVudD4uKiknCiAgICAtIGxhYmVsczoKICAgICAgICBuYW1lc3BhY2U6CiAgICAgICAgcG9kOgogICAgICAgIGNvbnRhaW5lcjoKICAgIC0gb3V0cHV0OgogICAgICAgIHNvdXJjZTogY29udGVudAogIC0gam9iX25hbWU6IGthaWwKICAgIHN0YXRpY19jb25maWdzOgogICAgICAtIHRhcmdldHM6CiAgICAgICAgICAtIGxvY2FsaG9zdAogICAgICAgIGxhYmVsczoKICAgICAgICAgIGpvYjoga2FpbAogICAgICAgICAgX19wYXRoX186IC9sb2tpL2xvZ3Mva2FpbAogICAgcGlwZWxpbmVfc3RhZ2VzOgogICAgLSByZWdleDoKICAgICAgICBleHByZXNzaW9uOiAnXnRpbWU9Iig/UDx0aW1lPi4qKSIgbGV2ZWw9KD9QPGxldmVsPi4qKSBtc2c9Iig/UDxjb250ZW50Pi4qKSIgY21wPSg/UDxjb21wb25lbnQ+LiopJwogICAgLSBsYWJlbHM6CiAgICAgICAgdGltZToKICAgICAgICBsZXZlbDoKICAgICAgICBjb21wb25lbnQ6CiAgICAtIHRpbWVzdGFtcDoKICAgICAgICBzb3VyY2U6IHRpbWUKICAgICAgICBmb3JtYXQ6IFJGQzMzMzkKICAgIC0gb3V0cHV0OgogICAgICAgIHNvdXJjZTogY29udGVudA=="
    },
    "kind": "Secret",
    "metadata": {
        "annotations": {
            "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Secret\",\"metadata\":{\"annotations\":{},\"name\":\"team-alpha-dev-log-recolector-config\",\"namespace\":\"team-alpha-dev\"},\"stringData\":{\"promtail.yaml\":\"server:\\n  http_listen_port: 9080\\n  grpc_listen_port: 0\\nclient:\\n  url: http://loki-multi-tenant-proxy.co-monitoring.svc.cluster.local:3100/api/prom/push\\n  basic_auth:\\n    username: alpha\\n    password: alpha\\nscrape_configs:\\n  - job_name: containers\\n    static_configs:\\n      - targets:\\n          - localhost\\n        labels:\\n          job: containers\\n          __path__: /loki/logs/containers\\n    pipeline_stages:\\n    - regex:\\n        expression: '^(?P\\u003cnamespace\\u003e.*)\\\\/(?P\\u003cpod\\u003e.*)\\\\[(?P\\u003ccontainer\\u003e.*)\\\\]: (?P\\u003ccontent\\u003e.*)'\\n    - labels:\\n        namespace:\\n        pod:\\n        container:\\n    - output:\\n        source: content\\n  - job_name: kail\\n    static_configs:\\n      - targets:\\n          - localhost\\n        labels:\\n          job: kail\\n          __path__: /loki/logs/kail\\n    pipeline_stages:\\n    - regex:\\n        expression: '^time=\\\"(?P\\u003ctime\\u003e.*)\\\" level=(?P\\u003clevel\\u003e.*) msg=\\\"(?P\\u003ccontent\\u003e.*)\\\" cmp=(?P\\u003ccomponent\\u003e.*)'\\n    - labels:\\n        time:\\n        level:\\n        component:\\n    - timestamp:\\n        source: time\\n        format: RFC3339\\n    - output:\\n        source: content\\n\"}}\n"
        },
        "creationTimestamp": "2020-12-28T13:19:27Z",
        "managedFields": [
            {
                "apiVersion": "v1",
                "fieldsType": "FieldsV1",
                "fieldsV1": {
                    "f:data": {
                        ".": {},
                        "f:promtail.yaml": {}
                    },
                    "f:metadata": {
                        "f:annotations": {
                            ".": {},
                            "f:kubectl.kubernetes.io/last-applied-configuration": {}
                        }
                    },
                    "f:type": {}
                },
                "manager": "kubectl-client-side-apply",
                "operation": "Update",
                "time": "2020-12-28T13:19:27Z"
            }
        ],
        "name": "team-alpha-dev-log-recolector-config",
        "namespace": "team-alpha-dev",
        "resourceVersion": "63853",
        "selfLink": "/api/v1/namespaces/team-alpha-dev/secrets/team-alpha-dev-log-recolector-config",
        "uid": "d6ca05ba-c265-4cb8-a23b-7a3b32453f72"
    },
    "type": "Opaque"
}
`
