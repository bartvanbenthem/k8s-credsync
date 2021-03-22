package grafana

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

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
