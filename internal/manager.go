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

func (m *Manager) addDownload(d *Download) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	q, exists := m.Queues[d.GetQueueName()]
	if !exists {
		return errors.New("queue does not exist")
	}

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

func (m *Manager) removeDownload(d *Download) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	d.Stop() // TODO: error handling

	for i, download := range m.Downloads {
		if download == d {
			m.Downloads = slices.Delete(m.Downloads, i, i+1)
			break
		}
	}

	log.Printf("removed download %q from queue %q\n", d.URL, d.GetQueueName())
	return nil
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
