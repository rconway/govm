package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

var goVersionsURL = "https://golang.org/dl/?mode=json"

func init() {
	log.Println("...init(main)...")
	log.Printf("os: %v, arch: %v\n", runtime.GOOS, runtime.GOARCH)
}

// GoVersion JSON structure for a go version
type GoVersion struct {
	Version *string `json:"version"`
	Stable  *bool   `json:"stable"`
	Files   []struct {
		Filename *string `json:"filename"`
		OS       *string `json:"os"`
		Arch     *string `json:"arch"`
		Version  *string `json:"version"`
		Sha256   *string `json:"sha256"`
		Size     *int    `json:"size"`
		Kind     *string `json:"kind"`
	} `json:"files"`
}

func hadAnError(err error) {
	log.Printf("ERROR: %v\n", err)
}

func main() {
	log.Println("...main...")

	// Create a web client with a connection/request timeout
	webClient := &http.Client{Timeout: 3 * time.Second}

	// Get the JSON version list from the endpoint
	resp, err := webClient.Get(goVersionsURL)
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
	for _, version := range data {
		if version.Version != nil {
			log.Println(*version.Version)
			// For each version print out the files
			for _, file := range version.Files {
				if file.Filename != nil {
					log.Printf("  %v (%v)\n", *file.Filename, *file.Size)
				}
			}
		}
	}

	// Deduce the latest version
	latestVersion, latestFilename := findLatestVersion(data)
	if len(latestVersion) > 0 {
		log.Printf("Latest version is: %v (%v)\n", latestVersion, latestFilename)
	} else {
		log.Printf("Cannot deduce latest version\n")
	}
}

func findLatestVersion(data []GoVersion) (string, string) {
	latestVersion := ""
	latestFilename := ""
	os := runtime.GOOS
	arch := runtime.GOARCH
	// Iterate over the versions
	for _, version := range data {
		// Check version is set
		if version.Version != nil {
			// If 'latest' is not yet set or if this version is higher (later)
			if len(latestVersion) == 0 || *version.Version > latestVersion {
				// Iterate through the files looking for those of type 'archive' that match our os/arch
				for _, file := range version.Files {
					if file.Kind != nil && *file.Kind == "archive" {
						if file.OS != nil && *file.OS == os && file.Arch != nil && *file.Arch == arch {
							// Match is found so set as latest
							latestVersion = *version.Version
							latestFilename = *file.Filename
						}
					}
				}
			}
		}
	}
	return latestVersion, latestFilename
}
