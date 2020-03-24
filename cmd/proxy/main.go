package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type server struct{}

func main() {
	server := mux.NewRouter()

	server.HandleFunc("/", index).Methods(http.MethodGet)
}
