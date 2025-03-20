package internal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
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
	ID                 int
	URL                string
	Destination        string
	OutputFileName     string
	path               string
	queueName          string
	headResp           *http.Response
	numberOfParts      int
	totalSize          int64
	lastDownloadedSize int64
	downloadedSize     int64
	downloadPercentage float64
	currentSpeed       float64
	lastUpdateTime     time.Time
	channel            chan error
	parts              []*Part
	mu                 sync.Mutex
	Status
}

func NewDownload(id int, url, destination, outputFileName, queueName string) *Download {
	return &Download{
		ID:                 id,
		URL:                url,
		Destination:        destination,
		OutputFileName:     outputFileName,
		path:               destination + "/" + outputFileName,
		queueName:          queueName,
		Status:             Pending,
		lastDownloadedSize: 0,
		downloadedSize:     0,
		downloadPercentage: 0,
	}
}

func (d *Download) setHttpResponse() error {
	req, err := http.NewRequest("HEAD", d.URL, nil)
	if err != nil {
		log.Printf("Error in getting HEAD of http request for downloadID = %d %v\n", d.ID, err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error performing http request for downloadID = %d: %v\n", d.ID, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error getting response from server for downloadID = %d: %v\n", d.ID, err)
		return errors.New("response status code is not OK")
	}

	d.headResp = resp
	return nil
}

func (d *Download) setTotalSize() {
	d.totalSize = d.headResp.ContentLength
}

func (d *Download) supportsPartialDownload() bool {
	if d.headResp.Header.Get("Accept-Ranges") == "" || d.headResp.Header.Get("Accept-Ranges") == "none" {
		log.Printf("downloadID = %d does not support partial downloading\n", d.ID)
		return false
	}

	return true
}

func (d *Download) downloadParts(bandwidthLimiter *BandwidthLimiter) error {
	partSize := d.totalSize / int64(d.numberOfParts)
	for i := range d.numberOfParts {
		req, err := http.NewRequest("GET", d.URL, nil)
		if err != nil {
			log.Printf("Error in GET http request for downloadID = %d: %v\n", d.ID, err)
			return err
		}

		p := Part{
			partIndex:       i,
			startIndex:      int64(i) * partSize,
			endIndex:        int64(i+1) * partSize,
			downloadedBytes: 0,
			Status:          Pending,
			req:             req,
		}
		if i == d.numberOfParts-1 {
			p.endIndex = d.totalSize - 1
		}
		p.rangeOfDownload = strconv.Itoa(int(p.startIndex+p.downloadedBytes)) + "-" + strconv.Itoa(int(p.endIndex))
		if p.path == "" {
			p.path = d.Destination + "/" + d.OutputFileName + p.rangeOfDownload + ".part"
		}

		d.parts[i] = &p
		fmt.Println("starting part", i)
		go d.parts[i].start(d.channel, bandwidthLimiter)
	}

	for range d.numberOfParts {
		err := <-d.channel
		if err != nil {
			d.Status = Failed
			return err
		}
	}
	return nil
}

func (d *Download) mergeParts() error {
	file, err := os.Create(d.path)
	if err != nil {
		log.Printf("Error creating merged file for downloadID = %d: %v\n", d.ID, err)
		return err
	}
	defer file.Close()

	for _, part := range d.parts {
		resp, err := os.Open(part.path)
		if err != nil {
			log.Printf("Error opening part file while merging for partId = %d: %v\n", part.partIndex, err)
			return err
		}

		_, err = io.Copy(file, resp)
		if err != nil {
			log.Printf("Error copying content from partId = %d to merged file in downloadID = %d: %v\n", part.partIndex, d.ID, err)
			return err
		}
		err = os.Remove(part.path)
		if err != nil {
			log.Printf("Error deleting .part file of partId = %d after merging in downloadID = %d: %v\n", part.partIndex, d.ID, err)
			return err
		}
	}
	return nil
}

func (d *Download) Start(bandwidthLimiter *BandwidthLimiter) error {
	d.path = d.Destination + "/" + d.OutputFileName

	d.lastUpdateTime = time.Now()
	go d.monitorProgress()

	err := d.setHttpResponse()
	if err != nil {
		return err
	}
	d.setTotalSize()

	if d.totalSize == 0 {
		d.setStatus(Failed)
		log.Printf("Content length in downloadID = %d is invalid\n", d.ID)
		return errors.New("content length is invalid")
	}

	d.setStatus(InProgress)
	log.Printf("Content length in downloadID = %d is %d\n", d.ID, d.totalSize)
	if d.supportsPartialDownload() {
		d.numberOfParts = NUMBER_OF_PARTS
	} else {
		d.numberOfParts = 1
	}
	d.parts = make([]*Part, d.numberOfParts)
	d.channel = make(chan error, d.numberOfParts)

	err = d.downloadParts(bandwidthLimiter)
	if err != nil {
		log.Printf("Error in downloadParts() function for downloadID = %d : %v\n", d.ID, err)
		return err
	}
	err = d.mergeParts()
	if err != nil {
		log.Printf("Error in mergeParts() function for downloadID = %d : %v\n", d.ID, err)
		return err
	}

	d.setStatus(Completed)
	return nil
}

func (d *Download) Pause() error {
	d.setStatus(Paused)
	for _, part := range d.parts {
		err := part.pause()
		if err != nil {
			log.Printf("Error while pausing download of partId %v", part.partIndex)
			return err
		}
	}
	return nil
}

func (d *Download) Pend() error {
	d.setStatus(Pending)
	for _, part := range d.parts {
		err := part.pend()
		if err != nil {
			log.Printf("Error while pending download of partId %v", part.partIndex)
			return err
		}
	}
	return nil
}

func (d *Download) Cancel() error {
	d.setStatus(Cancelled)
	for _, part := range d.parts {
		err := os.Remove(part.path)
		if err != nil {
			log.Printf("Error deleting .part file of partId = %d after canceling downloadID = %d: %v\n", part.partIndex, d.ID, err)
			return err
		}
	}
	return nil
}

func (d *Download) setStatus(status Status) {
	d.mu.Lock()
	d.Status = status
	d.mu.Unlock()
}

func (d *Download) monitorProgress() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	isActive := true
	for range ticker.C {
		<-ticker.C
		if d.Status != InProgress {
			if isActive {
				isActive = false
			} else {
				return
			}
		}

		d.mu.Lock()
		d.downloadedSize = 0
		for _, part := range d.parts {
			d.downloadedSize += part.downloadedBytes
		}

		now := time.Now()
		elapsed := now.Sub(d.lastUpdateTime).Seconds()
		bytesDownloaded := d.downloadedSize - d.lastDownloadedSize

		if d.Status != Paused && elapsed > 0 {
			d.currentSpeed = float64(bytesDownloaded) / elapsed
			d.lastDownloadedSize = d.downloadedSize
			d.lastUpdateTime = now
		} else if d.Status == Paused {
			d.currentSpeed = 0
		}

		percentage := float64(d.downloadedSize) / float64(d.totalSize) * 100
		d.downloadPercentage = percentage
		d.mu.Unlock()
	}
}

func (d *Download) GetQueueName() string {
	return d.queueName
}

func (d *Download) GetStatus() Status {
	return d.Status
}

func (d *Download) GetTransferRate() float64 {
	return d.currentSpeed
}

func (d *Download) GetProgress() float64 {
	return d.downloadPercentage
}
