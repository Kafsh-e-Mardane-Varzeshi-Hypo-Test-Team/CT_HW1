package models

import "sync"

type Manager struct {
	mu        sync.Mutex
	Downloads []*Download
	Queues    map[string]*Queue
}

func NewManager() *Manager {
	return &Manager{
		Queues: make(map[string]*Queue),
	}
}
