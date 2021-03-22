package tenant

import (
	"github.com/bartvanbenthem/k8s-ntenant/kube"
	"gopkg.in/yaml.v2"
)

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
		TenantID string `yaml:"tenant_id"`
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

// collects all tenant credentials
func AllTenantCredentials(tenantsecret string) ([]TenantCredential, error) {
	var err error
	// initiate kube client
	var kube kube.KubeCLient
	//create slice of tenant credential
	var tcreds []TenantCredential

	namespaces := kube.GetAllNamespaceNames(kube.CreateClientSet())
	for _, ns := range namespaces {
		var c TenantCredential
		s := kube.GetSecretData(kube.CreateClientSet(),
			ns, tenantsecret, "promtail.yaml")
		if len(s) != 0 {
			err = yaml.Unmarshal(s, &c)
			if err != nil {
				return nil, err
			}
			tcreds = append(tcreds, c)
		}
	}
	return tcreds, err
}
