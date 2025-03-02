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
