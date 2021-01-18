package proxy

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/bartvanbenthem/k8s-ntenant/kube"
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

func ReplaceProxySecret(namespace, datafield string, newc ProxyCredentials) error {
	// import required environment variables
	proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	var kube kube.KubeCLient

	credbyte, err := yaml.Marshal(&newc)
	if err != nil {
		log.Printf("Error encoding yaml: %v\n", err)
		return err
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
	// get/validate secret
	_ = kube.GetSecret(kube.CreateClientSet(), namespace, newsecret.Name)
	if err != nil {
		log.Printf("%v\n", err)
		return err
	}
	return err
}

// collects all proxy credentials
func AllProxyCredentials() (ProxyCredentials, error) {
	var err error
	// import environment variables
	proxysec := os.Getenv("K8S_PROXY_SECRET_NAME")
	proxyns := os.Getenv("K8S_PROXY_SECRET_NAMESPACE")
	// initiate kube client
	var kube kube.KubeCLient
	// get the proxy credentials
	proxycred, err := GetProxyCredentials(string(
		kube.GetSecretData(kube.CreateClientSet(),
			proxyns, proxysec, "authn.yaml")))
	if err != nil {
		return proxycred, err
	}
	return proxycred, err
}

func RestartProxy(namespace, podname string) {
	// initiate kube client
	var kube kube.KubeCLient
	// restart proxy pod by deleting pod
	// the replicaset will create a new pod with updated config
	pods := kube.GetAllPodNames(kube.CreateClientSet(), namespace)
	for _, p := range pods {
		if strings.Contains(p, podname) {
			kube.DeletePod(kube.CreateClientSet(), namespace, p)
		}
	}
}
