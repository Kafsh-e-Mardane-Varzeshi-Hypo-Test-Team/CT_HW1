package models

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Queue struct {
	Name         string
	downloadChan chan *Download
	done         chan struct{}
	wg           sync.WaitGroup

	mu            sync.Mutex
	SavePath      string
	NumConcurrent int
	NumRetries    int
	StartTime     time.Time
	EndTime       time.Time
	MaxBandwidth  int
	active        bool
}

func NewQueue(name, savePath string, numConcurrent, numRetries int, startTime, endTime time.Time, maxBandwidth int) *Queue {
	return &Queue{
		Name:          name,
		SavePath:      savePath,
		NumConcurrent: numConcurrent,
		NumRetries:    numRetries,
		StartTime:     startTime,
		EndTime:       endTime,
		MaxBandwidth:  maxBandwidth,
		active:        false,
	}
}

func (q *Queue) UpdateConfig(savePath string, numConcurrent, numRetries int, startTime, endTime time.Time, maxBandwidth int) {
	q.SavePath = savePath
	q.NumConcurrent = numConcurrent
	q.NumRetries = numRetries
	q.StartTime = startTime
	q.EndTime = endTime
	q.MaxBandwidth = maxBandwidth
}

func (q *Queue) AddDownload(d *Download) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.active {
		log.Printf("download %T did NOT add to queue %T downloadChan\n", d, q)
		return nil // TODO: error handling
	}

	select {
	case q.downloadChan <- d:
		log.Printf("download %q added to queue %q\n", d.URL, q.Name)
	default:
		log.Printf("failed to add downlaod %T to queue %T, too many downloads has beed added", d, q)
		return errors.New("failed to add to queue")
	}
	return nil
}
