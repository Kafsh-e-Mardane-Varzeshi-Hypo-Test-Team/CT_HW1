package internal

import (
	"io"
	"log"
	"net/http"
	"os"
)

type Part struct {
	partIndex       int
	startIndex      int
	endIndex        int
	downloadedBytes int64
	rangeOfDownload string
	path            string
	Status
}

func (p *Part) start(req *http.Request, channel chan error) {
	req.Header.Set("Range", "bytes="+p.rangeOfDownload)

	// TODO: Ask if it's better to save this client as a field in Download struct
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		channel <- err
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(p.path)
	if err != nil {
		log.Fatal(err)
		channel <- err
		return
	}
	defer file.Close()

	written, err := io.Copy(file, resp.Body)
	p.downloadedBytes += written
	if err != nil {
		p.Status = Failed
		log.Fatal(err)
		channel <- err
		return
	}

	channel <- nil
	p.Status = Completed
	log.Println("Downloaded ", p.rangeOfDownload)
}
