package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bartvanbenthem/k8s-ntenant/sync"
)

func main() {
	// get from environment variables
	address := os.Getenv("K8S_SERVER_ADDRESS")
	//cert := os.Getenv("K8S_SERVER_CERT")
	//key := os.Getenv("K8S_SERVER_KEY")

	// http handler functions
	http.HandleFunc("/", HandlerDefault)
	http.HandleFunc("/proxy/sync", HandlerProxySync())
	http.HandleFunc("/grafana/sync", HandlerGrafanaSync())

	// listen and serve http connections
	log.Printf("About to listen on http://%v/\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}

	/*
		// listen and serve https connections
		log.Printf("About to listen on https://%v/\n", address)
		err := http.ListenAndServeTLS(address, cert, key, nil)
		if err != nil {
			log.Fatal(err)
		}
	*/
}

// default handler
func HandlerDefault(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"server":"running"}`)
}

// handler for proxy synchronization service
func HandlerProxySync() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		err := sync.Proxy()
		if err != nil {
			log.Printf("Proxy sync completed with errors inspect log")
			io.WriteString(w, `{"proxy":"sync completed with errors inspect log"}`)
		} else {
			log.Printf("Proxy sync completed")
			io.WriteString(w, `{"proxy":" sync completed"}`)
		}
	})
}

// handler for grafana synchronization service
func HandlerGrafanaSync() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		err := sync.Grafana()
		if err != nil {
			log.Printf("Grafana sync completed with errors inspect log")
			io.WriteString(w, `{"grafana":"sync completed with errors inspect log"}`)
		} else {
			log.Printf("Grafana synchronization completed")
			io.WriteString(w, `{"grafana":" sync completed"}`)
		}
	})
}
