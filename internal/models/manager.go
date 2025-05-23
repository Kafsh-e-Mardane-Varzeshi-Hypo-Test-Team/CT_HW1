package models

import (
	"encoding/json"
	"errors"
	"log"
	"maps"
	"slices"
	"sort"
	"strings"
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

func (m *Manager) Start() {
	go m.monitorActiveHours()
}

func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, q := range m.Queues {
		m.pauseQueueDownloads(q.Name)
		q.Stop()
	}
}

func (m *Manager) AddDownload(url, outputFileName, queueName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	q, exists := m.Queues[queueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	if outputFileName == "" {
		split := strings.Split(url, "/")
		outputFileName = split[len(split)-1]
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

	err := d.Cancel()
	if err != nil {
		return err
	}

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

	err := d.Pause()
	if err != nil {
		return err
	}

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

	q, exists := m.Queues[d.QueueName]
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

func (m *Manager) AddQueue(qInfo QueueInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.Queues[qInfo.Name]; exists {
		return errors.New("queue already exists")
	}

	if err := checkQueueInfo(qInfo); err != nil {
		return err
	}

	q := NewQueue(
		qInfo.Name,
		qInfo.TargetDirectory,
		qInfo.MaxParallel,
		qInfo.NumRetries,
		qInfo.StartTime,
		qInfo.EndTime,
		qInfo.SpeedLimit,
	)
	m.Queues[qInfo.Name] = q
	log.Printf("added queue %q\n", q.Name)
	return nil
}

func (m *Manager) RemoveQueue(queueName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	q, exists := m.Queues[queueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	for _, d := range m.Downloads {
		if d.GetQueueName() == queueName {
			d.Cancel()
		}
	}
	m.Downloads = slices.DeleteFunc(m.Downloads, func(d *Download) bool {
		return d.GetQueueName() == queueName
	})
	q.Stop() // TODO: error handling

	delete(m.Queues, queueName)
	log.Printf("removed queue %q\n", queueName)
	return nil
}

func (m *Manager) UpdateQueue(qInfo QueueInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	q, exists := m.Queues[qInfo.Name]
	if !exists {
		return errors.New("queue does not exist")
	}

	if err := checkQueueInfo(qInfo); err != nil {
		return err
	}

	q.UpdateConfig(
		qInfo.TargetDirectory,
		qInfo.MaxParallel,
		qInfo.NumRetries,
		qInfo.StartTime,
		qInfo.EndTime,
		qInfo.SpeedLimit,
	)
	return nil
}

func checkQueueInfo(qInfo QueueInfo) error {
	if qInfo.MaxParallel < 1 {
		return errors.New("parallel count error")
	}
	if qInfo.NumRetries < 0 {
		return errors.New("retry count error")
	}
	return nil
}

func (m *Manager) GetQueueList() []*QueueInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	var list []*QueueInfo

	for q := range maps.Values(m.Queues) {
		list = append(list, &QueueInfo{
			Name:            q.Name,
			TargetDirectory: q.GetSavePath(),
			MaxParallel:     q.GetNumConcurrent(),
			SpeedLimit:      q.MaxBandwidth,
			NumRetries:      q.NumRetries,
			StartTime:       q.StartTime,
			EndTime:         q.EndTime,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

func (m *Manager) getQueuePendingDownloads(queueName string) []*Download {
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

	m.checkTimeAndActivate()

	for {
		select {
		case <-ticker.C:
			m.checkTimeAndActivate()
			// case <-m.done:
		}
	}
}

func (m *Manager) checkTimeAndActivate() {
	now := time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()

	for q := range maps.Values(m.Queues) {
		isActive := q.IsActive()
		checkTime := q.CheckActiveTime(now)
		if isActive && !checkTime {
			m.pauseQueueDownloads(q.Name)
			q.Stop()
		}
		if !isActive && checkTime {
			queuedDownloads := m.getQueuePendingDownloads(q.Name)
			q.Start(queuedDownloads)
		}
	}
}

func (m *Manager) pauseQueueDownloads(qName string) {
	for _, d := range m.Downloads {
		if d.GetQueueName() == qName {
			switch d.GetStatus() {
			case Pending, InProgress:
				d.Pend()
			case Paused:
				d.Pause()
			case Cancelled:
				d.Cancel()
			}
		}
	}
}

func (m *Manager) GetJson() ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

type DownloadInfo struct {
	ID           int
	URL          string
	QueueName    string
	TransferRate float64
	Progress     float64
	Status
}

type QueueInfo struct {
	Name            string
	TargetDirectory string
	MaxParallel     int
	SpeedLimit      int64
	NumRetries      int
	StartTime       time.Time
	EndTime         time.Time
}
