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

var (
	certFile = os.Getenv("K8S_GRAFANA_CERT_FILE")
	keyFile  = os.Getenv("K8S_GRAFANA_KEY_FILE")
	caFile   = os.Getenv("K8S_GRAFANA_CA_FILE")
)

func SwitchUserContext(org Organization) {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/user/using/%v", grafanapi, org.ID)
	o, err := json.Marshal(&org)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
	}
	log.Printf("Switching context to %v Organization\n", org.Name)
	data := RequestAUTH("POST", url, o)
	log.Printf("%v\n", string(data))
}

func GetOrganization(orgname string) Organization {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/orgs/name/%v", grafanapi, orgname)
	data := RequestAUTH("GET", url, []byte(""))
	var org Organization
	err := json.Unmarshal(data, &org)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
	}
	return org
}

func CreateOrganization(org Organization) {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/orgs", grafanapi)
	b, err := json.Marshal(&org)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
	}
	log.Printf("Create \"%v\" Grafana Organization\n", org.Name)
	data := RequestAUTH("POST", url, b)
	log.Printf("%v\n", string(data))
}

func GetDatasource(dsname string) Datasource {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/datasources/name/%v", grafanapi, dsname)
	data := RequestAUTH("GET", url, []byte(""))
	var ds Datasource
	err := json.Unmarshal(data, &ds)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
	}
	return ds
}

func CreateDatasource(ds Datasource) {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/datasources", grafanapi)
	b, err := json.Marshal(&ds)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
	}
	log.Printf("Create \"%v\" Grafana Datasource\n", ds.Name)
	data := RequestAUTH("POST", url, b)
	log.Printf("%v\n", string(data))
}

func UpdateDatasource(ds Datasource) {
	grafanapi := os.Getenv("K8S_GRAFANA_API_URL")
	url := fmt.Sprintf("%v/datasources/%v", grafanapi, ds.ID)
	b, err := json.Marshal(&ds)
	if err != nil {
		log.Printf("Error encoding yaml: %v", err)
	}
	log.Printf("\nUpdate \"%v\" Grafana Datasource\n", ds.Name)
	data := RequestAUTH("PUT", url, b)
	log.Printf("%v\n", string(data))
}

func RequestAUTH(method, url string, body []byte) []byte {
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
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
					// InsecureSkipVerify: true
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
	}
	// import grafana credentials from environment var
	// needed for basic authentication
	user := os.Getenv("K8S_GRAFANA_BA_USER")
	pass := os.Getenv("K8S_GRAFANA_BA_PASSWORD")
	req.SetBasicAuth(user, pass)
	response, err := client.Do(req)
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer response.Body.Close()
	// read response body
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("%v\n", err)
	}
	return data
}
