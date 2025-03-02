package internal

type Queue struct {
	Name            string
	downloads       []*Download
	savePath        string
	maxConcurrent   int
	activeDownloads int
	maxBandwidth    int
	maxRetries      int
	activeHours     string
}

func (queue *Queue) AddDownload(download *Download) {
	queue.downloads = append(queue.downloads, download)
	if queue.activeDownloads < queue.maxConcurrent {
		queue.activeDownloads++
		// Start download
	}
}
