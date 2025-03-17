package internal

import (
	"errors"
	"log"
	"sync"
)

type Queue struct {
	name          string
	savePath      string
	numConcurrent int
	numRetries    int
	activeHours   string
	maxBandwidth  int
	downloadChan  chan *Download
	done          chan struct{}

	mu        sync.Mutex
	isActive  bool
	downloads []*Download
}

func NewQueue(name, savePath string, numConcurrent, numRetries int, activeHours string, maxBandwidth int) *Queue {
	return &Queue{
		name:          name,
		savePath:      savePath,
		numConcurrent: numConcurrent,
		numRetries:    numRetries,
		activeHours:   activeHours,
		maxBandwidth:  maxBandwidth,
		isActive:      false,
	}
}

func (q *Queue) AddDownload(d *Download) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.isActive {
		q.downloads = append(q.downloads, d)
		log.Printf("download %T added to inactive queue %T\n", d, q)
		return nil
	}

	select {
	case q.downloadChan <- d:
		q.downloads = append(q.downloads, d)
		log.Printf("download %T added to queue %T\n", d, q)
	default:
		log.Printf("failed to add downlaod %T to queue %T, too many downloads has beed added", d, q)
		return errors.New("failed to add to queue")
	}
	return nil
}

func (q *Queue) Start() {
	q.downloadChan = make(chan *Download, 100)
	q.done = make(chan struct{})

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isActive {
		return
	}
	q.isActive = true

	for i := 0; i < q.numConcurrent; i++ {
		go q.downloader()
	}

	go q.addBufferedDownloads()
}

func (q *Queue) downloader() {
	for {
		select {
		case d := <-q.downloadChan:
			if d.GetQueueName() == q.name && d.GetStatus() == Pending {
				for i := 0; i < q.numRetries; i++ {
					err := d.Start()
					if err == nil {
						return
					}
					log.Println(err)

					if d.GetStatus() != Failed {
						return
					}
				}
			}
		case <-q.done:
			return
		}
	}
}

func (q *Queue) addBufferedDownloads() {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, d := range q.downloads {
		if d.GetStatus() == Pending {
			select {
			case q.downloadChan <- d:
				log.Printf("download %T added to queue %T\n", d, q)
			case <-q.done:
				return
			}
		}
	}
}

func (q *Queue) Stop() {
	close(q.done)
	close(q.downloadChan)
}
