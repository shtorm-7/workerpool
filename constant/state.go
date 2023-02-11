package constant

type State uint8

const (
	Pending State = iota
	Received
	Succeeded
)

func (s State) String() string {
	return [...]string{"pending", "received", "succeeded"}[s]
}
