package internal

import (
	"errors"
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
	parts              []Part
	isInitialized      bool
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
	for i := range d.parts {
		go d.parts[i].start(d.channel, bandwidthLimiter)
	}

	for range d.numberOfParts {
		err := <-d.channel
		if err != nil {
			d.setStatus(Failed)
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

	for i := range d.parts {
		part := &d.parts[i]
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

func (d *Download) initializeParts() error {
	partSize := d.totalSize / int64(d.numberOfParts)
	for i := range d.numberOfParts {
		req, err := http.NewRequest("GET", d.URL, nil)
		if err != nil {
			log.Printf("Error in GET http request for downloadID = %d: %v\n", d.ID, err)
			return err
		}

		d.parts[i] = Part{
			partIndex:       i,
			startIndex:      int64(i) * partSize,
			endIndex:        int64(i+1) * partSize,
			downloadedBytes: 0,
			Status:          Pending,
			req:             req,
			channel:         make(chan Status),
		}

		if i == d.numberOfParts-1 {
			d.parts[i].endIndex = d.totalSize - 1
		}
		d.parts[i].rangeOfDownload = strconv.Itoa(int(d.parts[i].startIndex+d.parts[i].downloadedBytes)) + "-" + strconv.Itoa(int(d.parts[i].endIndex))
		d.parts[i].path = d.Destination + "/" + d.OutputFileName + d.parts[i].rangeOfDownload + ".part"
	}

	return nil
}

func (d *Download) initializeDownload() error {
	d.isInitialized = true
	d.path = d.Destination + "/" + d.OutputFileName

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

	if d.supportsPartialDownload() {
		d.numberOfParts = NUMBER_OF_PARTS
	} else {
		d.numberOfParts = 1
	}
	d.parts = make([]Part, d.numberOfParts)

	err = d.initializeParts()
	if err != nil {
		log.Printf("Error in initializing parts for downloadID = %d", d.ID)
		return err
	}
	return nil
}

func (d *Download) Start(bandwidthLimiter *BandwidthLimiter) error {
	if !d.isInitialized {
		err := d.initializeDownload()
		if err != nil {
			log.Printf("Error while initializing downloadID = %d:%v", d.ID, err)
			return err
		}
	}
	d.channel = make(chan error, d.numberOfParts)
	d.setStatus(InProgress)
	log.Printf("Content length in downloadID = %d is %d\n", d.ID, d.totalSize)

	d.lastUpdateTime = time.Now()
	go d.monitorProgress()

	err := d.downloadParts(bandwidthLimiter)
	if err != nil {
		log.Printf("Error in downloadParts() function for downloadID = %d : %v\n", d.ID, err)
		return err
	}
	log.Printf("All parts downloaded successfully")
	err = d.mergeParts()
	if err != nil {
		log.Printf("Error in mergeParts() function for downloadID = %d : %v\n", d.ID, err)
		return err
	}

	d.setStatus(Completed)
	return nil
}

func (d *Download) Pause() error {
	log.Printf("Pausing downloadID = %d", d.ID)
	d.setStatus(Paused)
	for i := range d.parts {
		part := &d.parts[i]
		err := part.pause()
		if err != nil {
			log.Printf("Error while pausing download of partId %v", part.partIndex)
			return err
		}
	}
	return nil
}

func (d *Download) Pend() error {
	log.Printf("Pending downloadID = %d", d.ID)
	d.setStatus(Pending)
	for i := range d.parts {
		part := &d.parts[i]
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
	for i := range d.parts {
		part := &d.parts[i]
		err := part.cancel()
		if err != nil {
			log.Printf("Error canceling partId = %d while canceling downloadID = %d: %v\n", part.partIndex, d.ID, err)
			return err
		}

		if _, err := os.Stat(part.path); errors.Is(err, os.ErrNotExist) {
			log.Printf(".part file of partId = %d does not exists in downloadID = %d\n", part.partIndex, d.ID)
			continue
		}

		err = os.Remove(part.path)
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
	ticker := time.NewTicker(300 * time.Millisecond)
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
		for i := range d.parts {
			part := &d.parts[i]
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

		//log.Printf("monitoring :: %.2f%% (%.2f MB/%.2f MB) - %.2f MB/s\n",
		//	percentage, float64(d.downloadedSize)/1024/1024, float64(d.totalSize)/1024/1024, d.currentSpeed/1024/1024)
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
