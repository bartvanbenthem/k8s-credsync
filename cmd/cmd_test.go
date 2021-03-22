package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bartvanbenthem/k8s-ntenant/credential"
	"github.com/bartvanbenthem/k8s-ntenant/tenant"
)

func TestGetCredentials(t *testing.T) {
	// import environment variables
	csec := os.Getenv("K8S_CRED_SECRET_NAME")
	cns := os.Getenv("K8S_CRED_SECRET_NAMESPACE")
	// Prints the current proxycredentials
	pcurrent, err := credential.AllCredentials(cns, csec)
	if err != nil {
		log.Printf("%v\n", err)
	}
	fmt.Printf("\nProxy\n-----\n")
	for _, pc := range pcurrent.Users {
		fmt.Printf("User: %v Password: %v org: %v\n",
			pc.Username, pc.Password, pc.Orgid)
	}
}

// Prints the current tenant and proxy credentials
func TestMainFunctions(t *testing.T) {
	// import environment variables
	csec := os.Getenv("K8S_CRED_SECRET_NAME")
	cns := os.Getenv("K8S_CRED_SECRET_NAMESPACE")
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")
	// Prints the current tenant credentials
	tcurrent, err := tenant.AllTenantCredentials(tenantsec)
	fmt.Printf("\nTenant\n------\n")
	for _, tc := range tcurrent {
		fmt.Printf("User:%v Password:%v\n",
			tc.Client.BasicAuth.Username,
			tc.Client.BasicAuth.Password)
	}

	// Prints the current proxycredentials
	pcurrent, err := credential.AllCredentials(cns, csec)
	if err != nil {
		log.Printf("%v\n", err)
	}
	fmt.Printf("\nProxy\n-----\n")
	for _, pc := range pcurrent.Users {
		fmt.Printf("User:%v Password:%v org:%v\n",
			pc.Username, pc.Password, pc.Orgid)
	}
}
