package internal

type Status int

const (
	Pending Status = iota
	InProgress
	Paused
	Cancelled
	Completed
)

type Download struct {
	URL            string
	Destination    string
	OutputFileName string
	Queue
	Status
}
