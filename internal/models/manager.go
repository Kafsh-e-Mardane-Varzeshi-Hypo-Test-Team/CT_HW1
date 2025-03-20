package models

import (
	"slices"
	"sync"
)

var idCounter int = 0

type Manager struct {
	mu        sync.Mutex
	Downloads []*Download
	Queues    map[string]*Queue
}

type Status int

const (
	InProgress Status = iota
	Paused
	Completed
	Failed
)

type DownloadInfo struct {
	ID           int
	URL          string
	QueueName    string
	TransferRate int64
	Progress     float32
	Status
}

func NewManager() *Manager {
	return &Manager{
		Queues: make(map[string]*Queue),
	}
}

func (manager *Manager) AddDownload(
	URL string,
	OutputFileName string,
	queueName string,
) error {
	idCounter++
	manager.Downloads = append(manager.Downloads, &Download{
		ID:           idCounter + 1,
		Url:          URL,
		Queue:        queueName,
		TransferRate: 0,
		Progress:     0,
		Status:       InProgress,
	})
	return nil
}

func (m *Manager) GetDownloadList() []*DownloadInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	downloads := []*DownloadInfo{}
	for _, download := range m.Downloads {
		downloads = append(downloads, &DownloadInfo{
			ID:           download.ID,
			URL:          download.Url,
			QueueName:    download.Queue,
			TransferRate: download.TransferRate,
			Progress:     download.Progress,
			Status:       download.Status,
		})
	}
	return downloads
}

func (m *Manager) RemoveDownload(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, dl := range m.Downloads {
		if dl.ID == id {
			m.Downloads = slices.Delete(m.Downloads, i, i+1)
			break
		}
	}

	return nil
}

func (m *Manager) PauseDownload(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var d *Download
	for _, dl := range m.Downloads {
		if dl.ID == id {
			d = dl
			break
		}
	}

	d.Status = Paused

	return nil
}

func (m *Manager) ResumeDownload(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var d *Download
	for _, dl := range m.Downloads {
		if dl.ID == id {
			d = dl
			break
		}
	}

	d.Status = InProgress

	return nil
}
