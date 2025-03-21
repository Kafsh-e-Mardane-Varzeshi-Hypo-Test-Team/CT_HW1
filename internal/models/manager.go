package models

import (
	"errors"
	"log"
	"maps"
	"slices"
	"sync"
	"time"
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

type QueueInfo struct {
	Name            string
	TargetDirectory string
	MaxParallel     int
	SpeedLimit      int
	NumRetries      int
	StartTime       time.Time
	EndTime         time.Time
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
		URL:          URL,
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
			URL:          download.URL,
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

	_, exists := m.Queues[queueName]
	if !exists {
		return errors.New("queue does not exist")
	}

	delete(m.Queues, queueName)
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

	return list
}
func (q *Queue) GetSavePath() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.SavePath
}

func (q *Queue) GetNumConcurrent() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.NumConcurrent
}

func (q *Queue) GetMaxBandwidth() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.MaxBandwidth
}

func (q *Queue) GetStartTime() time.Time {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.StartTime
}

func (q *Queue) GetEndTime() time.Time {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.EndTime
}

func (q *Queue) IsActive() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.active
}
