package internal

import (
	"log"
	"net/http"
	"time"
)

type Status int

const (
	Pending Status = iota
	InProgress
	Paused
	Cancelled
	Completed
)

type Download struct {
	URL            string
	Destination    string
	OutputFileName string
	Queue          *Queue
	Status
	contentLength int64
}

func (d *Download) SupportsPartialDownload() bool {
	req, err := http.NewRequest("HEAD", d.URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to get response from server")
		return false
	}

	d.contentLength = resp.ContentLength
	if resp.Header.Get("Accept-Ranges") == "" || resp.Header.Get("Accept-Ranges") == "none" {
		log.Fatal("Server does not support partial download")
		return false
	}

	return true
}
