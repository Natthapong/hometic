package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func setupRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/pair-device", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"active"}`))
	}).Methods(http.MethodPost)

	server := http.Server{
		Addr:    "127.0.0.1:2009",
		Handler: r,
	}
	log.Println("...Starting...")
	log.Fatal(server.ListenAndServe()) // listen and serve on 127.0.0.0:8080
}

func main() {
	setupRouter()
}
