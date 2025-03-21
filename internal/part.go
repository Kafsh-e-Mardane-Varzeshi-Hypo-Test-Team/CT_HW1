package internal

import (
	"errors"
	"fmt"
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
	Status
}

func (p *Part) start(channel chan error, bandwidthLimiter *BandwidthLimiter) {
	if p.Status == Completed {
		channel <- nil
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
		channel <- err
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(p.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error opening part file with partId = %d: %v\n", p.partIndex, err)
		p.setStatus(Failed)
		channel <- err
		return
	}
	defer file.Close()

	buffer := make([]byte, 32*1024)
	fmt.Println(runtime.NumGoroutine(), "before for of partIndex", p.partIndex)
	for {
		if p.getStatus() == Completed {
			return
		}
		if p.getStatus() != InProgress {
			log.Printf("%d partIndex = %d ~ p.status != InProgress", runtime.NumGoroutine(), p.partIndex)
			channel <- errors.New("part " + strconv.Itoa(p.partIndex) + " status is not InProgress (it is " + strconv.Itoa(int(p.Status)) + ")")
			return
		}

		bandwidthLimiter.WaitForToken()
		n, err := resp.Body.Read(buffer)
		log.Printf("downloading partId = %d with n = %d and downloadedBytes = %d/%d", p.partIndex, n, p.downloadedBytes, p.endIndex - p.startIndex)
		n = min(n, int(p.endIndex - p.startIndex - p.downloadedBytes))
		if n > 0 {
			_, err := file.Write(buffer[:n])
			if err != nil {
				log.Printf("Error writing buffer to part file for partId = %d: %v\n", p.partIndex, err)
				p.setStatus(Failed)
				channel <- err
				return
			}

			p.addToDownloadedBytes(n)
		}
		if err == io.EOF {
			log.Printf("%d Downloaded partIndex = %d (bytes %d - %d)", runtime.NumGoroutine(), p.partIndex, p.startIndex, p.endIndex)
			p.setStatus(Completed)
			channel <- nil
			return
		}
		if err != nil {
			log.Printf("Error reading body of http request for partId = %d: %v\n", p.partIndex, err)
			p.setStatus(Failed)
			channel <- err
			return
		}
	}
}

func (p *Part) pause() error {
	log.Printf("Pausing download of part %d : %d bytes downloaded", p.partIndex, p.downloadedBytes)
	p.setStatus(Paused)
	return nil
}

func (p *Part) pend() error {
	log.Printf("Pending download of part %d : %d bytes downloaded", p.partIndex, p.downloadedBytes)
	p.setStatus(Pending)
	return nil
}

func (p *Part) cancel() error {
	log.Printf("Canceling download of part %d : %d bytes downloaded", p.partIndex, p.downloadedBytes)
	p.setStatus(Cancelled)
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
