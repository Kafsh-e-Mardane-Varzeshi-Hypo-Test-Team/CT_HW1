package internal

import (
	"errors"
	"sync"
)

type Queue struct {
	Name            string
	SavePath        string
	MaxConcurrent   int
	activeDownloads int
	MaxBandwidth    int
	MaxRetries      int
	ActiveHours     string
	downloadChan    chan *Download

	mu        sync.Mutex
	downloads []*Download
}

func NewQueue(name string, savePath string, maxConcurrent int, maxBandwidth int, maxRetries int, activeHours string) *Queue {
	return &Queue{
		Name:          name,
		SavePath:      savePath,
		MaxConcurrent: maxConcurrent,
		MaxBandwidth:  maxBandwidth,
		MaxRetries:    maxRetries,
		ActiveHours:   activeHours,
		downloadChan:  make(chan *Download, 100), // This may cause waiting for the channel to have capacity, check if there is a better way
	}
}

func (q *Queue) UpdateConfig(savePath string, maxConcurrent int, maxBandwidth int, maxRetries int, activeHours string) {
	q.SavePath = savePath
	q.MaxConcurrent = maxConcurrent
	q.MaxBandwidth = maxBandwidth
	q.MaxRetries = maxRetries
	q.ActiveHours = activeHours
}

func (q *Queue) AddDownload(d *Download) error {
	if d == nil {
		return errors.New("invalid download (nil)")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	q.downloads = append(q.downloads, d)
	q.downloadChan <- d

	return nil
}

func (q *Queue) CancelDownload(d *Download) error {
	if d == nil {
		return errors.New("invalid download (nil)")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	for i, dl := range q.downloads {
		if dl == d {
			if d.Status == InProgress || d.Status == Pending {
				d.Stop()
			}
			q.downloads = append(q.downloads[:i], q.downloads[i+1:]...)
			return nil
		}
	}

	return errors.New("download not found")
}

func (q *Queue) CancelAll() {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, d := range q.downloads {
		if d.Status == InProgress || d.Status == Pending {
			d.Stop()
		}
	}
}

func (q *Queue) Start() {
	for i := 0; i < q.MaxConcurrent; i++ {
		go func() {
			for d := range q.downloadChan {
				if d.Status == Pending {
					q.processDownload(d)
				}
			}
		}()
	}
}

func (q *Queue) processDownload(d *Download) {
	for i := 0; i < q.MaxRetries; i++ {
		err := d.Start()
		if err == nil {
			break
		}
	}
}
