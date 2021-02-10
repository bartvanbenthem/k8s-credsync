package grafana

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Datasource struct {
	Name           string `json:"name"`
	ID             int    `json:id`
	Type           string `json:"type"`
	URL            string `json:"url"`
	Access         string `json:"access"`
	OrgID          int    `json:"orgId"`
	BasicAuth      bool   `json:"basicAuth"`
	BasicAuthUser  string `json:"basicAuthUser"`
	SecureJSONData struct {
		BasicAuthPassword string `json:"basicAuthPassword"`
	} `json:"secureJsonData"`
}

type Organization struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address struct {
		Address1 string `json:"address1"`
		Address2 string `json:"address2"`
		City     string `json:"city"`
		ZipCode  string `json:"zipCode"`
		State    string `json:"state"`
		Country  string `json:"country"`
	} `json:"address"`
}

type Organizations []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func SwitchUserContext(org Organization) error {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/user/using/%v", grafanapi, org.ID)
	o, err := json.Marshal(&org)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
		return err
	}
	log.Printf("Switching context to %v Organization\n", org.Name)
	data, err := RequestAUTH("POST", url, o)
	if err != nil {
		return err
	}
	log.Printf("%v\n", string(data))
	return err
}

func GetOrganization(orgname string) (Organization, error) {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/orgs/name/%v", grafanapi, orgname)
	var org Organization
	data, err := RequestAUTH("GET", url, []byte(""))
	if err != nil {
		return org, err
	}
	err = json.Unmarshal(data, &org)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
		return org, err
	}
	return org, err
}

func CreateOrganization(org Organization) error {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/orgs", grafanapi)
	b, err := json.Marshal(&org)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
		return err
	}
	log.Printf("Create \"%v\" Grafana Organization\n", org.Name)
	data, err := RequestAUTH("POST", url, b)
	if err != nil {
		return err
	}
	log.Printf("%v\n", string(data))
	return err
}

func GetAllOrganizations() (Organizations, error) {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/orgs", grafanapi)
	var orgs Organizations
	data, err := RequestAUTH("GET", url, []byte(""))
	if err != nil {
		return orgs, err
	}
	err = json.Unmarshal(data, &orgs)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
		return orgs, err
	}
	return orgs, err
}

func GetDatasource(dsname string) (Datasource, error) {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/datasources/name/%v", grafanapi, dsname)
	var ds Datasource
	data, err := RequestAUTH("GET", url, []byte(""))
	if err != nil {
		return ds, err
	}
	err = json.Unmarshal(data, &ds)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
		return ds, err
	}
	return ds, err
}

func CreateDatasource(ds Datasource) error {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/datasources", grafanapi)
	b, err := json.Marshal(&ds)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
		return err
	}
	log.Printf("Create \"%v\" Grafana Datasource\n", ds.Name)
	data, err := RequestAUTH("POST", url, b)
	if err != nil {
		return err
	}
	log.Printf("%v\n", string(data))
	return err
}

func UpdateDatasource(ds Datasource) error {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/datasources/%v", grafanapi, ds.ID)
	b, err := json.Marshal(&ds)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
		return err
	}
	log.Printf("Update \"%v\" Grafana Datasource\n", ds.Name)
	data, err := RequestAUTH("PUT", url, b)
	if err != nil {
		return err
	}
	log.Printf("%v\n", string(data))
	return err
}

func DeleteDatasource(ds Datasource) error {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/datasources/%v", grafanapi, ds.ID)
	b, err := json.Marshal(&ds)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
		return err
	}
	log.Printf("Delete \"%v\" Grafana Datasource\n", ds.Name)
	data, err := RequestAUTH("DELETE", url, b)
	if err != nil {
		return err
	}
	log.Printf("%v\n", string(data))
	return err
}

func RequestAUTH(method, url string, body []byte) ([]byte, error) {
	caFile := os.Getenv("K8S_GRAFANA_CA_FILE")
	var err error
	var client *http.Client
	var req *http.Request
	if len(caFile) == 0 {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
		url := fmt.Sprintf("http://%v", url)
		req, err = http.NewRequest(method, url,
			bytes.NewBuffer(body))
	} else {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Printf("FATAL ERROR: %v\n", err)
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:            caCertPool,
					InsecureSkipVerify: true, // enable for self signed certificates
				},
			},
		}
		tls := fmt.Sprintf("https://%v", url)
		req, err = http.NewRequest(method, tls,
			bytes.NewBuffer(body))
	}

	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	// import grafana credentials from environment var
	// needed for basic authentication
	user := os.Getenv("K8S_GRAFANA_BA_USER")
	pass := os.Getenv("K8S_GRAFANA_BA_PASSWORD")
	req.SetBasicAuth(user, pass)
	response, err := client.Do(req)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	defer response.Body.Close()
	// read response body
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("%v\n", err)
		return data, err
	}
	return data, err
}
