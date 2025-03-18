package internal

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Queue struct {
	name         string
	downloadChan chan *Download
	done         chan struct{}
	wg           sync.WaitGroup

	mu            sync.Mutex
	savePath      string
	numConcurrent int
	numRetries    int
	startTime     time.Time
	endTime       time.Time
	maxBandwidth  int
	isActive      bool
}

func NewQueue(name, savePath string, numConcurrent, numRetries int, startTime, endTime time.Time, maxBandwidth int) *Queue {
	return &Queue{
		name:          name,
		savePath:      savePath,
		numConcurrent: numConcurrent,
		numRetries:    numRetries,
		startTime:     startTime,
		endTime:       endTime,
		maxBandwidth:  maxBandwidth,
		downloadChan:  make(chan *Download, 100),
		isActive:      false,
	}
}

func (q *Queue) UpdateConfig(savePath string, numConcurrent, numRetries int, startTime, endTime time.Time, maxBandwidth int) {
	q.savePath = savePath
	q.numConcurrent = numConcurrent
	q.numRetries = numRetries
	q.startTime = startTime
	q.endTime = endTime
	q.maxBandwidth = maxBandwidth
}

func (q *Queue) AddDownload(d *Download) error {
	select {
	case q.downloadChan <- d:
		log.Printf("download %T added to queue %T\n", d, q)
	default:
		log.Printf("failed to add downlaod %T to queue %T, too many downloads has beed added", d, q)
		return errors.New("failed to add to queue")
	}
	return nil
}

func (q *Queue) Start() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isActive {
		return
	}

	q.isActive = true
	q.done = make(chan struct{})

	q.wg.Add(q.numConcurrent)
	for i := 0; i < q.numConcurrent; i++ {
		go q.downloader()
	}
}

func (q *Queue) downloader() {
	defer q.wg.Done()

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

func (q *Queue) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.isActive = false

	close(q.done)
	q.wg.Wait()

	log.Printf("queue %T stopped", q)
}

func (q *Queue) IsActive() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.isActive
}

func (q *Queue) CheckActiveTime(now time.Time) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return now.After(q.startTime) && now.Before(q.endTime)
}
