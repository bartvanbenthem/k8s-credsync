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
	cert := os.Getenv("K8S_SERVER_CERT")
	key := os.Getenv("K8S_SERVER_KEY")
	// http handler functions
	http.HandleFunc("/", HandlerDefault)
	http.HandleFunc("/proxy/sync", HandlerProxySync())
	http.HandleFunc("/grafana/sync", HandlerGrafanaSync())
	// listen and serve https connections
	log.Printf("About to listen on https://%v/\n", address)
	err := http.ListenAndServeTLS(address, cert, key, nil)
	if err != nil {
		log.Fatal(err)
	}
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
		sync.Proxy()
		io.WriteString(w, `{"sync":"finished"}`)
	})
}

// handler for grafana synchronization service
func HandlerGrafanaSync() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		sync.Grafana()
		io.WriteString(w, `{"sync":"finished"}`)
	})
}
