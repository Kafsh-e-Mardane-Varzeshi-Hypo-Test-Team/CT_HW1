package models

type Download struct {
	ID           int
	Url          string
	Queue        string
	Status       string
	TransferRate string
	Progress     float64
}

func GetDownloads() []Download {
	return []Download{
		{
			Url:          "https://example.com",
			Queue:        "1",
			Status:       "Downloading",
			TransferRate: "1MB/s",
			Progress:     0.5,
		},
		{
			Url:          "https://downloadha.com",
			Queue:        "2",
			Status:       "Downloading",
			TransferRate: "2MB/s",
			Progress:     0.8,
		},
		{
			Url:          "https://google.com",
			Queue:        "3",
			Status:       "Downloading",
			TransferRate: "1.5MB/s",
			Progress:     0.2,
		},
		{
			Url:          "https://youtube.com",
			Queue:        "4",
			Status:       "Paused",
			TransferRate: "",
			Progress:     0.1,
		},
		{
			Url:          "https://facebook.com",
			Queue:        "5",
			Status:       "Completed",
			TransferRate: "",
			Progress:     1,
		},
		{
			Url:          "https://instagram.com",
			Queue:        "6",
			Status:       "Error",
			TransferRate: "",
			Progress:     0,
		},
		{
			Url:          "https://twitter.com",
			Queue:        "7",
			Status:       "Downloading",
			TransferRate: "0.5MB/s",
			Progress:     0.7,
		},
	}
}

func AddDownload(url string, filename string, queue string) (string, error) {
	// Add download to the database
	return "Download added successfully", nil
}
