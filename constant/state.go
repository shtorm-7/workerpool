package constant

type State uint8

const (
	Pending State = iota
	Received
	Complete
)

var States = []State{Pending, Received, Complete}

func (s State) String() string {
	return [...]string{"pending", "received", "complete"}[s]
}
