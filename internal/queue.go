package internal

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Queue struct {
	Name         string
	downloadChan chan *Download
	done         chan struct{}
	wg           sync.WaitGroup

	mu            sync.Mutex
	SavePath      string
	NumConcurrent int
	NumRetries    int
	StartTime     time.Time
	EndTime       time.Time
	MaxBandwidth  int
	isActive      bool
}

func NewQueue(name, savePath string, numConcurrent, numRetries int, startTime, endTime time.Time, maxBandwidth int) *Queue {
	return &Queue{
		Name:          name,
		SavePath:      savePath,
		NumConcurrent: numConcurrent,
		NumRetries:    numRetries,
		StartTime:     startTime,
		EndTime:       endTime,
		MaxBandwidth:  maxBandwidth,
		isActive:      false,
	}
}

func (q *Queue) UpdateConfig(savePath string, numConcurrent, numRetries int, startTime, endTime time.Time, maxBandwidth int) {
	q.SavePath = savePath
	q.NumConcurrent = numConcurrent
	q.NumRetries = numRetries
	q.StartTime = startTime
	q.EndTime = endTime
	q.MaxBandwidth = maxBandwidth
}

func (q *Queue) AddDownload(d *Download) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.isActive {
		log.Printf("download %T did NOT add to queue %T downloadChan\n", d, q)
		return nil // TODO: error handling
	}

	select {
	case q.downloadChan <- d:
		log.Printf("download %T added to queue %T\n", d, q)
	default:
		log.Printf("failed to add downlaod %T to queue %T, too many downloads has beed added", d, q)
		return errors.New("failed to add to queue")
	}
	return nil
}

func (q *Queue) Start(queuedDownloads []*Download) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isActive {
		return
	}
	q.isActive = true

	q.downloadChan = make(chan *Download, 100)
	q.done = make(chan struct{})

	bl := NewBandwidthLimiter(q.MaxBandwidth, q.done)

	q.wg.Add(q.NumConcurrent)
	for i := 0; i < q.NumConcurrent; i++ {
		go q.downloader(bl)
	}

	go q.addQueuedDownloads(queuedDownloads)
}

func (q *Queue) addQueuedDownloads(downloads []*Download) {
	for _, d := range downloads {
		select {
		case <-q.done:
			return
		default:
			q.AddDownload(d)
		}
	}
}

func (q *Queue) downloader(bl *BandwidthLimiter) {
	defer q.wg.Done()

	for {
		select {
		case d, ok := <-q.downloadChan:
			if !ok {
				return
			}
			if d.GetQueueName() == q.Name && d.GetStatus() == Pending {
				for i := 0; i < q.NumRetries; i++ {
					err := d.Start(bl)
					if err == nil {
						return
					}
					log.Println(err)

					if d.GetStatus() != Failed {
						return
					}
				}
			}
		case <-q.done:
			return
		}
	}
}

func (q *Queue) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.isActive = false

	close(q.downloadChan)
	close(q.done)
	q.wg.Wait()

	log.Printf("queue %T stopped", q)
}

func (q *Queue) GetSavePath() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.SavePath
}

func (q *Queue) GetNumConcurrent() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.SavePath
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
	return q.isActive
}

func (q *Queue) CheckActiveTime(now time.Time) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return now.After(q.StartTime) && now.Before(q.EndTime)
}

type BandwidthLimiter struct {
	rate   int
	tokens chan struct{}
	stop   chan struct{}
}

func NewBandwidthLimiter(rate int, stop chan struct{}) *BandwidthLimiter {
	bl := &BandwidthLimiter{
		rate:   rate,
		tokens: make(chan struct{}, rate/32), // why 32? like fps =)
		stop:   stop,
	}

	if rate > 0 {
		go bl.generateTokens()
	}

	return bl
}

func (bl *BandwidthLimiter) generateTokens() {
	ticker := time.NewTicker(time.Second / time.Duration(bl.rate/32)) // why 32? like fps =)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-bl.stop:
			return
		case bl.tokens <- struct{}{}:
		default:
		}
	}
}

func (bl *BandwidthLimiter) WaitForToken() {
	if bl.rate == 0 {
		return
	}
	<-bl.tokens
}
