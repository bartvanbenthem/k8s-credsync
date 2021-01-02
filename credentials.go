package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
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

func ReplaceProxySecret(namespace, datafield string, newc ProxyCredentials) {
	// import required environment variables
	proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	var kube KubeCLient

	credbyte, err := yaml.Marshal(&newc)
	if err != nil {
		fmt.Printf("\nerror encoding yaml: %v\n", err)
	}

	sec := kube.GetSecret(kube.CreateClientSet(), namespace, proxysec)
	sec.Data[datafield] = credbyte

	// create new secret object
	var newsecret v1.Secret
	newsecret.Kind = sec.Kind
	newsecret.APIVersion = sec.APIVersion
	newsecret.Data = map[string][]byte{datafield: credbyte}
	newsecret.Name = sec.Name
	newsecret.Namespace = sec.Namespace

	// delete secret
	kube.DeleteSecret(kube.CreateClientSet(), namespace, sec)

	// ONLY FOR TESTING PURPOSES
	// CREATE A SECRET EXISTS CHECK FUNCTION INSTEAD OF SLEEP !!!!!!!!
	time.Sleep(5 * time.Second)

	// create secret
	_ = kube.CreateSecret(kube.CreateClientSet(), namespace, &newsecret)

	_ = kube.GetSecret(kube.CreateClientSet(), namespace, newsecret.Name)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
