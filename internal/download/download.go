package download

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
	QueueName      string
	Destination    string
	OutputFileName string
	Status
}
