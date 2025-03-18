package internal

import (
	"errors"
	"log"
	"slices"
	"sync"
)

type Manager struct {
	mu        sync.Mutex
	downloads []*Download
	queues    map[string]*Queue
}

func NewManager() *Manager {
	return &Manager{
		queues: make(map[string]*Queue),
	}
}

func (m *Manager) Start(send, recieve chan interface{}) {

}

func (m *Manager) Stop() {

}

func (m *Manager) addDownload(d *Download) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	q, exists := m.queues[d.GetQueueName()]
	if !exists {
		return errors.New("queue does not exist")
	}

	err := q.AddDownload(d)
	if err != nil {
		return err
	}

	m.downloads = append(m.downloads, d)
	log.Printf("added download %q to queue %q\n", d.URL, d.GetQueueName())
	return nil
}

func (m *Manager) removeDownload(d *Download) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	d.Stop() // TODO: error handling

	for i, download := range m.downloads {
		if download == d {
			m.downloads = slices.Delete(m.downloads, i, i+1)
			break
		}
	}

	log.Printf("removed download %q from queue %q\n", d.URL, d.GetQueueName())
	return nil
}

func (m *Manager) addQueue(q *Queue) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.queues[q.name]; exists {
		return errors.New("queue already exists")
	}

	m.queues[q.name] = q
	log.Printf("Manager: Added queue %q\n", q.name)
	return nil
}

func (m *Manager) removeQueue(queueName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	q, exists := m.queues[queueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	q.Stop() // TODO: error handling

	delete(m.queues, queueName)
	log.Printf("Manager: Removed queue %q\n", queueName)
	return nil
}
