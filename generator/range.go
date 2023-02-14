package generator

func Range[V any](values []V) Generator[V] {
	i := 0
	return func() (V, bool) {
		if i != len(values) {
			value := values[i]
			i++
			return value, true
		}
		return *new(V), false
	}
}

func ChannelRange[V any](channel <-chan V) Generator[V] {
	return func() (V, bool) {
		if value, ok := <-channel; ok {
			return value, true
		}
		return *new(V), false
	}
}

func SequenceRange(start, stop int) Generator[int] {
	current := start
	return func() (int, bool) {
		if current < stop {
			value := current
			current++
			return value, true
		}
		return 0, false
	}
}
