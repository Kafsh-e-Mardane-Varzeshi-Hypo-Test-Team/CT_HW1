package internal

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
)

type Part struct {
	partIndex       int
	startIndex      int64
	endIndex        int64
	downloadedBytes int64
	rangeOfDownload string
	path            string
	req             *http.Request
	mu              sync.Mutex
	channel         chan Status
	Status
}

func (p *Part) start(commonChannelOfParts chan error, bandwidthLimiter *BandwidthLimiter) {
	if p.getStatus() == Completed {
		commonChannelOfParts <- nil
		return
	}
	p.setStatus(InProgress)
	p.rangeOfDownload = strconv.Itoa(int(p.startIndex+p.downloadedBytes)) + "-" + strconv.Itoa(int(p.endIndex))
	p.req.Header.Set("Range", "bytes="+p.rangeOfDownload)
	log.Printf("%d downloading part %d started (bytes %d - %d)", runtime.NumGoroutine(), p.partIndex, p.startIndex+p.downloadedBytes, p.endIndex)

	client := &http.Client{}
	resp, err := client.Do(p.req)
	if err != nil {
		log.Printf("Error performing http request for partId = %d: %v\n", p.partIndex, err)
		p.setStatus(Failed)
		commonChannelOfParts <- err
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(p.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error opening part file with partId = %d: %v\n", p.partIndex, err)
		p.setStatus(Failed)
		commonChannelOfParts <- err
		return
	}
	defer file.Close()

	buffer := make([]byte, 32*1024)
	for {
		select {
		case status := <-p.channel:
			log.Printf("Stop downloading partIndex = %d due to it's status = %d", p.partIndex, status)
			commonChannelOfParts <- errors.New("part " + strconv.Itoa(p.partIndex) + " has status = " + strconv.Itoa(int(status)))
			return
		default:
			bandwidthLimiter.WaitForToken()
			n, err := resp.Body.Read(buffer)
			// log.Printf("downloading partId = %d with n = %d and downloadedBytes = %d/%d", p.partIndex, n, p.downloadedBytes, p.endIndex - p.startIndex)
			if n > 0 {
				_, err := file.Write(buffer[:n])
				if err != nil {
					log.Printf("Error writing buffer to part file for partId = %d: %v\n", p.partIndex, err)
					p.fail()
					commonChannelOfParts <- err
					return
				}

				p.addToDownloadedBytes(n)
			}
			if err == io.EOF {
				log.Printf("%d Downloaded partIndex = %d (bytes %d - %d)", runtime.NumGoroutine(), p.partIndex, p.startIndex, p.endIndex)
				p.setStatus(Completed)
				commonChannelOfParts <- nil
				return
			}
			if err != nil {
				log.Printf("Error reading body of http request for partId = %d: %v\n", p.partIndex, err)
				p.fail()
				commonChannelOfParts <- err
				return
			}
		}
	}
}

func (p *Part) pause() error {
	p.channel <- Paused
	p.setStatus(Paused)
	log.Printf("Pausing download of part %d : %d bytes downloaded", p.partIndex, p.downloadedBytes)
	return nil
}

func (p *Part) pend() error {
	p.channel <- Pending
	p.setStatus(Pending)
	log.Printf("Pending download of part %d : %d bytes downloaded", p.partIndex, p.downloadedBytes)
	return nil
}

func (p *Part) cancel() error {
	p.channel <- Cancelled
	p.setStatus(Cancelled)
	log.Printf("Canceling download of part %d : %d bytes downloaded", p.partIndex, p.downloadedBytes)
	return nil
}

func (p *Part) fail() error {
	p.channel <- Failed
	p.setStatus(Failed)
	log.Printf("Failing download of part %d : %d bytes downloaded", p.partIndex, p.downloadedBytes)
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
	p.downloadedBytes += int64(n)
	p.mu.Unlock()
}
