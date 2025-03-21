package models

type Download struct {
	ID           int
	URL          string
	Queue        string
	Status       Status
	TransferRate int64
	Progress     float32
}
