package main

import (
	"bufio"
	"encoding/base64"
	"log"
	"math/rand"
	"strings"
	"time"

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

// input is a decoded yaml config file from the secret
func GetProxyCredentials(file string) (ProxyCredentials, error) {
	var err error
	var c ProxyCredentials
	// unmarshall entire tenant JSON into a map
	err = yaml.Unmarshal([]byte(file), &c)
	if err != nil {
		return c, err
	}
	return c, err
}

// input is a decoded yaml config file from the secret
func GetTenantCredential(file string) (TenantCredential, error) {
	var err error
	var c TenantCredential
	// unmarshall entire tenant JSON into a map
	err = yaml.Unmarshal([]byte(file), &c)
	if err != nil {
		return c, err
	}
	return c, err
}

func PasswordGenerator() string {
	var str string
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 12
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str = b.String()
	return str
}

func DecodeSecret(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatalf("Error decoding: %v", err)
	}
	return string(decoded)
}

// extract the encoded secret from the k8s json response
func GetEncodedSecret(jsonresponse, partial string) (string, error) {
	var err error
	var lines []string
	// Scan all the lines in sd byte slice
	// append every line to the lines slice of string
	scanner := bufio.NewScanner(strings.NewReader(jsonresponse))
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
