package internal

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/volume/service/opts"
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
	headResp *http.Response
	contentLength int
}

func (d *Download) setHttpResponse() {
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
		return
	}
	
	d.headResp = resp
}

func (d *Download) setContentLength() {
	d.contentLength = int(d.headResp.ContentLength)
}

func (d *Download) supportsPartialDownload() bool {
	if d.headResp.Header.Get("Accept-Ranges") == "" || d.headResp.Header.Get("Accept-Ranges") == "none" {
		log.Fatal("Server does not support partial download")
		return false
	}

	return true
}
