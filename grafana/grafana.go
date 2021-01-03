package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Datasource struct {
	Name           string `json:"name"`
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

func GetDatasource() {

}

func CreateDatasource() {

}

func GetOrganization(orgname string) Organization {
	url := fmt.Sprintf("http://grafana/api/orgs/name/%v", orgname)
	data := RequestAUTH("GET", url, []byte(""))
	var org Organization
	err := json.Unmarshal(data, &org)
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}
	return org
}

func CreateOrganization(org Organization) {
	url := fmt.Sprintf("http://grafana/api/orgs")
	b, err := json.Marshal(&org)
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}
	fmt.Printf("\nCreating %v Grafana Organization\n", org.Name)
	data := RequestAUTH("POST", url, b)
	fmt.Printf("%v\n", string(data))
}

func RequestAUTH(method, url string, body []byte) []byte {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}
	// import grafana credentials from environment var
	user := os.Getenv("K8S_GRAFANA_BA_USER")
	pass := os.Getenv("K8S_GRAFANA_BA_PASSWORD")
	req.SetBasicAuth(user, pass)
	response, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}
	return data
}