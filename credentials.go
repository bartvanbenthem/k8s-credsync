package main

import (
	"bufio"
	"encoding/base64"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type ProxyCredentials struct {
	Users []struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Orgid    string `yaml:"orgid"`
	} `yaml:"users"`
}

type TenantCredential struct {
	Client struct {
		URL       string `yaml:"url"`
		BasicAuth struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"basic_auth"`
	} `yaml:"client"`
}

func getProxyCredentials(file string) (ProxyCredentials, error) {
	var err error
	var c ProxyCredentials
	// unmarshall entire tenant JSON into a map
	err = yaml.Unmarshal([]byte(file), &c)
	if err != nil {
		return c, err
	}
	return c, err
}

func getTenantCredential(file string) (TenantCredential, error) {
	var err error
	var c TenantCredential
	// unmarshall entire tenant JSON into a map
	err = yaml.Unmarshal([]byte(file), &c)
	if err != nil {
		return c, err
	}
	return c, err
}

func decodeSecret(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatalf("Error decoding: %v", err)
	}
	return string(decoded)
}

func getEncodedSecret(json, partial string) (string, error) {
	var err error
	var lines []string
	// Scan all the lines in sd byte slice
	// append every line to the lines slice of string
	scanner := bufio.NewScanner(strings.NewReader(json))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err != nil {
		return "", err
	}
	// check every line on the given partial
	// split the line on :
	for _, line := range lines {
		if strings.Contains(line, partial) {
			lines = strings.Split(line, ":")
		}
	}
	// remove unwanted charachters and spaces
	str := lines[1]
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, " ", "")

	return str, err
}
