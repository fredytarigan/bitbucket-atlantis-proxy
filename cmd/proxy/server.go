package main

import "net/http"

/*
Home handler
URI Path : "/"
*/
func index(w http.ResponseWriter, r *http.Request) {

}

/*
Healthz handler
URI Path : "/healthz"
*/
func healthz(w http.ResponseWriter, r *http.Request) {

}

/*
Ping handler
URI Path : "/ping"
*/
func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"pong"}`))
}

/*
Hook handler
URI Path : "/hook"
*/
func hook(w http.ResponseWriter, r *http.Request) {

}
