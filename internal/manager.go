package internal

import "sync"

type Manager struct {
	mu        sync.Mutex
	downloads []*Download
	queues    map[string]*Queue
}

func (m *Manager) Start(send, recieve chan interface{}) {

}

func (m *Manager) Stop() {

}

func (m *Manager) addDownload(d *Download) error {
	return nil
}

func (m *Manager) removeDownload(d *Download) error {
	return nil
}
