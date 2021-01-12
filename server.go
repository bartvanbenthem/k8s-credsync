package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	address := os.Getenv("K8S_SERVER_ADDRESS")
	cert := os.Getenv("K8S_SERVER_CERT")
	key := os.Getenv("K8S_SERVER_KEY")

	http.HandleFunc("/", Handler)
	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	log.Printf("About to listen on 8443. Go to https://localhost:8443/")
	err := http.ListenAndServeTLS(address, cert, key, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}
