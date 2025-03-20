package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	p.setStatus(InProgress)
	if p.Status == Completed {
		return
	}
	fmt.Println("part", p.partIndex, "started")
	p.req.Header.Set("Range", "bytes="+p.rangeOfDownload)

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
	for {
		if p.Status != InProgress {
			return
		}

		bandwidthLimiter.WaitForToken()
		n, err := resp.Body.Read(buffer)
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
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading body of http request for partId = %d: %v\n", p.partIndex, err)
			p.setStatus(Failed)
			channel <- err
			return
		}
	}

	p.setStatus(Completed)
	channel <- nil
	log.Println("Downloaded ", p.rangeOfDownload)
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
	p.Status = status
	p.mu.Unlock()
}

func (p *Part) addToDownloadedBytes(n int) {
	p.mu.Lock()
	p.downloadedBytes += int64(n)
	p.mu.Unlock()
}
