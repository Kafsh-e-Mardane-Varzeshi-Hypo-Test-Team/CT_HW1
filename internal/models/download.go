package models

type Download struct {
	ID           int
	URL          string
	Queue        string
	Status       Status
	TransferRate float64
	Progress     float32
}
