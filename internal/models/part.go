package models

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Part struct {
	PartIndex       int
	StartIndex      int64
	EndIndex        int64
	DownloadedBytes int64
	RangeOfDownload string
	Path            string
	req             *http.Request
	mu              sync.Mutex
	channel         chan Status
	Status
}

func (p *Part) start(commonChannelOfParts chan connectionWithPart, bandwidthLimiter *BandwidthLimiter) {
	if p.getStatus() == Completed {
		commonChannelOfParts <- connectionWithPart{nil, Completed}
		return
	}
	p.setStatus(InProgress)

	startByte := p.StartIndex + p.DownloadedBytes
    p.RangeOfDownload = strconv.FormatInt(startByte, 10) + "-" + strconv.FormatInt(p.EndIndex, 10)
	p.req.Header.Set("Range", "bytes="+p.RangeOfDownload)
	log.Printf("downloading part %d started (bytes %d - %d)", p.PartIndex, p.StartIndex+p.DownloadedBytes, p.EndIndex)

	client := &http.Client{}
	resp, err := client.Do(p.req)
	if err != nil {
		log.Printf("Error performing http request for partId = %d: %v\n", p.PartIndex, err)
		p.setStatus(Failed)
		commonChannelOfParts <- connectionWithPart{err, Failed}
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(p.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error opening part file with partId = %d: %v\n", p.PartIndex, err)
		p.fail()
		commonChannelOfParts <- connectionWithPart{err, Failed}
		return
	}
	defer file.Close()

	buffer := make([]byte, 32*1024)
	for {
		select {
		case status := <-p.channel:
			log.Printf("Stop downloading partIndex = %d due to it's status = %d", p.PartIndex, status)
			commonChannelOfParts <- connectionWithPart{errors.New("part " + strconv.Itoa(p.PartIndex) + " has status = " + strconv.Itoa(int(status))), status}
			return
		default:
			bandwidthLimiter.WaitForToken()
			n, err := resp.Body.Read(buffer)
			// log.Printf("downloading partId = %d with n = %d and downloadedBytes = %d/%d", p.PartIndex, n, p.DownloadedBytes, p.EndIndex - p.StartIndex)
			if n > 0 {
				_, err := file.Write(buffer[:n])
				if err != nil {
					log.Printf("Error writing buffer to part file for partId = %d: %v\n", p.PartIndex, err)
					p.fail()
					commonChannelOfParts <- connectionWithPart{err, Failed}
					return
				}

				p.addToDownloadedBytes(n)
			}
			if err == io.EOF {
				log.Printf("Downloaded partIndex = %d (bytes %d - %d)", p.PartIndex, p.StartIndex, p.EndIndex)
				p.setStatus(Completed)
				commonChannelOfParts <- connectionWithPart{nil, Completed}
				return
			}
			if err != nil {
				log.Printf("Error reading body of http request for partId = %d: %v\n", p.PartIndex, err)
				p.fail()
				commonChannelOfParts <- connectionWithPart{err, Failed}
				return
			}
		}
	}
}

func (p *Part) pause() error {
	if p.getStatus() == InProgress {
		log.Printf("--------entering channel pause partId = %d", p.PartIndex)
		p.channel <- Paused
		log.Printf("--------read from channel pause partId = %d", p.PartIndex)
	}
	p.setStatus(Paused)
	log.Printf("Pausing download of part %d : %d bytes downloaded", p.PartIndex, p.DownloadedBytes)
	return nil
}

func (p *Part) pend() error {
	if p.getStatus() == InProgress {
		p.channel <- Pending
	}
	p.setStatus(Pending)
	log.Printf("Pending download of part %d : %d bytes downloaded", p.PartIndex, p.DownloadedBytes)
	return nil
}

func (p *Part) cancel() error {
	if p.getStatus() == InProgress {
		p.channel <- Cancelled
	}
	p.setStatus(Cancelled)
	log.Printf("Canceling download of part %d : %d bytes downloaded", p.PartIndex, p.DownloadedBytes)
	return nil
}

func (p *Part) fail() error {
	if p.getStatus() == InProgress {
		p.channel <- Failed
	}
	p.setStatus(Failed)
	log.Printf("Failing download of part %d : %d bytes downloaded", p.PartIndex, p.DownloadedBytes)
	return nil
}

func (p *Part) setStatus(status Status) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Status = status
}

func (p *Part) getStatus() Status {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.Status
}

func (p *Part) addToDownloadedBytes(n int) {
	p.mu.Lock()
	p.DownloadedBytes += int64(n)
	p.mu.Unlock()
}
