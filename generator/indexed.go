package generator

type IndexedValue[V any] struct {
	Index int
	Value V
}

func Indexed[V any](generator Generator[V]) Generator[IndexedValue[V]] {
	index := 0
	return func() (IndexedValue[V], bool) {
		if value, ok := generator(); ok {
			indexedValue := IndexedValue[V]{index, value}
			index++
			return indexedValue, true
		}
		return IndexedValue[V]{0, *new(V)}, false
	}
}
