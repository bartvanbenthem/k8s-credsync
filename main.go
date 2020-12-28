package main

import (
	"fmt"
	"log"
)

func main() {
	// get base64 encoded proxy secret
	proxy, err := getEncodedSecret(secretProxy, "\"authn.yaml\":")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// get base64 encoded tenant secret
	tenant, err := getEncodedSecret(secretTenantEmptyPassword, "\"promtail.yaml\":")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// test by printing decoded values
	//fmt.Printf("\nproxy\n-----\n%v\n", decodeSecret(proxy))
	//fmt.Printf("\ntenant\n------\n%v\n", decodeSecret(tenant))

	proxycred, err := getProxyCredentials(decodeSecret(proxy))
	if err != nil {
		log.Printf("error: %v", err)
	}

	tenantcred, err := getTenantCredential(decodeSecret(tenant))
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
	p := tenantcred.Client.BasicAuth.Password
	if p == "" || p == " " {
		tenantcred.Client.BasicAuth.Password = passwordGenerator()
	}
	fmt.Printf("username:%v password:%v\n",
		tenantcred.Client.BasicAuth.Username,
		tenantcred.Client.BasicAuth.Password)
}
