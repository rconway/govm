package main

type GoVersion []struct {
	Files []struct {
		Arch     string `json:"arch"`
		Filename string `json:"filename"`
		Kind     string `json:"kind"`
		Os       string `json:"os"`
		Sha256   string `json:"sha256"`
		Size     int64  `json:"size"`
		Version  string `json:"version"`
	} `json:"files"`
	Stable  bool   `json:"stable"`
	Version string `json:"version"`
}
