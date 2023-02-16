package constant

type Status uint8

const (
	Stopped Status = iota
	Starting
	Started
	Stopping
)

var Statuses = []Status{Stopped, Starting, Started, Stopping}

func (s Status) String() string {
	return [...]string{"stopped", "starting", "started", "stopping"}[s]
}
