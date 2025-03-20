package internal

import (
	"errors"
	"log"
	"maps"
	"slices"
	"sync"
	"time"
)

type Manager struct {
	mu        sync.Mutex
	LastID    int
	Downloads []*Download
	Queues    map[string]*Queue
}

func NewManager() *Manager {
	return &Manager{
		Queues: make(map[string]*Queue),
	}
}

func (m *Manager) Start(send, recieve chan interface{}) {

}

func (m *Manager) Stop() {

}

func (m *Manager) AddDownload(url, outputFileName, queueName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	q, exists := m.Queues[queueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	d := NewDownload(m.LastID, url, q.GetSavePath(), outputFileName, queueName)
	m.LastID++

	d.Pend()
	if q.IsActive() {
		err := q.AddDownload(d)
		if err != nil {
			return err
		}
	}

	m.Downloads = append(m.Downloads, d)
	log.Printf("added download %q to queue %q\n", d.URL, d.GetQueueName())
	return nil
}

func (m *Manager) RemoveDownload(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var d *Download
	for i, dl := range m.Downloads {
		if dl.ID == id {
			d = dl
			m.Downloads = slices.Delete(m.Downloads, i, i+1)
			break
		}
	}

	d.Cancel() // TODO: error handling, cancell

	log.Printf("removed download %q from queue %q\n", d.URL, d.GetQueueName())
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

	d.Pause() // TODO: error handling, paused

	log.Printf("paused download %q\n", d.URL)
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

	q, exists := m.Queues[d.queueName]
	if !exists {
		return errors.New("queue not found")
	}

	d.Pend()
	if q.IsActive() {
		err := q.AddDownload(d)
		if err != nil {
			return err
		}
	}

	log.Printf("resume download %q in queue %q\n", d.URL, d.GetQueueName())
	return nil
}

func (m *Manager) GetDownloadList() []*DownloadInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	var list []*DownloadInfo

	for _, d := range m.Downloads {
		list = append(list, &DownloadInfo{d.ID, d.URL, d.GetQueueName(), d.GetTransferRate(), d.GetProgress(), d.GetStatus()})
	}

	return list
}

func (m *Manager) addQueue(q *Queue) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.Queues[q.Name]; exists {
		return errors.New("queue already exists")
	}

	m.Queues[q.Name] = q
	log.Printf("Manager: Added queue %q\n", q.Name)
	return nil
}

func (m *Manager) removeQueue(queueName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	q, exists := m.Queues[queueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	q.Stop() // TODO: error handling

	delete(m.Queues, queueName)
	log.Printf("Manager: Removed queue %q\n", queueName)
	return nil
}

func (m *Manager) getQueuePendingDownloads(queueName string) []*Download {
	m.mu.Lock()
	defer m.mu.Unlock()

	var queuedDownloads []*Download
	for _, d := range m.Downloads {
		if d.GetQueueName() == queueName && d.GetStatus() == Pending {
			queuedDownloads = append(queuedDownloads, d)
		}
	}
	return queuedDownloads
}

func (m *Manager) monitorActiveHours() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			m.mu.Lock()
			for q := range maps.Values(m.Queues) {
				isActive := q.IsActive()
				checkTime := q.CheckActiveTime(now)
				if isActive && !checkTime {
					q.Stop()
				}
				if !isActive && checkTime {
					queuedDownloads := m.getQueuePendingDownloads(q.Name)
					q.Start(queuedDownloads)
				}
			}
			m.mu.Unlock()
		}
	}
}

type DownloadInfo struct {
	ID           int
	URL          string
	QueueName    string
	TransferRate int64
	Progress     float32
	Status
}

type QueueInfo struct {
	Name            string
	TargetDirectory string
	MaxParallel     int
	SpeedLimit      int
	StartTime       time.Time
	EndTime         time.Time
}
