package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Users struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Orgid    string `yaml:"orgid"`
}

type ProxyCredentials struct {
	Users []struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Orgid    string `yaml:"orgid"`
	} `yaml:"users"`
}

type TenantCredential struct {
	Server struct {
		HTTPListenPort int `yaml:"http_listen_port"`
		GrpcListenPort int `yaml:"grpc_listen_port"`
	} `yaml:"server"`
	Client struct {
		URL       string `yaml:"url"`
		BasicAuth struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"basic_auth"`
	} `yaml:"client"`
	ScrapeConfigs []struct {
		JobName       string `yaml:"job_name"`
		StaticConfigs []struct {
			Targets []string `yaml:"targets"`
			Labels  struct {
				Job  string `yaml:"job"`
				Path string `yaml:"__path__"`
			} `yaml:"labels"`
		} `yaml:"static_configs"`
		PipelineStages []struct {
			Regex struct {
				Expression string `yaml:"expression"`
			} `yaml:"regex,omitempty"`
			Labels struct {
				Namespace interface{} `yaml:"namespace"`
				Pod       interface{} `yaml:"pod"`
				Container interface{} `yaml:"container"`
			} `yaml:"labels,omitempty"`
			Output struct {
				Source string `yaml:"source"`
			} `yaml:"output,omitempty"`
		} `yaml:"pipeline_stages"`
	} `yaml:"scrape_configs"`
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

func UpdateProxySecret(namespace, datafield string, newc ProxyCredentials) {
	// import required environment variables
	proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	var kube KubeCLient

	credbyte, err := yaml.Marshal(&newc)
	if err != nil {
		fmt.Printf("\nerror encoding yaml: %v\n", err)
	}

	sec := kube.GetSecret(kube.CreateClientSet(), namespace, proxysec)
	sec.Data[datafield] = credbyte
	// update secret ignore output validation is in main
	_ = kube.UpdateSecret(kube.CreateClientSet(), namespace, sec)
}

// scan the tenant credentials
// create and update password when empty
func UpdateTenantSecret(namespace, datafield string) {
	// import required environment variables
	tenantsec := os.Getenv("K8S_TENANT_SECRET_NAME")
	var kube KubeCLient
	// create and set password in the tenant credential when password is empty
	cred, err := GetTenantCredential(string(kube.GetSecretData(kube.CreateClientSet(),
		namespace, tenantsec, datafield)))
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
	if len(cred.Client.BasicAuth.Username) != 0 {
		PasswordSetter(&cred)
	}

	credbyte, err := yaml.Marshal(&cred)
	if err != nil {
		fmt.Printf("\nerror encoding yaml: %v\n", err)
	}

	sec := kube.GetSecret(kube.CreateClientSet(), namespace, tenantsec)
	sec.Data[datafield] = credbyte
	// update secret ignore output validation is in main
	_ = kube.UpdateSecret(kube.CreateClientSet(), namespace, sec)

}

func PasswordSetter(t *TenantCredential) *TenantCredential {
	if len(t.Client.BasicAuth.Password) == 0 {
		t.Client.BasicAuth.Password = PasswordGenerator()
	}
	return t
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

////////////////////////////////////////////////////////////////////
// ONLY USED FOR RAW ENCODED JSON RESPONSE FROM THE KUBERNETES API /
////////////////////////////////////////////////////////////////////

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
