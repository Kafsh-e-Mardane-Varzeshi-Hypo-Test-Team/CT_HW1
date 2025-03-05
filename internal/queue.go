package internal

import (
	"errors"
	"sync"
)

type Queue struct {
	Name            string
	downloads       []*Download
	SavePath        string
	MaxConcurrent   int
	activeDownloads int
	MaxBandwidth    int
	MaxRetries      int
	ActiveHours     string
	// ch              chan int
	mu sync.Mutex
}

func NewQueue(name string, savePath string, maxConcurrent int, maxBandwidth int, maxRetries int, activeHours string) *Queue {
	return &Queue{
		Name:          name,
		SavePath:      savePath,
		MaxConcurrent: maxConcurrent,
		MaxBandwidth:  maxBandwidth,
		MaxRetries:    maxRetries,
		ActiveHours:   activeHours,
	}
}

func (q *Queue) AddDownload(d *Download) error {
	if d == nil {
		return errors.New("invalid download (nil)")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	q.downloads = append(q.downloads, d)

	return nil
}

func (q *Queue) RemoveDownload(d *Download) error {
	if d == nil {
		return errors.New("invalid download (nil)")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	for i, dl := range q.downloads {
		if dl == d {
			q.downloads = append(q.downloads[:i], q.downloads[i+1:]...)
			return nil
		}
	}

	return errors.New("download not found")
}
