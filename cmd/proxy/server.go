package main

import (
	"encoding/json"
	"io/ioutil"
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
	var atlantisURL string
	var requestBody []byte

	eventType := r.Header.Get(bitbucketEventTypeHeader)

	defer r.Body.Close()

	// store the body
	// this will be the data we sent to atlantis
	body, err := ioutil.ReadAll(r.Body)

	requestBody = body
	requestHeader := r.Header

	log.Printf("Request body %s", requestBody)
	log.Printf("Request header %s", requestHeader)

	// log.Printf("%s", body)

	err = json.NewDecoder(r.Body).Decode(&c)

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

	commitHash := c.PullRequest.Source.Commit.Hash

	log.Printf("Got webhook with event type : %s", eventType)
	log.Printf("Commit Hash : %s", commitHash)

	// Checkout to commit hash
	environment, err := gitClone(commitHash)
	if err != nil {
		log.Printf("Got error when cloning repository : %s", err)
	}

	log.Printf("%s", environment)

	// set url
	if environment == "dev" {
		atlantisURL = "https://atlantis.ext.bit-stack.net"
	} else if environment == "prd" {
		atlantisURL = "http://atlantis.ovo.co.id"
	} else {
		atlantisURL = ""
	}

	if atlantisURL != "" {
		log.Printf("Proxying bitbucket hook to atlantis server at %s", atlantisURL)
	} else {
		log.Printf("Cannot find atlantis URL for environment %s", environment)
	}

}
