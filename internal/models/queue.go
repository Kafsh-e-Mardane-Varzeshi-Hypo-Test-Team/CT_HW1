package models

import (
	"time"
)

type Queue struct {
	Name                 string
	TargetDirectory      string
	MaxParallelDownloads int
	SpeedLimit           string
	StartTime            time.Time
	EndTime              time.Time
}

func GetQueues() []Queue {
	return []Queue{
		{
			Name:                 "queue1",
			TargetDirectory:      "/home/user/Downloads",
			MaxParallelDownloads: 2,
			SpeedLimit:           "1MB/s",
			StartTime:            time.Now(),
			EndTime:              time.Now().Add(time.Hour),
		},
		{
			Name:                 "queue2",
			TargetDirectory:      "/home/user/Downloads",
			MaxParallelDownloads: 3,
			SpeedLimit:           "2MB/s",
			StartTime:            time.Now(),
			EndTime:              time.Now().Add(time.Hour),
		},
	}
}
