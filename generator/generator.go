package generator

type (
	Next[V any] func() (value V, ok bool)

	Generator[V any] Next[V]
)

func (g Generator[V]) Process(handler func(V)) {
	for value, ok := g(); ok; value, ok = g() {
		handler(value)
	}
}

func (g Generator[V]) ToSlice() []V {
	result := make([]V, 0)
	for value, ok := g(); ok; value, ok = g() {
		result = append(result, value)
	}
	return result
}
