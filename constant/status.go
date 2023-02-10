package constant

type Status int

const (
	Stopped Status = iota
	Starting
	Started
	Stopping
)

func (s Status) String() string {
	return [...]string{"stopped", "starting", "started", "stopping", "resizing"}[s]
}