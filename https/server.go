package main

import (
	"net/http"
	"fmt"
)

func main() {
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: nil,
	}
	http.HandleFunc("/", echoHandler)
	server.ListenAndServeTLS("cert.pem", "key.pem")
}

func echoHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)
}
