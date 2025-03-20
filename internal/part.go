package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Part struct {
	partIndex       int
	startIndex      int64
	endIndex        int64
	downloadedBytes int64
	rangeOfDownload string
	path            string
	req             *http.Request
	Status
}

func (p *Part) start(channel chan error) {
	if p.Status == Completed {
		return
	}
	fmt.Println("part", p.partIndex, "started")
	p.req.Header.Set("Range", "bytes="+p.rangeOfDownload)

	// TODO: Ask if it's better to save this client as a field in Download struct
	client := &http.Client{}
	resp, err := client.Do(p.req)
	if err != nil {
		log.Printf("Error performing http request for partId = %d: %v\n", p.partIndex, err)
		p.Status = Failed
		channel <- err
		return
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(p.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error opening part file with partId = %d: %v\n", p.partIndex, err)
		p.Status = Failed
		channel <- err
		return
	}
	defer file.Close()

	buffer := make([]byte, 32*1024)
	for {
		if p.Status == Paused {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if p.Status == Failed || p.Status == Cancelled {
			return
		}

		p.Status = InProgress
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, err := file.Write(buffer[:n])
			if err != nil {
				p.Status = Failed
				log.Printf("Error writing buffer to part file for partId = %d: %v\n", p.partIndex, err)
				log.Fatal(err)
				channel <- err
				return
			}

			p.downloadedBytes += int64(n)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			p.Status = Failed
			log.Printf("Error reading body of http request for partId = %d: %v\n", p.partIndex, err)
			channel <- err
			return
		}
	}

	channel <- nil
	p.Status = Completed
	log.Println("Downloaded ", p.rangeOfDownload)
}

func (p *Part) stop() {
	fmt.Println("downloaded bytes of part", p.partIndex, " :: ", p.downloadedBytes)
	p.Status = Paused
}
