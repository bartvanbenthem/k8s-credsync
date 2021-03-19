package credential

import (
	"log"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/kube"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
)

type Users struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Orgid    string `yaml:"orgid"`
	TenantID string `yaml:"tenantid"`
}

type Credentials struct {
	Users []struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Orgid    string `yaml:"orgid"`
		TenantID string `yaml:"tenantid"`
	} `yaml:"users"`
}

// input is a decoded yaml config file from the secret
func GetCredentials(file string) (Credentials, error) {
	var err error
	var c Credentials
	// unmarshall entire tenant JSON into a map
	err = yaml.Unmarshal([]byte(file), &c)
	if err != nil {
		return c, err
	}
	return c, err
}

func UpdateCredentialSecret(namespace, datafield string, newc Credentials) error {
	// import required environment variables
	csec := os.Getenv("K8S_CRED_SECRET_NAME")
	var kube kube.KubeCLient
	// marshall proxy credentials into byte slice
	credbyte, err := yaml.Marshal(&newc)
	if err != nil {
		log.Printf("Error encoding yaml: %v\n", err)
		return err
	}

	// get current secret
	sec := kube.GetSecret(kube.CreateClientSet(), namespace, csec)
	sec.Data[datafield] = credbyte
	// create new secret object
	var newsecret v1.Secret
	newsecret.Kind = sec.Kind
	newsecret.APIVersion = sec.APIVersion
	newsecret.Data = map[string][]byte{datafield: credbyte}
	newsecret.Name = sec.Name
	newsecret.Namespace = sec.Namespace
	// Update Credential secret with tenant credentials
	kube.UpdateSecret(kube.CreateClientSet(), namespace, &newsecret)
	// get/validate secret
	_ = kube.GetSecret(kube.CreateClientSet(), namespace, newsecret.Name)
	return err
}

// collects all credentials
func AllCredentials() (Credentials, error) {
	var err error
	// import environment variables
	csec := os.Getenv("K8S_CRED_SECRET_NAME")
	cns := os.Getenv("K8S_CRED_SECRET_NAMESPACE")
	// initiate kube client
	var kube kube.KubeCLient
	// get the credentials
	cred, err := GetCredentials(string(
		kube.GetSecretData(kube.CreateClientSet(),
			cns, csec, "authn.yaml")))
	if err != nil {
		return cred, err
	}
	return cred, err
}
