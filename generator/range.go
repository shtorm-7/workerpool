package generator

func Range[V any](values []V) Scheme[V] {
	return func() Generator[V] {
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
}

func ChannelRange[V any](channel <-chan V) Scheme[V] {
	return func() Generator[V] {
		return func() (V, bool) {
			if value, ok := <-channel; ok {
				return value, true
			}
			return *new(V), false
		}
	}
}

func SequenceRange(start, stop int) Scheme[int] {
	return func() Generator[int] {
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
}
