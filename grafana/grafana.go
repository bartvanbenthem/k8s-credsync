package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Datasource struct {
	Name          string `json:"name"`
	ID            int    `json:id`
	Type          string `json:"type"`
	URL           string `json:"url"`
	Access        string `json:"access"`
	OrgID         int    `json:"orgId"`
	BasicAuth     bool   `json:"basicAuth"`
	BasicAuthUser string `json:"basicAuthUser"`
	ReadOnly      bool   `json: "readOnly"`
	JSONData      struct {
		HTTPHeaderName1 string `json:"httpHeaderName1"`
	} `json:"jsonData"`
	SecureJSONData struct {
		HTTPHeaderValue1  string `json:"httpHeaderValue1"`
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

// switch to a specific grafana context
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

// get a specific grafana organization based on the org name
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

// get all the grafana organizations
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

// create a new grafana organizations
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

// get grafana datasource
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

// create a new grafana datasource
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

// update existing grafana datasource
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

// delete grafana datasource
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

// function for making web requests with basic auth
func RequestAUTH(method, url string, body []byte) ([]byte, error) {
	var err error
	var client *http.Client
	var req *http.Request
	client = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err = http.NewRequest(method, url,
		bytes.NewBuffer(body))

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
