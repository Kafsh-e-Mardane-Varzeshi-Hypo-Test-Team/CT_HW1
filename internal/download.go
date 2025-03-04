// PERSONAL NOTE: Don't forget to update download status
package internal

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	headResp      *http.Response
	contentLength int
	// TODO: Add array of size 'numberOfParts' for storing number of downloaded btyes from this part
	// TODO: Calculate download percentage using this array
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

func (d *Download) downloadThisPart(startIndex, endIndex int) bool {
	req, err := http.NewRequest("GET", d.URL, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}

	rangeOfDownload := strconv.Itoa(startIndex) + "-" + strconv.Itoa(endIndex)
	rangeHeader := "bytes=" + rangeOfDownload
	req.Header.Set("Range", rangeHeader)

	// TODO: Ask if it's better to save this client as a field in Download struct
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer resp.Body.Close()

	// Write the response to file
	file, err := os.Create(d.Destination + "/" + d.OutputFileName + rangeOfDownload + ".part")
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Println("Downloaded ", rangeOfDownload)
	return true
}

func (d *Download) downloadInParts() {
	numberOfParts := 2
	partSize := d.contentLength / numberOfParts

	for i := 0; i < numberOfParts; i++ {
		startIndex := i * partSize
		endIndex := (i + 1) * partSize
		if i == numberOfParts-1 {
			endIndex = d.contentLength - 1
		}
		ok := d.downloadThisPart(startIndex, endIndex)
		if !ok {
			d.Status = Cancelled
			return
		}
	}

	d.Status = Completed
}

func (d *Download) downloadInOneGo() {
	ok := d.downloadThisPart(0, d.contentLength-1)
	if !ok {
		d.Status = Cancelled
		return
	}
	d.Status = Completed
}

func (d *Download) Start() {
	d.setHttpResponse()
	d.setContentLength()

	if d.contentLength == 0 {
		d.Status = Cancelled
		log.Fatal("Content length is invalid")
	} else {
		d.Status = InProgress
		log.Println("Content length is", d.contentLength)
		if d.supportsPartialDownload() {
			d.downloadInParts()
		} else {
			d.downloadInOneGo()
		}
	}
}
