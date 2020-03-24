package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct{}

func main() {
	server := mux.NewRouter()

	server.HandleFunc("/", index).Methods(http.MethodGet)
	server.HandleFunc("/ping", ping).Methods(http.MethodGet)
	server.HandleFunc("/healthz", healthz).Methods(http.MethodGet)
	server.HandleFunc("/hook", hook).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", server))
}
