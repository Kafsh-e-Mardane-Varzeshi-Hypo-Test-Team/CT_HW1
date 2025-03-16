package internal

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const NUMBER_OF_PARTS int = 5

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
	headResp       *http.Response
	contentLength  int
	numberOfParts  int
	channel        chan error
	parts          []*Part
	Status
	// TODO: Add array of size 'numberOfParts' for storing number of downloaded bytes from this part
	// TODO: Calculate download percentage using this array
	// TODO: Don't forget to update download status
}

func (d *Download) setHttpResponse() error {
	req, err := http.NewRequest("HEAD", d.URL, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to get response from server")
		return errors.New("response status code is not OK")
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

func (d *Download) downloadParts() error {
	partSize := d.contentLength / d.numberOfParts
	for i := range d.numberOfParts {
		req, err := http.NewRequest("GET", d.URL, nil)
		if err != nil {
			log.Fatal(err)
			return err
		}

		p := Part{
			partIndex:       i,
			startIndex:      i * partSize,
			endIndex:        (i + 1) * partSize,
			downloadedBytes: 0,
			Status:          Pending,
			req:             req,
		}
		if i == d.numberOfParts-1 {
			p.endIndex = d.contentLength - 1
		}
		p.rangeOfDownload = strconv.Itoa(p.startIndex) + "-" + strconv.Itoa(p.endIndex)
		p.path = d.Destination + "/" + d.OutputFileName + p.rangeOfDownload + ".part"

		d.parts[i] = &p
		go d.parts[i].start(d.channel)
	}

	for range d.numberOfParts {
		select {
		case err := <-d.channel:
			if err != nil {
				d.Status = Failed
				return err
			}
		}
	}
	return nil
}

func (d *Download) mergeParts() error {
	file, err := os.Create(d.Destination + "/" + d.OutputFileName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	for _, part := range d.parts {
		resp, err := os.Open(part.path)
		if err != nil {
			log.Fatal(err)
			return err
		}

		_, err = io.Copy(file, resp)
		if err != nil {
			log.Fatal(err)
			return err
		}
		err = os.Remove(part.path)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func (d *Download) Start() error {
	err := d.setHttpResponse()
	if err != nil {
		return err
	}
	d.setContentLength()

	if d.contentLength == 0 {
		d.Status = Failed
		log.Fatal("Content length is invalid")
		return errors.New("content length is invalid")
	}

	d.Status = InProgress
	log.Println("Content length is", d.contentLength)
	if d.supportsPartialDownload() {
		d.numberOfParts = NUMBER_OF_PARTS
	} else {
		d.numberOfParts = 1
	}
	d.parts = make([]*Part, d.numberOfParts)
	d.channel = make(chan error, d.numberOfParts)

	err = d.downloadParts()
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = d.mergeParts()
	if err != nil {
		log.Fatal(err)
		return err
	}

	d.Status = Completed
	return nil
}

func (d *Download) Stop() {
	for _, part := range d.parts {
		part.stop()
	}
	d.Status = Paused
}
