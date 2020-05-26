package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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
	var err error
	var atlantisURL string

	eventType := r.Header.Get(bitbucketEventTypeHeader)

	requestBody, err := ioutil.ReadAll(r.Body)
	requestHeader := r.Header

	defer r.Body.Close()

	log.Printf("request header %s", requestHeader)

	err = json.Unmarshal(requestBody, &c)

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
		atlantisURL = os.Getenv("ATLANTIS_DEV")
	} else if environment == "stg" {
		atlantisURL = os.Getenv("ATLANTIS_STG")
	} else if environment == "prd" {
		atlantisURL = os.Getenv("ATLANTIS_PRD")
	} else {
		atlantisURL = ""
	}

	if atlantisURL != "" {
		log.Printf("Proxying bitbucket hook to atlantis server at %s", atlantisURL)
		timeout := time.Duration(5 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		log.Printf("Request Body : %s", requestBody)

		if err != nil {
			log.Printf("unable to marshal request body %s", err)
		}

		request, err := http.NewRequest("POST", atlantisURL+"/events", bytes.NewBuffer(requestBody))

		for key, values := range requestHeader {
			for _, value := range values {
				request.Header.Set(key, value)
			}
		}

		if err != nil {
			log.Printf("Error when sending request to atlantis %s", err)
		}

		resp, err := client.Do(request)

		if err != nil {
			log.Printf("Got error from response %s", err)
		}

		defer resp.Body.Close()

		log.Printf("Response status from atlantis %v", resp.StatusCode)

	} else {
		log.Printf("Cannot find atlantis URL for environment %s", environment)
	}

}
