package generator

type (
	Generator[V any] func() (value V, ok bool)

	Scheme[V any] func() Generator[V]
)

func (sc Scheme[V]) Process(handler func(V)) {
	generator := sc()
	for value, ok := generator(); ok; value, ok = generator() {
		handler(value)
	}
}

func (sc Scheme[V]) ToSlice() []V {
	generator := sc()
	result := make([]V, 0)
	for value, ok := generator(); ok; value, ok = generator() {
		result = append(result, value)
	}
	return result
}
