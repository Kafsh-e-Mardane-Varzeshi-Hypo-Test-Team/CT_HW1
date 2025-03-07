// PERSONAL NOTE: Don't forget to update download status
package internal

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const numberOfParts int = 2

type Status int

const (
	Pending Status = iota
	InProgress
	Paused
	Cancelled
	Failed
	Completed
)

type Download struct {
	URL            string
	Destination    string
	OutputFileName string
	Queue          *Queue
	Status
	headResp               *http.Response
	contentLength          int
	indexOfDownloadedBytes [numberOfParts]int64
	// TODO: Add array of size 'numberOfParts' for storing number of downloaded bytes from this part
	// TODO: Calculate download percentage using this array
}

func (d *Download) setHttpResponse() error {
	req, err := http.NewRequest("HEAD", d.URL, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to get response from server")
		return errors.New("Response status code is not OK!")
	}

	d.headResp = resp
	return nil
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

func (d *Download) downloadThisPart(index, startIndex, endIndex int) error {
	req, err := http.NewRequest("GET", d.URL, nil)
	if err != nil {
		log.Fatal(err)
		return err
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
		return err
	}
	defer resp.Body.Close()

	// Write the response to file
	file, err := os.Create(d.Destination + "/" + d.OutputFileName + rangeOfDownload + ".part")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	written, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}
	d.indexOfDownloadedBytes[index] = written

	log.Println("Downloaded ", rangeOfDownload)
	return nil
}

func (d *Download) downloadInParts() error {
	file, err := os.Create(d.Destination + "/" + d.OutputFileName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	partSize := d.contentLength / numberOfParts
	for i := 0; i < numberOfParts; i++ {
		startIndex := i * partSize
		endIndex := (i + 1) * partSize
		if i == numberOfParts-1 {
			endIndex = d.contentLength - 1
		}
		err := d.downloadThisPart(i, startIndex, endIndex)
		if err != nil {
			d.Status = Failed
			return err
		}

		resp, err := os.Open(d.Destination + "/" + d.OutputFileName + strconv.Itoa(startIndex) + "-" + strconv.Itoa(endIndex) + ".part")
		if err != nil {
			log.Fatal(err)
			return err
		}

		_, err = io.Copy(file, resp)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	d.Status = Completed
	return nil
}

func (d *Download) downloadInOneGo() error {
	err := d.downloadThisPart(0, 0, d.contentLength-1)
	if err != nil {
		d.Status = Failed
		return err
	}

	d.Status = Completed
	return nil
}

func (d *Download) Start() error {
	err := d.setHttpResponse()
	if err != nil {
		return err
	}

	d.setContentLength()

	if d.contentLength == 0 {
		d.Status = Cancelled
		log.Fatal("Content length is invalid")
		return errors.New("Content length is invalid")
	} else {
		d.Status = InProgress
		log.Println("Content length is", d.contentLength)
		if d.supportsPartialDownload() {
			err = d.downloadInParts()
		} else {
			err = d.downloadInOneGo()
		}
	}
	return err
}

func (d *Download) Stop() {

}
