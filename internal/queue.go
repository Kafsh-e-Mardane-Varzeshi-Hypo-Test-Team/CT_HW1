package internal

import "sync"

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
