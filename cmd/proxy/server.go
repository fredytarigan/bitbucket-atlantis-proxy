package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type StandardResponse struct {
	Message    string
	StatusCode string
}

type CommitHashID struct {
	PullRequest PullRequest `json:"pullrequest"`
}

type PullRequest struct {
	Source Source `json:"source"`
}

type Source struct {
	Commit Commit `json:"commit"`
}

type Commit struct {
	Hash string `json:"hash"`
}

const bitbucketEventTypeHeader = "X-Event-Key"
const bitbucketCloudRequestIDHeader = "X-Request-UUID"

/*
Home handler
URI Path : "/"
*/
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := StandardResponse{
		Message:    "Welcome !!!",
		StatusCode: "200",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonData)
}

/*
Healthz handler
URI Path : "/healthz"
*/
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := StandardResponse{
		Message:    "Application is healthy",
		StatusCode: "200",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonData)
}

/*
Ping handler
URI Path : "/ping"
*/
func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"pong"`))
}

/*
Hook handler
URI Path : "/hook"
*/
func hook(w http.ResponseWriter, r *http.Request) {
	var c CommitHashID

	eventType := r.Header.Get(bitbucketEventTypeHeader)

	defer r.Body.Close()
	// body, err := ioutil.ReadAll(r.Body)

	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		data := StandardResponse{
			Message:    "Unable to read data body: ",
			StatusCode: "504",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		jsonData, _ := json.Marshal(data)
		w.Write(jsonData)
		log.Fatal(err)
	}

	log.Printf("Got webhook with event type : %s", eventType)
	log.Printf("Commit Hash : %s", c.PullRequest.Source.Commit.Hash)
}
