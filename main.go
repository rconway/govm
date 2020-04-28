package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func init() {
	log.Println("...init(main)...")
}

// GoVersion JSON structure for a go version
type GoVersion struct {
	Version *string `json:"version"`
}

func hadAnError(err error) {
	log.Printf("ERROR: %v\n", err)
}

func main() {
	log.Println("...main...")

	// Create a web client with a connection/request timeout
	webClient := &http.Client{Timeout: 3 * time.Second}

	// Get the JSON version list from the endpoint
	resp, err := webClient.Get("https://golang.org/dl/?mode=json")
	if err != nil {
		hadAnError(fmt.Errorf("Problem retrieving metadata: %w", err))
		return
	}
	defer resp.Body.Close()

	// Read the data bytes from the response body
	responseBodyBytes, err := ioutil.ReadAll(resp.Body)

	// Dump some info from the response
	log.Println("Response status:", resp.Status)
	log.Println("Response body:", string(responseBodyBytes))

	// Parse/decode the json response into our data structure (an array of go versions)
	var data []GoVersion
	err = json.NewDecoder(bytes.NewBuffer(responseBodyBytes)).Decode(&data)
	if err != nil {
		hadAnError(fmt.Errorf("Problem parsing metadata: %w", err))
		return
	}

	// Print out the versions found
	for _, v := range data {
		log.Println(*v.Version)
	}
}
