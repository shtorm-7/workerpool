package constant

type State int

const (
	Pending State = iota
	Received
	Succeeded
)

func (s State) String() string {
	return [...]string{"pending", "received", "succeeded"}[s]
}
