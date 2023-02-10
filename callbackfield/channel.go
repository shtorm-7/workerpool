package callbackfield

var ClosedChannel chan struct{} = make(chan struct{})

func init() {
	close(ClosedChannel)
}
